package i18n

// LocalisableError is an error that is translate-able (Localisable)
// this has to be modified to implement the error interface.
type LocalisableError struct {
	Data  Localisable
	Inner error // the core error
}

func (le LocalisableError) Error() string {
	return le.Inner.Error()
}
