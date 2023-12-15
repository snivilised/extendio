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

	case FilterTypePolyEn:
		filter = newPolyFilter(def.Poly)

	default:
		panic(fmt.Sprintf("Filter definition for '%v' is missing the Type field", def.Description))
	}

	if def.Type != FilterTypePolyEn {
		filter.Validate()
	}

	return filter
}

func newPolyFilter(polyDef *PolyFilterDef) TraverseFilter {
	// lets enforce the correct filter scopes
	//
	polyDef.File.Scope.Set(ScopeFileEn)     // file scope must be set for files
	polyDef.File.Scope.Clear(ScopeFolderEn) // folder scope must NOT be set for files

	polyDef.Folder.Scope.Set(ScopeFolderEn) // folder scope must be set for folders
	polyDef.Folder.Scope.Clear(ScopeFileEn) // file scope must NOT be set for folders

	filter := &PolyFilter{
		File:   newNodeFilter(&polyDef.File),
		Folder: newNodeFilter(&polyDef.Folder),
	}

	return filter
}

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
	case FilterTypePolyEn:
	}

	filter.Validate()

	return filter
}
