package nav

import "fmt"

func newStateMarshaler(o *TraverseOptions, state *persistState) stateMarshaller {
	if o.Persist.Restore == nil {
		panic(MISSING_RESTORER_FN_L_ERR)
	}

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
