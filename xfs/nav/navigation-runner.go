package nav

import (
	"fmt"
	"runtime"
	"time"

	"github.com/snivilised/extendio/internal/lo"
	"github.com/snivilised/lorax/boost"
)

type CreateNewRunnerWith int

const (
	RunnerDefault    CreateNewRunnerWith = 0
	RunnerWithResume CreateNewRunnerWith = 1
	RunnerWithPool   CreateNewRunnerWith = 2
)

type Acceleration struct {
	WgAn            boost.WaitGroupAn
	RoutineName     boost.GoRoutineName
	NoW             int
	JobsChOut       TraverseItemJobStream
	JobResultsCh    boost.JobOutputStream[TraverseOutput]
	OutputChTimeout time.Duration
}

type RunnerInfo struct {
	ResumeInfo       *Resumption
	PrimeInfo        *Prime
	AccelerationInfo *Acceleration
}

const (
	MinNoWorkers = 1
	MaxNoWorkers = 100
)

func CreateTraverseOutputCh(outputChSize int) boost.JobOutputStream[TraverseOutput] {
	return lo.TernaryF(outputChSize > 0,
		func() boost.JobOutputStream[TraverseOutput] {
			return make(boost.JobOutputStream[TraverseOutput], outputChSize)
		},
		func() boost.JobOutputStream[TraverseOutput] {
			return nil
		},
	)
}

type Runnable interface {
	Run(args ...any) (*TraverseResult, error)
}

type AccelerationOperators interface {
	Runnable
	NoW(now int) AccelerationOperators
	Consume(outputCh boost.JobOutputStream[TraverseOutput], timeout time.Duration) AccelerationOperators
}

type SessionRunner interface {
	With(with CreateNewRunnerWith, info *RunnerInfo) NavigationRunner
	Primary(info *Prime) NavigationRunner
	Resume(info *Resumption) NavigationRunner
}

type NavigationRunner interface {
	AccelerationOperators
	WithPool(ai *AsyncInfo) AccelerationOperators
	Save(path string) error
}

func New() SessionRunner {
	return &runner{}
}

type runner struct {
	session TraverseSession
	sync    *acceleratedSync
}

func (r *runner) With(with CreateNewRunnerWith, info *RunnerInfo) NavigationRunner {
	lo.TernaryF(with&RunnerWithResume == 0,
		func() NavigationRunner {
			return r.Primary(&Prime{
				Path:            info.PrimeInfo.Path,
				OptionsFn:       info.PrimeInfo.OptionsFn,
				ProvidedOptions: info.PrimeInfo.ProvidedOptions,
			})
		},
		func() NavigationRunner {
			return r.Resume(&Resumption{
				RestorePath: info.ResumeInfo.RestorePath,
				Restorer:    info.ResumeInfo.Restorer,
				Strategy:    info.ResumeInfo.Strategy,
			})
		},
	)

	lo.TernaryF(with&RunnerWithPool == 0,
		func() AccelerationOperators {
			return r
		},
		func() AccelerationOperators {
			if info.AccelerationInfo == nil {
				// As this is not a user facing issue (ie programming error),
				// it does not have to be i18n error
				//
				panic("internal: acceleration info missing from runner info")
			}

			if info.AccelerationInfo.JobsChOut == nil {
				panic("internal: job channel not set on acceleration info")
			}

			return r.WithPool(&AsyncInfo{
				NavigatorRoutineName: info.AccelerationInfo.RoutineName,
				WaitAQ:               info.AccelerationInfo.WgAn,
				JobsChanOut:          info.AccelerationInfo.JobsChOut,
			})
		},
	)

	if info.AccelerationInfo != nil && with&RunnerWithPool > 0 {
		if info.AccelerationInfo.JobResultsCh != nil {
			r.Consume(info.AccelerationInfo.JobResultsCh, info.AccelerationInfo.OutputChTimeout)
		}

		if info.AccelerationInfo.NoW > 0 {
			r.NoW(info.AccelerationInfo.NoW)
		}
	}

	return r
}

func IfWithPoolUseContext(with CreateNewRunnerWith, args ...any) []any {
	return lo.TernaryF(with&RunnerWithPool > 0,
		func() []any {
			return args
		},
		func() []any {
			return []any{}
		},
	)
}

func (r *runner) Primary(info *Prime) NavigationRunner {
	r.session = &Primary{
		Path:            info.Path,
		OptionFn:        info.OptionsFn,
		ProvidedOptions: info.ProvidedOptions,
	}

	return r
}

func (r *runner) Resume(info *Resumption) NavigationRunner {
	r.session = &Resume{
		RestorePath: info.RestorePath,
		Restorer:    info.Restorer,
		Strategy:    info.Strategy,
	}

	return r
}

func (r *runner) WithPool(ai *AsyncInfo) AccelerationOperators {
	r.sync = &acceleratedSync{
		ai:        ai,
		noWorkers: runtime.NumCPU(),
	}

	return r
}

func (r *runner) Save(path string) error {
	return r.session.Save(path)
}

func (r *runner) NoW(now int) AccelerationOperators {
	if now >= MinNoWorkers && now <= MaxNoWorkers {
		r.sync.noWorkers = now
	} else {
		// TODO: turn this into an i18n error
		panic(fmt.Errorf("no of workers requested (%v) is out of range ('%v' - '%v')",
			now, MinNoWorkers, MaxNoWorkers),
		)
	}

	return r
}

func (r *runner) Consume(
	outputCh boost.JobOutputStream[TraverseOutput],
	timeout time.Duration,
) AccelerationOperators {
	r.sync.outputChOut = outputCh
	r.sync.outputChTimeout = timeout

	return r
}

func (r *runner) Run(args ...any) (*TraverseResult, error) {
	sync := lo.TernaryF(r.sync == nil,
		func() NavigationSync {
			return &inlineSync{}
		},
		func() NavigationSync {
			return r.sync
		},
	)

	return r.session.run(sync, args...)
}
