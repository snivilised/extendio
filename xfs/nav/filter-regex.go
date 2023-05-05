package nav

import (
	"io/fs"
	"regexp"

	"github.com/samber/lo"
)

// RegexFilter ================================================================

// RegexFilter regex filter.
type RegexFilter struct {
	Filter
	rex *regexp.Regexp
}

// Validate ensures the filter definition is valid, panics when invalid
func (f *RegexFilter) Validate() {
	f.Filter.Validate()
	f.rex = regexp.MustCompile(f.pattern)
}

// IsMatch
func (f *RegexFilter) IsMatch(item *TraverseItem) bool {
	if f.IsApplicable(item) {
		return f.invert(f.rex.MatchString(item.Extension.Name))
	}

	return f.ifNotApplicable
}

// CompoundRegexFilter ========================================================

type CompoundRegexFilter struct {
	CompoundFilter
	rex *regexp.Regexp
}

func (f *CompoundRegexFilter) Validate() {
	f.rex = regexp.MustCompile(f.Pattern)
}

func (f *CompoundRegexFilter) Matching(children []fs.DirEntry) []fs.DirEntry {
	return lo.Filter(children, func(entry fs.DirEntry, index int) bool {
		return f.invert(f.rex.MatchString(entry.Name()))
	})
}
