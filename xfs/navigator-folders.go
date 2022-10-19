package xfs

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/samber/lo"
)

type foldersNavigator struct {
	navigator
}

func (n *foldersNavigator) top(root string) *LocalisableError {
	info, err := os.Lstat(root)
	var le *LocalisableError = nil
	if err != nil {
		item := TraverseItem{Path: root, Info: info, Error: &LocalisableError{Inner: err}}
		le = n.options.Callback(&item)
	} else {

		if info.IsDir() {
			item := TraverseItem{Path: root, Info: info}
			le = n.traverse(&item)
		} else {
			item := TraverseItem{Path: root, Info: info, Error: &LocalisableError{Inner: errors.New("Not a directory")}}
			le = n.options.Callback(&item)
		}
	}
	if (le != nil) && (le.Inner == fs.SkipDir) {
		return nil
	}
	return le
}

func (n *foldersNavigator) traverse(currentItem *TraverseItem) *LocalisableError {
	if le := n.options.Callback(currentItem); le != nil || (currentItem.Entry != nil && !currentItem.Entry.IsDir()) {
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
		if le := n.options.Callback(item); le != nil {
			if err == fs.SkipDir && (currentItem.Entry != nil && currentItem.Entry.IsDir()) {
				err = nil
			}
			return &LocalisableError{Inner: err}
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
		le := lo.Ternary(err == nil, nil, &LocalisableError{Inner: err})
		childItem := TraverseItem{Path: childPath, Info: info, Entry: childEntry, Error: le}

		if childLe := n.traverse(&childItem); childLe != nil {
			if childLe.Inner == fs.SkipDir {
				break
			}
			return childLe
		}
	}
	return nil
}
