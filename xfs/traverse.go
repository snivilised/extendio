package xfs

import (
	"io/fs"
)

// ExtendedItem provides extended information if the client requests
// it by setting the Extend boolean in the traverse options.
type ExtendedItem struct {
	Depth     int             // traversal depth relative to the root
	IsLeaf    bool            // defines whether this node a leaf node
	Parent    string          // derived as the directory from filepath.Split
	NodeScope FilterScopeEnum // type of folder corresponding to the Filter Scope
	Custom    any             // to be set and used by the client
}

// TraverseItem info provided for each file system entity encountered
// during traversal. The root item does not have a DirEntry because it is
// not created as a result of a readDir invoke. Therefore, the client has
// to know that when its function is called back, they will be no DirEntry
// for the root entity.
type TraverseItem struct {
	Path      string
	Entry     fs.DirEntry   // contains a FileInfo via Info() function
	Info      fs.FileInfo   // optional file info instance
	Extension *ExtendedItem // extended information about the file system node, if requested
	Error     *LocalisableError
}

// Clone makes shallow copy of TraverseItem (except the error)
func (ti *TraverseItem) Clone() *TraverseItem {

	return &TraverseItem{
		Path: ti.Path, Entry: ti.Entry, Info: ti.Info, Extension: ti.Extension,
	}
}

// TraverseSubscription type to define traversal subscription (for which file system
// items the client defined callback are invoked for)
type TraverseSubscription uint

const (
	_                TraverseSubscription = iota
	SubscribeAny                          // invoke callback for files and folders
	SubscribeFolders                      // invoke callback for folders only
	SubscribeFiles                        // invoke callback for files only
)

// TraverseCallback defines traversal callback function signature
type TraverseCallback func(item *TraverseItem) *LocalisableError

// TraverseOptions customise way a directory tree in traversed
type TraverseOptions struct {
	Subscription    TraverseSubscription
	IsCaseSensitive bool             // case sensitive traversal order
	Extend          bool             // request an extended result
	WithMetrics     bool             // request metrics in TraversalResult
	Callback        TraverseCallback // traversal callback (universal, folders, files)
	OnDescend       TraverseCallback // callback to invoke as a folder is descended (before children)
	OnAscend        TraverseCallback // callback to invoke as a folder is ascended (after children)
}
type TraverseOptionFn func(o *TraverseOptions) // functional traverse options

// TraverseResult the result of the traversal process
type TraverseResult struct {
	Error *LocalisableError
}

// TraverseNavigator interface to the main traverse instance
type TraverseNavigator interface {
	Walk(root string) *TraverseResult
}

type navigatorSubject interface {
	top(root string) *LocalisableError
	traverse(currentItem *TraverseItem) *LocalisableError
}

func composeTraverseOptions(fn ...TraverseOptionFn) TraverseOptions {
	options := TraverseOptions{
		Subscription:    SubscribeAny,
		IsCaseSensitive: false,
		Extend:          false,
	}

	for _, functionalOption := range fn {
		functionalOption(&options)
	}
	return options
}
