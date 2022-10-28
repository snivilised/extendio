package nav

import (
	"errors"

	. "github.com/snivilised/extendio/translate"
)

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
	ctrl := &navigatorController{
		impl: impl,
	}

	return ctrl
}
