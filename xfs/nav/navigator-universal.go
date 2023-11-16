package nav

import "io/fs"

type universalNavigator struct {
	navigator
}

func (n *universalNavigator) top(frame *navigationFrame, root string) (*TraverseResult, error) {
	return n.agent.top(&agentTopParams{
		impl:  n,
		frame: frame,
		top:   root,
	})
}

func (n *universalNavigator) traverse(params *traverseParams) (*TraverseItem, error) {
	defer func() {
		n.ascend(&NavigationInfo{
			Options: n.o,
			Item:    params.item,
			frame:   params.frame},
		)
	}()

	navi := &NavigationInfo{
		Options: n.o,
		Item:    params.item,
		frame:   params.frame,
	}
	n.descend(navi)

	var (
		entries  *DirectoryEntries
		contents []fs.DirEntry
		readErr  error
		isDir    = params.item.IsDir()
	)

	if isDir {
		entries, readErr = n.agent.read(params.item.Path)

		// Files and Folders need to be sorted independently to preserve the navigation order
		// stipulated by .Behaviours.Sort.DirectoryEntryOrder
		//
		entries.sort(entries.Files)
		entries.sort(entries.Folders)

		contents = entries.All()

		if n.o.isSamplingActive() {
			n.o.Sampler.Fn(entries)
			contents = entries.All()
		}
	} else {
		entries = newEmptyDirectoryEntries(n.o)
	}

	n.o.Hooks.Extend(navi, entries)

	if le := params.frame.proxy(params.item, nil); le != nil {
		return nil, le
	}

	if skip, err := n.agent.notify(&agentNotifyParams{
		frame:    params.frame,
		item:     params.item,
		contents: contents,
		readErr:  readErr,
	}); skip == SkipTraversalAllEn {
		return nil, err
	} else if skip == SkipTraversalDirEn {
		return params.item.Parent, err
	}

	return n.agent.traverse(&agentTraverseParams{
		impl:     n,
		contents: contents,
		parent:   params.item,
		frame:    params.frame,
	})
}
