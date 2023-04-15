package nav

import (
	"github.com/snivilised/extendio/internal/log"
	"github.com/snivilised/extendio/xfs/utils"

	xi18n "github.com/snivilised/extendio/i18n"
)

type navigatorFactory struct{}

func (f navigatorFactory) new(fn ...TraverseOptionFn) TraverseNavigator {
	o := composeTraverseOptions(fn...)

	if o.Callback.Fn == nil {
		panic(xi18n.NewMissingCallbackError())
	}

	impl := navigatorImplFactory{}.new(o)
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

func (f navigatorImplFactory) new(o *TraverseOptions) navigatorImpl {
	var impl navigatorImpl

	doInvoke := o.Store.Subscription != SubscribeFiles
	deFactory := directoryEntriesFactory{}
	agent := agentFactory{}.new(&agentFactoryParams{
		o:         o,
		doInvoke:  doInvoke,
		deFactory: deFactory,
	})
	logger := f.makeLogger(o)

	switch o.Store.Subscription {
	case SubscribeAny:
		impl = &universalNavigator{
			navigator: navigator{o: o,
				agent: agent,
				log:   logger,
			},
		}

	case SubscribeFolders, SubscribeFoldersWithFiles:
		impl = &foldersNavigator{
			navigator: navigator{o: o,
				agent: agent,
				log:   logger,
			},
		}

	case SubscribeFiles:
		impl = &filesNavigator{
			navigator: navigator{o: o,
				agent: agent,
				log:   logger,
			},
		}
	}

	return impl
}

func (f navigatorImplFactory) makeLogger(o *TraverseOptions) utils.RoProp[log.Logger] {
	return log.NewLogger(&log.LoggerInfo{
		Rotation: log.Rotation{
			Filename:       o.Store.Logging.Path,
			MaxSizeInMb:    o.Store.Logging.Rotation.MaxSizeInMb,
			MaxNoOfBackups: o.Store.Logging.Rotation.MaxNoOfBackups,
			MaxAgeInDays:   o.Store.Logging.Rotation.MaxAgeInDays,
		},
		Enabled:         o.Store.Logging.Enabled,
		Path:            o.Store.Logging.Path,
		TimeStampFormat: o.Store.Logging.TimeStampFormat,
		Level:           o.Store.Logging.Level,
	})
}
