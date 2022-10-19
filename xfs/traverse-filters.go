package xfs

import (
	"math"

	"github.com/samber/lo"
)

// FilterScopeEnum allows client to define which node types should be filtered.
// Filters can be applied to multiple node types by bitwise or-ing the XXXNodes
// definitions.
//

type FilterScopeEnum uint32

const (
	// For directories, any node that has no sub folders. For files, any node
	// that appears under a leaf directory node
	//
	LeafNodes FilterScopeEnum = 1 << iota

	// Any node that is a direct descendent of the root node
	//
	TopNodes

	// IntermediateNodes apply filter to nodes which are neither leaf or top nodes
	//
	IntermediateNodes

	// FolderNodes apply filter to folder nodes
	//
	FolderNodes

	// FileNodes apply filter to file nodes
	//
	FileNodes

	// CustomNodes apply filter to node using client defined categorisation
	// (yet to be confirmed)
	//
	CustomNodes

	// AllNodes apply the filter to any node type
	//
	AllNodes = math.MaxInt32
)

var filterScopeStrings map[FilterScopeEnum]string = map[FilterScopeEnum]string{
	LeafNodes:         "Leaf",
	TopNodes:          "Top",
	IntermediateNodes: "Intermediate",
	FolderNodes:       "Folder",
	FileNodes:         "File",
	CustomNodes:       "Custom",
	AllNodes:          "All",
}

// String
func (f FilterScopeEnum) String() string {
	result := filterScopeStrings[f]
	return lo.Ternary(result == "", "[multi]", result)
}

// TraverseFilter filter that can be applied to file system entries. When specified,
// the callback will only be invoked for file system nodes that pass the filter.
type TraverseFilter interface {
	IsMatch(name string, scope FilterScopeEnum) bool
}

// Filter base filter struct
type Filter struct {
	Scope FilterScopeEnum // defines which file system nodes the filter should be applied to
}

// RegexFilter regex filter
type RegexFilter struct {
	Filter
}

// IsMatch
func (f *RegexFilter) IsMatch(name string, scope FilterScopeEnum) bool {

	return false
}

// GlobFilter wildcard filter
type GlobFilter struct {
	Filter
}

// IsMatch
func (f *GlobFilter) IsMatch(name string, scope FilterScopeEnum) bool {

	return false
}

// CustomFilter is not a real filter, it represents a filter that would be defined by the client
type CustomFilter struct {
	Filter
}

func (f *CustomFilter) IsMatch(name string, scope FilterScopeEnum) bool {
	panic("IsMatch not defined for custom filter")
}
