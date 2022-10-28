package nav

type navigatorController struct {
	impl navigatorImpl
}

func (n *navigatorController) Walk(root string) *TraverseResult {
	n.impl.options().Notify.OnBegin(root)
	frame := navigationFrame{
		Root: root,
	}

	result := &TraverseResult{
		Error: n.impl.top(&frame),
	}
	n.impl.options().Notify.OnEnd(result)

	return result
}
