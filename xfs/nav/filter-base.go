package nav

import "github.com/samber/lo"

// Filter =====================================================================

// Filter base filter struct.
type Filter struct {
	name            string
	pattern         string
	scope           FilterScopeBiEnum // defines which file system nodes the filter should be applied to
	negate          bool              // select to define a negative match
	ifNotApplicable bool
}

// Description description of the filter
func (f *Filter) Description() string {
	return f.name
}

// Source text defining the filter
func (f *Filter) Source() string {
	return f.pattern
}

func (f *Filter) IsApplicable(item *TraverseItem) bool {
	return (f.scope & item.Extension.NodeScope) > 0
}

func (f *Filter) Scope() FilterScopeBiEnum {
	return f.scope
}

func (f *Filter) invert(result bool) bool {
	return lo.Ternary(f.negate, !result, result)
}

func (f *Filter) Validate() {
	if f.scope == ScopeUndefinedEn {
		f.scope = ScopeAllEn
	}
}

// CompoundFilter =============================================================

// CompoundFilter filter used when subscription is FoldersWithFiles
type CompoundFilter struct {
	Name    string
	Pattern string
	Negate  bool
}

func (f *CompoundFilter) Description() string {
	return f.Name
}

func (f *CompoundFilter) Validate() {}

func (f *CompoundFilter) Source() string {
	return f.Pattern
}

func (f *CompoundFilter) invert(result bool) bool {
	return lo.Ternary(f.Negate, !result, result)
}
