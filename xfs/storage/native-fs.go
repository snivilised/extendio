package storage

import (
	"io/fs"
	"os"
)

type nativeFS struct {
	backend VirtualBackend
}

func UseNativeFS() VirtualFS {
	return &nativeFS{
		backend: "native",
	}
}

func (ns *nativeFS) Backend() VirtualBackend {
	return ns.backend
}

// interface ExistsInFS

func (ns *nativeFS) FileExists(path string) bool {
	result := false
	if info, err := os.Lstat(path); err == nil {
		result = !info.IsDir()
	}

	return result
}

func (ns *nativeFS) DirectoryExists(path string) bool {
	result := false
	if info, err := os.Lstat(path); err == nil {
		result = info.IsDir()
	}

	return result
}

// end: interface ExistsInFS

// interface ReadOnlyVirtualFS

func (ns *nativeFS) Lstat(path string) (fs.FileInfo, error) {
	return os.Lstat(path)
}

func (ns *nativeFS) Stat(path string) (fs.FileInfo, error) {
	return os.Stat(path)
}

func (ns *nativeFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (ns *nativeFS) ReadDir(name string) ([]os.DirEntry, error) {
	return os.ReadDir(name)
}

// end: interface ReadOnlyVirtualFS

func (ns *nativeFS) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(name, mode)
}

func (ns *nativeFS) Chown(name string, uid, gid int) error {
	return os.Chown(name, uid, gid)
}

func (ns *nativeFS) Create(name string) (*os.File, error) {
	return os.Create(name)
}

// interface WriteToFS

func (ns *nativeFS) Link(oldname, newname string) error {
	return os.Link(oldname, newname)
}

func (ns *nativeFS) Mkdir(name string, perm fs.FileMode) error {
	return os.Mkdir(name, perm)
}

func (ns *nativeFS) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (ns *nativeFS) Remove(name string) error {
	return os.Remove(name)
}

func (ns *nativeFS) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (ns *nativeFS) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func (ns *nativeFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

// end: interface WriteToFS
