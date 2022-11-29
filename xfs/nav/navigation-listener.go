package nav

import (
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
	// currently used for composeListenStates and listener.decorate
	//
	decorated *LabelledTraverseCallback // only use me for composeListenStates
	lo        *ListenOptions
	o         *TraverseOptions
	frame     *navigationFrame
	detach    func()
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

func (l *navigationListener) composeListenStates(params *listenStatesParams) {

	l.states = navigationListeningStates{

		// Just use the original unadulterated (filtered) client
		// (this depends on filter-init assigning to raw)
		//
		ListenDeaf: params.frame.raw,

		ListenFastward: LabelledTraverseCallback{
			Label: "ListenFastward decorator",
			Fn: func(item *TraverseItem) *LocalisableError {
				// fast forwarding to resume point
				//
				if params.frame.listener.lo.Stop.IsMatch(item) {
					if params.detach != nil {
						params.detach()
					} else {
						panic("listen-state(fastward): missing detacher function from listenStatesParams")
					}
				}
				return nil
			},
		},

		ListenPending: LabelledTraverseCallback{
			Label: "ListenPending decorator",
			Fn: func(item *TraverseItem) *LocalisableError {
				// listening not yet started
				//
				if params.frame.listener.lo.Start.IsMatch(item) {
					params.frame.listener.transition(ListenActive)
					params.o.Notify.OnStart(params.frame.listener.lo.Start.Description())

					if params.o.Store.Behaviours.Listen.InclusiveStart {
						return params.decorated.Fn(item)
					}
					return nil
				}
				return nil
			},
		},

		ListenActive: LabelledTraverseCallback{
			Label: "ListenActive decorator",
			Fn: func(item *TraverseItem) *LocalisableError {
				// listening
				//
				if params.frame.listener.lo.Stop.IsMatch(item) {
					params.frame.listener.transition(ListenRetired)
					params.o.Notify.OnStop(params.frame.listener.lo.Stop.Description())

					if params.o.Store.Behaviours.Listen.InclusiveStop {
						return params.decorated.Fn(item)
					}
					return nil
				}
				return params.decorated.Fn(item)
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
	// Since we know these, listenStatesParams does not have to include
	// the decorated member. (TODO: may be we should not repurpose
	// listenStatesParams for multiple scenarios)
	//
	decorator := LabelledTraverseCallback{
		Label: "listener decorator",
		Fn: func(item *TraverseItem) *LocalisableError {
			return l.current.Fn(item)
		},
	}
	decorated := params.frame.decorate("listener 🎀", decorator)

	l.composeListenStates(&listenStatesParams{
		decorated: decorated, o: params.o, frame: params.frame,
		detach: func() {
			l.detach()
		},
	})
	l.lo = params.lo
	l.resumeStack.Push(l.lo)
	l.init()
}

func (l *navigationListener) attach(options *ListenOptions, state ListeningState) {
	// ??? TODO: don't we have to also rebuild the states with composeListenStates?
	//
	l.lo = options
	l.resumeStack.Push(options)

	if state != ListenUndefined {
		l.transition(state)
	}
}

func (l *navigationListener) detach() *ListenOptions {
	previous, _ := l.resumeStack.Pop()
	l.lo, _ = l.resumeStack.Current()

	return previous
}
