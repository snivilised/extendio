package utils

// Must
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
