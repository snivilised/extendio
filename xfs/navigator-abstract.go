package xfs

type navigator struct {
	o     *TraverseOptions
	agent *childAgent
}

func (n *navigator) options() *TraverseOptions {
	return n.o
}

func (n *navigator) descend(navi *NavigationParams) {
	navi.Frame.Depth++
	n.o.OnDescend(navi.Item)
}

func (n *navigator) ascend(navi *NavigationParams) {
	navi.Frame.Depth--
	n.o.OnAscend(navi.Item)
}
