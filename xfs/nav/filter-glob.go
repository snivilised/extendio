package nav

import (
	"io/fs"
	"path/filepath"

	"github.com/samber/lo"
)

// GlobFilter wildcard filter.
type GlobFilter struct {
	Filter
}

// IsMatch
func (f *GlobFilter) IsMatch(item *TraverseItem) bool {
	if f.IsApplicable(item) {
		matched, _ := filepath.Match(f.Pattern, item.Extension.Name)
		return f.invert(matched)
	}
	return f.IfNotApplicable
}

// CompoundGlobFilter =========================================================

type CompoundGlobFilter struct {
	CompoundFilter
}

func (f *CompoundGlobFilter) Matching(children []fs.DirEntry) []fs.DirEntry {
	return lo.Filter(children, func(entry fs.DirEntry, index int) bool {
		matched, _ := filepath.Match(f.Pattern, entry.Name())
		return f.invert(matched)
	})
}
