package nav

type universalNavigator struct {
	navigator
}

func (n *universalNavigator) init(ns *NavigationState) {
	n.navigator.init(ns)

	if n.samplingActive {
		n.samplingCtrl.initInspector(n)
	}
}

func (n *universalNavigator) top(frame *navigationFrame, root string) (*TraverseResult, error) {
	return n.agent.top(&agentTopParams{
		impl:  n,
		frame: frame,
		top:   root,
	})
}

func (n *universalNavigator) inspect(params *traverseParams) *inspection {
	stash := &inspection{
		current: params.current,
		isDir:   params.current.IsDirectory(),
	}

	if stash.isDir {
		stash.contents, stash.readErr = n.agent.read(params.current.Path)

		stash.contents.sort(stash.contents.Files)
		stash.contents.sort(stash.contents.Folders)

		stash.entries = stash.contents.All()
	} else {
		stash.clearContents(n.o)
	}

	n.o.Hooks.Extend(params.navi, stash.contents)

	return stash
}

func (n *universalNavigator) traverse(params *traverseParams) (*TraverseItem, error) {
	navi := &NavigationInfo{
		Options: n.o,
		Item:    params.current,
		frame:   params.frame,
	}

	params.navi = navi
	descended := n.descend(navi)

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
	entries := stash.entries

	if stash.isDir {
		if n.samplingActive {
			n.samplingCtrl.sample(stash.contents, navi, params)
			entries = stash.contents.All()
		}
	}

	if le := params.frame.proxy(params.current, nil); le != nil {
		return nil, le
	}

	if skip, err := n.agent.notify(&agentNotifyParams{
		frame:   params.frame,
		current: params.current,
		entries: entries,
		readErr: stash.readErr,
	}); skip == SkipAllTraversalEn {
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
