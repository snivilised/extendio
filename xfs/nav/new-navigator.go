package nav

type NavigatorFactory struct{}

// Create navigator factory function which uses the functional
// options pattern.
func (f *NavigatorFactory) Create(fn ...TraverseOptionFn) TraverseNavigator {
	o := composeTraverseOptions(fn...)

	if o.Callback.Fn == nil {
		panic(MISSING_CALLBACK_FN_L_ERR)
	}

	impl := (&navigatorImplFactory{}).create(o)
	ctrl := &navigatorController{
		impl: impl,
	}

	booter := bootstrapper{
		o:  o,
		nc: ctrl,
	}
	booter.init()

	return ctrl
}

type navigatorImplFactory struct{}

func (f *navigatorImplFactory) create(o *TraverseOptions) navigatorImpl {
	var impl navigatorImpl

	switch o.Store.Subscription {
	case SubscribeAny:
		impl = &universalNavigator{
			navigator: navigator{o: o, agent: &agent{
				o: o, DO_INVOKE: true,
			}},
		}

	case SubscribeFolders, SubscribeFoldersWithFiles:
		impl = &foldersNavigator{
			navigator: navigator{o: o, agent: &agent{
				o: o, DO_INVOKE: true,
			}},
		}

	case SubscribeFiles:
		impl = &filesNavigator{
			navigator: navigator{o: o, agent: &agent{
				o: o, DO_INVOKE: false,
			}},
		}
	}

	return impl
}
