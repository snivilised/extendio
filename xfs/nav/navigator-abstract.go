package nav

import (
	"github.com/snivilised/extendio/internal/log"
	"github.com/snivilised/extendio/xfs/utils"
)

type navigator struct {
	o     *TraverseOptions
	agent *navigationAgent
	log   utils.RoProp[log.Logger]
}

func (n *navigator) options() *TraverseOptions {
	return n.o
}

func (n *navigator) logger() log.Logger {
	return n.log.Get()
}

func (n *navigator) descend(navi *NavigationInfo) {
	navi.Frame.periscope.descend()
	navi.Frame.notifiers.descend.invoke(navi.Item)
}

func (n *navigator) ascend(navi *NavigationInfo) {
	navi.Frame.periscope.ascend()
	navi.Frame.notifiers.ascend.invoke(navi.Item)
}

func (n *navigator) finish() error {
	return n.log.Get().Sync()
}
