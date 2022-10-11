package xfs

import (
	"io/fs"
	"os"
)

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
