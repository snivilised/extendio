package nav

// o.Store.Behaviours.Listen
type fastwardStrategy struct {
	client struct {
		state      ListeningState
		behaviours ListenBehaviour
	}
	lo *ListenOptions

	// to be replaced with spawn/fastward
}

func (s *fastwardStrategy) init(params *listenerInitParams) {
	s.client.state = params.state
	s.client.behaviours = params.o.Store.Behaviours.Listen

	if params.listener == nil {

		// resumeListener()

		// &navigationListener{
		// 	state:       initialState,
		// 	resumeStack: collections.NewStack[*ListenOptions](),
		// }

		s.lo = &ListenOptions{}
	} else {
		s.lo = &ListenOptions{}
	}
	s.lo = &ListenOptions{
		Start: params.listener.lo.Start,
		Stop:  params.listener.lo.Stop,
	}

	params.state = ListenFastward

	// now set up synthetic listener
}
