package xfs

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/samber/lo"
)

func composeFileOptions(fn ...FileOptionFn) FileOptions {
	options := FileOptions{
		GenericOptions: GenericOptions{CaseSensitive: false, Extend: false},
	}

	for _, functionalOption := range fn {
		functionalOption(&options)
	}
	return options
}

func TraverseFiles(root string, fn ...FileOptionFn) *LocalisableError {
	options := composeFileOptions(fn...)

	if options.Callback == nil {
		return &LocalisableError{Error: errors.New("missing callback function")}
	}

	info, err := os.Lstat(root)
	var le *LocalisableError = nil
	if err != nil {
		le = &LocalisableError{Error: err}
	} else {

		if info.IsDir() {
			item := TraverseItem{Path: root, Info: info}
			le = traverseFiles(&options, &item)
		} else {
			le = &LocalisableError{Error: errors.New("Not a directory")}
		}
	}
	if (le != nil) && (le.Error == fs.SkipDir) {
		return nil
	}
	return le
}

// The main difference between traverseFiles and traverseFolders is the directory entry
// filter function so this should be re-factored out. The other main difference when
// the callback is invoked, which we do only for file entries, but still recurse
// on the folder entries.
// For files, the registered callback will only be invoked for file entries. This means
// that the client will have no way to skip the descending of a particular directory. In
// this case, the client should use the OnDescend callback (yet to be implemented) and
// return SkipDir from there.
func traverseFiles(options *FileOptions, currentItem *TraverseItem) *LocalisableError {

	if (currentItem.Entry != nil) && !(currentItem.Entry.IsDir()) {
		return options.Callback(currentItem)
	}

	entries, err := readDir(currentItem.Path)
	if err != nil {
		return &LocalisableError{Error: err}
	}

	for _, childEntry := range entries {
		childPath := filepath.Join(currentItem.Path, childEntry.Name())
		info, err := childEntry.Info()
		le := lo.Ternary(err == nil, nil, &LocalisableError{Error: err})
		childItem := TraverseItem{Path: childPath, Info: info, Entry: childEntry, Error: le}

		if childLe := traverseFiles(options, &childItem); childLe != nil {
			if childLe.Error == fs.SkipDir {
				break
			}
			return childLe
		}
	}
	return nil
}
