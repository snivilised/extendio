package xfs

import "errors"

// NewNavigator navigator factory function which uses the functional
// options pattern.
func NewNavigator(fn ...TraverseOptionFn) TraverseNavigator {
	o := composeTraverseOptions(fn...)

	if o.Callback == nil {
		panic(LocalisableError{
			Inner: errors.New("missing callback function"),
		})
	}

	var impl navigatorImpl

	switch o.Subscription {
	case SubscribeAny:
		impl = &universalNavigator{
			navigator: navigator{o: o, agent: &childAgent{
				o: o, DO_INVOKE: true,
			}},
		}

	case SubscribeFolders:
		impl = &foldersNavigator{
			navigator: navigator{o: o, agent: &childAgent{
				o: o, DO_INVOKE: true,
			}},
		}

	case SubscribeFiles:
		impl = &filesNavigator{
			navigator: navigator{o: o, agent: &childAgent{
				o: o, DO_INVOKE: false,
			}},
		}
	}
	nav := &navigatorController{
		impl: impl,
	}

	return nav
}
