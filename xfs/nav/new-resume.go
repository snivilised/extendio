package nav

import (
	"fmt"
)

type NewResumerInfo struct {
	Path     string
	Restore  PersistenceRestorer
	Strategy ResumeStrategyEnum
}

func NewResumer(info *NewResumerInfo) (Resumer, error) {
	marshaller := stateMarshallerJSON{
		restore: info.Restore,
	}
	err := marshaller.unmarshal(info.Path)

	if err != nil {
		return nil, err
	}
	o := marshaller.o

	impl := newImpl(o)
	strategy := newResumeStrategy(o, info.Strategy, marshaller.ps.Active)
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
	booter.resume(o, marshaller.ps)

	return resumerCtrl, nil
}

func newResumeStrategy(o *TraverseOptions, strategyEn ResumeStrategyEnum, active *ActiveState) resumeStrategy {

	var strategy resumeStrategy

	switch strategyEn {
	case ResumeStrategyFastwardEn:
		strategy = &fastwardStrategy{
			baseStrategy: baseStrategy{
				o:      o,
				active: active,
			},
		}
	default:
		panic(fmt.Errorf("*** newResumeStrategy: unsupported strategy: '%v'", strategyEn))
	}

	return strategy
}
