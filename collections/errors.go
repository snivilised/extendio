package collections

import (
	"fmt"
)

// ❌ Stack Is Empty (internal error)

// NewStackIsEmptyNativeError creates an untranslated error to
// indicate stack is empty (internal error)
func NewStackIsEmptyNativeError() error {
	return fmt.Errorf("internal: stack is empty")
}
