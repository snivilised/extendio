package nav

import (
	"fmt"

	xi18n "github.com/snivilised/extendio/i18n"
)

type marshallerFactory struct{}

func (m *marshallerFactory) new(o *TraverseOptions, state *persistState) stateMarshaller {
	var marshaller stateMarshaller

	if o.Persist.Format != PersistInJSONEn {
		panic(xi18n.NewUnknownMarshalFormatError(
			fmt.Sprintf("%v", o.Persist.Format), "Options/Persist/Format",
		))
	}

	marshaller = &stateMarshallerJSON{
		o:  o,
		ps: state,
	}

	return marshaller
}
