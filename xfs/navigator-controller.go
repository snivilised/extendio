package xfs

import "fmt"

type navigatorController struct {
	subject navigatorSubject
}

func (n navigatorController) Walk(root string) *TraverseResult {
	fmt.Printf("---> 🛡️ [navigatorController]::Walk, root: '%v'\n", root)

	n.subject.top(root)
	return nil
}
