package nav

import (
	"fmt"

	. "github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/xfs/utils"
)

func NewNodeFilter(def *FilterDef) TraverseFilter {
	var filter TraverseFilter
	switch def.Type {
	case FilterTypeRegexEn:
		filter = &RegexFilter{
			Filter: Filter{
				Name:            def.Description,
				RequiredScope:   def.Scope,
				Pattern:         def.Pattern,
				Negate:          def.Negate,
				IfNotApplicable: def.IfNotApplicable,
			},
		}

	case FilterTypeGlobEn:
		filter = &GlobFilter{
			Filter: Filter{
				Name:            def.Description,
				RequiredScope:   def.Scope,
				Pattern:         def.Pattern,
				Negate:          def.Negate,
				IfNotApplicable: def.IfNotApplicable,
			},
		}

	case FilterTypeCustomEn:
		if utils.IsNil(def.Custom) {
			panic(NewMissingCustomFilterDefinitionError("Options/Store/FilterDefs/Node/Custom"))
		}

		filter = def.Custom

	default:
		panic(fmt.Sprintf("Filter definition for '%v' is missing the Type field", def.Description))
	}

	filter.Validate()
	return filter
}

func NewRegexNodeFilter(def *FilterDef) TraverseFilter {
	def.Type = FilterTypeRegexEn
	scope := def.Scope
	if scope == ScopeUndefinedEn {
		scope = ScopeAllEn
	}

	filter := RegexFilter{
		Filter: Filter{
			Name:            def.Description,
			RequiredScope:   scope,
			Pattern:         def.Pattern,
			Negate:          def.Negate,
			IfNotApplicable: def.IfNotApplicable,
		},
	}
	filter.Validate()

	return &filter
}

func NewGlobNodeFilter(def *FilterDef) TraverseFilter {
	def.Type = FilterTypeGlobEn
	scope := def.Scope
	if scope == ScopeUndefinedEn {
		scope = ScopeAllEn
	}

	filter := GlobFilter{
		Filter: Filter{
			Name:            def.Description,
			RequiredScope:   scope,
			Pattern:         def.Pattern,
			Negate:          def.Negate,
			IfNotApplicable: def.IfNotApplicable,
		},
	}
	filter.Validate()

	return &filter
}

// NewCompoundFilter exported for testing purposes only (do not use)
func NewCompoundFilter(def *CompoundFilterDef) CompoundTraverseFilter {
	var filter CompoundTraverseFilter

	switch def.Type {
	case FilterTypeRegexEn:
		filter = &CompoundRegexFilter{
			CompoundFilter: CompoundFilter{
				Name:    def.Description,
				Pattern: def.Pattern,
				Negate:  def.Negate,
			},
		}

	case FilterTypeGlobEn:
		filter = &CompoundGlobFilter{
			CompoundFilter: CompoundFilter{
				Name:    def.Description,
				Pattern: def.Pattern,
				Negate:  def.Negate,
			},
		}

	case FilterTypeCustomEn:
		if utils.IsNil(def.Custom) {
			panic(NewMissingCustomFilterDefinitionError("Options/Store/FilterDefs/Children/Custom"))
		}

		filter = def.Custom
	}

	return filter
}
