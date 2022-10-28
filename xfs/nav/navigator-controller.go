package nav

import (
	"fmt"

	. "github.com/snivilised/extendio/translate"
)

type navigatorController struct {
	impl navigatorImpl
}

func (c *navigatorController) Walk(root string) *TraverseResult {
	o := c.impl.options()
	o.Notify.OnBegin(root)
	frame := navigationFrame{
		Root: root,
	}
	initialState := c.initial(&o.Listen)

	if initialState != ListenInactive {
		frame.listener = &navigationListener{
			listen: initialState,
		}
		frame.listener.states = *c.listeners(&frame)
		frame.listener.init()

		// NB: the original client has already been preserved so overwriting
		// it here is not a problem.
		//
		o.Callback = func(item *TraverseItem) *LocalisableError {
			return frame.listener.current(item)
		}

		if o.Notify.OnStart == nil {
			o.Notify.OnStart = func(description string) {
				fmt.Printf("===> Start Listening: '%v'\n", description)
			}
		}

		if o.Notify.OnStop == nil {
			o.Notify.OnStop = func(description string) {
				fmt.Printf("===> Stopped Listening: '%v'\n", description)
			}
		}
	}

	result := &TraverseResult{
		Error: c.impl.top(&frame),
	}
	o.Notify.OnEnd(result)

	return result
}

func (c *navigatorController) initial(o *ListenOptions) ListeningState {
	initialState := ListenInactive

	switch {
	case (o.Start != nil) && (o.Stop != nil):
		initialState = ListenPending

	case o.Start != nil:
		initialState = ListenPending
		o.Stop = &ListenerFn{
			Name: "no-op: run to completion (don't stop early)",
			Fn: func(item *TraverseItem) bool {
				return false
			},
		}

	case o.Stop != nil:
		initialState = ListenActive
		o.Start = &ListenerFn{
			Name: "no-op: start listening straight away",
			Fn: func(item *TraverseItem) bool {
				return true
			},
		}
	}

	return initialState
}

func (c *navigatorController) listeners(frame *navigationFrame) *navigationListeners {
	// should this go into listener?
	//
	o := c.impl.options()
	client := o.Callback

	return &navigationListeners{
		ListenPending: func(item *TraverseItem) *LocalisableError {
			// listening not yet started
			//
			if c.impl.options().Listen.Start.IsMatch(item) {
				frame.listener.transition(ListenActive)
				o.Notify.OnStart(o.Listen.Start.Description())

				if o.Behaviours.Listen.InclusiveStart {
					return client(item)
				}
				return nil
			}
			return nil
		},

		ListenActive: func(item *TraverseItem) *LocalisableError {
			// listening
			//
			if c.impl.options().Listen.Stop.IsMatch(item) {
				frame.listener.transition(ListenRetired)
				o.Notify.OnStop(o.Listen.Stop.Description())

				if o.Behaviours.Listen.InclusiveStop {
					return client(item)
				}
				return nil
			}
			return client(item)
		},

		ListenRetired: func(item *TraverseItem) *LocalisableError {
			// done with listening, lets finish early
			// TODO: define a new return error that denotes finished
			//
			return nil
		},
	}
}
