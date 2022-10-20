package xfs

// LocalisableError is an error that is translate-able (Localisable)
// this has to be modified to implement the error interface. IE there
// should be no distinction between LocalisableError and error
type LocalisableError struct {
	Inner error // the core error
}

// TODO: implement translation
func (le LocalisableError) Error() string {
	return le.Inner.Error()
}
