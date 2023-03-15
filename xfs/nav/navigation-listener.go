package nav

import (
	"github.com/snivilised/extendio/collections"
	. "github.com/snivilised/extendio/i18n"
)

// ListenHandler
type ListenHandler func(description string)

// ListenTriggers
type ListenTriggers struct {
	Start TraverseFilter
	Stop  TraverseFilter
}

type ListenBehaviour struct {
	InclusiveStart bool
	InclusiveStop  bool
}

// ListeningState denotes whether user defined callback is being invoked.
type ListeningState uint

const (
	ListenUndefined ListeningState = iota
	ListenDeaf                     // listen not active, callback always invoked (subject to filtering)
	ListenFastward                 // listen used to resume by fast-forwarding
	ListenPending                  // conditional listening is awaiting activation
	ListenActive                   // conditional listening is active (callback is invoked)
	ListenRetired                  // conditional listening is now deactivated
)

type navigationListeningStates map[ListeningState]LabelledTraverseCallback

type listenStatesParams struct {
	// currently used for makeStates and listener.decorate
	//
	triggers *ListenTriggers
	o        *TraverseOptions
	frame    *navigationFrame
	detacher resumeDetacher
}

type navigationListener struct {
	state       ListeningState
	states      navigationListeningStates
	current     LabelledTraverseCallback
	resumeStack *collections.Stack[*ListenTriggers]
	triggers    *ListenTriggers
}

func (l *navigationListener) init() {
	l.transition(l.state)
}

func (l *navigationListener) makeStates(params *listenStatesParams) {

	// The listen states are aware of the raw callback, because frame.client
	// denotes the decorated client which may incorporate the listener callback.
	// If the client simply called frame.client, then there would be an infinite
	// loop if listening is active. Elsewhere, frame.client is acceptable to call,
	// so when listen is active, it is routed through the listener callback embedded
	// into frame.client. The listener callback simply delegates to the current
	// listener state. When an attachment occurs for the purposes of resume, the state
	// machine takes account of required change in behaviour, ie we don't have to
	// re-decorate the client. The only thing required in this scenario is the modification
	// of the resume stack which is updated with the resume specific ListenTriggers and
	// reverted at a later point via detach (resume stack pop).
	//
	l.states = navigationListeningStates{

		// Just use the original unadulterated (filtered) client
		// (this depends on filter-init assigning to raw)
		//
		ListenDeaf: params.frame.raw,

		ListenFastward: LabelledTraverseCallback{
			Label: "ListenFastward decorator",
			Fn: func(item *TraverseItem) error {
				// fast forwarding to resume point
				//
				if params.frame.listener.triggers.Stop.IsMatch(item) {
					if params.detacher != nil {
						// detach performs state transition
						//
						params.detacher.detach(params.frame)

						// NB: ok to call the client here without concern over causing an infinite
						// loop because the detach has performed a state transition.
						//
						return params.frame.client.Fn(item)
					} else {
						panic(NewMissingListenDetacherFunctionNativeError("fastward"))
					}
				} else {
					item.skip = true
				}
				return nil
			},
		},

		ListenPending: LabelledTraverseCallback{
			Label: "ListenPending decorator",
			Fn: func(item *TraverseItem) error {
				// listening not yet started
				//
				if params.frame.listener.triggers.Start.IsMatch(item) {
					params.frame.listener.transition(ListenActive)
					params.frame.notifiers.start.invoke(params.frame.listener.triggers.Start.Description())

					if params.o.Store.Behaviours.Listen.InclusiveStart {
						return params.frame.raw.Fn(item)
					}
					return nil
				}
				return nil
			},
		},

		ListenActive: LabelledTraverseCallback{
			Label: "ListenActive decorator",
			Fn: func(item *TraverseItem) error {
				// listening
				//
				if params.frame.listener.triggers.Stop.IsMatch(item) {
					params.frame.listener.transition(ListenRetired)
					params.frame.notifiers.stop.invoke(params.frame.listener.triggers.Stop.Description())

					if params.o.Store.Behaviours.Listen.InclusiveStop {
						return params.frame.raw.Fn(item)
					}
					return nil
				}
				return params.frame.raw.Fn(item)
			},
		},

		ListenRetired: LabelledTraverseCallback{
			Label: "ListenRetired decorator",
			Fn: func(item *TraverseItem) error {
				return NewTerminateTraverseError()
			},
		},
	}
}

func (l *navigationListener) transition(state ListeningState) {
	l.state = state
	l.current = l.states[state]
}

func (l *navigationListener) decorate(params *listenStatesParams) {
	// decorator: is the listen state machine, ie l.current.
	// decorated: is frame.client, what is returned from frame.decorate.
	// Since we know these, listenStatesParams does not have to include
	// the decorated member.
	//

	decorator := &LabelledTraverseCallback{
		Label: "listener decorator",
		Fn: func(item *TraverseItem) error {
			return l.current.Fn(item)
		},
	}
	params.frame.decorate("listener 🎀", decorator)

	l.triggers = params.triggers
	l.resumeStack.Push(l.triggers)
	l.init()
}

type initialListenerState struct {
	initialState ListeningState
	Listen       ListenTriggers
}

func backfill(defs *ListenDefinitions) *initialListenerState {

	state := initialListenerState{
		initialState: ListenDeaf,
	}

	startAt := FilterDef{
		Type:            FilterTypeGlobEn,
		Description:     "start listening straight away",
		Pattern:         "*",
		Scope:           ScopeAllEn,
		IfNotApplicable: TriStateBoolTrueEn,
	}
	stopAt := FilterDef{
		Type:        FilterTypeGlobEn,
		Description: "run to completion, don't stop early",
		// We don't want this to match, so the stop trigger is never fired.
		// "/" is prohibited within the name of a file-system item on
		// all OSs.
		//
		Pattern:         "/",
		Scope:           ScopeRootEn,
		IfNotApplicable: TriStateBoolFalseEn,
	}

	switch {
	case (defs.StartAt != nil) && (defs.StopAt != nil):
		state.initialState = ListenPending

	case defs.StartAt != nil:
		state.initialState = ListenPending
		defs.StopAt = &stopAt

	case defs.StopAt != nil:
		state.initialState = ListenActive
		defs.StartAt = &startAt

	default:
		defs.StartAt = &startAt
		defs.StopAt = &stopAt
	}

	state.Listen.Start = newNodeFilter(defs.StartAt)
	state.Listen.Stop = newNodeFilter(defs.StopAt)

	return &state
}

func (l *navigationListener) dispose() *ListenTriggers {

	previous, _ := l.resumeStack.Pop()
	l.triggers, _ = l.resumeStack.Current()

	return previous
}
