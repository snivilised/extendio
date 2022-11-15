package nav

import (
	"io/fs"

	. "github.com/snivilised/extendio/translate"
)

type universalNavigator struct {
	navigator
}

func (n *universalNavigator) top(frame *navigationFrame) *LocalisableError {

	return n.agent.top(&agentTopParams{
		impl:  n,
		frame: frame,
	})
}

func (n *universalNavigator) traverse(currentItem *TraverseItem, frame *navigationFrame) *LocalisableError {
	defer func() {
		n.ascend(&NavigationInfo{Options: n.o, Item: currentItem, Frame: frame})
	}()
	navi := &NavigationInfo{Options: n.o, Item: currentItem, Frame: frame}
	n.descend(navi)
	entries, readErr := n.agent.read(currentItem, n.o.Store.Behaviours.Sort.DirectoryEntryOrder)
	// Files and Folders need to be sorted independently to preserve the navigation order
	// stipulated by .Behaviours.Sort.DirectoryEntryOrder
	//
	entries.sort(&entries.Files)
	entries.sort(&entries.Folders)
	sorted := entries.all()
	n.o.Hooks.Extend(navi, *sorted)

	if le := n.agent.proxy(currentItem, frame); le != nil || (currentItem.Entry != nil && !currentItem.Entry.IsDir()) {
		if le != nil && le.Inner == fs.SkipDir && currentItem.Entry.IsDir() {
			// Successfully skipped directory
			//
			le = nil
		}
		return le
	}

	if exit, err := n.agent.notify(&agentNotifyParams{
		frame: frame, item: currentItem, entries: *sorted, readErr: readErr,
	}); exit {
		return err
	} else {

		return n.agent.traverse(&agentTraverseParams{
			impl:     n,
			contents: sorted,
			parent:   currentItem,
			frame:    frame,
		})
	}
}
