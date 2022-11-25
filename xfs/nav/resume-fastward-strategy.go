package nav

import "fmt"

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
		state      ListeningState
		behaviours ListenBehaviour
	}
	lo *ListenOptions
}

func (s *fastwardStrategy) init(params *strategyInitParams) {
	s.client.state = params.state
	s.client.behaviours = s.o.Store.Behaviours.Listen

	listener := &fastwardListener{
		target: s.active.NodePath,
	}
	s.lo = &ListenOptions{
		Start: nil,      // we want to start listening immediately
		Stop:  listener, // stop when we get the the resume point
	}

	// force the state into fast forward
	//
	params.state = ListenFastward
}

func (s *fastwardStrategy) listenOptions() *ListenOptions {
	return s.lo
}
