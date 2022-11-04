package nav

import (
	"errors"

	. "github.com/snivilised/extendio/translate"
)

var TERMINATE_ERR = errors.New("terminate traverse")

var NOT_DIRECTORY_ERR = errors.New("Not a directory")
var NOT_DIRECTORY_L_ERR = LocalisableError{Inner: NOT_DIRECTORY_ERR}

var SORT_ERR = errors.New("sort function failed")
var SORT_L_ERR = LocalisableError{Inner: SORT_ERR}

var MISSING_CALLBACK_FN_ERR = errors.New("missing callback function")
var MISSING_CALLBACK_FN_L_ERR = LocalisableError{Inner: MISSING_CALLBACK_FN_ERR}

var PATTERN_NOT_DEFINED_ERR = errors.New("pattern not defined")
var PATTERN_NOT_DEFINED_L_ERR = LocalisableError{Inner: PATTERN_NOT_DEFINED_ERR}
