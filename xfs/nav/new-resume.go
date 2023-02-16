package nav

import (
	"fmt"
)

// ResumerInfo
type ResumerInfo struct {
	RestorePath string
	Restorer    PersistenceRestorer
	Strategy    ResumeStrategyEnum
}

type resumerFactory struct{}

func (f resumerFactory) new(info *ResumerInfo) (*resumeController, error) {
	marshaller := stateMarshallerJSON{
		restore: info.Restorer,
	}
	err := marshaller.unmarshal(info.RestorePath)

	if err != nil {
		return nil, err
	}
	o := marshaller.o

	impl := navigatorImplFactory{}.new(o)
	navigator := &navigatorController{
		impl: impl,
	}

	strategy := strategyFactory{}.new(&createStrategyParams{
		o:          o,
		strategyEn: info.Strategy,
		ps:         marshaller.ps,
		nc:         navigator,
	})

	resumerCtrl := &resumeController{
		navigator: navigator,
		ps:        marshaller.ps,
		strategy:  strategy,
	}

	booter := bootstrapper{
		o:  o,
		nc: navigator,
		rc: resumerCtrl,
	}
	booter.init()
	booter.initResume(o, marshaller.ps)

	return resumerCtrl, nil
}

type strategyFactory struct{}

type createStrategyParams struct {
	o          *TraverseOptions
	strategyEn ResumeStrategyEnum
	ps         *persistState
	nc         *navigatorController
}

func (f strategyFactory) new(params *createStrategyParams) resumeStrategy {
	var strategy resumeStrategy
	deFactory := &directoryEntriesFactory{}

	switch params.strategyEn {

	case ResumeStrategySpawnEn:
		strategy = &spawnStrategy{
			baseStrategy: baseStrategy{
				o:         params.o,
				ps:        params.ps,
				nc:        params.nc,
				deFactory: deFactory,
			},
		}
	case ResumeStrategyFastwardEn:
		strategy = &fastwardStrategy{
			baseStrategy: baseStrategy{
				o:         params.o,
				ps:        params.ps,
				nc:        params.nc,
				deFactory: deFactory,
			},
		}

	default:
		panic(fmt.Errorf("*** newResumeStrategy: unsupported strategy: '%v'", params.strategyEn))
	}

	return strategy
}
