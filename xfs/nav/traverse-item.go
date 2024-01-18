package nav

import (
	"io/fs"

	"github.com/snivilised/extendio/xfs/utils"
)

// ExtendedItem provides extended information if the client requests
// it by setting the DoExtend boolean in the traverse options.
type ExtendedItem struct {
	Depth     int               // traversal depth relative to the root
	IsLeaf    bool              // defines whether this node a leaf node
	Name      string            // derived as the leaf segment from filepath.Split
	Parent    string            // derived as the directory from filepath.Split
	SubPath   string            // represents the path between the root and the current item
	NodeScope FilterScopeBiEnum // type of folder corresponding to the Filter Scope
	Custom    any               // to be set and used by the client
}

// TraverseItem info provided for each file system entity encountered
// during traversal. The root item does not have a DirEntry because it is
// not created as a result of a readDir invoke. Therefore, the client has
// to know that when its function is called back, they will be no DirEntry
// for the root entity.
type TraverseItem struct {
	Path        string
	Entry       fs.DirEntry  // contains a FileInfo via Info() function
	Info        fs.FileInfo  // optional file info instance
	Extension   ExtendedItem // extended information about the file system node, if requested
	Error       error
	Children    []fs.DirEntry
	filteredOut bool
	Parent      *TraverseItem
	admit       bool
	dir         bool
}

func isDir(item *TraverseItem) bool {
	if !utils.IsNil(item.Entry) {
		return item.Entry.IsDir()
	} else if !utils.IsNil(item.Info) {
		return item.Info.IsDir()
	}
	// only get here in error scenario, because neither Entry or Info is set
	//
	return false
}

func newTraverseItem(
	path string, entry fs.DirEntry, info fs.FileInfo, parent *TraverseItem, err error,
) *TraverseItem {
	item := &TraverseItem{
		Path:     path,
		Entry:    entry,
		Info:     info,
		Parent:   parent,
		Children: []fs.DirEntry{},
		Error:    err,
	}
	item.dir = isDir(item)

	return item
}

// clone makes shallow copy of TraverseItem (except the error).
func (ti *TraverseItem) clone() *TraverseItem {
	return &TraverseItem{
		Path:      ti.Path,
		Entry:     ti.Entry,
		Info:      ti.Info,
		Extension: ti.Extension,
		Children:  ti.Children,
	}
}

func (ti *TraverseItem) IsDirectory() bool {
	return ti.dir
}

func (ti *TraverseItem) filtered() {
	// 📚 filtered is used by sampling functions to mark an item as already having
	// been filtered. Sampling functions require the ability to 'Preview' an item
	// so that it can be filtered, but doing so means there is potential for a
	// child item to be double filtered. By marking an item is being pre-filtered,
	// when the navigation process reaches the child entry in its own right (as
	// opposed to being previewed), the filter can be bypassed and the client
	// callback for this item can be invoked; ie if an item passes filtering at
	// the preview stage, it does not needed to be filtered again.
	//
	ti.admit = true
}

func (ti *TraverseItem) key() string {
	return ti.Extension.SubPath
}
