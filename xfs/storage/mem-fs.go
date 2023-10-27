package storage

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/avfs/avfs/vfs/memfs"
	"github.com/pkg/errors"
)

type memFS struct {
	backend VirtualBackend
	mfs     *memfs.MemFS
}

func UseMemFS() VirtualFS {
	return &memFS{
		backend: "mem",
		mfs:     memfs.New(),
	}
}

func (ms *memFS) Backend() VirtualBackend {
	return ms.backend
}

// interface ExistsInFS

func (ms *memFS) FileExists(path string) bool {
	result := false
	if info, err := ms.mfs.Lstat(path); err == nil {
		result = !info.IsDir()
	}

	return result
}

func (ms *memFS) DirectoryExists(path string) bool {
	result := false
	if info, err := ms.mfs.Lstat(path); err == nil {
		result = info.IsDir()
	}

	return result
}

// end: interface ExistsInFS

// interface ReadOnlyVirtualFS

func (ms *memFS) Lstat(path string) (fs.FileInfo, error) {
	return ms.mfs.Lstat(path)
}

func (ms *memFS) Stat(path string) (fs.FileInfo, error) {
	return ms.mfs.Stat(path)
}

func (ms *memFS) ReadFile(name string) ([]byte, error) {
	return ms.mfs.ReadFile(name)
}

func (ms *memFS) ReadDir(name string) ([]os.DirEntry, error) {
	return ms.mfs.ReadDir(name)
}

// end: interface ReadOnlyVirtualFS

// interface WriteToFS

func (ms *memFS) Chmod(name string, mode os.FileMode) error {
	return ms.mfs.Chmod(name, mode)
}

func (ms *memFS) Chown(name string, uid, gid int) error {
	return ms.mfs.Chown(name, uid, gid)
}

func (ms *memFS) Create(name string) (*os.File, error) {
	f, err := ms.mfs.Create(name)

	if file, ok := f.(*os.File); ok {
		return file, err
	}

	return nil, errors.Wrap(err,
		fmt.Sprintf("file '%v' creation in '%v' failed", name, ms.backend),
	)
}

func (ms *memFS) Link(oldname, newname string) error {
	return ms.mfs.Link(oldname, newname)
}

func (ms *memFS) Mkdir(name string, perm fs.FileMode) error {
	return ms.mfs.Mkdir(name, perm)
}

func (ms *memFS) MkdirAll(path string, perm os.FileMode) error {
	return ms.mfs.MkdirAll(path, perm)
}

func (ms *memFS) Remove(name string) error {
	return ms.mfs.Remove(name)
}

func (ms *memFS) RemoveAll(path string) error {
	return ms.mfs.RemoveAll(path)
}

func (ms *memFS) Rename(oldpath, newpath string) error {
	return ms.mfs.Rename(oldpath, newpath)
}

func (ms *memFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	return ms.mfs.WriteFile(name, data, perm)
}

// end: interface WriteToFS
