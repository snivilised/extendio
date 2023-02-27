package nav

import (
	"fmt"

	. "github.com/snivilised/extendio/i18n"
)

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
		// NewUnknownMarshalFormatError()
		panic(NewUnknownMarshalFormatError(
			fmt.Sprintf("%v", o.Persist.Format), "Options/Persist/Format",
		))
	}

	return marshaller

}
