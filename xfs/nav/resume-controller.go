package nav

type resumeController struct {
	navigator *navigatorController
	ps        *persistState
}

func (c *resumeController) init(initialiser resumeInit) {
	c.navigator.resume(c.ps, initialiser)
}

func (c *resumeController) Walk() *TraverseResult {
	if c.navigator.frame.Root == "" {
		panic("resumeController:Walk 'Root' not set")
	}
	return c.navigator.Walk(c.navigator.frame.Root)
}

func (c *resumeController) Save(path string) error {
	return c.navigator.Save(path)
}
