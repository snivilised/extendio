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
	state ListeningState
	frame *navigationFrame
	pci   *preserveClientInfo
}

type Resumer interface {
	Continue() *TraverseResult
	Save(path string) error
}

type resumeStrategy interface {
	init(params *strategyInitParams)
	listenOptions() *ListenOptions
	preservedClient() *preserveClientInfo
}

type baseStrategy struct {
	o      *TraverseOptions
	active *ActiveState
	pci    *preserveClientInfo
}

func (s *baseStrategy) preservedClient() *preserveClientInfo {
	return s.pci
}
