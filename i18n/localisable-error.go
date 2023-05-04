package i18n

import (
	"reflect"
)

// LocalisableError is an error that is translate-able (Localisable)
type LocalisableError struct {
	Data Localisable
}

func (le LocalisableError) Error() string {
	return Text(le.Data)
}

func QueryGeneric[T any](method string, target error) bool {
	if target == nil {
		return false
	}

	nativeIf, ok := target.(T)

	if !ok {
		return false
	}

	none := []reflect.Value{}

	return reflect.ValueOf(&nativeIf).Elem().MethodByName(method).Call(none)[0].Bool()
}
