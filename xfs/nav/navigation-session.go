package nav

import (
	xi18n "github.com/snivilised/extendio/i18n"
)

type traverseSession interface {
	run() (*TraverseResult, error)
}

type NavigationRunner interface {
	Run() (*TraverseResult, error)
}

type sessionRunner struct {
	session traverseSession
}

type primaryRunner struct {
	sessionRunner
}

// Run invokes the traversal run for a primary session
func (r *primaryRunner) Run() (*TraverseResult, error) {
	return r.session.run()
}

type resumeRunner struct {
	sessionRunner
}

// Run invokes the traversal run for a resume session
func (r *resumeRunner) Run() (*TraverseResult, error) {
	return r.session.run()
}

type PrimarySession struct {
	Path      string
	navigator TraverseNavigator
}

// Configure is the pre run stage for a primary session
func (s *PrimarySession) Configure(fn ...TraverseOptionFn) NavigationRunner {
	s.navigator = navigatorFactory{}.new(fn...)

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

func (s *PrimarySession) run() (*TraverseResult, error) {
	defer s.finish()

	return s.navigator.walk(s.Path)
}

func (s *PrimarySession) finish() {
	_ = s.navigator.finish()
}

type ResumeSession struct {
	Path     string
	Strategy ResumeStrategyEnum
	resumer  *resumeController
}

// Configure is the pre run stage for a resume session
func (s *ResumeSession) Configure(restore func(o *TraverseOptions, active *ActiveState)) NavigationRunner {
	info := &ResumerInfo{
		RestorePath: s.Path,
		Restorer:    restore,
		Strategy:    s.Strategy,
	}

	var err error
	s.resumer, err = resumerFactory{}.new(info)

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

func (s *ResumeSession) run() (*TraverseResult, error) {
	defer s.finish()

	return s.resumer.Continue()
}

func (s *ResumeSession) finish() {
	_ = s.resumer.finish()
}
