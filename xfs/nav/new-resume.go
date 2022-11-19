package nav

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
		strategy := &dummyResumeStrategy{}
		ctrl := &resumeController{
			navigator: &navigatorController{
				impl: impl,
			},
			ps:       marshaller.ps,
			strategy: strategy,
		}
		ctrl.init()

		return ctrl, nil
	} else {
		return nil, err
	}
}
