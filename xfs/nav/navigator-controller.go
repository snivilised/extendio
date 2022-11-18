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

func (c *navigatorController) init(root string) {
	o := c.impl.options()
	c.frame = &navigationFrame{
		Root:   root,
		client: o.Callback,
	}

	bootstrapFilter(o, c.frame)
	bootstrapListener(o, c.frame)

	c.ns = NavigationState{Root: root, Filters: c.frame.filters}
	o.Notify.OnBegin(&c.ns)
}

func (c *navigatorController) resume(ps *persistState, strategy resumeStrategy) {

	if ps.Active.Listen != ListenUndefined {
		if ps.Active.Listen == ListenDefault {
			// TODO: what do we do here?
			//
		} else {
			if c.frame.listener == nil {
				panic("navigatorController: 🔥 listener has not been set!")
			}
			// this state transition is the default, but can be overridden
			// by the resume initialiser
			//
			c.frame.listener.transition(ps.Active.Listen)
			// TODO: don't take this seriously, its speculative at the
			// moment and wont be fixed in this issue (#59)
			//
			strategy.init(&listenerInitParams{
				state:    ps.Active.Listen,
				listener: c.frame.listener,
			})
		}
	}
}

func (c *navigatorController) Walk(root string) *TraverseResult {
	o := c.impl.options()

	if c.frame == nil {
		c.init(root)
	}

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
		Active: &activeState{
			Root:   c.frame.Root,
			Listen: listen,
		},
	}
	marshaller := newStateMarshaler(o, state)
	return marshaller.marshal(path)
}
