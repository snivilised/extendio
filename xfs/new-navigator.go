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

	var core navigatorCore

	switch o.Subscription {
	case SubscribeAny:
		core = &universalNavigator{
			navigator: navigator{options: o, children: &childAgent{
				options: o, DO_INVOKE: true,
			}},
		}

	case SubscribeFolders:
		core = &foldersNavigator{
			navigator: navigator{options: o, children: &childAgent{
				options: o, DO_INVOKE: true,
			}},
		}

	case SubscribeFiles:
		core = &filesNavigator{
			navigator: navigator{options: o, children: &childAgent{
				options: o, DO_INVOKE: false,
			}},
		}
	}
	nav := &navigatorController{
		core: core,
	}

	return nav
}
