package nav

type resumeController struct {
	navigator *navigatorController
	ps        *persistState
	strategy  resumeStrategy
}

func (c *resumeController) Continue() *TraverseResult {
	// TODO: can't call walk here, resume instead
	// calling top for resume is not a problem, in fact it should, but that is ...
	// strategy dependent ...
	//
	return c.navigator.resume(c.ps, c.strategy)
}

func (c *resumeController) Save(path string) error {
	return c.navigator.Save(path)
}

// func (c *resumeController) detach() *ListenOptions {
// 	o := c.navigator.impl.options()
// 	frame := c.navigator.frame
// 	pci := c.strategy.preservedClient()

// 	o.Listen = *pci.lo
// 	o.Store.Behaviours.Listen = pci.behaviours
// 	o.Notify = pci.notify

// 	frame.listener.composeListenStates(&listenStatesParams{
// 		decorated: &frame.raw, o: o, frame: frame,
// 	})
// 	frame.listener.transition(c.ps.Active.Listen)

// 	detached := frame.listener.detach()

// 	return detached
// }
