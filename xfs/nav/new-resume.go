package nav

import (
	"fmt"
)

type NewResumerInfo struct {
	RestorePath string
	Restorer    PersistenceRestorer
	Strategy    ResumeStrategyEnum
}

func newResumer(info *NewResumerInfo) (resumer, error) {
	marshaller := stateMarshallerJSON{
		restore: info.Restorer,
	}
	err := marshaller.unmarshal(info.RestorePath)

	if err != nil {
		return nil, err
	}
	o := marshaller.o

	impl := newImpl(o)
	strategy := newStrategy(o, info.Strategy, marshaller.ps)
	navigator := &navigatorController{
		impl: impl,
	}

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

func newStrategy(o *TraverseOptions, strategyEn ResumeStrategyEnum, ps *persistState) resumeStrategy {

	var strategy resumeStrategy

	switch strategyEn {
	case ResumeStrategyFastwardEn:
		strategy = &fastwardStrategy{
			baseStrategy: baseStrategy{
				o:  o,
				ps: ps,
			},
		}
	default:
		panic(fmt.Errorf("*** newResumeStrategy: unsupported strategy: '%v'", strategyEn))
	}

	return strategy
}
