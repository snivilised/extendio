package xfs

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

// func composeOptions[O struct{ Options }](fn ...func(o *O)) O {

// 	options := O{
// 		Options: Options{CaseSensitive: false, Extend: false},
// 	}

// 	for _, functionalOption := range fn {
// 		functionalOption(&options)
// 	}
// 	return options
// }

// turn into generic?
func composeAnyOptions(fn ...AnyOptionFn) AnyOptions {
	options := AnyOptions{
		GenericOptions: GenericOptions{CaseSensitive: false, Extend: false},
	}

	for _, functionalOption := range fn {
		functionalOption(&options)
	}
	return options
}

func TraverseAny(root string, fn ...AnyOptionFn) *LocalisableError {
	options := composeAnyOptions(fn...)
	// options := composeOptions[AnyOptions](fn...)

	if options.Callback == nil {
		return &LocalisableError{Inner: errors.New("missing callback function")}
	}

	info, err := os.Lstat(root)
	var le *LocalisableError = nil
	if err != nil {
		item := TraverseItem{Path: root, Error: &LocalisableError{Inner: err}}
		le = options.Callback(&item)
	} else {
		item := TraverseItem{Path: root, Info: info}
		le = traverseAny(&options, &item)
	}
	if (le != nil) && (le.Inner == fs.SkipDir) {
		return nil
	}
	return le
}

func traverseAny(options *AnyOptions, currentItem *TraverseItem) *LocalisableError {

	if le := options.Callback(currentItem); le != nil || (currentItem.Entry != nil && !currentItem.Entry.IsDir()) {
		if le != nil && le.Inner == fs.SkipDir && currentItem.Entry.IsDir() {
			// Successfully skipped directory
			//
			le = nil
		}
		return le
	}

	entries, err := readDir(currentItem.Path)
	if err != nil {
		item := currentItem.Clone()
		item.Error = &LocalisableError{Inner: err}

		// Second call, to report ReadDir error
		//
		if le := options.Callback(item); le != nil {
			if err == fs.SkipDir && (currentItem.Entry != nil && currentItem.Entry.IsDir()) {
				err = nil
			}
			return &LocalisableError{Inner: err}
		}
	}

	for _, childEntry := range entries {
		childPath := filepath.Join(currentItem.Path, childEntry.Name())
		childItem := TraverseItem{Path: childPath, Entry: childEntry}

		if childLe := traverseAny(options, &childItem); childLe != nil {
			if childLe.Inner == fs.SkipDir {
				break
			}
			return childLe
		}
	}
	return nil
}
