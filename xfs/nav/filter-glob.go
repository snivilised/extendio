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

// IsMatch does this item match the filter
func (f *GlobFilter) IsMatch(item *TraverseItem) bool {
	if f.IsApplicable(item) {
		matched, _ := filepath.Match(f.pattern, item.Extension.Name)
		return f.invert(matched)
	}

	return f.ifNotApplicable
}

// CompoundGlobFilter =========================================================

type CompoundGlobFilter struct {
	CompoundFilter
}

// Matching returns the collection of files contained within this
// item's folder that matches this filter.
func (f *CompoundGlobFilter) Matching(children []fs.DirEntry) []fs.DirEntry {
	return lo.Filter(children, func(entry fs.DirEntry, _ int) bool {
		matched, _ := filepath.Match(f.Pattern, entry.Name())
		return f.invert(matched)
	})
}
