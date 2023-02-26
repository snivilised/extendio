package nav

import (
	"github.com/samber/lo"
	"github.com/snivilised/extendio/xfs/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	. "github.com/snivilised/extendio/i18n"
)

type navigatorFactory struct{}

func (f navigatorFactory) new(fn ...TraverseOptionFn) TraverseNavigator {
	o := composeTraverseOptions(fn...)

	if o.Callback.Fn == nil {
		panic(NewMissingCallbackError())
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
	doInvoke := o.Store.Subscription != SubscribeFiles

	var impl navigatorImpl
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

func (f navigatorImplFactory) makeLogger(o *TraverseOptions) utils.RoProp[*zap.Logger] {

	return utils.NewRoProp(lo.TernaryF(o.Store.Logging.Enabled,
		func() *zap.Logger {
			if o.Store.Logging.Path == "" {
				panic(NewInvalidConfigEntryError(o.Store.Logging.Path, "Store/Logging/Path"))
			}
			ws := zapcore.AddSync(&lumberjack.Logger{
				Filename:   o.Store.Logging.Path,
				MaxSize:    o.Store.Logging.Rotation.MaxSizeInMb,
				MaxBackups: o.Store.Logging.Rotation.MaxNoOfBackups,
				MaxAge:     o.Store.Logging.Rotation.MaxAgeInDays,
			})
			config := zap.NewProductionEncoderConfig()
			config.EncodeTime = zapcore.TimeEncoderOfLayout(o.Store.Logging.TimeStampFormat)
			core := zapcore.NewCore(
				zapcore.NewJSONEncoder(config),
				ws,
				o.Store.Logging.Level,
			)
			return zap.New(core)
		}, func() *zap.Logger {
			return zap.NewNop()
		}))
}
