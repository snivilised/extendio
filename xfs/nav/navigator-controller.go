package nav

import (
	"github.com/samber/lo"
	"github.com/snivilised/extendio/internal/log"
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
		metrics:     navigationMetricsFactory{}.new(),
	}

	return c.frame
}

func (c *navigatorController) init() {
	c.ns = &NavigationState{
		Filters: c.frame.filters,
		Root:    &c.frame.root,
		Logger:  utils.NewRoProp[ClientLogger](c.impl.logger()),
	}
}

func (c *navigatorController) logger() log.Logger {
	return c.impl.logger()
}

func (c *navigatorController) walk(root string) (*TraverseResult, error) {
	c.frame.root.Set(root)
	c.impl.logger().Info("walk", log.String("root", root))

	c.frame.notifiers.begin.invoke(c.ns)

	result, err := c.impl.top(c.frame, root)

	fields := []log.Field{}
	for _, m := range result.Metrics.collection {
		fields = append(fields, log.Uint(m.Name, m.Count))
	}

	c.impl.logger().Info("Result", fields...)
	c.frame.notifiers.end.invoke(result)

	return result, err
}

func (c *navigatorController) save(path string) error {
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

	marshaller := (&marshallerFactory{}).new(o, state)

	return marshaller.marshal(path)
}

func (c *navigatorController) finish() error {
	return c.impl.finish()
}
