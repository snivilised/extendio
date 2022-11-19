package nav

type ResumeStrategyEnum uint

const (
	ResumeStrategyUndefinedEn ResumeStrategyEnum = iota
	ResumeStrategySpawnEn
	ResumeStrategyFastwardEn
)

type listenerInitParams struct {
	state    ListeningState
	listener *navigationListener
}

type Resumer interface {
	Walk() *TraverseResult
	Save(path string) error
}

type resumeStrategy interface {
	init(params *listenerInitParams)
}

type dummyResumeStrategy struct {
	// to be replaced with spawn/fastward
}

func (s *dummyResumeStrategy) init(params *listenerInitParams) {

}
