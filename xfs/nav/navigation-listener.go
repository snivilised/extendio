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
	ListenDefault
	ListenFastward // resume fast-forwarding
	ListenPending  // denotes conditional listening is awaiting activation
	ListenActive   // denotes conditional listening is active (callback is invoked)
	ListenRetired  // denoted conditional listening is now deactivated
)

type navigationListeningStates map[ListeningState]TraverseCallback

func bootstrapListener(o *TraverseOptions, frame *navigationFrame) {

	if o.Listen.Start == nil && o.Listen.Stop == nil {
		return
	}
	frame.listener = newListener(o, frame)
	frame.listener.attach(&o.Listen)
}

func newListener(o *TraverseOptions, frame *navigationFrame) *navigationListener {

	initialState := backfillListenState(&o.Listen)

	listener := &navigationListener{
		state:       initialState,
		resumeStack: collections.NewStack[*ListenOptions](),
	}

	decorated := frame.client
	decorator := func(item *TraverseItem) *LocalisableError {
		return listener.current(item)
	}
	frame.decorate("listener 🎀", decorator)

	listener.states = *listenStates(&listenStatesParams{
		decorated: decorated, o: o, frame: frame,
	})
	listener.init()

	if o.Notify.OnStart == nil {
		o.Notify.OnStart = func(description string) {}
	}

	if o.Notify.OnStop == nil {
		o.Notify.OnStop = func(description string) {}
	}

	return listener
}

// func resumeListener(o *TraverseOptions, frame *navigationFrame) *navigationListener {

// 	return nil
// }

func backfillListenState(lo *ListenOptions) ListeningState {
	initialState := ListenDefault

	switch {
	case (lo.Start != nil) && (lo.Stop != nil):
		initialState = ListenPending

	case lo.Start != nil:
		initialState = ListenPending
		lo.Stop = &ListenBy{
			Name: "run to completion (don't stop early)",
			Fn: func(item *TraverseItem) bool {
				return false
			},
		}

	case lo.Stop != nil:
		initialState = ListenActive
		lo.Start = &ListenBy{
			Name: "start listening straight away",
			Fn: func(item *TraverseItem) bool {
				return true
			},
		}
	}

	return initialState
}

type listenStatesParams struct {
	decorated TraverseCallback
	o         *TraverseOptions
	frame     *navigationFrame
}

func listenStates(params *listenStatesParams) *navigationListeningStates {

	return &navigationListeningStates{

		ListenFastward: func(item *TraverseItem) *LocalisableError {
			// fast forwarding to resume point
			//
			return nil
		},

		ListenPending: func(item *TraverseItem) *LocalisableError {
			// listening not yet started
			//
			if params.frame.listener.lo.Start.IsMatch(item) {
				params.frame.listener.transition(ListenActive)
				params.o.Notify.OnStart(params.frame.listener.lo.Start.Description())

				if params.o.Store.Behaviours.Listen.InclusiveStart {
					return params.decorated(item)
				}
				return nil
			}
			return nil
		},

		ListenActive: func(item *TraverseItem) *LocalisableError {
			// listening
			//
			if params.frame.listener.lo.Stop.IsMatch(item) {
				params.frame.listener.transition(ListenRetired)
				params.o.Notify.OnStop(params.frame.listener.lo.Stop.Description())

				if params.o.Store.Behaviours.Listen.InclusiveStop {
					return params.decorated(item)
				}
				return nil
			}
			return params.decorated(item)
		},

		ListenRetired: func(item *TraverseItem) *LocalisableError {
			return &LocalisableError{Inner: TERMINATE_ERR}
		},
	}
}

type navigationListener struct {
	state       ListeningState
	states      navigationListeningStates
	current     TraverseCallback
	resumeStack *collections.Stack[*ListenOptions]
	lo          *ListenOptions
}

func (l *navigationListener) init() {
	l.transition(l.state)
}

func (l *navigationListener) transition(state ListeningState) {
	l.state = state
	l.current = l.states[state]
}

func (l *navigationListener) attach(options *ListenOptions) {
	l.lo = options
	l.resumeStack.Push(options)
}

// to be called by ListenFastward state:
//
// func (l *navigationListener) detach() {
// 	_, _ = l.resumeStack.Pop()
// 	l.currentLo, _ = l.resumeStack.Current()
// }
