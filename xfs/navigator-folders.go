package xfs

import (
	"errors"
	"io/fs"
	"os"

	"github.com/samber/lo"
)

type foldersNavigator struct {
	navigator
}

func (n *foldersNavigator) top(frame *navigationFrame) *LocalisableError {
	info, err := os.Lstat(frame.Root)
	var le *LocalisableError = nil
	if err != nil {
		item := &TraverseItem{Path: frame.Root, Info: info, Error: &LocalisableError{Inner: err}}
		n.options.Hooks.Extend(&navigationInfo{
			options: n.options, item: item, frame: frame,
		}, []fs.DirEntry{})
		le = n.options.Callback(item)
	} else {

		if info.IsDir() {
			item := &TraverseItem{Path: frame.Root, Info: info}
			le = n.traverse(item, frame)
		} else {
			item := &TraverseItem{
				Path: frame.Root, Info: info, Error: &LocalisableError{Inner: errors.New("not a directory")},
			}
			n.options.Hooks.Extend(&navigationInfo{
				options: n.options, item: item, frame: frame,
			}, []fs.DirEntry{})
			le = n.options.Callback(item)
		}
	}
	if (le != nil) && (le.Inner == fs.SkipDir) {
		return nil
	}
	return le
}

func (n *foldersNavigator) traverse(currentItem *TraverseItem, frame *navigationFrame) *LocalisableError {
	defer func() {
		n.ascend(&navigationInfo{options: n.options, item: currentItem, frame: frame})
	}()
	navi := &navigationInfo{options: n.options, item: currentItem, frame: frame}
	n.descend(navi)
	entries, readErr := n.children.read(currentItem)
	n.options.Hooks.Extend(navi, entries)

	if le := n.options.Callback(currentItem); le != nil || (currentItem.Entry != nil && !currentItem.Entry.IsDir()) {
		if le != nil && le.Inner == fs.SkipDir && currentItem.Entry.IsDir() {
			// Successfully skipped directory
			//
			le = nil
		}
		return le
	}

	if exit, err := n.children.notify(&notifyInfo{
		item: currentItem, entries: entries, readErr: readErr,
	}); exit {
		return err
	} else {
		dirs := lo.Filter(entries, func(de fs.DirEntry, i int) bool {
			return de.Type().IsDir()
		})

		var err error
		if err = n.options.Hooks.Sort(dirs); err != nil {
			panic(LocalisableError{
				Inner: errors.New("folder navigator sort function failed"),
			})
		}

		return n.children.traverse(&agentTraverseInfo{
			core:    n,
			entries: dirs,
			parent:  currentItem,
			frame:   frame,
		})
	}
}
