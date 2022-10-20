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

	var subject navigatorSubject

	switch o.Subscription {
	case SubscribeAny:
		subject = &universalNavigator{
			navigator: navigator{options: o},
		}

	case SubscribeFolders:
		subject = &foldersNavigator{
			navigator: navigator{options: o},
		}

	case SubscribeFiles:
		subject = &filesNavigator{
			navigator: navigator{options: o},
		}
	}
	nav := &navigatorController{
		subject: subject,
	}

	return nav
}
