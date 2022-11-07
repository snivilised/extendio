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

func (f *RegexFilter) Validate() {
	f.Filter.Validate()
	if f.Pattern == "" {
		panic(PATTERN_NOT_DEFINED_L_ERR)
	}
	f.rex = regexp.MustCompile(f.Pattern)
}

// IsMatch
func (f *RegexFilter) IsMatch(item *TraverseItem) bool {
	if f.IsApplicable(item) {
		return f.invert(f.rex.Match([]byte(item.Extension.Name)))
	}
	return f.IfNotApplicable
}

// CompoundRegexFilter ========================================================

type CompoundRegexFilter struct {
	CompoundFilter
	rex *regexp.Regexp
}

func (f *CompoundRegexFilter) Validate() {
	if f.Pattern == "" {
		panic(PATTERN_NOT_DEFINED_L_ERR)
	}
	f.rex = regexp.MustCompile(f.Pattern)
}

func (f *CompoundRegexFilter) Matching(children []fs.DirEntry) []fs.DirEntry {
	return lo.Filter(children, func(entry fs.DirEntry, index int) bool {
		return f.invert(f.rex.Match([]byte(entry.Name())))
	})
}
