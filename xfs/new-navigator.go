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
			navigator: navigator{options: o},
		}

	case SubscribeFolders:
		core = &foldersNavigator{
			navigator: navigator{options: o},
		}

	case SubscribeFiles:
		core = &filesNavigator{
			navigator: navigator{options: o},
		}
	}
	nav := &navigatorController{
		core: core,
	}

	return nav
}
