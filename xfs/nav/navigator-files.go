package nav

import (
	. "github.com/snivilised/extendio/translate"
)

type filesNavigator struct {
	navigator
}

func (n *filesNavigator) top(frame *navigationFrame) *LocalisableError {

	return n.agent.top(&agentTopParams{
		impl:  n,
		frame: frame,
	})
}

func (n *filesNavigator) traverse(currentItem *TraverseItem, frame *navigationFrame) *LocalisableError {
	//
	// For files, the registered callback will only be invoked for file entries. This means
	// that the client will have no way to skip the descending of a particular directory. In
	// this case, the client should use the OnDescend callback (yet to be implemented) and
	// return SkipDir from there.

	defer func() {
		n.ascend(&NavigationInfo{Options: n.o, Item: currentItem, Frame: frame})
	}()
	navi := &NavigationInfo{Options: n.o, Item: currentItem, Frame: frame}
	n.descend(navi)
	entries, readErr := n.agent.read(currentItem, n.o.Behaviours.Sort.DirectoryEntryOrder)
	// Files and Folders need to be sorted independently to preserve the navigation order
	// stipulated by .Behaviours.Sort.DirectoryEntryOrder
	//
	entries.sort(&entries.Files)
	entries.sort(&entries.Folders)
	sorted := entries.all()

	if (currentItem.Entry != nil) && !(currentItem.Entry.IsDir()) {
		n.o.Hooks.Extend(navi, *sorted)

		// Effectively, this is the file only filter
		//
		return n.agent.proxy(currentItem, frame)
	}

	if exit, err := n.agent.notify(&agentNotifyParams{
		frame: frame, item: currentItem, entries: *sorted, readErr: readErr,
	}); exit || err != nil {
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
