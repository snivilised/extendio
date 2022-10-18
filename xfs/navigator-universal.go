package xfs

import "fmt"

type universalNavigator struct {
	navigator
}

func (n *universalNavigator) top(root string) *TraverseResult {
	fmt.Printf("---> 🚀 [universalNavigator]::top\n")
	return nil
}

func (n *universalNavigator) traverse(currentItem *TraverseItem) *TraverseResult {
	fmt.Printf("---> 🚀 [universalNavigator]::traverse\n")
	return nil
}
