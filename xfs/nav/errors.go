package nav

import (
	"fmt"
	"io/fs"
)

// ❌ Invalid Notification Mute Requested

// NewInvalidNotificationMuteRequestedNativeError creates an untranslated error to
// indicate invalid notification mute requested (internal error)
func NewInvalidNotificationMuteRequestedNativeError(value string) error {
	return fmt.Errorf("internal: invalid notification mute requested (%v)", value)
}

// ❌ Invalid Resume State Transition Detected

// NewItemAlreadyExtendedNativeError creates an untranslated error to
// indicate invalid resume state transition occurred (internal error)
func NewInvalidResumeStateTransitionNativeError(state string) error {
	return fmt.Errorf("internal: invalid resume state transition detected (%v)", state)
}

// ❌ Item already extended

// NewItemAlreadyExtendedNativeError creates an untranslated error to
// indicate traverse-item already extended (internal error)
func NewItemAlreadyExtendedNativeError(path string) error {
	return fmt.Errorf("internal: item already extended for item at: '%v'", path)
}

// ❌ Missing listen detacher function

// NewMissingListenDetacherFunctionNativeError creates an untranslated error to
// indicate invalid resume state transition occurred (internal error)
func NewMissingListenDetacherFunctionNativeError(state string) error {
	return fmt.Errorf("internal: missing listen detacher function (%v)", state)
}

// ❌ Invalid Periscope Root Path

// NewInvalidPeriscopeRootPathNativeError creates an untranslated error to
// indicate invalid resume state transition occurred (internal error)
func NewInvalidPeriscopeRootPathNativeError(root, current string) error {
	return fmt.Errorf("internal: root path '%v' can't be longer than current '%v'", root, current)
}

// ❌ Resume controller not set

// NewResumeControllerNotSetNativeError creates an untranslated error to
// indicate resume controller not set (internal error)
func NewResumeControllerNotSetNativeError(from string) error {
	return fmt.Errorf("internal: resume controller not set (from: '%v')", from)
}

// ❌ SkipDir

// QuerySkipDirError query if error is the fs SkipDir error
func QuerySkipDirError(target error) bool {
	return target != nil && target == fs.SkipDir
}
