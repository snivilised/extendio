package nav

import (
	"fmt"

	xi18n "github.com/snivilised/extendio/i18n"
)

type resumerFactory struct{}

func (f resumerFactory) new(info *Resumption) (*resumeStrategyController, error) {
	marshaller := stateMarshallerJSON{
		restore: info.Restorer,
	}
	err := marshaller.unmarshal(info.RestorePath)

	if err != nil {
		return nil, err
	}

	o := marshaller.o
	impl := navigatorImplFactory{}.new(o)
	nc := &navigationController{
		impl: impl,
	}

	strategy := strategyFactory{}.new(&createStrategyParams{
		o:          o,
		strategyEn: info.Strategy,
		ps:         marshaller.ps,
		nc:         nc,
	})

	rc := &resumeStrategyController{
		nc:       nc,
		ps:       marshaller.ps,
		strategy: strategy,
	}

	booter := bootstrapper{
		o:  o,
		nc: nc,
		rc: rc,
	}
	booter.init()
	booter.initResume(marshaller.ps)

	return rc, nil
}

type strategyFactory struct{}

type createStrategyParams struct {
	o          *TraverseOptions
	strategyEn ResumeStrategyEnum
	ps         *persistState
	nc         *navigationController
}

func (f strategyFactory) new(params *createStrategyParams) resumeStrategy {
	var strategy resumeStrategy

	switch params.strategyEn { //nolint:exhaustive // default case is present
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
		panic(xi18n.NewInvalidResumeStrategyError(fmt.Sprintf("%v", params.strategyEn)))
	}

	return strategy
}
