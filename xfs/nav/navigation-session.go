package nav

import (
	"errors"
	"time"

	xi18n "github.com/snivilised/extendio/i18n"
)

type Session interface {
	StartedAt() time.Time
	Elapsed() time.Duration
}

type TraverseSession interface {
	Session
	run(sync NavigationSync, args ...any) (*TraverseResult, error)
	Save(path string) error
}

type sessionCallback func() (*TraverseResult, error)

type session struct {
	startAt  time.Time
	duration time.Duration
}

func (s *session) StartedAt() time.Time {
	return s.startAt
}

func (s *session) Elapsed() time.Duration {
	return s.duration
}

func (s *session) start() {
	s.startAt = time.Now()
}

func (s *session) finish(result *TraverseResult, _ error) {
	s.duration = time.Since(s.startAt)

	if result != nil {
		result.Session = s
	}
}

// Primary
type Primary struct {
	session
	Path            string
	OptionFn        TraverseOptionFn
	ProvidedOptions *TraverseOptions
	navigator       TraverseNavigator
}

// Save persists the current state for a primary session, that allows
// a subsequent run to complete the resume.
func (s *Primary) Save(path string) error {
	return s.navigator.save(path)
}

func (s *Primary) init() {
	switch {
	case s.OptionFn != nil:
		s.navigator = navigatorFactory{}.fromOptionsFn(s.OptionFn)

	case s.ProvidedOptions != nil:
		s.navigator = navigatorFactory{}.fromProvidedOptions(s.ProvidedOptions)

	default:
		panic(errors.New("missing traverse options"))
	}
}

func (s *Primary) run(sync NavigationSync, args ...any) (*TraverseResult, error) {
	s.start()
	s.init()

	result, err := sync.Run(
		func() (*TraverseResult, error) {
			return s.navigator.walk(s.Path)
		},
		s.navigator,
		args...,
	)

	s.finish(result, err)

	return result, err
}

func (s *Primary) finish(result *TraverseResult, err error) {
	if s.navigator != nil {
		_ = s.navigator.finish()
	}

	s.session.finish(result, err)
}

// Resume represents a traversal that is invoked as a result
// of the user needing to resume a previously interrupted navigation
// session.
type Resume struct {
	session
	RestorePath string
	Restorer    func(o *TraverseOptions, active *ActiveState)
	Strategy    ResumeStrategyEnum
	rsc         *resumeStrategyController
}

// Save persists the current state for a resume session, that allows
// a subsequent run to complete the resume.
func (s *Resume) Save(path string) error {
	return s.rsc.nc.save(path)
}

func (s *Resume) init() {
	var err error

	s.rsc, err = resumerFactory{}.new(&Resumption{
		RestorePath: s.RestorePath,
		Restorer:    s.Restorer,
		Strategy:    s.Strategy,
	})

	if err != nil {
		panic(xi18n.NewFailedToResumeFromFileError(s.RestorePath, err))
	}
}

func (s *Resume) run(sync NavigationSync, args ...any) (*TraverseResult, error) {
	s.start()
	s.init()

	result, err := sync.Run(
		func() (*TraverseResult, error) {
			return s.rsc.run()
		},
		s.rsc.nc,
		args...,
	)

	s.finish(result, err)

	return result, err
}

func (s *Resume) finish(result *TraverseResult, err error) {
	if s.rsc != nil {
		_ = s.rsc.finish()
	}

	s.session.finish(result, err)
}
