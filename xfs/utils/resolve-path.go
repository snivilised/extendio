package utils

import (
	"os"
	"path/filepath"

	"github.com/samber/lo"
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
type HomeUserFunc func() (string, error)

// Home function invoker, allows a function to be used in place where
// an instance of an interface would be expected.
func (f HomeUserFunc) Home() (string, error) {
	return f()
}

// ResolveMocks, used to override the internal functions used
// to resolve the home path (os.UserHomeDir) and the abs path
// (filepath.Abs). In normal usage, these do not need to be provided,
// just used for testing purposes.
type ResolveMocks struct {
	HomeFunc HomeUserFunc
	AbsFunc  AbsFunc
}

// ResolvePath performs 2 forms of path resolution. The first is resolving a
// home path reference, via the ~ character; ~ is replaced by the user's
// home path. The second resolves ./ or ../ relative path. The overrides
// do not need to be provided.
func ResolvePath(path string, mocks ...ResolveMocks) string {
	result := path

	if len(mocks) > 0 {
		m := mocks[0]
		result = lo.TernaryF(result[0] == '~',
			func() string {
				if h, err := m.HomeFunc(); err == nil {
					return filepath.Join(h, result[1:])
				}

				return path
			},
			func() string {
				if a, err := m.AbsFunc(result); err == nil {
					return a
				}

				return path
			},
		)
	} else {
		result = lo.TernaryF(result[0] == '~',
			func() string {
				if h, err := os.UserHomeDir(); err == nil {
					return filepath.Join(h, result[1:])
				}

				return path
			},
			func() string {
				if a, err := filepath.Abs(result); err == nil {
					return a
				}

				return path
			},
		)
	}

	return result
}
