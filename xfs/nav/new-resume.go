package nav

func NewResumer(path string, restore PersistenceRestorer) (Resumer, error) {
	marshaller := stateMarshallerJSON{
		restore: restore,
	}

	if err := marshaller.unmarshal(path); err == nil {

		impl := newImpl(marshaller.o)
		ctrl := &resumeController{
			navigator: &navigatorController{
				impl: impl,
			},
			ps: marshaller.ps,
		}
		ctrl.init(func(params *listenerInitParams) {})

		return ctrl, nil
	} else {
		return nil, err
	}
}
