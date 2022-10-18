package xfs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type universalNavigator struct {
	navigator
}

func (n *universalNavigator) top(root string) *LocalisableError {
	fmt.Printf("---> 🚀 [universalNavigator]::top\n")

	if n.options.Callback == nil {
		return &LocalisableError{Inner: errors.New("missing callback function")}
	}

	info, err := os.Lstat(root)
	var le *LocalisableError = nil
	if err != nil {
		item := TraverseItem{Path: root, Error: &LocalisableError{Inner: err}}
		le = n.options.Callback(&item)
	} else {
		item := TraverseItem{Path: root, Info: info}
		le = n.traverse(&item)
	}
	if (le != nil) && (le.Inner == fs.SkipDir) {
		return nil
	}
	return le
}

func (n *universalNavigator) traverse(currentItem *TraverseItem) *LocalisableError {
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
				return nil
			}
			return &LocalisableError{Inner: err}
		}
	}

	for _, childEntry := range entries {
		childPath := filepath.Join(currentItem.Path, childEntry.Name())
		childItem := TraverseItem{Path: childPath, Entry: childEntry}

		if childLe := n.traverse(&childItem); childLe != nil {
			if childLe.Inner == fs.SkipDir {
				break
			}
			return childLe
		}
	}
	return nil
}
