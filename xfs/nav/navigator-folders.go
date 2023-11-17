package nav

import (
	"io/fs"

	"github.com/samber/lo"
)

type foldersNavigator struct {
	navigator
}

func (n *foldersNavigator) init(ns *NavigationState) {
	n.navigator.init(ns)

	if n.samplingActive {
		n.samplingCtrl.initInspector(n)
	}
}

func (n *foldersNavigator) top(frame *navigationFrame, root string) (*TraverseResult, error) {
	return n.agent.top(&agentTopParams{
		impl:  n,
		frame: frame,
		top:   root,
	})
}

func (n *foldersNavigator) inspect(params *traverseParams) *inspection {
	stash := &inspection{
		current: params.current,
		isDir:   params.current.IsDir(),
	}

	// for the folders navigator, we ignore the user defined setting in
	// n.o.Store.Behaviours.Sort.DirectoryEntryOrder, as we're only interested in
	// folders and therefore force to use DirectoryEntryOrderFoldersFirstEn instead
	//
	stash.contents, stash.readErr = n.agent.read(params.current.Path)
	stash.entries = stash.contents.Folders
	stash.contents.sort(stash.entries)

	if n.o.Store.Subscription == SubscribeFoldersWithFiles {
		var files []fs.DirEntry

		allFilesCount := len(stash.contents.Files)
		filteredIn := allFilesCount

		if params.frame.filters == nil {
			files = stash.contents.Files
		} else {
			files = lo.TernaryF(params.frame.filters.Children == nil,
				func() []fs.DirEntry { return stash.contents.Files },
				func() []fs.DirEntry { return params.frame.filters.Children.Matching(stash.contents.Files) },
			)
			filteredIn = len(files)
		}

		stash.compoundCounts = &compoundCounters{
			filteredIn:  uint(filteredIn),
			filteredOut: uint(allFilesCount - filteredIn),
		}

		stash.contents.sort(files)
		params.current.Children = files
	}

	n.o.Hooks.Extend(params.navi, stash.contents)

	return stash
}

func (n *foldersNavigator) traverse(params *traverseParams) (*TraverseItem, error) {
	navi := &NavigationInfo{
		Options: n.o,
		Item:    params.current,
		frame:   params.frame,
	}
	defer func() {
		if n.samplingFilterActive {
			delete(n.agent.cache, params.current.key())
		}

		n.ascend(navi)
	}()

	params.navi = navi
	n.descend(navi)

	stash := n.inspect(params)
	entries := stash.entries

	if n.samplingActive {
		n.samplingCtrl.sample(stash.contents, navi, params)
		entries = stash.contents.Folders
	}

	if le := params.frame.proxy(params.current, stash.compoundCounts); le != nil {
		return nil, le
	}

	if skip, err := n.agent.notify(&agentNotifyParams{
		frame:   params.frame,
		current: params.current,
		entries: entries,
		readErr: stash.readErr,
	}); skip == SkipAllTraversalEn {
		return nil, err
	}

	return n.agent.traverse(&agentTraverseParams{
		impl:    n,
		entries: entries,
		parent:  params.current,
		frame:   params.frame,
	})
}
