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
//

type TraverseItem struct {
	Path      string
	Entry     fs.DirEntry   // contains a FileInfo via Info() function
	Info      fs.FileInfo   // optional file info instance
	Extension *ExtendedItem // make this an interface?, no
	Error     *LocalisableError
}

func (ti *TraverseItem) Clone() *TraverseItem {

	return &TraverseItem{
		Path: ti.Path, Entry: ti.Entry, Info: ti.Info, Extension: ti.Extension, Error: ti.Error,
	}
}

type TraverseSubscription uint

const (
	_ TraverseSubscription = iota
	SubscribeAny
	SubscribeFolders
	SubscribeFiles
)

type TraverseCallback func(item *TraverseItem) *LocalisableError

type TraverseOptions struct {
	Subscription    TraverseSubscription
	IsCaseSensitive bool // case sensitive traversal order
	Extend          bool // request an extended result
	WithMetrics     bool
	Callback        TraverseCallback
	OnDescend       TraverseCallback
	OnAscend        TraverseCallback
}
type TraverseOptionFn func(o *TraverseOptions)

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

type TraverseResult struct {
	Error *LocalisableError
}

type TraverseNavigator interface {
	Walk(root string) *TraverseResult
}

type navigatorSubject interface {
	top(root string) *LocalisableError
	traverse(currentItem *TraverseItem) *LocalisableError
}
