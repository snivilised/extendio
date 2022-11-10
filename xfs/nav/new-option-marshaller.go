package nav

import "fmt"

func newOptionMarshaller(options *TraverseOptions) OptionMarshaller {

	if options.Persist.Restorer == nil {
		panic(MISSING_RESTORER_FN_L_ERR)
	}

	var marshaller OptionMarshaller
	switch options.Persist.Format {
	case PersistInJSONEn:
		marshaller = &OptionMarshallerJSON{
			Options: options,
		}

	default:
		panic(fmt.Errorf("unknown marshal format: '%v'", options.Persist.Format))
	}

	return marshaller
}
