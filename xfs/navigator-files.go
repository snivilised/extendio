package xfs

import (
	"errors"
	"io/fs"
	"os"
)

type filesNavigator struct {
	navigator
}

func (n *filesNavigator) top(frame *navigationFrame) *LocalisableError {
	info, err := os.Lstat(frame.Root)
	var le *LocalisableError = nil
	if err != nil {
		le = &LocalisableError{Inner: err}
	} else {

		if info.IsDir() {
			item := &TraverseItem{Path: frame.Root, Info: info}
			le = n.traverse(item, frame)
		} else {
			le = &LocalisableError{Inner: errors.New("Not a directory")}
		}
	}
	if (le != nil) && (le.Inner == fs.SkipDir) {
		return nil
	}
	return le
}

func (n *filesNavigator) traverse(currentItem *TraverseItem, frame *navigationFrame) *LocalisableError {
	//
	// For files, the registered callback will only be invoked for file entries. This means
	// that the client will have no way to skip the descending of a particular directory. In
	// this case, the client should use the OnDescend callback (yet to be implemented) and
	// return SkipDir from there.

	defer func() {
		_ = n.ascend(&navigationInfo{options: n.options, item: currentItem, frame: frame})
	}()
	navi := &navigationInfo{options: n.options, item: currentItem, frame: frame}
	_ = n.descend(navi)

	entries, readErr := n.children.read(currentItem)
	if (currentItem.Entry != nil) && !(currentItem.Entry.IsDir()) {
		_ = n.options.Hooks.Extend(navi, entries)

		// Effectively, this is the file only filter
		//
		return n.options.Callback(currentItem)
	}

	if exit, err := n.children.notify(&notifyInfo{
		item: currentItem, entries: entries, readErr: readErr,
	}); exit || err != nil {
		return err
	} else {
		var err error
		if err = n.options.Hooks.Sort(entries); err != nil {
			panic(LocalisableError{
				Inner: errors.New("files navigator sort function failed"),
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
