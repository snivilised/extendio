package utils

import (
	"path/filepath"
)

func SplitParent(path string) (string, string) {

	d := filepath.Dir(path)
	f := filepath.Base(path)

	return d, f
}
