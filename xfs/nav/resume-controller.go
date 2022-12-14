package nav

type resumeController struct {
	navigator *navigatorController
	ps        *persistState
	strategy  resumeStrategy
}

func (c *resumeController) Continue() *TraverseResult {

	return c.strategy.resume(&strategyResumeInfo{
		nc: c.navigator,
		ps: c.ps,
	})
}

func (c *resumeController) Save(path string) error {
	return c.navigator.Save(path)
}

func (c *resumeController) detach(frame *navigationFrame) {
	c.strategy.detach(frame)
}
