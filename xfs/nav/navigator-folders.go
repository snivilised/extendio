package nav

import (
	"io/fs"

	"github.com/samber/lo"
	. "github.com/snivilised/extendio/translate"
)

type foldersNavigator struct {
	navigator
}

func (n *foldersNavigator) top(frame *navigationFrame, root string) *TraverseResult {

	return n.agent.top(&agentTopParams{
		impl:  n,
		frame: frame,
		top:   root,
	})
}

func (n *foldersNavigator) traverse(params *traverseParams) *LocalisableError {
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
	// for the folders navigator, we ignore the user defined setting in
	// n.o.Store.Behaviours.Sort.DirectoryEntryOrder, as we're only interested in
	// folders and therefore force to use DirectoryEntryOrderFoldersFirstEn instead
	//
	entries, readErr := n.agent.read(params.currentItem.Path,
		DirectoryEntryOrderFoldersFirstEn,
	)
	folders := entries.Folders
	entries.sort(&folders)

	if n.o.Store.Subscription == SubscribeFoldersWithFiles {

		var files []fs.DirEntry
		if params.frame.filters == nil {
			files = entries.Files
		} else {
			files = lo.TernaryF(params.frame.filters.Compound == nil,
				func() []fs.DirEntry { return entries.Files },
				func() []fs.DirEntry { return params.frame.filters.Compound.Matching(entries.Files) },
			)
		}

		entries.sort(&files)
		params.currentItem.Children = files
	}

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
		entries: folders,
		readErr: readErr,
	}); exit {
		return err
	} else {

		return n.agent.traverse(&agentTraverseParams{
			impl:     n,
			contents: &folders,
			parent:   params.currentItem,
			frame:    params.frame,
		})
	}
}
