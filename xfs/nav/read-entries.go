package nav

import (
	"io/fs"
	"os"
)

// ReadEntries reads the contents of a directory. The resulting
// slice is left un-sorted
func ReadEntries(dirname string) ([]fs.DirEntry, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dirs, err := f.ReadDir(-1)
	if err != nil {
		return nil, err
	}

	return dirs, nil
}
