package utils

import (
	"os"
	"path/filepath"
)

// AbsFunc signature of function used to obtain the absolute representation of
// a path.
type AbsFunc func(path string) (string, error)

// Abs function invoker, allows a function to be used in place where
// an instance of an interface would be expected.
func (f AbsFunc) Abs(path string) (string, error) {
	return f(path)
}

// HomeUserFunc signature of function used to obtain the user's home directory.
type HomeUserFunc func() string

// Home function invoker, allows a function to be used in place where
// an instance of an interface would be expected.
func (f HomeUserFunc) Home() string {
	return f()
}

// ResolveOverrides, used to override the internal functions used
// to resolve the home path (os.UserHomeDir) and the abs path
// (filepath.Abs). In normal usage, these do not need to be overridden,
// just used for testing purposes.
type ResolveOverrides struct {
	HomeFunc HomeUserFunc
	AbsFunc  AbsFunc
}

// ResolvePath performs 2 forms of path resolution. The first is resolving a
// home path reference, via the ~ character; ~ is replaced by the user's
// home path. The second resolves ./ or ../ relative path. The overrides
// does not need to be provided
func ResolvePath(path string, overrides ...ResolveOverrides) string {
	result := path

	if len(overrides) > 0 {
		override := overrides[0]
		if result[0] == '~' {
			result = filepath.Join(override.HomeFunc(), result[1:])
		} else {
			a, err := override.AbsFunc(result)

			if err == nil {
				result = a
			}
		}
	} else {
		if result[0] == '~' {
			if h, err := os.UserHomeDir(); err == nil {
				result = filepath.Join(h, result[1:])
			}
		} else {
			if a, err := filepath.Abs(result); err == nil {
				result = a
			}
		}
	}

	return result
}
