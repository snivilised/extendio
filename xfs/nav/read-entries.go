package nav

import (
	"io/fs"
	"os"

	"github.com/samber/lo"
)

// ReadEntriesHookFn reads the contents of a directory. The resulting
// slice is left un-sorted
func ReadEntriesHookFn(dirname string) ([]fs.DirEntry, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	contents, err := f.ReadDir(-1)
	if err != nil {
		return nil, err
	}

	return lo.Filter(contents, func(item fs.DirEntry, _ int) bool {
		return item.Name() != ".DS_Store"
	}), nil
}
