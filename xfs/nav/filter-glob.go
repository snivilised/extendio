package nav

import "path/filepath"

// GlobFilter wildcard filter.
type GlobFilter struct {
	Filter
}

// IsMatch
func (f *GlobFilter) IsMatch(item *TraverseItem) bool {
	if f.IsApplicable(item) {
		matched, _ := filepath.Match(f.Pattern, item.Extension.Name)
		return f.Invert(matched)
	}
	return f.IfNotApplicable
}
