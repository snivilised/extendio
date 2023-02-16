package nav

import "fmt"

type marshallerFactory struct{}

func (m *marshallerFactory) new(o *TraverseOptions, state *persistState) stateMarshaller {
	var marshaller stateMarshaller
	switch o.Persist.Format {
	case PersistInJSONEn:
		marshaller = &stateMarshallerJSON{
			o:  o,
			ps: state,
		}

	default:
		panic(fmt.Errorf("unknown marshal format: '%v'", o.Persist.Format))
	}

	return marshaller

}
