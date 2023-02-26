package nav

type filesNavigator struct {
	navigator
}

func (n *filesNavigator) top(frame *navigationFrame, root string) *TraverseResult {

	return n.agent.top(&agentTopParams{
		impl:  n,
		frame: frame,
		top:   root,
	})
}

func (n *filesNavigator) traverse(params *traverseParams) error {
	//
	// For files, the registered callback will only be invoked for file entries. This means
	// that the client will have no way to skip the descending of a particular directory. In
	// this case, the client should use the OnDescend callback (yet to be implemented) and
	// return SkipDir from there.

	defer func() {
		n.ascend(&NavigationInfo{
			Options: n.o,
			Item:    params.item,
			Frame:   params.frame},
		)
	}()
	navi := &NavigationInfo{
		Options: n.o,
		Item:    params.item,
		Frame:   params.frame,
	}

	var (
		entries *DirectoryEntries
		readErr error
	)

	if params.item.Info.IsDir() {
		entries, readErr = n.agent.read(
			params.item.Path,
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

	if (params.item.Info != nil) && !(params.item.Info.IsDir()) {
		n.o.Hooks.Extend(navi, entries)

		// Effectively, this is the file only filter
		//
		return n.agent.proxy(params.item, params.frame)
	}

	if exit, err := n.agent.notify(&agentNotifyParams{
		frame:   params.frame,
		item:    params.item,
		entries: *sorted,
		readErr: readErr,
	}); exit || err != nil {
		return err
	} else {

		return n.agent.traverse(&agentTraverseParams{
			impl:     n,
			contents: sorted,
			parent:   params.item,
			frame:    params.frame,
		})
	}
}
