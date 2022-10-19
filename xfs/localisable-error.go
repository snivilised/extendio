package xfs

// LocalisableError is an error that is translate-able (Localisable)
type LocalisableError struct {
	Inner error // the core error
}
