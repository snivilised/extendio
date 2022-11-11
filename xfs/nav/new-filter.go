package nav

import "reflect"

func isNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}

// NewCurrentFilter exported for testing purposes only (do not use)
func NewCurrentFilter(info *FilterDef) TraverseFilter {

	var filter TraverseFilter
	switch info.Type {
	case FilterTypeRegexEn:
		filter = &RegexFilter{
			Filter: Filter{
				Name:            info.Description,
				RequiredScope:   info.Scope,
				Pattern:         info.Source,
				Negate:          info.Negate,
				IfNotApplicable: info.IfNotApplicable,
			},
		}

	case FilterTypeGlobEn:
		filter = &GlobFilter{
			Filter: Filter{
				Name:            info.Description,
				RequiredScope:   info.Scope,
				Pattern:         info.Source,
				Negate:          info.Negate,
				IfNotApplicable: info.IfNotApplicable,
			},
		}

	case FilterTypeCustomEn:
		if isNil(info.Custom) {
			panic("missing custom filter")
		}

		filter = info.Custom
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
		if isNil(def.Custom) {
			panic("missing custom compound filter")
		}

		filter = def.Custom
	}

	return filter
}
