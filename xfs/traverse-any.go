package xfs

import "errors"

func TraverseAny(options ...AnyOptionFn) *TranslateError {
	option := AnyOptions{
		Options: Options{CaseSensitive: false, Extend: false},
	}

	for _, functionalOption := range options {
		functionalOption(&option)
	}

	if option.Fn == nil {
		return &TranslateError{Error: errors.New("missing callback function")}
	}

	return nil
}
