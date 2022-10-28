package xfs

import (
	"io/fs"

	"github.com/samber/lo"
	. "github.com/snivilised/extendio/translate"
)

// ExtendedItem provides extended information if the client requests
// it by setting the DoExtend boolean in the traverse options.
type ExtendedItem struct {
	Depth     uint            // traversal depth relative to the root
	IsLeaf    bool            // defines whether this node a leaf node
	Name      string          // derived as the leaf segment from filepath.Split
	Parent    string          // derived as the directory from filepath.Split
	SubPath   string          // represents the path between the root and the current item
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

// Clone makes shallow copy of TraverseItem (except the error).
func (ti *TraverseItem) Clone() *TraverseItem {

	return &TraverseItem{
		Path: ti.Path, Entry: ti.Entry, Info: ti.Info, Extension: ti.Extension,
	}
}

func (ti *TraverseItem) IsDir() bool {
	return lo.TernaryF(ti.Entry != nil,
		func() bool { return ti.Entry.IsDir() },
		func() bool { return ti.Info.IsDir() },
	)
}

// TraverseSubscription type to define traversal subscription (for which file system
// items the client defined callback are invoked for).
type TraverseSubscription uint

const (
	_                TraverseSubscription = iota
	SubscribeAny                          // invoke callback for files and folders
	SubscribeFolders                      // invoke callback for folders only
	SubscribeFiles                        // invoke callback for files only
)

// TraverseCallback defines traversal callback function signature.
type TraverseCallback func(item *TraverseItem) *LocalisableError

// AscendancyHandler defines the signatures of ascend/descend handlers
type AscendancyHandler func(item *TraverseItem)

// BeginHandler life cycle event handler, invoked before start of traversal
type BeginHandler func(root string)

// EndHandler life cycle event handler, invoked at end of traversal
type EndHandler func(result *TraverseResult)

// TraverseResult the result of the traversal process.
type TraverseResult struct {
	Error error
}

// TraverseNavigator interface to the main traverse instance.
type TraverseNavigator interface {
	Walk(root string) *TraverseResult
}

type navigationFrame struct {
	Root  string
	Depth uint
}

type navigatorImpl interface {
	options() *TraverseOptions
	top(frame *navigationFrame) *LocalisableError
	traverse(currentItem *TraverseItem, frame *navigationFrame) *LocalisableError
}

type NavigationParams struct {
	Options *TraverseOptions
	Item    *TraverseItem
	Frame   *navigationFrame
}

type SubPathInfo struct {
	Root string
	Item *TraverseItem
}
