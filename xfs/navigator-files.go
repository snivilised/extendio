package xfs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/samber/lo"
)

type filesNavigator struct {
	navigator
}

func (n *filesNavigator) top(root string) *LocalisableError {
	fmt.Printf("---> 🛩️ [filesNavigator]::top\n")

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
	fmt.Printf("---> 🛩️ [filesNavigator]::traverse\n")

	if (currentItem.Entry != nil) && !(currentItem.Entry.IsDir()) {
		return n.options.Callback(currentItem)
	}

	entries, err := readDir(currentItem.Path)
	if err != nil {
		return &LocalisableError{Inner: err}
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
