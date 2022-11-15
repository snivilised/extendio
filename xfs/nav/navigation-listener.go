package nav

import (
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
	ListenPending // denotes conditional listening is awaiting activation
	ListenActive  // denotes conditional listening is active (callback is invoked)
	ListenRetired // denoted conditional listening is now deactivated
)

type navigationListeningStates map[ListeningState]TraverseCallback

func bootstrapListener(o *TraverseOptions, frame *navigationFrame) {

	if o.Listen.Start == nil && o.Listen.Stop == nil {
		return
	}
	frame.listener = newListener(o, frame)
}

func newListener(o *TraverseOptions, frame *navigationFrame) *navigationListener {

	initialState := backfillListenState(&o.Listen)

	listener := &navigationListener{
		listen: initialState,
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

func backfillListenState(o *ListenOptions) ListeningState {
	initialState := ListenDefault

	switch {
	case (o.Start != nil) && (o.Stop != nil):
		initialState = ListenPending

	case o.Start != nil:
		initialState = ListenPending
		o.Stop = &ListenBy{
			Name: "run to completion (don't stop early)",
			Fn: func(item *TraverseItem) bool {
				return false
			},
		}

	case o.Stop != nil:
		initialState = ListenActive
		o.Start = &ListenBy{
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
		ListenPending: func(item *TraverseItem) *LocalisableError {
			// listening not yet started
			//
			if params.o.Listen.Start.IsMatch(item) {
				params.frame.listener.transition(ListenActive)
				params.o.Notify.OnStart(params.o.Listen.Start.Description())

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
			if params.o.Listen.Stop.IsMatch(item) {
				params.frame.listener.transition(ListenRetired)
				params.o.Notify.OnStop(params.o.Listen.Stop.Description())

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
	listen  ListeningState
	states  navigationListeningStates
	current TraverseCallback
}

func (l *navigationListener) init() {
	l.transition(l.listen)
}

func (l *navigationListener) transition(state ListeningState) {
	l.listen = state
	l.current = l.states[state]
}
