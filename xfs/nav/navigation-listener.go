package nav

import (
	. "github.com/snivilised/extendio/translate"
)

// Listener
type Listener interface {
	Description() string
	IsMatch(item *TraverseItem) bool
}

type ListenerBy struct {
	Fn   ListenPredicate
	Name string
}

func (f *ListenerBy) Description() string {
	return f.Name
}

func (f *ListenerBy) IsMatch(item *TraverseItem) bool {
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
	ListenDefault ListeningState = iota
	ListenPending                // denotes conditional listening is awaiting activation
	ListenActive                 // denotes conditional listening is active (callback is invoked)
	ListenRetired                // denoted conditional listening is now deactivated
)

var listeningStateStrings map[ListeningState]string = map[ListeningState]string{
	ListenDefault: "Listen Default",
	ListenPending: "Listen Pending",
	ListenActive:  "Listen Active",
	ListenRetired: "Listen Retired",
}

func (s ListeningState) String() string {
	return listeningStateStrings[s]
}

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
	listener.states = *listenStates(o, frame)
	listener.init()

	o.Callback = func(item *TraverseItem) *LocalisableError {
		return listener.current(item)
	}

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
		o.Stop = &ListenerBy{
			Name: "run to completion (don't stop early)",
			Fn: func(item *TraverseItem) bool {
				return false
			},
		}

	case o.Stop != nil:
		initialState = ListenActive
		o.Start = &ListenerBy{
			Name: "start listening straight away",
			Fn: func(item *TraverseItem) bool {
				return true
			},
		}
	}

	return initialState
}

func listenStates(o *TraverseOptions, frame *navigationFrame) *navigationListeningStates {

	return &navigationListeningStates{
		ListenPending: func(item *TraverseItem) *LocalisableError {
			// listening not yet started
			//
			if o.Listen.Start.IsMatch(item) {
				frame.listener.transition(ListenActive)
				o.Notify.OnStart(o.Listen.Start.Description())

				if o.Behaviours.Listen.InclusiveStart {
					return frame.client(item)
				}
				return nil
			}
			return nil
		},

		ListenActive: func(item *TraverseItem) *LocalisableError {
			// listening
			//
			if o.Listen.Stop.IsMatch(item) {
				frame.listener.transition(ListenRetired)
				o.Notify.OnStop(o.Listen.Stop.Description())

				if o.Behaviours.Listen.InclusiveStop {
					return frame.client(item)
				}
				return nil
			}
			return frame.client(item)
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
