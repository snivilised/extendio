package xfs

type navigatorController struct {
	impl navigatorImpl
}

func (n *navigatorController) Walk(root string) *TraverseResult {
	n.impl.options().OnBegin(root)
	frame := navigationFrame{
		Root: root,
	}

	result := &TraverseResult{
		Error: n.impl.top(&frame),
	}
	n.impl.options().OnEnd(result)

	return result
}
