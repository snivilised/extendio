package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/snivilised/extendio/xfs/storage"
)

// EnsurePathAt ensures that the specified path exists. Given a path and a default
// filename, the specified path is created in the following manner:
// - If the path denotes a file (path does not end is a directory separator), the
// the parent folder is created if it doesn't exist on the file-system provided.
// - If the path denotes a directory, then that directory is created
//
// The returned string represents the file, so if the path specified was a path, then
// the defaultFilename provided is joined the the path and returned, otherwise
// the original path is returned un-modified.
// Note: filepath.Join does not preserve trailing separator, therefore to make sure
// a path is interpreted as a directory and not a file, then the separator has
// to be appended manually onto the end of the path.
// If vfs is not provided, then the path is ensured directly on the native file
// system.
func EnsurePathAt(path, defaultFilename string, perm int, vfs ...storage.VirtualFS) (at string, err error) {
	var (
		directory, file string
	)

	if strings.HasSuffix(path, string(os.PathSeparator)) {
		directory = path
		file = defaultFilename
	} else {
		directory, file = filepath.Split(path)
	}

	if len(vfs) > 0 {
		if !vfs[0].DirectoryExists(directory) {
			const perm = 0o766
			err = vfs[0].MkdirAll(directory, os.FileMode(perm))
		}
	} else {
		err = os.MkdirAll(directory, os.FileMode(perm))
	}

	return filepath.Clean(filepath.Join(directory, file)), err
}
