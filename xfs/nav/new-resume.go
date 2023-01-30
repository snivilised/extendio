package nav

import (
	"fmt"
)

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

	if info.Strategy == ResumeStrategySpawnEn {
		o.Resume.Spawn = true
	}

	/*
		impl := lo.TernaryF(info.Strategy == ResumeStrategySpawnEn,
			func() navigatorImpl {
				return &spawnerImpl{
					navigator: navigator{o: o, agent: &agent{
						o: o, DO_INVOKE: true,
					}},
					ps: marshaller.ps,
				}
			},
			func() navigatorImpl {
				return (&navigatorImplFactory{}).create(o)
			},
		)
	*/
	impl := (&navigatorImplFactory{}).create(o)
	strategy := (&strategyFactory{}).create(o, info.Strategy, marshaller.ps)

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

type strategyFactory struct{}

func (f *strategyFactory) create(o *TraverseOptions, strategyEn ResumeStrategyEnum, ps *persistState) resumeStrategy {
	var strategy resumeStrategy

	switch strategyEn {

	case ResumeStrategySpawnEn:
		strategy = &spawnStrategy{
			baseStrategy: baseStrategy{
				o:  o,
				ps: ps,
			},
		}
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
