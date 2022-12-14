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
		top:   frame.root,
	})
}

func (n *filesNavigator) traverse(params *traverseParams) *LocalisableError {
	//
	// For files, the registered callback will only be invoked for file entries. This means
	// that the client will have no way to skip the descending of a particular directory. In
	// this case, the client should use the OnDescend callback (yet to be implemented) and
	// return SkipDir from there.

	defer func() {
		n.ascend(&NavigationInfo{Options: n.o, Item: params.currentItem, Frame: params.frame})
	}()
	navi := &NavigationInfo{Options: n.o, Item: params.currentItem, Frame: params.frame}
	n.descend(navi)
	entries, readErr := n.agent.read(params.currentItem.Path, n.o.Store.Behaviours.Sort.DirectoryEntryOrder)
	// Files and Folders need to be sorted independently to preserve the navigation order
	// stipulated by .Behaviours.Sort.DirectoryEntryOrder
	//
	entries.sort(&entries.Files)
	entries.sort(&entries.Folders)
	sorted := entries.all()

	if (params.currentItem.Entry != nil) && !(params.currentItem.Entry.IsDir()) {
		n.o.Hooks.Extend(navi, *sorted)

		// Effectively, this is the file only filter
		//
		return n.agent.proxy(params.currentItem, params.frame)
	}

	if exit, err := n.agent.notify(&agentNotifyParams{
		frame: params.frame, item: params.currentItem, entries: *sorted, readErr: readErr,
	}); exit || err != nil {
		return err
	} else {

		return n.agent.traverse(&agentTraverseParams{
			impl:     n,
			contents: sorted,
			parent:   params.currentItem,
			frame:    params.frame,
		})
	}
}

func (n *filesNavigator) spawn(params *spawnParams) *LocalisableError {

	return nil
}

func (n *filesNavigator) seed(params *seedParams) *LocalisableError {
	return nil
}
