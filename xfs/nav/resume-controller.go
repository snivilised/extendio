package nav

type resumeController struct {
	navigator *navigatorController
	ps        *persistState
	strategy  resumeStrategy
}

func (c *resumeController) init() {
	c.navigator.resume(c.ps, c.strategy)
}

func (c *resumeController) Walk() *TraverseResult {
	return c.navigator.Walk(c.navigator.frame.Root)
}

func (c *resumeController) Save(path string) error {
	return c.navigator.Save(path)
}
