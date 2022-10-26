package xfs

import "fmt"

type navigatorController struct {
	impl navigatorImpl
}

func (n *navigatorController) Walk(root string) *TraverseResult {
	fmt.Printf("---> 🛡️ [navigatorController]::Walk, root: '%v'\n", root)
	frame := navigationFrame{
		Root: root,
	}

	return &TraverseResult{
		Error: n.impl.top(&frame),
	}
}
