package nav

import (
	"errors"
	"fmt"
)

type traverseSession interface {
	run() *TraverseResult
}

type NavigationRunner interface {
	Run() *TraverseResult
}

type sessionRunner struct {
	session traverseSession
}

type primaryRunner struct {
	sessionRunner
}

// Run invokes the traversal run for a primary session
func (r *primaryRunner) Run() *TraverseResult {
	fmt.Println("🧊🧊 Primary: Run 🧊🧊")

	return r.session.run()
}

type resumeRunner struct {
	sessionRunner
}

// Run invokes the traversal run for a resume session
func (r *resumeRunner) Run() *TraverseResult {
	fmt.Println("🧊🧊 Resume: Run 🧊🧊")

	return r.session.run()
}

type session struct {
	// Path string not defined here because that would mean the client
	// would have to use go's awkward embedding syntax for specifying literal
	// struct.
}

func (s *session) finish() {
	fmt.Println("🍧🍧 finish 🍧🍧")
}

type PrimarySession struct {
	session
	Path      string
	navigator TraverseNavigator
}

// Configure is the pre run stage for a primary session
func (s *PrimarySession) Configure(fn ...TraverseOptionFn) NavigationRunner {

	s.navigator = navigatorFactory{}.construct(fn...)
	return &primaryRunner{
		sessionRunner: sessionRunner{
			session: s,
		},
	}
}

// Save persists the current state for a primary session, that allows
// a subsequent run to complete the resume.
func (s *PrimarySession) Save(path string) error {
	return s.navigator.Save(path)
}

func (s *PrimarySession) run() *TraverseResult {
	defer s.finish()

	return s.navigator.Walk(s.Path)
}

type ResumeSession struct {
	session
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

	// TODO: resolve the ignored error _
	//
	s.resumer, _ = resumerFactory{}.construct(info)

	return &resumeRunner{
		sessionRunner: sessionRunner{
			session: s,
		},
	}
}

// Save persists the current state for a resume session, that allows
// a subsequent run to complete the resume.
func (s *ResumeSession) Save(path string) error {
	panic(errors.New("ResumeSession.Save: NOT-IMPL"))
}

func (s *ResumeSession) run() *TraverseResult {
	defer s.finish()

	return s.resumer.Continue()
}
