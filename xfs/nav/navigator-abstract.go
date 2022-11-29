package nav

type navigator struct {
	o     *TraverseOptions
	agent *agent
}

func (n *navigator) options() *TraverseOptions {
	return n.o
}

func (n *navigator) descend(navi *NavigationInfo) {
	navi.Frame.Depth++
	navi.Frame.notifiers.descend.invoke(navi.Item)
}

func (n *navigator) ascend(navi *NavigationInfo) {
	navi.Frame.Depth--
	navi.Frame.notifiers.ascend.invoke(navi.Item)
}
