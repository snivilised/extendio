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
	FilterTypeRegexEn
	FilterTypeGlobEn
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

// String
func (f FilterScopeBiEnum) String() string {
	result := filterScopeStrings[f]
	return lo.Ternary(result == "", "[multi]", result)
}

// TraverseFilter filter that can be applied to file system entries. When specified,
// the callback will only be invoked for file system nodes that pass the filter.
type TraverseFilter interface {
	Description() string
	Validate()
	Source() string
	IsMatch(item *TraverseItem) bool
	IsApplicable(item *TraverseItem) bool
	Scope() FilterScopeBiEnum
}

type FilterDef struct {
	Type            FilterTypeEnum
	Description     string
	Pattern         string
	Scope           FilterScopeBiEnum
	Negate          bool
	IfNotApplicable bool
	Custom          TraverseFilter `json:"-"`
}

// CompoundTraverseFilter filter that can be applied to a folder's collection of entries
// when subscription is
type CompoundTraverseFilter interface {
	Description() string
	Validate()
	Source() string
	Matching(children []fs.DirEntry) []fs.DirEntry
}

type CompoundFilterDef struct {
	Type        FilterTypeEnum
	Description string
	Pattern     string
	Negate      bool
	Custom      CompoundTraverseFilter `json:"-"`
}
