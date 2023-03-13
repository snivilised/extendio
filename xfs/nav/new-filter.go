package nav

import (
	"fmt"

	. "github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/xfs/utils"
)

func newNodeFilter(def *FilterDef) TraverseFilter {
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
			panic(NewMissingCustomFilterDefinitionError("Options/Store/FilterDefs/Children/Custom"))
		}

		filter = def.Custom
	}
	filter.Validate()

	return filter
}
