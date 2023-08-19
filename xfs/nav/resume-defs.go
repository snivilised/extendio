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
	ps       *persistState
	frame    *navigationFrame
	rc       *resumeController
	triggers *ListenTriggers
}

type strategyResumeInfo struct {
	ps *persistState
	nc *navigatorController
	ai *AsyncInfo
}

type resumeAttachParams struct {
	o        *TraverseOptions
	frame    *navigationFrame
	triggers *ListenTriggers
}

type resumeStrategy interface {
	init(_ *strategyInitParams)
	attach(_ *resumeAttachParams)
	detach(_ *navigationFrame)
	resume(_ *strategyResumeInfo) (*TraverseResult, error)
	finish() error
}

type baseStrategy struct {
	o         *TraverseOptions
	ps        *persistState
	nc        *navigatorController
	deFactory *directoryEntriesFactory
}

func (s *baseStrategy) ensync(ai *AsyncInfo) {
	if ai != nil {
		s.nc.impl.ensync(s.nc.frame, ai)
	}
}

func (s *baseStrategy) attach(_ *resumeAttachParams) {}
func (s *baseStrategy) detach(_ *navigationFrame)    {}
func (s *baseStrategy) finish() error {
	return s.nc.finish()
}

type resumeDetacher interface {
	detach(_ *navigationFrame)
}
