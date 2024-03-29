package nav

import (
	"fmt"

	"github.com/snivilised/extendio/i18n"
)

type marshallerFactory struct{}

func (m *marshallerFactory) new(o *TraverseOptions, state *persistState) stateMarshaller {
	var marshaller stateMarshaller

	if o.Persist.Format != PersistInJSONEn {
		panic(i18n.NewUnknownMarshalFormatError(
			fmt.Sprintf("%v", o.Persist.Format), "Options/Persist/Format",
		))
	}

	marshaller = &stateMarshallerJSON{
		o:  o,
		ps: state,
	}

	return marshaller
}
