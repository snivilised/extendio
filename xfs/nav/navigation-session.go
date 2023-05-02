package nav

import (
	xi18n "github.com/snivilised/extendio/i18n"
)

type TraverseSession interface {
	Init() NavigationRunner
	Run() (*TraverseResult, error)
}

type NavigationRunner interface {
	Run() (*TraverseResult, error)
}

type sessionRunner struct {
	session TraverseSession
}

// Run executes the traversal session
func (r *sessionRunner) Run() (*TraverseResult, error) {
	return r.session.Run()
}

type primaryRunner struct {
	sessionRunner
}

type resumeRunner struct {
	sessionRunner
}

// PrimarySession
type PrimarySession struct {
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

func (s *PrimarySession) Run() (*TraverseResult, error) {
	defer s.finish()

	return s.navigator.walk(s.Path)
}

func (s *PrimarySession) finish() {
	_ = s.navigator.finish()
}

// ResumeSession represents a traversal that is invoked as a result
// of the user needing to resume a previously interrupted navigation
// session.
type ResumeSession struct {
	Path     string
	Restorer func(o *TraverseOptions, active *ActiveState)
	Strategy ResumeStrategyEnum
	resumer  *resumeController
}

func (s *ResumeSession) Init() NavigationRunner {
	var err error

	s.resumer, err = resumerFactory{}.new(&ResumerInfo{
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

	s.resumer, err = resumerFactory{}.new(&ResumerInfo{
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
	return s.resumer.navigator.save(path)
}

func (s *ResumeSession) Run() (*TraverseResult, error) {
	defer s.finish()

	return s.resumer.Continue()
}

func (s *ResumeSession) finish() {
	_ = s.resumer.finish()
}
