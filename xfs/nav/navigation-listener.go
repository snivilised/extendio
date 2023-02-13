package nav

import (
	"fmt"

	"github.com/snivilised/extendio/collections"
	. "github.com/snivilised/extendio/translate"
)

// Listener
type Listener interface {
	Description() string
	IsMatch(item *TraverseItem) bool
}

type ListenBy struct {
	Fn   ListenPredicate
	Name string
}

func (f *ListenBy) Description() string {
	return f.Name
}

func (f *ListenBy) IsMatch(item *TraverseItem) bool {
	return f.Fn(item)
}

// ListenPredicate
type ListenPredicate func(item *TraverseItem) bool

// ListenHandler
type ListenHandler func(description string)

// ListenOptions
type ListenOptions struct {
	Start Listener
	Stop  Listener
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
	// currently used for buildStates and listener.decorate
	//
	lo       *ListenOptions
	o        *TraverseOptions
	frame    *navigationFrame
	detacher resumeDetacher
}

type navigationListener struct {
	state       ListeningState
	states      navigationListeningStates
	current     LabelledTraverseCallback
	resumeStack *collections.Stack[*ListenOptions]
	lo          *ListenOptions
}

func (l *navigationListener) init() {
	l.transition(l.state)
}

func (l *navigationListener) buildStates(params *listenStatesParams) {

	// The listen states are aware of the raw callback, because frame.client
	// denotes the decorated client which may incorporate the listener callback.
	// If the client simply called frame.client, then there would be an infinite
	// loop if listening is active. Elsewhere, frame.client is acceptable to call,
	// so when listen is active, it is routed through the listener callback embedded
	// into frame.client. The listener callback simply delegates to the current
	// listener state. When an attachment occurs for the purposes of resume, the state
	// machine takes account of required change in behaviour, ie we don't have to
	// re-decorate the client. The only thing required in this scenario is the modification
	// of the resume stack which is updated with the resume specific ListenOptions and
	// reverted at a later point via detach (resume stack pop).
	//
	l.states = navigationListeningStates{

		// Just use the original unadulterated (filtered) client
		// (this depends on filter-init assigning to raw)
		//
		ListenDeaf: params.frame.raw,

		ListenFastward: LabelledTraverseCallback{
			Label: "ListenFastward decorator",
			Fn: func(item *TraverseItem) *LocalisableError {
				fmt.Printf(">>> 🤎 ListenFastward decorator, item: '%s'\n", item.Path)
				// fast forwarding to resume point
				//
				if params.frame.listener.lo.Stop.IsMatch(item) {
					fmt.Printf(">>> 🤎🤎🤎 DETACHING-AT: ListenFastward decorator, item: '%s'\n", item.Path)
					if params.detacher != nil {
						// detach performs state transition
						//
						params.detacher.detach(params.frame)

						// NB: ok to call the client here without concern over causing an infinite
						// loop because the detach has performed a state transition.
						//
						return params.frame.client.Fn(item)
					} else {
						panic("listen-state(fastward): missing detacher function from listenStatesParams")
					}
				} else {
					item.skip = true
				}
				return nil
			},
		},

		ListenPending: LabelledTraverseCallback{
			Label: "ListenPending decorator",
			Fn: func(item *TraverseItem) *LocalisableError {
				fmt.Printf(">>> 💙 ListenPending decorator, item: '%s'\n", item.Path)
				// listening not yet started
				//
				if params.frame.listener.lo.Start.IsMatch(item) {
					params.frame.listener.transition(ListenActive)
					params.frame.notifiers.start.invoke(params.frame.listener.lo.Start.Description())

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
			Fn: func(item *TraverseItem) *LocalisableError {
				fmt.Printf(">>> 💜ListenActive decorator, item: '%s'\n", item.Path)
				// listening
				//
				if params.frame.listener.lo.Stop.IsMatch(item) {
					params.frame.listener.transition(ListenRetired)
					params.frame.notifiers.stop.invoke(params.frame.listener.lo.Start.Description())

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
			Fn: func(item *TraverseItem) *LocalisableError {
				return &LocalisableError{Inner: TERMINATE_ERR}
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
	// TODO: Since we know these, listenStatesParams does not have to include
	// the decorated member. (TODO: may be we should not repurpose
	// listenStatesParams for multiple scenarios)
	//

	decorator := &LabelledTraverseCallback{
		Label: "listener decorator",
		Fn: func(item *TraverseItem) *LocalisableError {
			// fmt.Printf(">>> ❤️ listener decorator, item: '%s'\n", item.Path)
			return l.current.Fn(item)
		},
	}
	params.frame.decorate("listener 🎀", decorator)

	l.lo = params.lo
	l.resumeStack.Push(l.lo)
	l.init()
}

func backfill(lo *ListenOptions) ListeningState {

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

func (l *navigationListener) dispose() *ListenOptions {

	previous, _ := l.resumeStack.Pop()
	l.lo, _ = l.resumeStack.Current()

	return previous
}
