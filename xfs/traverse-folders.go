package xfs

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/samber/lo"
)

func composeFolderOptions(fn ...FolderOptionFn) FolderOptions {
	options := FolderOptions{
		GenericOptions: GenericOptions{CaseSensitive: false, Extend: false},
	}

	for _, functionalOption := range fn {
		functionalOption(&options)
	}
	return options
}

// return a traversal result containing traversal stats if requested?
func TraverseFolders(root string, fn ...FolderOptionFn) *LocalisableError {
	// navigator.Top
	// navigator.traverse
	options := composeFolderOptions(fn...)

	if options.Callback == nil {
		return &LocalisableError{Error: errors.New("missing callback function")}
	}

	info, err := os.Lstat(root)
	var le *LocalisableError = nil
	if err != nil {
		item := TraverseItem{Path: root, Info: info, Error: &LocalisableError{Error: err}}
		le = options.Callback(&item)
	} else {

		if info.IsDir() {
			item := TraverseItem{Path: root, Info: info}
			le = traverseFolders(&options, &item)
		} else {
			item := TraverseItem{Path: root, Info: info, Error: &LocalisableError{Error: errors.New("Not a directory")}}
			le = options.Callback(&item)
		}
	}
	if (le != nil) && (le.Error == fs.SkipDir) {
		return nil
	}
	return le
}

func traverseFolders(options *FolderOptions, currentItem *TraverseItem) *LocalisableError {

	if le := options.Callback(currentItem); le != nil || (currentItem.Entry != nil && !currentItem.Entry.IsDir()) {
		if le != nil && le.Error == fs.SkipDir && currentItem.Entry.IsDir() {
			// Successfully skipped directory
			//
			le = nil
		}
		return le
	}

	entries, err := readDir(currentItem.Path)
	if err != nil {
		item := currentItem.Clone()
		item.Error = &LocalisableError{Error: err}

		// Second call, to report ReadDir error
		//
		if le := options.Callback(item); le != nil {
			if err == fs.SkipDir && (currentItem.Entry != nil && currentItem.Entry.IsDir()) {
				err = nil
			}
			return &LocalisableError{Error: err}
		}
	}

	// this should be extracted away into a directory-entry filter
	//
	filtered := lo.Filter(entries, func(de fs.DirEntry, i int) bool {
		return de.Type().IsDir()
	})

	for _, childEntry := range filtered {
		childPath := filepath.Join(currentItem.Path, childEntry.Name())
		info, err := childEntry.Info()
		le := lo.Ternary(err == nil, nil, &LocalisableError{Error: err})
		childItem := TraverseItem{Path: childPath, Info: info, Entry: childEntry, Error: le}

		if childLe := traverseFolders(options, &childItem); childLe != nil {
			if childLe.Error == fs.SkipDir {
				break
			}
			return childLe
		}
	}
	return nil
}
