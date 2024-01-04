package nav

import (
	"github.com/snivilised/extendio/i18n"
)

type navigatorFactory struct{}

func (f navigatorFactory) new(o *TraverseOptions) TraverseNavigator {
	impl := navigatorImplFactory{}.new(o)
	nc := &navigationController{
		impl: impl,
	}

	booter := bootstrapper{
		o:  o,
		nc: nc,
	}
	booter.init()

	return nc
}

func (f navigatorFactory) fromOptionsFn(fn ...TraverseOptionFn) TraverseNavigator {
	o := composeTraverseOptions(fn...)

	if o.Callback.Fn == nil {
		panic(i18n.NewMissingCallbackError())
	}

	return f.new(o)
}

func (f navigatorFactory) fromProvidedOptions(o *TraverseOptions) TraverseNavigator {
	nav := f.new(o)

	o.afterUserOptions()

	return nav
}

type navigatorImplFactory struct{}

func (f navigatorImplFactory) new(o *TraverseOptions) navigatorImpl {
	var (
		impl                 navigatorImpl
		samplingActive       = o.Store.Sampling.SampleType != SampleTypeUnsetEn
		filteringActive      = o.isFilteringActive()
		samplingFilterActive = samplingActive && filteringActive
		doInvoke             = o.Store.Subscription != SubscribeFiles
		agent                = newAgent(&newAgentParams{
			o:                    o,
			doInvoke:             doInvoke,
			handler:              &notifyCallbackErrorHandler{},
			samplingFilterActive: samplingFilterActive,
		})
		n = navigator{
			o:                    o,
			agent:                agent,
			samplingActive:       samplingActive,
			filteringActive:      filteringActive,
			samplingFilterActive: samplingFilterActive,
		}
	)

	switch o.Store.Subscription {
	case SubscribeAny:
		impl = &universalNavigator{
			navigator: n,
		}

	case SubscribeFolders, SubscribeFoldersWithFiles:
		impl = &foldersNavigator{
			navigator: n,
		}

	case SubscribeFiles:
		impl = &filesNavigator{
			navigator: n,
		}
	default:
		panic(ErrUndefinedSubscriptionType)
	}

	return impl
}
