package nav

import "regexp"

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
		return f.Invert(f.rex.Match([]byte(item.Extension.Name)))
	}
	return f.IfNotApplicable
}
