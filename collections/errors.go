package collections

import (
	"errors"

	. "github.com/snivilised/extendio/translate"
)

var STACK_IS_EMPTY_ERR = errors.New("stack is empty")
var STACK_IS_EMPTY_L_ERR = LocalisableError{Inner: STACK_IS_EMPTY_ERR}
