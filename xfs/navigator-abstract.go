package xfs

import (
	"fmt"
)

type navigator struct {
	options *TraverseOptions
}

func (n *navigator) top(root string) *LocalisableError {
	fmt.Printf("---> 🚁 [navigator]::top\n")

	return nil
}

func (n *navigator) traverse(currentItem *TraverseItem) *LocalisableError {
	fmt.Printf("---> 🚁 [navigator]::traverse\n")
	return nil
}
