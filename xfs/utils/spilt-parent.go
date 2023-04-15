package utils

import (
	"path/filepath"
)

func SplitParent(path string) (d, f string) {
	d = filepath.Dir(path)
	f = filepath.Base(path)

	return d, f
}
