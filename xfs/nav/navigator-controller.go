package nav

import (
	"github.com/samber/lo"
	_ "github.com/snivilised/extendio/translate"
)

type navigatorController struct {
	impl  navigatorImpl
	frame *navigationFrame
	ns    *NavigationState
}

func (c *navigatorController) init() *navigationFrame {
	o := c.impl.options()
	c.frame = &navigationFrame{
		client: o.Callback,
		raw:    o.Callback,
	}
	return c.frame
}

func (c *navigatorController) navState(fn ...func() *NavigationState) *NavigationState {

	if len(fn) == 0 {
		return c.ns
	}
	c.ns = fn[0]()
	return nil
}

// THIS (resume) MAY NOT BE REQUIRED (well not in its current form)
// Actually, it should probably implement the equivalent of Walk
// without calling top.
//

func (c *navigatorController) resume(ps *persistState, strategy resumeStrategy) *TraverseResult {
	c.root(func() string {
		return ps.Active.Root
	})

	// // this functionality is all bogus
	// if ps.Active.Listen != ListenUndefined {
	// 	c.setRoot(ps.Active.Root)
	// 	initParams := &strategyInitParams{
	// 		state: ps.Active.Listen,
	// 		// listener: c.frame.listener,
	// 		frame: c.frame,
	// 	}

	// 	if ps.Active.Listen == ListenDeaf {
	// 		// TODO: what else do we do here?
	// 		//
	// 		strategy.init(initParams)
	// 	} else {
	// 		// if c.frame.listener == nil {
	// 		// 	panic("navigatorController.resume: 🔥 listener has not been set!")
	// 		// }
	// 		strategy.init(initParams)
	// 		c.frame.listener.transition(ps.Active.Listen)
	// 	}
	// } else {
	// 	panic("navigatorController.resume: 🔥 listen state invalid (undefined)")
	// }

	return &TraverseResult{}
}

func (c *navigatorController) Walk(root string) *TraverseResult {
	o := c.impl.options()

	c.root(func() string {
		return root
	})
	o.Notify.OnBegin(c.ns)

	result := &TraverseResult{
		Error: c.impl.top(c.frame),
	}
	o.Notify.OnEnd(result)

	return result
}

func (c *navigatorController) Save(path string) error {
	o := c.impl.options()

	listen := lo.TernaryF(c.frame.listener == nil,
		func() ListeningState {
			return ListenUndefined
		},
		func() ListeningState {
			return c.frame.listener.state
		},
	)

	state := &persistState{
		Store: &o.Store,
		Active: &ActiveState{
			Root:     c.frame.Root,
			NodePath: c.frame.NodePath,
			Listen:   listen,
		},
	}
	marshaller := newStateMarshaler(o, state)
	return marshaller.marshal(path)
}

func (c *navigatorController) root(fn ...func() string) string {
	if len(fn) == 0 {
		return c.frame.Root
	}
	c.ns.Root = fn[0]()
	c.frame.Root = c.ns.Root
	return ""
}
