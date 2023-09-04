package nav

import (
	"time"

	xi18n "github.com/snivilised/extendio/i18n"
)

type TraverseSession interface {
	Init() NavigationRunner
	Run(ai ...*AsyncInfo) (*TraverseResult, error)
	StartedAt() time.Time
	Elapsed() time.Duration
}

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
	defer s.finish(result, err)

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

func (s *PrimarySession) finish(result *TraverseResult, err error) {
	defer s.session.finish(result, err)

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
	defer s.finish(result, err)

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

func (s *ResumeSession) finish(result *TraverseResult, err error) {
	defer s.session.finish(result, err)

	_ = s.rc.finish()
}
