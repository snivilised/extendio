package xfs

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/samber/lo"
)

type filesNavigator struct {
	navigator
}

func (n *filesNavigator) top(root string) *LocalisableError {
	info, err := os.Lstat(root)
	var le *LocalisableError = nil
	if err != nil {
		le = &LocalisableError{Inner: err}
	} else {

		if info.IsDir() {
			item := TraverseItem{Path: root, Info: info}
			le = n.traverse(&item)
		} else {
			le = &LocalisableError{Inner: errors.New("Not a directory")}
		}
	}
	if (le != nil) && (le.Inner == fs.SkipDir) {
		return nil
	}
	return le
}

func (n *filesNavigator) traverse(currentItem *TraverseItem) *LocalisableError {
	//
	// For files, the registered callback will only be invoked for file entries. This means
	// that the client will have no way to skip the descending of a particular directory. In
	// this case, the client should use the OnDescend callback (yet to be implemented) and
	// return SkipDir from there.

	if (currentItem.Entry != nil) && !(currentItem.Entry.IsDir()) {
		// Effectively, this is the file only filter
		//
		return n.options.Callback(currentItem)
	}

	entries, err := n.options.Hooks.ReadDirectory(currentItem.Path)
	if err != nil {
		return &LocalisableError{Inner: err}
	}

	if entries, err = n.options.Hooks.Sort(entries); err != nil {
		panic(LocalisableError{
			Inner: errors.New("files navigator sort function failed"),
		})
	}

	for _, childEntry := range entries {
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
