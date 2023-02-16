package nav

import (
	"io/fs"

	. "github.com/snivilised/extendio/translate"
)

type universalNavigator struct {
	navigator
}

func (n *universalNavigator) top(frame *navigationFrame, root string) *TraverseResult {

	return n.agent.top(&agentTopParams{
		impl:  n,
		frame: frame,
		top:   root,
	})
}

func (n *universalNavigator) traverse(params *traverseParams) *LocalisableError {
	defer func() {
		n.ascend(&NavigationInfo{
			Options: n.o,
			Item:    params.currentItem,
			Frame:   params.frame},
		)
	}()
	navi := &NavigationInfo{
		Options: n.o,
		Item:    params.currentItem,
		Frame:   params.frame,
	}
	n.descend(navi)

	var (
		entries *DirectoryEntries
		readErr error
	)

	if params.currentItem.Info.IsDir() {
		entries, readErr = n.agent.read(
			params.currentItem.Path,
			n.o.Store.Behaviours.Sort.DirectoryEntryOrder,
		)

		// Files and Folders need to be sorted independently to preserve the navigation order
		// stipulated by .Behaviours.Sort.DirectoryEntryOrder
		//
		entries.sort(&entries.Files)
		entries.sort(&entries.Folders)
	} else {
		entries = &DirectoryEntries{}
	}
	sorted := entries.all()

	n.o.Hooks.Extend(navi, entries)

	if le := n.agent.proxy(params.currentItem, params.frame); le != nil ||
		(params.currentItem.Entry != nil && !params.currentItem.Entry.IsDir()) {
		if le != nil && le.Inner == fs.SkipDir && params.currentItem.Entry.IsDir() {
			// Successfully skipped directory
			//
			le = nil
		}
		return le
	}

	if exit, err := n.agent.notify(&agentNotifyParams{
		frame:   params.frame,
		item:    params.currentItem,
		entries: *sorted,
		readErr: readErr,
	}); exit {
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
