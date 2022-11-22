package nav

import (
	"fmt"
)

type NewResumerInfo struct {
	Path     string
	Restore  PersistenceRestorer
	Strategy ResumeStrategyEnum
}

func NewResumer(info NewResumerInfo) (Resumer, error) {
	marshaller := stateMarshallerJSON{
		restore: info.Restore,
	}

	if err := marshaller.unmarshal(info.Path); err == nil {

		impl := newImpl(marshaller.o)
		strategy := newResumeStrategy(info.Strategy)
		navigator := &navigatorController{
			impl: impl,
		}
		navigator.init()

		resumerCtrl := &resumeController{
			navigator: navigator,
			ps:        marshaller.ps,
			strategy:  strategy,
		}
		resumerCtrl.init()

		return resumerCtrl, nil
	} else {
		return nil, err
	}
}

func newResumeStrategy(strategyEn ResumeStrategyEnum) resumeStrategy {

	var strategy resumeStrategy

	switch strategyEn {
	case ResumeStrategyFastwardEn:
		strategy = &fastwardStrategy{}
	default:
		panic(fmt.Errorf("*** newResumeStrategy: unsupported strategy: '%v'", strategyEn))
	}

	return strategy
}
