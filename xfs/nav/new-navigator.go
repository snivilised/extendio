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
	doInvoke := o.Store.Subscription != SubscribeFiles

	var impl navigatorImpl
	agtFactory := &agentFactory{}
	deFactory := directoryEntriesFactory{}
	agent := agtFactory.construct(&agentFactoryParams{
		o:         o,
		doInvoke:  doInvoke,
		deFactory: &deFactory,
	})

	switch o.Store.Subscription {
	case SubscribeAny:
		impl = &universalNavigator{
			navigator: navigator{o: o, agent: agent},
		}

	case SubscribeFolders, SubscribeFoldersWithFiles:
		impl = &foldersNavigator{
			navigator: navigator{o: o, agent: agent},
		}

	case SubscribeFiles:
		impl = &filesNavigator{
			navigator: navigator{o: o, agent: agent},
		}
	}

	return impl
}
