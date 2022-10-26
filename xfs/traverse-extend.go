package xfs

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/samber/lo"
)

// DefaultExtendHookFn is the default extend hook function. The client can choose to
// override this by setting the custom function on options.Hooks.Extend. If the client
// wishes to augment the default behaviour rather than replace it, they can call
// this function from inside the custom function.
func DefaultExtendHookFn(ei *NavigationParams, descendants []fs.DirEntry) {

	if ei.Item.Extension != nil {
		panic(LocalisableError{
			Inner: fmt.Errorf("extend: item for path '%v' already extended", ei.Item.Path),
		})
	}

	grouped := lo.GroupBy(descendants, func(item fs.DirEntry) bool {
		return item.IsDir()
	})

	isLeaf := len(grouped[true]) == 0

	scope := IntermediateScopeEn
	if ei.Frame.Depth == 1 {
		scope = TopScopeEn
	} else if isLeaf {
		scope = LeafScopeEn
	}

	parent, name := filepath.Split(ei.Item.Path)
	ei.Item.Extension = &ExtendedItem{
		Depth:     ei.Frame.Depth,
		IsLeaf:    isLeaf,
		Name:      name,
		Parent:    parent,
		NodeScope: scope,
	}
	// fmt.Printf("💥 extend> depth: '%v', name: '%v', scope: '%v'\n", ei.frame.Depth, name, scope)
}

func nullExtendHookFn(ei *NavigationParams, descendants []fs.DirEntry) {}
