package nav

import (
	"errors"

	"github.com/snivilised/extendio/collections"
)

type nullDetacher struct{}

func (d *nullDetacher) detach(frame *navigationFrame) {

}

type bootstrapper struct {
	o        *TraverseOptions
	nc       *navigatorController
	rc       *resumeController
	detacher resumeDetacher
}

func (b *bootstrapper) init() {
	b.detacher = &nullDetacher{}

	b.nc.frame = b.nc.makeFrame()
	b.initFilters()
	b.initNotifiers()
	b.initListener()
	b.nc.init()
	b.nc.ns = &NavigationState{Filters: b.nc.frame.filters, Root: &b.nc.frame.root}
}

func (b *bootstrapper) initFilters() {
	b.o.Hooks.InitFilters(b.o, b.nc.frame)
}

func (b *bootstrapper) initNotifiers() {
	if b.o.Notify.OnStart == nil {
		b.o.Notify.OnStart = func(description string) {}
	}

	if b.o.Notify.OnStop == nil {
		b.o.Notify.OnStop = func(description string) {}
	}

	b.nc.frame.notifiers.init(&b.o.Notify)
}

func (b *bootstrapper) initListener() {
	initialState := backfill(&b.o.Listen)

	b.nc.frame.listener = &navigationListener{
		state:       initialState,
		resumeStack: collections.NewStack[*ListenOptions](),
	}

	b.nc.frame.listener.makeStates(&listenStatesParams{
		o: b.o, frame: b.nc.frame,
		detacher: b,
	})

	b.nc.frame.listener.decorate(&listenStatesParams{
		lo: &b.o.Listen, frame: b.nc.frame,
	})
}

type preserveClientInfo struct {
	lo         *ListenOptions
	behaviours ListenBehaviour
}

type overrideClientInfo struct {
	lo *ListenOptions
}

type overrideListenerInfo struct {
	client   *preserveClientInfo
	override *overrideClientInfo
	ps       *persistState
}

func (b *bootstrapper) initResume(o *TraverseOptions, ps *persistState) {

	if b.rc == nil {
		b.nc.impl.logger().Error("bootstrapper.initResume: resume controller not set")
		panic(errors.New("bootstrapper.initResume: resume controller not set"))
	}

	strategyParams := &strategyInitParams{
		ps:    ps,
		frame: b.nc.frame,
		rc:    b.rc,
	}
	b.nc.frame.metrics.load(ps.Active)

	b.rc.strategy.init(strategyParams)
	b.detacher = b.rc
}

func (b *bootstrapper) detach(frame *navigationFrame) {
	b.detacher.detach(b.nc.frame)
}
