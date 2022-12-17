package nav

type spawnStrategy struct {
	baseStrategy
}

func (s *spawnStrategy) init(params *strategyInitParams) {
	// TODO: set the depth and other appropriate properties on the frame
	//
	params.frame.depth = params.ps.Active.Depth
}

func (s *spawnStrategy) attach(params *resumeAttachParams) {}
func (s *spawnStrategy) detach(frame *navigationFrame)     {}
func (s *spawnStrategy) resume(info *strategyResumeInfo) *TraverseResult {
	// Implementing spawning here
	//
	return nil
}
