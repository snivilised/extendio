package nav

import "github.com/samber/lo"

// Filter base filter struct.
type Filter struct {
	Name          string
	Pattern       string
	RequiredScope FilterScopeEnum // defines which file system nodes the filter should be applied to
	Negate        bool            // select to define a negative match
}

func (f *Filter) Description() string {
	return f.Name
}

func (f *Filter) Source() string {
	return f.Pattern
}

func (f *Filter) IsApplicable(item *TraverseItem) bool {
	return (f.RequiredScope | item.Extension.NodeScope) > 0
}

func (f *Filter) Invert(result bool) bool {
	return lo.Ternary(f.Negate, !result, result)
}

func (f *Filter) Scope() FilterScopeEnum {
	return f.RequiredScope
}

type FilterInfo struct {
	Filter      TraverseFilter
	ActualScope FilterScopeEnum
}
