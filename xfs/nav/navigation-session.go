package nav

import (
	"fmt"
	"runtime"
	"time"

	"github.com/samber/lo"
	xi18n "github.com/snivilised/extendio/i18n"
	"github.com/snivilised/lorax/boost"
)

const (
	MinNoWorkers = 1
	MaxNoWorkers = 100
)

type TraverseSession interface {
	Init() NavigationRunner
	Run(ai ...*AsyncInfo) (*TraverseResult, error)
	StartedAt() time.Time
	Elapsed() time.Duration
}

type RunnerOperators interface {
	WithCPUPool() NavigationRunner
	WithPool(now int) NavigationRunner
	Consume(outputCh boost.OutputStream[TraverseOutput]) NavigationRunner
}

type NavigationRunner interface {
	RunnerOperators
	Run(ai ...*AsyncInfo) (*TraverseResult, error)
}

type sessionRunner struct {
	session     TraverseSession
	accelerator navigationAccelerator
}

// Run executes the traversal session
func (r *sessionRunner) Run(ai ...*AsyncInfo) (*TraverseResult, error) {
	if r.accelerator.active {
		r.accelerator.start(ai[0])
		return r.session.Run(ai[0])
	}

	return r.session.Run()
}

func (r *sessionRunner) WithPool(now int) NavigationRunner {
	r.accelerator.active = true
	if now >= MinNoWorkers && now <= MaxNoWorkers {
		r.accelerator.noWorkers = now
	} else {
		panic(fmt.Errorf("no of workers requested (%v) is out of range ('%v' - '%v')",
			now, MinNoWorkers, MaxNoWorkers),
		)
	}

	return r
}

func (r *sessionRunner) WithCPUPool() NavigationRunner {
	r.accelerator.active = true
	r.accelerator.noWorkers = runtime.NumCPU()

	return r
}

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

func (r *sessionRunner) Consume(outputCh boost.OutputStream[TraverseOutput]) NavigationRunner {
	if !r.accelerator.active {
		panic(fmt.Errorf(
			"worker pool acceleration not active; ensure With(CPU)Pool specified before Consume",
		))
	}

	r.accelerator.outputChOut = outputCh

	return r
}

type primaryRunner struct {
	sessionRunner
}

type resumeRunner struct {
	sessionRunner
}

type session struct {
	startAt  time.Time
	duration time.Duration
}

func (s *session) start() {
	s.startAt = time.Now()
}

func (s *session) finish(_ *TraverseResult, _ error, ai ...*AsyncInfo) {
	defer func() {
		if len(ai) > 0 {
			fmt.Printf("---> observable navigator 😈😈😈 defer session.finish (CLOSE(JobsChanOut)/QUIT)\n")
			close(ai[0].JobsChanOut) // ⚠️ fastward: intermittent panic on close
			ai[0].WaitAQ.Done(ai[0].NavigatorRoutineName)
		}
	}()

	s.duration = time.Since(s.startAt)
}

// PrimarySession
type PrimarySession struct {
	session
	Path      string
	OptionFn  TraverseOptionFn
	navigator TraverseNavigator
}

func (s *PrimarySession) Init() NavigationRunner {
	s.navigator = navigatorFactory{}.new(s.OptionFn)

	return &primaryRunner{
		sessionRunner: sessionRunner{
			session: s,
		},
	}
}

// Save persists the current state for a primary session, that allows
// a subsequent run to complete the resume.
func (s *PrimarySession) Save(path string) error {
	return s.navigator.save(path)
}

func (s *PrimarySession) Run(ai ...*AsyncInfo) (result *TraverseResult, err error) {
	defer s.finish(result, err, ai...)

	s.session.start()

	if len(ai) > 0 {
		s.navigator.ensync(ai[0])
	}

	return s.navigator.walk(s.Path)
}

func (s *PrimarySession) StartedAt() time.Time {
	return s.startAt
}

func (s *PrimarySession) Elapsed() time.Duration {
	return s.duration
}

func (s *PrimarySession) finish(result *TraverseResult, err error, ai ...*AsyncInfo) {
	defer s.session.finish(result, err, ai...)

	_ = s.navigator.finish()
}

// ResumeSession represents a traversal that is invoked as a result
// of the user needing to resume a previously interrupted navigation
// session.
type ResumeSession struct {
	session
	Path     string
	Restorer func(o *TraverseOptions, active *ActiveState)
	Strategy ResumeStrategyEnum
	rc       *resumeController
}

func (s *ResumeSession) Init() NavigationRunner {
	var err error

	s.rc, err = resumerFactory{}.new(&ResumerInfo{
		RestorePath: s.Path,
		Restorer:    s.Restorer,
		Strategy:    s.Strategy,
	})

	if err != nil {
		panic(xi18n.NewFailedToResumeFromFileError(s.Path, err))
	}

	return &resumeRunner{
		sessionRunner: sessionRunner{
			session: s,
		},
	}
}

// Restore is the pre run stage for a resume session
func (s *ResumeSession) Restore(restore func(o *TraverseOptions, active *ActiveState)) NavigationRunner {
	var err error

	s.rc, err = resumerFactory{}.new(&ResumerInfo{
		RestorePath: s.Path,
		Restorer:    restore,
		Strategy:    s.Strategy,
	})

	if err != nil {
		panic(xi18n.NewFailedToResumeFromFileError(s.Path, err))
	}

	return &resumeRunner{
		sessionRunner: sessionRunner{
			session: s,
		},
	}
}

// Save persists the current state for a resume session, that allows
// a subsequent run to complete the resume.
func (s *ResumeSession) Save(path string) error {
	return s.rc.navigator.save(path)
}

func (s *ResumeSession) Run(ai ...*AsyncInfo) (result *TraverseResult, err error) {
	defer s.finish(result, err, ai...)

	s.session.start()

	if len(ai) > 0 {
		s.rc.navigator.ensync(ai[0])
	}

	return s.rc.Continue(ai...)
}

func (s *ResumeSession) StartedAt() time.Time {
	return s.startAt
}

func (s *ResumeSession) Elapsed() time.Duration {
	return s.duration
}

func (s *ResumeSession) finish(result *TraverseResult, err error, ai ...*AsyncInfo) {
	defer s.session.finish(result, err, ai...)

	_ = s.rc.finish()
}
