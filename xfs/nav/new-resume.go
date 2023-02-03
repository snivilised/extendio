package nav

import (
	"fmt"
)

// TODO: I don't like the name NewResumerInfo, it would be better called
// just ResumerInfo
type NewResumerInfo struct {
	RestorePath string
	Restorer    PersistenceRestorer
	Strategy    ResumeStrategyEnum
}

type resumerFactory struct{}

func (f *resumerFactory) create(info *NewResumerInfo) (resumer, error) {
	marshaller := stateMarshallerJSON{
		restore: info.Restorer,
	}
	err := marshaller.unmarshal(info.RestorePath)

	if err != nil {
		return nil, err
	}
	o := marshaller.o

	impl := (&navigatorImplFactory{}).create(o)
	navigator := &navigatorController{
		impl: impl,
	}

	strategy := (&strategyFactory{}).create(&createStrategyParams{
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

func (f *strategyFactory) create(params *createStrategyParams) resumeStrategy {
	var strategy resumeStrategy

	switch params.strategyEn {

	case ResumeStrategySpawnEn:
		strategy = &spawnStrategy{
			baseStrategy: baseStrategy{
				o:  params.o,
				ps: params.ps,
				nc: params.nc,
			},
		}
	case ResumeStrategyFastwardEn:
		strategy = &fastwardStrategy{
			baseStrategy: baseStrategy{
				o:  params.o,
				ps: params.ps,
				nc: params.nc,
			},
		}

	default:
		panic(fmt.Errorf("*** newResumeStrategy: unsupported strategy: '%v'", params.strategyEn))
	}

	return strategy
}
