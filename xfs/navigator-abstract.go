package xfs

import (
	"fmt"
)

type navigator struct {
	options  *TraverseOptions
	children *childAgent
}

func (n *navigator) top(frame *navigationFrame) *LocalisableError {
	fmt.Printf("---> 🚁 [navigator]::top\n")

	return nil
}

func (n *navigator) traverse(item *TraverseItem, frame *navigationFrame) *LocalisableError {
	fmt.Printf("---> 🚁 [navigator]::traverse\n")
	return nil
}

func (n *navigator) descend(navi *navigationInfo) *LocalisableError {
	navi.frame.Depth++
	return n.options.OnDescend(navi.item)
}

func (n *navigator) ascend(navi *navigationInfo) *LocalisableError {
	navi.frame.Depth--
	return n.options.OnAscend(navi.item)
}
