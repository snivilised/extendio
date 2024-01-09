package nav

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"
	"github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/xfs/utils"
)

func fromExtendedGlobPattern(pattern string) (segments, suffixes []string, err error) {
	if !strings.Contains(pattern, "|") {
		return []string{}, []string{},
			errors.New("invalid extended glob filter definition; pattern is missing separator")
	}

	segments = strings.Split(pattern, "|")
	suffixes = strings.Split(segments[1], ",")

	suffixes = lo.Reject(suffixes, func(item string, index int) bool {
		return item == ""
	})

	return segments, suffixes, nil
}

func newNodeFilter(def *FilterDef) TraverseFilter {
	var (
		filter             TraverseFilter
		ifNotApplicable    = true
		err                error
		segments, suffixes []string
	)

	switch def.IfNotApplicable {
	case TriStateBoolTrueEn:
		ifNotApplicable = true

	case TriStateBoolFalseEn:
		ifNotApplicable = false

	case TriStateBoolUnsetEn:
	}

	switch def.Type {
	case FilterTypeExtendedGlobEn:
		if segments, suffixes, err = fromExtendedGlobPattern(def.Pattern); err != nil {
			panic(err)
		}

		filter = &ExtendedGlobFilter{
			Filter: Filter{
				name:            def.Description,
				scope:           def.Scope,
				pattern:         def.Pattern,
				negate:          def.Negate,
				ifNotApplicable: ifNotApplicable,
			},
			baseGlob: strings.ToLower(segments[0]),
			suffixes: lo.Map(suffixes, func(s string, _ int) string {
				return strings.ToLower(strings.TrimPrefix(strings.TrimSpace(s), "."))
			}),
			anyExtension: slices.Contains(suffixes, "*"),
		}

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

	case FilterTypeUndefinedEn:
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
	var (
		filter CompoundTraverseFilter
	)

	switch def.Type {
	case FilterTypeExtendedGlobEn:
		var (
			err                error
			segments, suffixes []string
		)

		if segments, suffixes, err = fromExtendedGlobPattern(def.Pattern); err != nil {
			panic(errors.New("invalid incase filter definition; pattern is missing separator"))
		}

		filter = &CompoundExtendedGlobFilter{
			CompoundFilter: CompoundFilter{
				Name:    def.Description,
				Pattern: def.Pattern,
				Negate:  def.Negate,
			},
			baseGlob: strings.ToLower(segments[0]),
			suffixes: lo.Map(suffixes, func(s string, _ int) string {
				return strings.ToLower(strings.TrimPrefix(strings.TrimSpace(s), "."))
			}),
			anyExtension: slices.Contains(suffixes, "*"),
		}

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
