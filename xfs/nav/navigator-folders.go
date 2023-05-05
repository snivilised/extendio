package nav

import (
	"io/fs"

	"github.com/samber/lo"
)

type foldersNavigator struct {
	navigator
}

func (n *foldersNavigator) top(frame *navigationFrame, root string) (*TraverseResult, error) {
	return n.agent.top(&agentTopParams{
		impl:  n,
		frame: frame,
		top:   root,
	})
}

func (n *foldersNavigator) traverse(params *traverseParams) error {
	defer func() {
		n.ascend(&NavigationInfo{
			Options: n.o,
			Item:    params.item,
			Frame:   params.frame},
		)
	}()

	var compoundCounts *compoundCounters

	navi := &NavigationInfo{
		Options: n.o,
		Item:    params.item,
		Frame:   params.frame,
	}
	n.descend(navi)

	// for the folders navigator, we ignore the user defined setting in
	// n.o.Store.Behaviours.Sort.DirectoryEntryOrder, as we're only interested in
	// folders and therefore force to use DirectoryEntryOrderFoldersFirstEn instead
	//
	entries, readErr := n.agent.read(params.item.Path,
		DirectoryEntryOrderFoldersFirstEn,
	)
	folders := entries.Folders
	entries.sort(&folders)

	if n.o.Store.Subscription == SubscribeFoldersWithFiles {
		var files []fs.DirEntry

		allFilesCount := len(entries.Files)
		filteredIn := 0

		if params.frame.filters == nil {
			files = entries.Files
		} else {
			files = lo.TernaryF(params.frame.filters.Children == nil,
				func() []fs.DirEntry { return entries.Files },
				func() []fs.DirEntry { return params.frame.filters.Children.Matching(entries.Files) },
			)
			filteredIn = len(files)
		}

		compoundCounts = &compoundCounters{
			filteredIn:  uint(filteredIn),
			filteredOut: uint(allFilesCount - filteredIn),
		}

		entries.sort(&files)
		params.item.Children = files
	}

	n.o.Hooks.Extend(navi, entries)

	if le := n.agent.proxy(&agentProxyParams{
		item:           params.item,
		frame:          params.frame,
		compoundCounts: compoundCounts,
	}); le != nil ||
		(params.item.Entry != nil && !params.item.Entry.IsDir()) {
		if QuerySkipDirError(le) && params.item.Entry.IsDir() {
			// Successfully skipped directory
			//
			le = nil
		}

		return le
	}

	if exit, err := n.agent.notify(&agentNotifyParams{
		frame:   params.frame,
		item:    params.item,
		entries: folders,
		readErr: readErr,
	}); exit {
		return err
	}

	return n.agent.traverse(&agentTraverseParams{
		impl:     n,
		contents: &folders,
		parent:   params.item,
		frame:    params.frame,
	})
}
