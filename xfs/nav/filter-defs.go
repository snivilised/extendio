package nav

import (
	"io/fs"
	"math"

	"github.com/samber/lo"
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

	// ScopeCustomEn, client defined categorisation (yet to be confirmed)
	//
	ScopeCustomEn

	// ScopeAllEn, any node type
	//
	ScopeAllEn = math.MaxUint32
)

type FilterTypeEnum uint

const (
	FilterTypeUndefinedEn FilterTypeEnum = iota

	// FilterTypeRegexEn regex filter
	FilterTypeRegexEn

	// FilterTypeGlobEn glob filter
	FilterTypeGlobEn

	// FilterTypeCustomEn client definable filter
	FilterTypeCustomEn
)

var filterScopeStrings = map[FilterScopeBiEnum]string{
	ScopeUndefinedEn:    "Undefined",
	ScopeRootEn:         "Root",
	ScopeTopEn:          "Top",
	ScopeLeafEn:         "Leaf",
	ScopeIntermediateEn: "Intermediate",
	ScopeCustomEn:       "Custom",
	ScopeAllEn:          "All",
}

// String converts enum value to a string
func (f FilterScopeBiEnum) String() string {
	result := filterScopeStrings[f]
	return lo.Ternary(result == "", "[multi]", result)
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
