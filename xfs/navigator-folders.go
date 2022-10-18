package xfs

import "fmt"

type foldersNavigator struct {
	navigator
}

func (n *foldersNavigator) top(root string) *TraverseResult {
	fmt.Printf("---> ✈️ [foldersNavigator]::top\n")
	return nil
}

func (n *foldersNavigator) traverse(currentItem *TraverseItem) *TraverseResult {
	fmt.Printf("---> ✈️ [foldersNavigator]::traverse\n")
	return nil
}
