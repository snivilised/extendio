package nav

import "github.com/samber/lo"

// Filter =====================================================================

// Filter base filter struct.
type Filter struct {
	Name            string
	Pattern         string
	RequiredScope   FilterScopeBiEnum // defines which file system nodes the filter should be applied to
	Negate          bool              // select to define a negative match
	IfNotApplicable bool
}

func (f *Filter) Description() string {
	return f.Name
}

func (f *Filter) Validate() {
	if f.RequiredScope == ScopeUndefinedEn {
		f.RequiredScope = ScopeAllEn
	}
}

func (f *Filter) Source() string {
	return f.Pattern
}

func (f *Filter) IsApplicable(item *TraverseItem) bool {
	return (f.RequiredScope & item.Extension.NodeScope) > 0
}

func (f *Filter) Scope() FilterScopeBiEnum {
	return f.RequiredScope
}

func (f *Filter) invert(result bool) bool {
	return lo.Ternary(f.Negate, !result, result)
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
