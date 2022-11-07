package nav

import (
	"io/fs"
	"math"

	"github.com/samber/lo"
	. "github.com/snivilised/extendio/translate"
)

// InitFilter is the default filter initialiser. This can be overridden or extended
// by the client if the need arises. To extend this behaviour rather than replace it,
// call this function from inside the custom function set on o.Hooks.Filter. To
// replace the default functionality, do note the following points:
// - the original client callback is defined as frame.client, this should be referred to
// from outside the custom function (ie in the closure) as is performed here in the default.
// This will allow the custom function to invoke the core callback as appropriate.
// - The filters defined here in extendio make use of some extended fields, so if the client
// needs to define a custom function that is compatible with the native filters, then make
// sure the DoExtend value is set to true in the options, otherwise a panic will occur due to the
// filter attempting to de-reference the Extension on the TraverseItem.
func InitFilter(o *TraverseOptions, frame *navigationFrame) {
	if o.Filters.Current != nil {
		o.Filters.Current.Validate()
		o.DoExtend = true
		decorated := frame.client
		decorator := func(item *TraverseItem) *LocalisableError {
			if o.Filters.Current.IsMatch(item) {
				return decorated(item)
			}
			return nil
		}
		frame.decorate("init-current-filter 🎁", decorator)
	}

	if o.Filters.Children != nil {
		o.Filters.Children.Validate()
		o.DoExtend = true
	}
}

func bootstrapFilter(o *TraverseOptions, frame *navigationFrame) {
	o.Hooks.Filter(o, frame)
}

// FilterScopeEnum allows client to define which node types should be filtered.
// Filters can be applied to multiple node types by bitwise or-ing the XXXNodes
// definitions. A node may have multiple scope designations, eg a node may be top
// level and leaf if the top level directory does not itself contain further
// sub-directories thereby making it also a leaf.
// It should be noted a file is only a leaf node all of its siblings are all files
// only (TODO: write tests to ensure this characteristic).
type FilterScopeEnum uint32

const (
	ScopeUndefinedEn FilterScopeEnum = 0

	// ScopeRootEn, the Root scope
	//
	ScopeRootEn FilterScopeEnum = 1 << (iota - 1)

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

var filterScopeStrings map[FilterScopeEnum]string = map[FilterScopeEnum]string{
	ScopeUndefinedEn:    "Undefined",
	ScopeRootEn:         "Root",
	ScopeTopEn:          "Top",
	ScopeLeafEn:         "Leaf",
	ScopeIntermediateEn: "Intermediate",
	ScopeCustomEn:       "Custom",
	ScopeAllEn:          "All",
}

// String
func (f FilterScopeEnum) String() string {
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
	Scope() FilterScopeEnum
}

// CompoundTraverseFilter filter that can be applied to a folder's collection of entries
// when subscription is
type CompoundTraverseFilter interface {
	Description() string
	Validate()
	Source() string
	Matching(children []fs.DirEntry) []fs.DirEntry
}
