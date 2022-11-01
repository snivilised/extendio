package nav

import (
	"errors"

	. "github.com/snivilised/extendio/translate"
)

var TERMINATE_ERR = errors.New("terminate traverse")

var NOT_DIRECTORY_ERR = errors.New("Not a directory")
var NOT_DIRECTORY_L_ERR = LocalisableError{Inner: NOT_DIRECTORY_ERR}

var FILES_NAV_SORT_ERR = errors.New("files navigator sort function failed")
var FILES_NAV_SORT_L_ERR = LocalisableError{Inner: FILES_NAV_SORT_ERR}

var FOLDERS_NAV_SORT_ERR = errors.New("folders navigator sort function failed")
var FOLDERS_NAV_SORT_L_ERR = LocalisableError{Inner: FOLDERS_NAV_SORT_ERR}

var UNIVERSAL_NAV_SORT_ERR = errors.New("universal navigator sort function failed")
var UNIVERSAL_NAV_SORT_L_ERR = LocalisableError{Inner: UNIVERSAL_NAV_SORT_ERR}

var MISSING_CALLBACK_FN_ERR = errors.New("missing callback function")
var MISSING_CALLBACK_FN_L_ERR = LocalisableError{Inner: MISSING_CALLBACK_FN_ERR}

var PATTERN_NOT_DEFINED_ERR = errors.New("pattern not defined")
var PATTERN_NOT_DEFINED_L_ERR = LocalisableError{Inner: PATTERN_NOT_DEFINED_ERR}
