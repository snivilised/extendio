package nav

import (
	"github.com/snivilised/extendio/collections"
	"github.com/snivilised/extendio/xfs/utils"
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
	b.nc.ns = &NavigationState{
		Filters: b.nc.frame.filters,
		Root:    &b.nc.frame.root,
		Logger:  utils.NewRoProp[ClientLogger](b.nc.impl.logger()),
	}
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
	state := backfill(&b.o.Store.ListenDefs)

	b.nc.frame.listener = &navigationListener{
		state:       state.initialState,
		resumeStack: collections.NewStack[*ListenTriggers](),
	}

	b.nc.frame.listener.makeStates(&listenStatesParams{
		o: b.o, frame: b.nc.frame,
		detacher: b,
	})

	b.nc.frame.listener.decorate(&listenStatesParams{
		triggers: &state.Listen, frame: b.nc.frame,
	})
}

type preserveClientInfo struct {
	triggers   *ListenTriggers
	behaviours ListenBehaviour
}

type overrideClientInfo struct {
	triggers *ListenTriggers
}

type overrideListenerInfo struct {
	client   *preserveClientInfo
	override *overrideClientInfo
	ps       *persistState
}

func (b *bootstrapper) initResume(ps *persistState) {
	if b.rc == nil {
		err := NewResumeControllerNotSetNativeError("bootstrapper.initResume")
		b.nc.impl.logger().Error(err.Error())
		panic(err)
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
