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
	n.o.Notify.OnDescend(navi.Item)
}

func (n *navigator) ascend(navi *NavigationInfo) {
	navi.Frame.Depth--
	n.o.Notify.OnAscend(navi.Item)
}
