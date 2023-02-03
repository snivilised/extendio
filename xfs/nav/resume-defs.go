package nav

type ResumeStrategyEnum uint

// If these enum definitions change, the test data (eg, resume-fastward.json) also needs
// to be updated.

const (
	ResumeStrategyUndefinedEn ResumeStrategyEnum = iota
	ResumeStrategySpawnEn
	ResumeStrategyFastwardEn
)

type strategyInitParams struct {
	ps    *persistState
	frame *navigationFrame
	rc    *resumeController
}

type resumer interface {
	Continue() *TraverseResult
	Save(path string) error
}

type strategyResumeInfo struct {
	ps *persistState
	nc *navigatorController
}

type resumeAttachParams struct {
	o     *TraverseOptions
	frame *navigationFrame
	lo    *ListenOptions
}

type resumeStrategy interface {
	init(params *strategyInitParams)
	attach(params *resumeAttachParams)
	detach(frame *navigationFrame)
	resume(info *strategyResumeInfo) *TraverseResult
}

type baseStrategy struct {
	o         *TraverseOptions
	ps        *persistState
	nc        *navigatorController
	deFactory *directoryEntriesFactory
}

func (s *baseStrategy) attach(params *resumeAttachParams) {}
func (s *baseStrategy) detach(frame *navigationFrame)     {}

type resumeDetacher interface {
	detach(frame *navigationFrame)
}
