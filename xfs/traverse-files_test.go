package xfs_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/snivilised/extendio/xfs"
)

func fileCallbackWithErrorCheck(o *xfs.FileOptions) {
	o.Callback = func(item *xfs.TraverseItem) *xfs.LocalisableError {
		GinkgoWriter.Printf("---> 🍩 ON-CALLBACK: '%v' ...\n", item.Path)

		if item.Error != nil {
			GinkgoWriter.Printf("---> 🔥 ON-CALLBACK (error): '%s' ...\n", item.Error.Inner)
		}
		Expect(item.Info.IsDir()).To(BeFalse())

		return nil
	}
}

var _ = Describe("TraverseFiles", Ordered, func() {
	var root, heavy string

	BeforeAll(func() {
		if current, err := os.Getwd(); err == nil {
			parent, _ := filepath.Split(current)
			root = filepath.Join(parent, "Test", "data", "MUSICO")
			heavy = filepath.Join(root, "rock", "metal", "dark", "HEAVY-METAL")
		}
	})

	Context("Path is leaf", func() {
		Context("and: is Folder", func() {
			It("🧪 should: should visit all files only", func() {
				const relative = "Mötley Crüe/Theatre of Pain"
				path := path(heavy, relative)
				Expect(xfs.FolderExists(path)).To(BeTrue())

				xfs.TraverseFiles(path, func(o *xfs.FileOptions) {
					o.Callback = func(item *xfs.TraverseItem) *xfs.LocalisableError {
						GinkgoWriter.Printf("---> 🍰 ON-CALLBACK: '%v' ...\n", item.Path)
						Expect(item.Info.IsDir()).To(BeFalse())

						return nil
					}
				})
			})
		})

		Context("Path contains folders", func() {
			It("🧪 should: visit all files only", func() {
				const relative = "Mötley Crüe"
				path := path(heavy, relative)
				Expect(xfs.FolderExists(path)).To(BeTrue())

				xfs.TraverseFiles(path, fileCallbackWithErrorCheck)
			})
		})

		Context("Path contains folders (large)", func() {
			It("🧪 should: visit all files only", func() {
				xfs.TraverseFiles(root, fileCallbackWithErrorCheck)
			})
		})
	})
})
