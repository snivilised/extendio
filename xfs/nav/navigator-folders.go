package nav

import (
	"io/fs"

	"github.com/samber/lo"
	. "github.com/snivilised/extendio/translate"
)

type foldersNavigator struct {
	navigator
}

func (n *foldersNavigator) top(frame *navigationFrame) *LocalisableError {

	return n.agent.top(&agentTopParams{
		impl:  n,
		frame: frame,
	})
}

func (n *foldersNavigator) traverse(currentItem *TraverseItem, frame *navigationFrame) *LocalisableError {
	defer func() {
		n.ascend(&NavigationInfo{Options: n.o, Item: currentItem, Frame: frame})
	}()
	navi := &NavigationInfo{Options: n.o, Item: currentItem, Frame: frame}
	n.descend(navi)
	// for the folders navigator, we ignore the user defined setting in
	// n.o.Behaviours.Sort.DirectoryEntryOrder, as we're only interested in
	// folders and therefore force to use DirectoryEntryOrderFoldersFirstEn instead
	//
	entries, readErr := n.agent.read(currentItem, DirectoryEntryOrderFoldersFirstEn)
	folders := entries.Folders
	entries.sort(&folders)

	if n.o.Subscription == SubscribeFoldersWithFiles {
		files := lo.TernaryF(n.o.Filters.Children == nil,
			func() []fs.DirEntry { return entries.Files },
			func() []fs.DirEntry { return n.o.Filters.Children.Matching(entries.Files) },
		)

		entries.sort(&files)
		currentItem.Children = files
	}

	n.o.Hooks.Extend(navi, folders)

	if le := n.agent.proxy(currentItem, frame); le != nil || (currentItem.Entry != nil && !currentItem.Entry.IsDir()) {
		if le != nil && le.Inner == fs.SkipDir && currentItem.Entry.IsDir() {
			// Successfully skipped directory
			//
			le = nil
		}
		return le
	}

	if exit, err := n.agent.notify(&agentNotifyParams{
		frame: frame, item: currentItem, entries: folders, readErr: readErr,
	}); exit {
		return err
	} else {

		return n.agent.traverse(&agentTraverseParams{
			impl:     n,
			contents: &folders,
			parent:   currentItem,
			frame:    frame,
		})
	}
}
