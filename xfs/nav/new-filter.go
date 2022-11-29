package nav

import "reflect"

func IsNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}

// NewNodeFilter exported for testing purposes only (do not use)
func NewNodeFilter(def *FilterDef) TraverseFilter {

	var filter TraverseFilter
	switch def.Type {
	case FilterTypeRegexEn:
		filter = &RegexFilter{
			Filter: Filter{
				Name:            def.Description,
				RequiredScope:   def.Scope,
				Pattern:         def.Source,
				Negate:          def.Negate,
				IfNotApplicable: def.IfNotApplicable,
			},
		}

	case FilterTypeGlobEn:
		filter = &GlobFilter{
			Filter: Filter{
				Name:            def.Description,
				RequiredScope:   def.Scope,
				Pattern:         def.Source,
				Negate:          def.Negate,
				IfNotApplicable: def.IfNotApplicable,
			},
		}

	case FilterTypeCustomEn:
		if IsNil(def.Custom) {
			panic("missing custom filter")
		}

		filter = def.Custom
	}

	return filter
}

// NewCompoundFilter exported for testing purposes only (do not use)
func NewCompoundFilter(def *CompoundFilterDef) CompoundTraverseFilter {
	var filter CompoundTraverseFilter

	switch def.Type {
	case FilterTypeRegexEn:
		filter = &CompoundRegexFilter{
			CompoundFilter: CompoundFilter{
				Name:    def.Description,
				Pattern: def.Source,
				Negate:  def.Negate,
			},
		}

	case FilterTypeGlobEn:
		filter = &CompoundGlobFilter{
			CompoundFilter: CompoundFilter{
				Name:    def.Description,
				Pattern: def.Source,
				Negate:  def.Negate,
			},
		}

	case FilterTypeCustomEn:
		if IsNil(def.Custom) {
			panic("missing custom compound filter")
		}

		filter = def.Custom
	}

	return filter
}
