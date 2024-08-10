package nav

import (
	"context"
	"log/slog"

	"github.com/snivilised/extendio/internal/lo"
	"github.com/snivilised/extendio/xfs/utils"
)

type navigationController struct {
	impl  navigatorImpl
	frame *navigationFrame
	ns    *NavigationState
}

func (nc *navigationController) makeFrame() *navigationFrame {
	o := nc.impl.options()
	nc.frame = &navigationFrame{
		root:        utils.VarProp[string]{},
		currentPath: utils.VarProp[string]{},
		client:      o.Callback,
		raw:         o.Callback,
		notifiers:   notificationsSink{},
		periscope:   &navigationPeriscope{},
		metrics:     navigationMetricsFactory{}.new(),
	}

	return nc.frame
}

func (nc *navigationController) init() {
	nc.ns = &NavigationState{
		Filters: nc.frame.filters,
		Root:    &nc.frame.root,
		Logger:  nc.logger(),
	}
	nc.impl.init(nc.ns)
}

func (nc *navigationController) logger() *slog.Logger {
	return nc.impl.logger()
}

func (nc *navigationController) ensync(ctx context.Context, cancel context.CancelFunc, ai *AsyncInfo) {
	nc.impl.ensync(ctx, cancel, nc.frame, ai)
}

func (nc *navigationController) walk(root string) (*TraverseResult, error) {
	nc.frame.root.Set(root)
	nc.impl.logger().Info("walk", slog.String("root", root))

	nc.frame.notifiers.begin.invoke(nc.ns)

	result, err := nc.impl.top(nc.frame, root)

	fields := []any{}
	for _, m := range result.Metrics.collection {
		fields = append(fields, slog.Int(m.Name, int(m.Count)))
	}

	nc.impl.logger().Info("Result", fields...)
	nc.frame.notifiers.end.invoke(result)

	return result, err
}

func (nc *navigationController) save(path string) error {
	o := nc.impl.options()

	listen := lo.TernaryF(nc.frame.listener == nil,
		func() ListeningState {
			return ListenUndefined
		},
		func() ListeningState {
			return nc.frame.listener.state
		},
	)

	active := &ActiveState{
		Listen: listen,
	}
	nc.frame.save(active)

	state := &persistState{
		Store:  &o.Store,
		Active: active,
	}

	marshaller := (&marshallerFactory{}).new(o, state)

	return marshaller.marshal(path)
}

func (nc *navigationController) finish() error {
	return nc.impl.finish()
}
