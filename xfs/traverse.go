package xfs

import (
	"io/fs"
	"os"
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
// during traversal.
//

type TraverseItem struct {
	Path      string
	Entry     *fs.DirEntry  // contains a FileInfo via Info() function
	Extension *ExtendedItem // make this an interface?
	Error     error
}

// these will be similar to filepath.WalkFunc, defined as:
// func(path string, info fs.FileInfo, err error) error, except they
// will use TraverseItem instead of path string, info fs.FileInfo
// So far all these functions appear to be the same, so this may eventually
// reduced to just a single entity.
type FolderFn func(item *TraverseItem) error
type FileFn func(item *TraverseItem) error
type AnyFn func(item *TraverseItem) error

type Options struct {
	CaseSensitive bool // case sensitive traversal order
	Extend        bool // request an extended response
}

type FolderOptions struct {
	Options
	Fn FolderFn
}
type FolderOptionFn func(o *FolderOptions)

type FileOptions struct {
	Options
	Fn FileFn
}
type FileOptionFn func(o *FileOptions)

type AnyOptions struct {
	Options
	Fn AnyFn
}
type AnyOptionFn func(o *AnyOptions)

// FakeTraverse walks the file tree rooted at root, calling fn for each file or
// directory in the tree, including root.
//
// All errors that arise visiting files and directories are filtered by fn:
// see the fs.WalkDirFunc documentation for details.
//
// The files are walked in lexical order, which makes the output deterministic
// but requires FakeTraverse to read an entire directory into memory before proceeding
// to walk that directory.
//
// FakeTraverse does not follow symbolic links.
func FakeTraverse(root string, fn fs.WalkDirFunc) error {
	info, err := os.Lstat(root)
	if err != nil {
		err = fn(root, nil, err)
	} else {
		err = walkDir(root, &statDirEntry{info}, fn)
	}
	if err == fs.SkipDir {
		return nil
	}
	return err
}
