package xfs

import "fmt"

type navigatorController struct {
	subject navigatorSubject
}

func (n *navigatorController) Walk(root string) *TraverseResult {
	fmt.Printf("---> 🛡️ [navigatorController]::Walk, root: '%v'\n", root)

	return &TraverseResult{
		Error: n.subject.top(root),
	}
}

// func (n navigatorController) Dummy() bool {
// 	return false
// }
