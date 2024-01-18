package nav

type filesNavigator struct {
	navigator
}

func (n *filesNavigator) init(ns *NavigationState) {
	n.navigator.init(ns)

	if n.samplingActive {
		n.samplingCtrl.initInspector(n)
	}
}

func (n *filesNavigator) top(frame *navigationFrame, root string) (*TraverseResult, error) {
	return n.agent.top(&agentTopParams{
		impl:  n,
		frame: frame,
		top:   root,
	})
}

func (n *filesNavigator) inspect(params *traverseParams) *inspection {
	stash := &inspection{
		current: params.current,
		isDir:   params.current.IsDirectory(),
	}

	if stash.isDir {
		stash.contents, stash.readErr = n.agent.read(params.current.Path)
		stash.contents.sort(stash.contents.Files)
		stash.contents.sort(stash.contents.Folders)
	} else {
		stash.clearContents(n.o)
	}

	n.o.Hooks.Extend(params.navi, stash.contents)

	return stash
}

func (n *filesNavigator) traverse(params *traverseParams) (*TraverseItem, error) {
	navi := &NavigationInfo{
		Options: n.o,
		Item:    params.current,
		frame:   params.frame,
	}
	params.navi = navi
	descended := n.descend(navi)

	//
	// For files, the registered callback will only be invoked for file entries. This means
	// that the client will have no way to skip the descending of a particular directory. In
	// this case, the client should use the OnDescend callback (yet to be implemented) and
	// return SkipDir from there.
	defer func(permit bool) {
		if n.samplingFilterActive {
			delete(n.agent.cache, params.current.key())
		}

		n.ascend(navi, permit)
	}(descended)

	if !descended {
		return nil, nil
	}

	stash := n.inspect(params)

	if !stash.isDir {
		// Effectively, this is the file only filter
		//
		return nil, params.frame.proxy(params.current, nil)
	}

	if n.samplingActive {
		n.samplingCtrl.sample(stash.contents, navi, params)
	}

	entries := stash.contents.All()

	if skip, err := n.agent.notify(&agentNotifyParams{
		frame:   params.frame,
		current: params.current,
		entries: entries,
		readErr: stash.readErr,
	}); skip == SkipAllTraversalEn || err != nil {
		return nil, err
	} else if skip == SkipDirTraversalEn {
		return params.current.Parent, err
	}

	return n.agent.traverse(&agentTraverseParams{
		impl:    n,
		entries: entries,
		parent:  params.current,
		frame:   params.frame,
	})
}
