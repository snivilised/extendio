package nav

import (
	"github.com/samber/lo"
	_ "github.com/snivilised/extendio/translate"
)

type navigatorController struct {
	impl  navigatorImpl
	frame navigationFrame
	ns    NavigationState
}

func (c *navigatorController) Walk(root string) *TraverseResult {
	o := c.impl.options()
	c.frame = navigationFrame{
		Root:   root,
		client: o.Callback,
	}

	bootstrapFilter(o, &c.frame)
	bootstrapListener(o, &c.frame)

	c.ns = NavigationState{Root: root, Filters: c.frame.filters}
	o.Notify.OnBegin(&c.ns)

	result := &TraverseResult{
		Error: c.impl.top(&c.frame),
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
			return c.frame.listener.listen
		},
	)

	state := &persistState{
		Store: &o.Store,
		Active: activeState{
			Listen: listen,
		},
	}
	marshaller := newStateMarshaler(o, state)
	return marshaller.marshal(path)
}
