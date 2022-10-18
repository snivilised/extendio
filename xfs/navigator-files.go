package xfs

import "fmt"

type filesNavigator struct {
	navigator
}

func (n *filesNavigator) top(root string) *TraverseResult {
	fmt.Printf("---> 🛩️ [filesNavigator]::top\n")
	return nil
}

func (n *filesNavigator) traverse(currentItem *TraverseItem) *TraverseResult {
	fmt.Printf("---> 🛩️ [filesNavigator]::traverse\n")
	return nil
}
