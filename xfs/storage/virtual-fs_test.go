package storage_test

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/xfs/storage"
)

type virtualTE struct {
	message string
	should  string
	fn      func(vfs storage.VirtualFS, isNative bool)
}

func (v *virtualTE) action(vfs storage.VirtualFS, isNative bool) {
	v.fn(vfs, isNative)
}

var (
	faydeaudeau = os.FileMode(0o777)
	beezledub   = os.FileMode(0o666)
)

func reason(backend storage.VirtualBackend, message string, actual, expected any) string {
	return fmt.Sprintf("🔥 [%v:%v] expected '%v' to be '%v'",
		backend, message, actual, expected,
	)
}

type setupFile struct {
	path string
	data []byte
}

func setupDirectory(fs storage.VirtualFS, directoryPath string) {
	if e := fs.MkdirAll(directoryPath, faydeaudeau); e != nil {
		Fail(e.Error())
	}
}

func setupFiles(fs storage.VirtualFS, directoryPath string, files ...*setupFile) {
	setupDirectory(fs, directoryPath)

	for _, f := range files {
		if e := fs.WriteFile(f.path, f.data, beezledub); e != nil {
			Fail(e.Error())
		}
	}
}

var _ = Describe("virtual-fs", Ordered, func() {
	var (
		mfs           storage.VirtualFS
		nfs           storage.VirtualFS
		root, requiem string
	)

	BeforeAll(func() {
		if current, err := os.Getwd(); err == nil {
			resolved := filepath.Join(
				current, "..", "..", "Test", "data", "storage", "Nephilim", "Mourning Sun",
			)

			var err error
			root, err = filepath.Abs(resolved)

			if err != nil {
				Fail("failed to resolve root path")
			}

			requiem = filepath.Join(root, "info.requiem.txt")
		}
	})

	BeforeEach(func() {
		mfs = storage.UseMemFS()
		nfs = storage.UseNativeFS()

		if err := Use(func(o *UseOptions) {
			o.Tag = DefaultLanguage.Get()
		}); err != nil {
			Fail(err.Error())
		}
	})

	DescribeTable("vfs",
		func(entry *virtualTE) {
			entry.action(mfs, false)
			entry.action(nfs, true)
		},

		func(entry *virtualTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'",
				entry.message, entry.should,
			)
		},

		// --- ExistsInFS

		Entry(nil, &virtualTE{
			message: "FileExists",
			should:  "return correct existence status",
			fn: func(vfs storage.VirtualFS, isNative bool) {
				if !isNative {
					setupFiles(vfs, root, &setupFile{
						path: requiem,
						data: []byte("foo-bar"),
					})
				}
				actual := vfs.FileExists(requiem)

				Expect(actual).To(BeTrue(),
					reason(vfs.Backend(), "file exists return error", actual, true),
				)
			},
		}),

		Entry(nil, &virtualTE{
			message: "DirectoryExists",
			should:  "return correct existence status",
			fn: func(vfs storage.VirtualFS, isNative bool) {
				if !isNative {
					setupDirectory(vfs, root)
				}

				actual := vfs.DirectoryExists(root)

				Expect(actual).To(BeTrue(),
					reason(vfs.Backend(), "directory exists return error", actual, true),
				)
			},
		}),

		// --- end: ExistsInFS

		// --- ReadOnlyVirtualFS

		Entry(nil, &virtualTE{
			message: "Lstat",
			should:  "return correct file info",
			fn: func(vfs storage.VirtualFS, isNative bool) {
				if !isNative {
					setupFiles(vfs, root, &setupFile{
						path: requiem,
						data: []byte("requiem-content"),
					})
				}
				info, err := vfs.Lstat(requiem)
				Expect(err).Error().To(BeNil())

				expected := "info.requiem.txt"
				actual := info.Name()
				Expect(actual).To(Equal(expected),
					reason(vfs.Backend(), "lstat return correct name", actual, expected),
				)
			},
		}),

		Entry(nil, &virtualTE{
			message: "Stat",
			should:  "return correct file info",
			fn: func(vfs storage.VirtualFS, isNative bool) {
				if !isNative {
					setupFiles(vfs, root, &setupFile{
						path: requiem,
						data: []byte("requiem-content"),
					})
				}
				info, err := vfs.Stat(requiem)
				Expect(err).Error().To(BeNil())

				expected := "info.requiem.txt"
				actual := info.Name()
				Expect(actual).To(Equal(expected),
					reason(vfs.Backend(), "lstat return correct name", actual, expected),
				)
			},
		}),

		Entry(nil, &virtualTE{
			message: "ReadFile",
			should:  "return correct file content",
			fn: func(vfs storage.VirtualFS, isNative bool) {
				expected := "requiem-content"

				if !isNative {
					setupFiles(vfs, root, &setupFile{
						path: requiem,
						data: []byte(expected),
					})
				}
				content, err := vfs.ReadFile(requiem)
				actual := string(content)

				Expect(actual).To(Equal(expected),
					reason(vfs.Backend(), "read file return content", actual, expected),
				)
				Expect(err).Error().To(BeNil())
			},
		}),

		Entry(nil, &virtualTE{
			message: "ReadDir",
			should:  "return correct read status",
			fn: func(vfs storage.VirtualFS, isNative bool) {
				expected := "requiem-content"

				if !isNative {
					setupFiles(vfs, root, &setupFile{
						path: requiem,
						data: []byte(expected),
					})
				}
				actual, err := vfs.ReadDir(root)

				Expect(actual).To(HaveLen(1),
					reason(vfs.Backend(), "read directory return content", actual, expected),
				)
				Expect(err).Error().To(BeNil())
			},
		}),

		// --- end: ReadOnlyVirtualFS

		// --- WriteToFS

		Entry(nil, &virtualTE{
			message: "Create",
			should:  "create file",
			fn: func(vfs storage.VirtualFS, isNative bool) {
				path := filepath.Join(root, "shroud.txt")

				if !isNative {
					setupDirectory(vfs, root)
				}

				file, err := vfs.Create(path)
				if err == nil {
					defer file.Close()
				}

				defer func() {
					_ = vfs.Remove(path)
				}()

				Expect(err).Error().To(BeNil(),
					reason(vfs.Backend(), "create file return error", err, nil),
				)
			},
		}),

		// Chmod
		// Chown
		// Link

		Entry(nil, &virtualTE{
			message: "Mkdir",
			should:  "create all directory segments in path",
			fn: func(vfs storage.VirtualFS, isNative bool) {
				if isNative {
					return // bypass due to potential of access denied in native-fs
				}

				setupDirectory(vfs, root)

				path := filepath.Join(root, "__A")
				actual := vfs.Mkdir(path, beezledub)

				Expect(actual).Error().To(BeNil(),
					reason(vfs.Backend(), "Mkdir return error", actual, nil),
				)
				Expect(vfs.DirectoryExists(path)).To(BeTrue())
			},
		}),

		Entry(nil, &virtualTE{
			message: "MkdirAll",
			should:  "create all directory segments in path",
			fn: func(vfs storage.VirtualFS, isNative bool) {
				if isNative {
					return // bypass due to potential of access denied in native-fs
				}

				setupDirectory(vfs, root)

				path := filepath.Join(root, "__A", "__B", "__C")
				actual := vfs.MkdirAll(path, beezledub)

				Expect(actual).Error().To(BeNil(),
					reason(vfs.Backend(), "MkdirAll return error", actual, nil),
				)
				Expect(vfs.DirectoryExists(path)).To(BeTrue())
			},
		}),

		Entry(nil, &virtualTE{
			message: "Remove",
			should:  "remove file at path",
			fn: func(vfs storage.VirtualFS, isNative bool) {
				path := filepath.Join(root, "shroud.txt")
				setupFiles(vfs, root, &setupFile{
					path: path,
					data: []byte("foo-bar"),
				})

				actual := vfs.Remove(path)

				Expect(actual).Error().To(BeNil(),
					reason(vfs.Backend(), "remove file return error", actual, nil),
				)
			},
		}),

		Entry(nil, &virtualTE{
			message: "RemoveAll",
			should:  "remove all at path",
			fn: func(vfs storage.VirtualFS, isNative bool) {
				path := filepath.Join(root, "__A")

				setupFiles(vfs, path,
					&setupFile{
						path: filepath.Join(path, "x.txt"),
						data: []byte("x-content"),
					},
					&setupFile{
						path: filepath.Join(path, "y.txt"),
						data: []byte("y-content"),
					},
				)

				actual := vfs.RemoveAll(path)

				Expect(actual).Error().To(BeNil(),
					reason(vfs.Backend(), "remove all at path return error", actual, nil),
				)
			},
		}),

		Entry(nil, &virtualTE{
			message: "Rename",
			should:  "rename file at path",
			fn: func(vfs storage.VirtualFS, isNative bool) {
				path := filepath.Join(root, "shroud.txt")
				destination := filepath.Join(root, "renamed-shroud.txt")
				setupFiles(vfs, root, &setupFile{
					path: path,
					data: []byte("foo-bar"),
				})

				actual := vfs.Rename(path, destination)

				if isNative {
					defer func() {
						_ = vfs.Remove(destination)
					}()
				}

				Expect(actual).Error().To(BeNil(),
					reason(vfs.Backend(), "rename return error", actual, nil),
				)
				Expect(vfs.FileExists(destination)).To(BeTrue())
			},
		}),

		Entry(nil, &virtualTE{
			message: "WriteFile",
			should:  "write file to path",
			fn: func(vfs storage.VirtualFS, isNative bool) {
				setupDirectory(vfs, root)
				path := filepath.Join(root, "shroud.txt")

				content := []byte("Mourning Sun")
				actual := vfs.WriteFile(path, content, beezledub)

				if isNative {
					defer func() {
						_ = vfs.Remove(path)
					}()
				}

				Expect(actual).Error().To(BeNil(),
					reason(vfs.Backend(), "write file return error", actual, nil),
				)
			},
		}),

		// --- end: WriteToFS
	)
})
