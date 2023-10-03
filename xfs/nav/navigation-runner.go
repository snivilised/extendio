package nav

import (
	"fmt"
	"runtime"

	"github.com/samber/lo"
	"github.com/snivilised/lorax/boost"
)

const (
	MinNoWorkers = 1
	MaxNoWorkers = 100
)

func CreateTraverseOutputCh(outputChSize int) boost.OutputStream[TraverseOutput] {
	return lo.TernaryF(outputChSize > 0,
		func() boost.OutputStream[TraverseOutput] {
			return make(boost.OutputStream[TraverseOutput], outputChSize)
		},
		func() boost.OutputStream[TraverseOutput] {
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
	Consume(outputCh boost.OutputStream[TraverseOutput]) AccelerationOperators
}

type SessionRunner interface {
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

func (r *runner) Primary(info *Prime) NavigationRunner {
	r.session = &Primary{
		Path:     info.Path,
		OptionFn: info.OptionsFn,
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

func (r *runner) Consume(outputCh boost.OutputStream[TraverseOutput]) AccelerationOperators {
	r.sync.outputChOut = outputCh

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
