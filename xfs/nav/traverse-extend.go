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
func DefaultExtendHookFn(navi *NavigationInfo, descendants []fs.DirEntry) {

	if navi.Item.Extension != nil {
		panic(LocalisableError{
			Inner: fmt.Errorf("extend: item for path '%v' already extended", navi.Item.Path),
		})
	}
	isLeaf := false
	scope := ScopeIntermediateEn

	if navi.Item.IsDir() {
		// TODO: you shouldn't have to work this out again as it has already been performed:
		// the descendent passed in should be replaced with an Entries instance
		//
		grouped := lo.GroupBy(descendants, func(item fs.DirEntry) bool {
			return item.IsDir()
		})
		isLeaf = len(grouped[true]) == 0

		// TODO: eventually, the scope/depth designation will be put into an abstraction,
		// perhaps a scope->depth map. This will support resume, where these designations
		// will have to be adjusted
		//
		// Root=1
		// Top=2
		//

		switch {
		case isLeaf && navi.Frame.Depth == 1:
			scope = ScopeRootEn | ScopeLeafEn
		case navi.Frame.Depth == 1:
			scope = ScopeRootEn
		case isLeaf && navi.Frame.Depth == 2:
			scope = ScopeTopEn | ScopeLeafEn
		case navi.Frame.Depth == 2:
			scope = ScopeTopEn
		case isLeaf:
			scope = ScopeLeafEn
		}
	} else {
		scope = ScopeLeafEn
	}

	parent, name := filepath.Split(navi.Item.Path)
	navi.Item.Extension = &ExtendedItem{
		Depth:     navi.Frame.Depth,
		IsLeaf:    isLeaf,
		Name:      name,
		Parent:    parent,
		NodeScope: scope,
	}

	spInfo := &SubPathInfo{
		Root:      navi.Frame.Root,
		Item:      navi.Item,
		Behaviour: &navi.Options.Store.Behaviours.SubPath,
	}
	subpath := lo.TernaryF(navi.Item.IsDir(),
		func() string { return navi.Options.Hooks.FolderSubPath(spInfo) },
		func() string { return navi.Options.Hooks.FileSubPath(spInfo) },
	)

	subpath = lo.TernaryF(navi.Options.Store.Behaviours.SubPath.KeepTrailingSep,
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

	navi.Item.Extension.SubPath = subpath
}

func nullExtendHookFn(params *NavigationInfo, descendants []fs.DirEntry) {}
