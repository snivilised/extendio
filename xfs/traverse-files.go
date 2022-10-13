package xfs

import "errors"

func TraverseFiles(options ...FileOptionFn) *TranslateError {
	option := FileOptions{
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
