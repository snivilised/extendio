package storage

import (
	"io/fs"
	"os"
)

type filepathAPI interface {
	// Intended only for those filepath methods that actually affect the
	// filesystem. Eg: there is no point in replicating methods like
	// filepath.Join here they are just path helpers that do not read/write
	// to the filesystem.
	// Currently, there is no requirement for using any filepath methods
	// with the golang generator, hence nothing is defined here. We may
	// want to replicate this filesystem model in other contexts, so this
	// will serve as a reminder in the intended use of this interface.
}

// ExistsInFS contains methods that check the existence of file system items.
type ExistsInFS interface {
	// FileExists does file exist at the path specified
	FileExists(path string) bool

	// DirectoryExists does directory exist at the path specified
	DirectoryExists(path string) bool
}

type ReadFromFS interface {
	// Lstat, see https://pkg.go.dev/os#Lstat
	Lstat(path string) (fs.FileInfo, error)

	// Lstat, see https://pkg.go.dev/os#Stat
	Stat(path string) (fs.FileInfo, error)

	// ReadFile, see https://pkg.go.dev/os#ReadFile
	ReadFile(name string) ([]byte, error)

	// ReadDir, see https://pkg.go.dev/os#ReadDir
	ReadDir(name string) ([]os.DirEntry, error)
}

// WriteToFS contains methods that perform mutative operations on the file system.
type WriteToFS interface {

	// Chmod, see https://pkg.go.dev/os#Chmod
	Chmod(name string, mode os.FileMode) error

	// Chown, https://pkg.go.dev/os#Chown
	Chown(name string, uid, gid int) error

	// Create, see https://pkg.go.dev/os#Create
	Create(name string) (*os.File, error)

	// Link, see https://pkg.go.dev/os#Link
	Link(oldname, newname string) error

	// Mkdir, see https://pkg.go.dev/os#Mkdir
	Mkdir(name string, perm fs.FileMode) error

	// MkdirAll, see https://pkg.go.dev/os#MkdirAll
	MkdirAll(path string, perm os.FileMode) error

	// Remove, see https://pkg.go.dev/os#Remove
	Remove(name string) error

	// RemoveAll, see https://pkg.go.dev/os#RemoveAll
	RemoveAll(path string) error

	// Rename, see https://pkg.go.dev/os#Rename
	Rename(oldpath, newpath string) error

	// WriteFile, see https://pkg.go.dev/os#WriteFile
	WriteFile(name string, data []byte, perm os.FileMode) error
}

// ReadOnlyVirtualFS provides read-only access to the file system.
type ReadOnlyVirtualFS interface {
	filepathAPI
	ExistsInFS
	ReadFromFS
}

// VirtualFS is a facade over the native file system, which include read
// and write access.
type VirtualFS interface {
	filepathAPI
	ExistsInFS
	ReadFromFS
	WriteToFS

	Backend() VirtualBackend
}

type VirtualBackend string
