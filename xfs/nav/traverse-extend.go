package nav

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
	. "github.com/snivilised/extendio/translate"
)

// DefaultExtendHookFn is the default extend hook function. The client can choose to
// override this by setting the custom function on options.Hooks.Extend. If the client
// wishes to augment the default behaviour rather than replace it, they can call
// this function from inside the custom function.
func DefaultExtendHookFn(params *NavigationParams, descendants []fs.DirEntry) {

	if params.Item.Extension != nil {
		panic(LocalisableError{
			Inner: fmt.Errorf("extend: item for path '%v' already extended", params.Item.Path),
		})
	}

	grouped := lo.GroupBy(descendants, func(item fs.DirEntry) bool {
		return item.IsDir()
	})
	isLeaf := len(grouped[true]) == 0

	// TODO(issue #38): this scope designation is not correct. We also need to define
	// a child scope which are the immediate descendants of the top level.
	//
	scope := IntermediateScopeEn
	if params.Frame.Depth == 1 {
		scope = TopScopeEn
	} else if isLeaf {
		scope = LeafScopeEn
	}

	parent, name := filepath.Split(params.Item.Path)
	params.Item.Extension = &ExtendedItem{
		Depth:     params.Frame.Depth,
		IsLeaf:    isLeaf,
		Name:      name,
		Parent:    parent,
		NodeScope: scope,
	}

	spInfo := &SubPathInfo{Root: params.Frame.Root, Item: params.Item}
	subpath := lo.TernaryF(params.Item.IsDir(),
		func() string { return params.Options.Hooks.FolderSubPath(spInfo) },
		func() string { return params.Options.Hooks.FileSubPath(spInfo) },
	)

	subpath = lo.TernaryF(params.Options.Behaviours.SubPath.KeepTrailingSep,
		func() string { return subpath },
		func() string {
			result := subpath
			sep := string(filepath.Separator)
			if strings.HasSuffix(subpath, sep) {
				result = subpath[:strings.LastIndex(subpath, sep)]
			}
			return result
		},
	)

	params.Item.Extension.SubPath = subpath
}

func nullExtendHookFn(params *NavigationParams, descendants []fs.DirEntry) {}
