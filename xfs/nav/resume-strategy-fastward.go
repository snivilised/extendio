package nav

import (
	"fmt"

	"github.com/snivilised/extendio/internal/log"
)

type fastwardListener struct {
	target string
}

func (l *fastwardListener) Description() string {
	return fmt.Sprintf(">>> fast forwarding >>> to: '%v'", l.target)
}

func (l *fastwardListener) Validate() {}

func (l *fastwardListener) Source() string {
	return l.target
}

func (l *fastwardListener) IsMatch(item *TraverseItem) bool {
	return item.Path == l.target
}

func (l *fastwardListener) IsApplicable(_ *TraverseItem) bool {
	return true
}

func (l *fastwardListener) Scope() FilterScopeBiEnum {
	return ScopeAllEn
}

type fastwardStrategy struct {
	baseStrategy
	client struct {
		state ListeningState
	}
	overrideInfo *overrideListenerInfo
}

func (s *fastwardStrategy) init(params *strategyInitParams) {
	s.nc.logger().Info("fastward resume", log.String("path", params.ps.Active.Root))
	// this is the start we revert back to when we get back to the resume point
	//
	s.client.state = params.ps.Active.Listen

	pci := &preserveClientInfo{
		triggers:   params.triggers,
		behaviours: s.o.Store.Behaviours.Listen,
	}

	s.overrideInfo = &overrideListenerInfo{
		client: pci,
		override: &overrideClientInfo{
			triggers: &ListenTriggers{
				Stop: &fastwardListener{
					target: s.ps.Active.NodePath,
				},
			},
		},
		ps: params.ps,
	}

	s.attach(&resumeAttachParams{
		o:        s.o,
		frame:    params.frame,
		triggers: s.overrideInfo.override.triggers,
	})
}

func (s *fastwardStrategy) attach(params *resumeAttachParams) {
	// for the resume scenario, we don't have to attach a new decorator, rather the
	// change of state accounts for the new behaviour. frame.attach is only used
	// to setup the navigation-listener. resume sits atop the listener but a new item
	// still has to be pushed onto the resumeStack.
	//
	params.frame.listener.triggers = params.triggers
	params.frame.listener.resumeStack.Push(params.triggers)
	params.frame.listener.transition(ListenFastward)
	params.frame.notifiers.mute(notificationAllEn)
}

func (s *fastwardStrategy) detach(frame *navigationFrame) {
	frame.listener.dispose()

	// now restore the client state (preserveClientInfo)
	//
	s.o.Store.Behaviours.Listen = s.overrideInfo.client.behaviours

	frame.notifiers.mute(notificationAllEn, false)

	if s.ps.Active.Listen == ListenFastward {
		panic(NewInvalidResumeStateTransitionNativeError("ListenFastward"))
	}

	frame.listener.transition(s.ps.Active.Listen)
}

func (s *fastwardStrategy) resume(info *strategyResumeInfo) (*TraverseResult, error) {
	resumeAt := info.ps.Active.NodePath
	s.nc.logger().Info("fastward resume",
		log.String("root-path", info.ps.Active.Root),
		log.String("resume-at-path", resumeAt),
	)

	// fast-forward doesn't need to restore the entire state, eg, the
	// Depth can begin as per usual, without being restored.
	//
	return info.nc.walk(info.ps.Active.Root)
}
