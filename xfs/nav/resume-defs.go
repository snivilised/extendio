package nav

type ResumeStrategyEnum uint

const (
	_ ResumeStrategyEnum = iota
	ResumeStrategySpawnEn
	ResumeStrategyFastwardEn
)

type listenerInitParams struct {
	state    ListeningState
	listener *navigationListener
}

type resumeInit func(params *listenerInitParams)

type Resumer interface {
	Walk() *TraverseResult
	Save(path string) error
}
