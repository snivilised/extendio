package xfs

import (
	"io/fs"
	"os"
)

// ReadEntries readds the contents of a directory. The resulting
// slice is left un-sorted
func ReadEntries(dirname string) ([]fs.DirEntry, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	dirs, err := f.ReadDir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}

	return dirs, nil
}
