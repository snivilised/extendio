package nav

import (
	"github.com/samber/lo"
)

type navigatorController struct {
	impl  navigatorImpl
	frame *navigationFrame
	ns    *NavigationState
}

func (c *navigatorController) makeFrame() *navigationFrame {
	o := c.impl.options()
	c.frame = &navigationFrame{
		client:    o.Callback,
		raw:       o.Callback,
		notifiers: notificationsSink{},
		periscope: &navigationPeriscope{},
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

func (c *navigatorController) Walk(root string) *TraverseResult {
	c.root(func() string {
		return root
	})
	c.frame.notifiers.begin.invoke(c.ns)

	result := &TraverseResult{
		Error: c.impl.top(c.frame, root),
	}
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

func (c *navigatorController) root(fn ...func() string) string {
	if len(fn) == 0 {
		return c.frame.root
	}
	c.ns.Root = fn[0]()
	c.frame.root = c.ns.Root
	return ""
}
