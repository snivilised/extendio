package nav

import (
	"github.com/samber/lo"
	"github.com/snivilised/extendio/xfs/utils"
)

type navigatorController struct {
	impl  navigatorImpl
	frame *navigationFrame
	ns    *NavigationState
}

func (c *navigatorController) makeFrame() *navigationFrame {

	o := c.impl.options()
	c.frame = &navigationFrame{
		root:        utils.VarProp[string]{},
		currentPath: utils.VarProp[string]{},
		client:      o.Callback,
		raw:         o.Callback,
		notifiers:   notificationsSink{},
		periscope:   &navigationPeriscope{},
		metrics:     navigationMetricsFactory{}.construct(),
	}
	return c.frame
}

func (c *navigatorController) init() {
	c.ns = &NavigationState{Filters: c.frame.filters, Root: &c.frame.root}
}

func (c *navigatorController) Walk(root string) *TraverseResult {
	c.frame.root.Set(root)
	c.frame.notifiers.begin.invoke(c.ns)

	result := c.impl.top(c.frame, root)
	c.frame.notifiers.end.invoke(result)

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

	active := &ActiveState{
		Listen: listen,
	}
	c.frame.save(active)

	state := &persistState{
		Store:  &o.Store,
		Active: active,
	}

	marshaller := (&marshallerFactory{}).create(o, state)
	return marshaller.marshal(path)
}
