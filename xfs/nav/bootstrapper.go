package nav

import (
	"errors"

	"github.com/snivilised/extendio/collections"
)

type bootstrapper struct {
	o  *TraverseOptions
	oi *overrideListenerInfo
	nc *navigatorController
	rc *resumeController
}

func (b *bootstrapper) init() {
	b.nc.frame = b.nc.init()
	b.initFilters()
	b.initNotifiers()
	b.initListener()
	b.nc.navState(func() *NavigationState {
		return &NavigationState{Filters: b.nc.frame.filters}
	})
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
}

func (b *bootstrapper) initListener() {
	initialState := b.backfill(&b.o.Listen)

	listener := &navigationListener{
		state:       initialState,
		resumeStack: collections.NewStack[*ListenOptions](),
	}
	listener.decorate(&listenStatesParams{
		o: b.o, lo: &b.o.Listen, frame: b.nc.frame,
	})
	b.nc.frame.listener = listener
}

// func (b *bootstrapper) decorate(decorator LabelledTraverseCallback, listener *navigationListener) *LabelledTraverseCallback {
// 	decorated := &b.nc.frame.client
// 	b.nc.frame.decorate("listener 🎀", decorator)

// 	listener.composeListenStates(&listenStatesParams{
// 		decorated: decorated, o: b.o, frame: b.nc.frame,
// 		detach: func() {
// 			b.rc.detach()
// 		},
// 	})

// 	return decorated
// }

func (b *bootstrapper) backfill(lo *ListenOptions) ListeningState {

	initialState := ListenDeaf

	start := func(item *TraverseItem) bool {
		return false
	}
	stop := func(item *TraverseItem) bool {
		return true
	}

	switch {
	case (lo.Start != nil) && (lo.Stop != nil):
		initialState = ListenPending

	case lo.Start != nil:
		initialState = ListenPending
		lo.Stop = &ListenBy{
			Name: "run to completion, don't stop early",
			Fn:   start,
		}

	case lo.Stop != nil:
		initialState = ListenActive
		lo.Start = &ListenBy{
			Name: "start listening straight away",
			Fn:   stop,
		}

	default:
		lo.Stop = &ListenBy{
			Name: "dormant listener, don't stop early",
			Fn:   start,
		}
		lo.Start = &ListenBy{
			Name: "dormant listener, start listening straight away",
			Fn:   stop,
		}
	}

	return initialState
}

type preserveClientInfo struct {
	lo         *ListenOptions
	behaviours ListenBehaviour
	notify     Notifications
}

type overrideClientInfo struct {
	lo *ListenOptions
}

type overrideListenerInfo struct {
	client   *preserveClientInfo
	override *overrideClientInfo
	ps       *persistState
}

func (b *bootstrapper) resume(o *TraverseOptions, ps *persistState) {

	if b.rc == nil {
		panic(errors.New("bootstrapper.resume: resume controller not set"))
	}

	pci := &preserveClientInfo{
		lo:         &o.Listen,
		behaviours: b.o.Store.Behaviours.Listen,
		notify:     b.o.Notify,
	}

	initParams := &strategyInitParams{
		state: ps.Active.Listen,
		frame: b.nc.frame,
		pci:   pci,
	}

	// the state defaults to the state restored as the persist state,
	// but the strategy is free to modify it in the params to suit
	// its needs.
	//
	b.rc.strategy.init(initParams)

	b.oi = &overrideListenerInfo{
		client: pci,
		override: &overrideClientInfo{
			lo: b.rc.strategy.listenOptions(),
		},
		ps: ps,
	}

	// we can ignore the returned listen state from backfill, because
	// the strategy knows the correct initial state.
	//
	_ = b.backfill(b.oi.override.lo)
	b.nc.frame.listener.attach(b.oi.override.lo, initParams.state)
}
