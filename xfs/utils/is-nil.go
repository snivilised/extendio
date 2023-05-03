package utils

import (
	"reflect"
)

func IsNil(i interface{}) bool {
	value := reflect.ValueOf(i)
	if !value.IsValid() {
		return true
	}

	kind := value.Kind()

	switch kind { //nolint:exhaustive // default case IS present to handle all other cases
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.Interface, reflect.Slice:
		return reflect.ValueOf(i).IsNil()

	default:
		return false
	}
}
