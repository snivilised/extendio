package xfs

import (
	"errors"
	"io/fs"
	"os"
)

type universalNavigator struct {
	navigator
}

func (n *universalNavigator) top(frame *navigationFrame) *LocalisableError {

	info, err := os.Lstat(frame.Root)
	var le *LocalisableError = nil
	if err != nil {
		item := &TraverseItem{Path: frame.Root, Info: info, Error: &LocalisableError{Inner: err}}

		le = n.options.Callback(item)
	} else {
		item := &TraverseItem{Path: frame.Root, Info: info}
		le = n.traverse(item, frame)
	}
	if (le != nil) && (le.Inner == fs.SkipDir) {
		return nil
	}
	return le
}

func (n *universalNavigator) traverse(currentItem *TraverseItem, frame *navigationFrame) *LocalisableError {
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
		var err error
		if err = n.options.Hooks.Sort(entries); err != nil {
			panic(LocalisableError{
				Inner: errors.New("universal navigator sort function failed"),
			})
		}

		return n.children.traverse(&agentTraverseInfo{
			core:    n,
			entries: entries,
			parent:  currentItem,
			frame:   frame,
		})
	}
}
