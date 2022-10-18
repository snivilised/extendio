package xfs

import "fmt"

type navigator struct {
	options *TraverseOptions
}

func (n *navigator) Traverse() *TraverseResult {
	return nil
}

func (n *navigator) top(root string) *TraverseResult {
	fmt.Printf("---> 🚁 [navigator]::top\n")
	return nil
}

func (n *navigator) traverse(currentItem *TraverseItem) *TraverseResult {
	fmt.Printf("---> 🚁 [navigator]::traverse\n")
	return nil
}
