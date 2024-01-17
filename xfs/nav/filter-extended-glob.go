package nav

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

type ExtendedGlobFilter struct {
	Filter
	baseGlob     string
	suffixes     []string
	anyExtension bool
	exclusion    string
}

func filterFileByExtendedGlob(name, base, exclusion string, suffixes []string, anyExtension bool) bool {
	extension := filepath.Ext(name)
	baseName := strings.ToLower(strings.TrimSuffix(name, extension))

	if baseMatch, _ := filepath.Match(base, baseName); !baseMatch {
		return false
	}

	if excluded, _ := filepath.Match(exclusion, baseName); excluded {
		return false
	}

	return lo.TernaryF(anyExtension,
		func() bool {
			return true
		},
		func() bool {
			return lo.TernaryF(extension == "",
				func() bool {
					return len(suffixes) == 0
				},
				func() bool {
					return lo.Contains(
						suffixes, strings.ToLower(strings.TrimPrefix(extension, ".")),
					)
				},
			)
		},
	)
}

// IsMatch does this item match the filter
func (f *ExtendedGlobFilter) IsMatch(item *TraverseItem) bool {
	if f.IsApplicable(item) {
		result := lo.TernaryF(item.IsDir(),
			func() bool {
				result, _ := filepath.Match(f.baseGlob, strings.ToLower(item.Extension.Name))

				return result
			},
			func() bool {
				return filterFileByExtendedGlob(
					item.Extension.Name, f.baseGlob, f.exclusion, f.suffixes, f.anyExtension,
				)
			},
		)

		return f.invert(result)
	}

	return f.ifNotApplicable
}

// CompoundExtendedGlobFilter =======================================================

type CompoundExtendedGlobFilter struct {
	CompoundFilter
	baseGlob     string
	exclusion    string
	suffixes     []string
	anyExtension bool
}

func (f *CompoundExtendedGlobFilter) Matching(children []fs.DirEntry) []fs.DirEntry {
	return lo.Filter(children, func(entry fs.DirEntry, _ int) bool {
		name := entry.Name()

		return f.invert(filterFileByExtendedGlob(
			name, f.baseGlob, f.exclusion, f.suffixes, f.anyExtension,
		))
	})
}
