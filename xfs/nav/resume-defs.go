package nav

type ResumeStrategyEnum uint

// If these enum definitions change, the data (eg, resume-fastward.json) also needs
// to be updated.

const (
	ResumeStrategyUndefinedEn ResumeStrategyEnum = iota
	ResumeStrategySpawnEn
	ResumeStrategyFastwardEn
)

type listenerInitParams struct {
	o        *TraverseOptions
	state    ListeningState
	listener *navigationListener
	frame    *navigationFrame
}

type Resumer interface {
	Walk() *TraverseResult
	Save(path string) error
}

type resumeStrategy interface {
	init(params *listenerInitParams)
}
