package nav

// Listener
type Listener interface {
	Description() string
	IsMatch(item *TraverseItem) bool
}

type ListenerFn struct {
	Fn   ListenPredicate
	Name string
}

func (f *ListenerFn) Description() string {
	return f.Name
}

func (f *ListenerFn) IsMatch(item *TraverseItem) bool {
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
	ListenInactive ListeningState = iota
	ListenPending                 // denotes conditional listening is awaiting activation
	ListenActive                  // denotes conditional listening is active (callback is invoked)
	ListenRetired                 // denoted conditional listening is now deactivated
)

var listeningStateStrings map[ListeningState]string = map[ListeningState]string{
	ListenInactive: "Listening Inactive",
	ListenPending:  "Listening Pending",
	ListenActive:   "Listening Active",
	ListenRetired:  "Listening Retired",
}

func (s ListeningState) String() string {
	return listeningStateStrings[s]
}

type navigationListeners map[ListeningState]TraverseCallback

type navigationListener struct {
	listen  ListeningState
	states  navigationListeners
	current TraverseCallback
}

func (l *navigationListener) init() {
	l.transition(l.listen)
}

func (l *navigationListener) transition(state ListeningState) {
	l.listen = state
	l.current = l.states[state]
}
