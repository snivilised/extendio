package nav

type navigator struct {
	o     *TraverseOptions
	agent *navigationAgent
}

func (n *navigator) options() *TraverseOptions {
	return n.o
}

func (n *navigator) descend(navi *NavigationInfo) {
	navi.Frame.periscope.descend()
	navi.Frame.notifiers.descend.invoke(navi.Item)
}

func (n *navigator) ascend(navi *NavigationInfo) {
	navi.Frame.periscope.ascend()
	navi.Frame.notifiers.ascend.invoke(navi.Item)
}
