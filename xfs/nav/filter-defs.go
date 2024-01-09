package nav

import (
	"io/fs"
	"math"
	"strings"

	"github.com/snivilised/extendio/collections"
	"golang.org/x/exp/maps"
)

// FilterScopeBiEnum allows client to define which node types should be filtered.
// Filters can be applied to multiple node types by bitwise or-ing the XXXNodes
// definitions. A node may have multiple scope designations, eg a node may be top
// level and leaf if the top level directory does not itself contain further
// sub-directories thereby making it also a leaf.
// It should be noted a file is only a leaf node all of its siblings are all files
// only
type FilterScopeBiEnum uint32

const (
	ScopeUndefinedEn FilterScopeBiEnum = 0

	// ScopeRootEn, the Root scope
	//
	ScopeRootEn FilterScopeBiEnum = 1 << (iota - 1)

	// ScopeTopEn, any node that is a direct descendent of the root node
	//
	ScopeTopEn

	// ScopeLeafEn, for directories, any node that has no sub folders. For files, any node
	// that appears under a leaf directory node
	//
	ScopeLeafEn

	// ScopeIntermediateEn, apply filter to nodes which are neither leaf or top nodes
	//
	ScopeIntermediateEn

	// ScopeFileEn attributed to file nodes
	//
	ScopeFileEn

	// ScopeFolderEn attributed to directory nodes
	//
	ScopeFolderEn

	// ScopeCustomEn, client defined categorisation (yet to be confirmed)
	//
	ScopeCustomEn

	// ScopeAllEn represents any node type
	//
	ScopeAllEn = math.MaxUint32
)

type FilterTypeEnum uint

const (
	FilterTypeUndefinedEn FilterTypeEnum = iota

	// FilterTypeExtendedGlobEn is the preferred filter type as it the most
	// user friendly. The base part of the name is filtered by a glob
	// and the suffix is filtered by a list of defined extensions. The pattern
	// for the extended filter type is composed of 2 parts; the first is a
	// glob, which is applied to the base part of the name. The second part
	// is a csv of required extensions to filter for. The pattern is specified
	// in the form: "<base-glob>|ext1,ext2...". Each extension may include a
	// a leading dot. An example pattern definition would be:
	// "cover.*|.jpg,jpeg"
	FilterTypeExtendedGlobEn

	// FilterTypeRegexEn regex filter
	FilterTypeRegexEn

	// FilterTypeGlobEn glob filter
	FilterTypeGlobEn

	// FilterTypeCustomEn client definable filter
	FilterTypeCustomEn

	// FilterTypePolyEn poly filter
	FilterTypePolyEn
)

type allOrderedFilterScopeEnums collections.OrderedKeysMap[FilterScopeBiEnum, string]

var filterScopeStrings = allOrderedFilterScopeEnums{
	ScopeUndefinedEn:    "Undefined",
	ScopeRootEn:         "Root",
	ScopeTopEn:          "Top",
	ScopeLeafEn:         "Leaf",
	ScopeIntermediateEn: "Intermediate",
	ScopeFileEn:         "File",
	ScopeFolderEn:       "Folder",
	ScopeCustomEn:       "Custom",
	ScopeAllEn:          "All",
}

var filterScopeKeys = maps.Keys(filterScopeStrings)

// String converts enum value to a string
func (f FilterScopeBiEnum) String() string {
	result := make([]string, 0, len(filterScopeKeys))

	for _, en := range filterScopeKeys {
		if en == ScopeAllEn {
			continue
		}

		if (en & f) > 0 {
			result = append(result, filterScopeStrings[en])
		}
	}

	return strings.Join(result, "|")
}

// Set sets the bit position indicated by mask
func (f *FilterScopeBiEnum) Set(mask FilterScopeBiEnum) {
	*f |= mask
}

// Clear clears the bit position indicated by mask
func (f *FilterScopeBiEnum) Clear(mask FilterScopeBiEnum) {
	*f &^= mask
}

// TraverseFilter filter that can be applied to file system entries. When specified,
// the callback will only be invoked for file system nodes that pass the filter.
type TraverseFilter interface {
	// Description describes filter
	Description() string

	// Validate ensures the filter definition is valid, panics when invalid
	Validate()

	// Source, filter definition (comes from filter definition Pattern)
	Source() string

	// IsMatch does this item match the filter
	IsMatch(item *TraverseItem) bool

	// IsApplicable is this filter applicable to this item's scope
	IsApplicable(item *TraverseItem) bool

	// Scope, what items this filter applies to
	Scope() FilterScopeBiEnum
}

// FilterDef defines a filter to be used filtering or listening features.
type FilterDef struct {
	// Type specifies the type of filter (mandatory)
	Type FilterTypeEnum

	// Description describes filter (optional)
	Description string

	// Pattern filter definition (mandatory)
	Pattern string

	// Scope which file system entries this filter applies to (defaults
	// to ScopeAllEn)
	Scope FilterScopeBiEnum

	// Negate, reverses the applicability of the filter (Defaults to false)
	Negate bool

	// IfNotApplicable, when the filter does not apply to a directory entry,
	// this value determines whether the callback is invoked for this entry
	// or not (defaults to TriStateBoolTrueEn/true).
	IfNotApplicable TriStateBoolEnum

	// Custom client define-able filter. When restoring for resume feature,
	// its the client's responsibility to restore this themselves (see
	// PersistenceRestorer)
	Custom TraverseFilter `json:"-"`

	// Poly allows for the definition of a PolyFilter which contains separate
	// filters that target files and folders separately. If present, then
	// all other fields are redundant, since the filter definitions inside
	// Poly should be referred to instead.
	Poly *PolyFilterDef
}

type PolyFilterDef struct {
	File   FilterDef
	Folder FilterDef
}

// CompoundTraverseFilter filter that can be applied to a folder's collection of entries
// when subscription is
type CompoundTraverseFilter interface {
	// Description describes filter
	Description() string

	// Validate ensures the filter definition is valid, panics when invalid
	Validate()

	// Source, filter definition (comes from filter definition Pattern)
	Source() string

	// Matching returns the collection of files contained within this
	// item's folder that matches this filter.
	Matching(children []fs.DirEntry) []fs.DirEntry
}

type CompoundFilterDef struct {
	// Type specifies the type of filter (mandatory)
	Type FilterTypeEnum

	// Description describes filter (optional)
	Description string

	// Pattern filter definition (mandatory)
	Pattern string

	// Negate, reverses the applicability of the filter (Defaults to false)
	Negate bool

	// Custom client define-able filter. When restoring for resume feature,
	// its the client's responsibility to restore this themselves (see
	// PersistenceRestorer)
	Custom CompoundTraverseFilter `json:"-"`
}

type compoundCounters struct {
	filteredIn  uint
	filteredOut uint
}

var BenignNodeFilterDef = FilterDef{
	Type:        FilterTypeRegexEn,
	Description: "benign allow all",
	Pattern:     ".",
	Scope:       ScopeRootEn,
}
