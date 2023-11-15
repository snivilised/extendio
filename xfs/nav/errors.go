package nav

import (
	"errors"
	"fmt"
)

// Errors defined here are internal errors that are of no value to end
// users (hence not l10n). There are usually programming errors which
// means they only have meaning for client developers.

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

// static errors, identifiable with errors.Is

// ErrUndefinedSubscriptionType indicates client has not set the navigation
// subscription type at /Options.Store.Subscription.
var ErrUndefinedSubscriptionType = errors.New(
	"undefined subscription type; please set in traverse options (/Options.Store.Subscription)",
)
