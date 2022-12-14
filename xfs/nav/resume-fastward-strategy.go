package nav

import (
	"fmt"
)

type fastwardListener struct {
	target string
}

func (l *fastwardListener) Description() string {
	return fmt.Sprintf(">>> fast forwarding >>> to: '%v'", l.target)
}

func (l *fastwardListener) IsMatch(item *TraverseItem) bool {
	return item.Path == l.target
}

type fastwardStrategy struct {
	baseStrategy
	client struct {
		state ListeningState
	}
	overrideInfo *overrideListenerInfo
}

func (s *fastwardStrategy) init(params *strategyInitParams) {

	// this is the start we revert back to when we get back to the resume point
	//
	s.client.state = params.ps.Active.Listen

	pci := &preserveClientInfo{
		lo:         &s.o.Listen,
		behaviours: s.o.Store.Behaviours.Listen,
	}

	s.overrideInfo = &overrideListenerInfo{
		client: pci,
		override: &overrideClientInfo{
			lo: &ListenOptions{
				Stop: &fastwardListener{
					target: s.ps.Active.NodePath,
				},
			},
		},
		ps: params.ps,
	}
	backfill(s.overrideInfo.override.lo)

	s.attach(&resumeAttachParams{
		o:     s.o,
		frame: params.frame,
		lo:    s.overrideInfo.override.lo,
	})
}

func (s *fastwardStrategy) attach(params *resumeAttachParams) {
	// for the resume scenario, we don't have to attach a new decorator, rather the
	// change of state accounts for the new behaviour. frame.attach is only used
	// to setup the navigation-listener. resume sits atop the listener but a new item
	// still has to be pushed onto the resumeStack.
	//
	params.frame.listener.lo = params.lo
	params.frame.listener.resumeStack.Push(params.lo)
	params.frame.listener.transition(ListenFastward)
	params.frame.notifiers.mute(notificationAllEn)
}

func (s *fastwardStrategy) detach(frame *navigationFrame) {
	frame.listener.dispose()
	fmt.Printf("==>⚠️⚠️⚠️ fastwardStrategy: detach\n")

	// now restore the client state (preserveClientInfo)
	//
	s.o.Store.Behaviours.Listen = s.overrideInfo.client.behaviours

	frame.notifiers.mute(notificationAllEn, false)

	if s.ps.Active.Listen == ListenFastward {
		panic("invalid state transition detected (ListenFastward)")
	}
	frame.listener.transition(s.ps.Active.Listen)
}

func (s *fastwardStrategy) resume(info *strategyResumeInfo) *TraverseResult {
	// fast-forward doesn't need to restore the entire state, eg, the
	// Depth can begin as per usual, without being restored.
	//
	info.nc.root(func() string {
		return info.ps.Active.Root
	})

	return info.nc.Walk(info.ps.Active.Root)
}
