package nav

import (
	"time"

	xi18n "github.com/snivilised/extendio/i18n"
)

type TraverseSession interface {
	StartedAt() time.Time
	Elapsed() time.Duration
	run(sync NavigationSync, args ...any) (*TraverseResult, error)
	Save(path string) error
}

type sessionCallback func() (*TraverseResult, error)

type session struct {
	startAt  time.Time
	duration time.Duration
}

func (s *session) start() {
	s.startAt = time.Now()
}

func (s *session) finish(_ *TraverseResult, _ error) {
	s.duration = time.Since(s.startAt)
}

// Primary
type Primary struct {
	session
	Path      string
	OptionFn  TraverseOptionFn
	navigator TraverseNavigator
}

func (s *Primary) init() {
	s.navigator = navigatorFactory{}.new(s.OptionFn)
}

// Save persists the current state for a primary session, that allows
// a subsequent run to complete the resume.
func (s *Primary) Save(path string) error {
	return s.navigator.save(path)
}

func (s *Primary) run(sync NavigationSync, args ...any) (result *TraverseResult, err error) {
	defer s.finish(result, err)

	s.start()
	s.init()

	return sync.Run(
		func() (*TraverseResult, error) {
			return s.navigator.walk(s.Path)
		},
		s.navigator,
		args...,
	)
}

func (s *Primary) StartedAt() time.Time {
	return s.startAt
}

func (s *Primary) Elapsed() time.Duration {
	return s.duration
}

func (s *Primary) finish(result *TraverseResult, err error) {
	defer s.session.finish(result, err)

	if s.navigator != nil {
		_ = s.navigator.finish()
	}
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

// Save persists the current state for a resume session, that allows
// a subsequent run to complete the resume.
func (s *Resume) Save(path string) error {
	return s.rsc.nc.save(path)
}

func (s *Resume) run(sync NavigationSync, args ...any) (result *TraverseResult, err error) {
	defer s.finish(result, err)

	s.init()
	s.start()

	return sync.Run(
		func() (*TraverseResult, error) {
			return s.rsc.run()
		},
		s.rsc.nc,
		args...,
	)
}

func (s *Resume) StartedAt() time.Time {
	return s.startAt
}

func (s *Resume) Elapsed() time.Duration {
	return s.duration
}

func (s *Resume) finish(result *TraverseResult, err error) {
	defer s.session.finish(result, err)

	_ = s.rsc.finish()
}
