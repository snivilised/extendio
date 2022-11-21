package nav

import (
	"github.com/samber/lo"
	_ "github.com/snivilised/extendio/translate"
)

type navigatorController struct {
	impl  navigatorImpl
	frame *navigationFrame
	ns    NavigationState
}

func (c *navigatorController) init() {
	o := c.impl.options()
	c.frame = &navigationFrame{
		client: o.Callback,
	}
	bootstrapFilter(o, c.frame)
	bootstrapListener(o, c.frame)

	c.ns = NavigationState{Filters: c.frame.filters}
}

func (c *navigatorController) resume(ps *persistState, strategy resumeStrategy) {
	if ps.Active.Listen != ListenUndefined {
		c.setRoot(ps.Active.Root)
		initParams := &listenerInitParams{
			o:     c.impl.options(),
			state: ps.Active.Listen,
			// listener: c.frame.listener,
			frame: c.frame,
		}

		if ps.Active.Listen == ListenDefault {
			// TODO: what else do we do here?
			//
			strategy.init(initParams)
		} else {
			// if c.frame.listener == nil {
			// 	panic("navigatorController.resume: 🔥 listener has not been set!")
			// }
			strategy.init(initParams)
			c.frame.listener.transition(ps.Active.Listen)
		}
	} else {
		panic("navigatorController.resume: 🔥 listen state invalid (undefined)")
	}
}

func (c *navigatorController) Walk(root string) *TraverseResult {
	o := c.impl.options()

	c.setRoot(root)
	o.Notify.OnBegin(&c.ns)

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
			Root:   c.frame.Root,
			Listen: listen,
		},
	}
	marshaller := newStateMarshaler(o, state)
	return marshaller.marshal(path)
}

func (c *navigatorController) setRoot(root string) {
	c.ns.Root = root
	c.frame.Root = root
}
