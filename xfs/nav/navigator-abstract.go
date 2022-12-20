package nav

type navigator struct {
	o     *TraverseOptions
	agent *navigationAgent
}

func (n *navigator) options() *TraverseOptions {
	return n.o
}

func (n *navigator) descend(navi *NavigationInfo) {
	navi.Frame.depth++
	navi.Frame.notifiers.descend.invoke(navi.Item)
}

func (n *navigator) ascend(navi *NavigationInfo) {
	navi.Frame.depth--
	navi.Frame.notifiers.ascend.invoke(navi.Item)
}
