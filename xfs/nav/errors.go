package nav

import (
	"errors"

	. "github.com/snivilised/extendio/translate"
)

// This file should only contain user facing error messages. Other errors/panics
// maybe generated in this lib, but if they are internal errors or situations
// that occur due to programming mistakes, they they do not need to be defined here
// and can be constructed on the fly.
//

var TERMINATE_ERR = errors.New("terminate traverse")

var NOT_DIRECTORY_ERR = errors.New("Not a directory")
var NOT_DIRECTORY_L_ERR = LocalisableError{Inner: NOT_DIRECTORY_ERR}

var SORT_ERR = errors.New("sort function failed")
var SORT_L_ERR = LocalisableError{Inner: SORT_ERR}

var MISSING_CALLBACK_FN_ERR = errors.New("missing callback function")
var MISSING_CALLBACK_FN_L_ERR = LocalisableError{Inner: MISSING_CALLBACK_FN_ERR}
