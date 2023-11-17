package nav

import (
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

// DefaultExtendHookFn is the default extend hook function. The client can choose to
// override this by setting the custom function on options.Hooks.Extend. If the client
// wishes to augment the default behaviour rather than replace it, they can call
// this function from inside the custom function.
func DefaultExtendHookFn(navi *NavigationInfo, entries *DirectoryContents) {
	if navi.Item.Extension != nil {
		return
	}

	var (
		scope  FilterScopeBiEnum
		isLeaf bool
	)

	if navi.Item.IsDir() {
		isLeaf = len(entries.Folders) == 0
		scope = navi.frame.periscope.scope(isLeaf)
		scope |= ScopeFolderEn
	} else {
		scope = ScopeLeafEn
		scope |= ScopeFileEn
	}

	parent, name := filepath.Split(navi.Item.Path)
	navi.Item.Extension = &ExtendedItem{
		Depth:     navi.frame.periscope.depth(),
		IsLeaf:    isLeaf,
		Name:      name,
		Parent:    parent,
		NodeScope: scope,
	}

	spInfo := &SubPathInfo{
		Root:      navi.frame.root.Get(),
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

func nullExtendHookFn(_ *NavigationInfo, _ *DirectoryContents) {}
