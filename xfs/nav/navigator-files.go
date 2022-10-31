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
		n.ascend(&NavigationParams{Options: n.o, Item: currentItem, Frame: frame})
	}()
	navi := &NavigationParams{Options: n.o, Item: currentItem, Frame: frame}
	n.descend(navi)
	entries, readErr := n.agent.read(currentItem)

	if (currentItem.Entry != nil) && !(currentItem.Entry.IsDir()) {
		n.o.Hooks.Extend(navi, entries)

		// Effectively, this is the file only filter
		//
		return n.o.Callback(currentItem)
	}

	if exit, err := n.agent.notify(&agentNotifyParams{
		item: currentItem, entries: entries, readErr: readErr,
	}); exit || err != nil {
		return err
	} else {
		var err error
		if err = n.o.Hooks.Sort(entries); err != nil {
			panic(FILES_NAV_SORT_L_ERR)
		}

		return n.agent.traverse(&agentTraverseParams{
			impl:    n,
			entries: entries,
			parent:  currentItem,
			frame:   frame,
		})
	}
}
