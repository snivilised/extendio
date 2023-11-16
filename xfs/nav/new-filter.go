package nav

import (
	"fmt"

	"github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/xfs/utils"
)

func newNodeFilter(def *FilterDef) TraverseFilter {
	var (
		filter          TraverseFilter
		ifNotApplicable = true
	)

	switch def.IfNotApplicable { //nolint:exhaustive // already accounted for
	case TriStateBoolTrueEn:
		ifNotApplicable = true

	case TriStateBoolFalseEn:
		ifNotApplicable = false
	}

	switch def.Type { //nolint:exhaustive // default case is present
	case FilterTypeRegexEn:
		filter = &RegexFilter{
			Filter: Filter{
				name:            def.Description,
				scope:           def.Scope,
				pattern:         def.Pattern,
				negate:          def.Negate,
				ifNotApplicable: ifNotApplicable,
			},
		}

	case FilterTypeGlobEn:
		filter = &GlobFilter{
			Filter: Filter{
				name:            def.Description,
				scope:           def.Scope,
				pattern:         def.Pattern,
				negate:          def.Negate,
				ifNotApplicable: ifNotApplicable,
			},
		}

	case FilterTypeCustomEn:
		if utils.IsNil(def.Custom) {
			panic(i18n.NewMissingCustomFilterDefinitionError("Options/Store/FilterDefs/Node/Custom"))
		}

		filter = def.Custom

	default:
		panic(fmt.Sprintf("Filter definition for '%v' is missing the Type field", def.Description))
	}

	filter.Validate()

	return filter
}

// newCompoundFilter exported for testing purposes only (do not use)
func newCompoundFilter(def *CompoundFilterDef) CompoundTraverseFilter {
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
			panic(i18n.NewMissingCustomFilterDefinitionError("Options/Store/FilterDefs/Children/Custom"))
		}

		filter = def.Custom

	case FilterTypeUndefinedEn:
	}

	filter.Validate()

	return filter
}
