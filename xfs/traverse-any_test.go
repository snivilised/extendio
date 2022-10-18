package xfs_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/snivilised/extendio/xfs"
)

func anyCallbackWithErrorCheck(o *xfs.AnyOptions) {
	o.Callback = func(item *xfs.TraverseItem) *xfs.LocalisableError {
		GinkgoWriter.Printf("---> 🥭 ON-CALLBACK: '%v' ...\n", item.Path)

		if item.Error != nil {
			GinkgoWriter.Printf("---> 🔥 ON-CALLBACK (error): '%s' ...\n", item.Error.Inner)
		}

		return nil
	}
}

var _ = Describe("TraverseAny", Ordered, func() {
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
			It("🧪 should: should visit all files in directory", func() {
				const relative = "Mötley Crüe/Theatre of Pain"
				path := path(heavy, relative)
				Expect(xfs.FolderExists(path)).To(BeTrue())

				xfs.TraverseAny(path, func(o *xfs.AnyOptions) {
					o.Callback = func(item *xfs.TraverseItem) *xfs.LocalisableError {
						GinkgoWriter.Printf("---> 🥭 ON-CALLBACK: '%v' ...\n", item.Path)

						return nil
					}
				})
			})
		})

		Context("and: is File", func() {
			It("🧪 should: double invoke callback (2nd with error)", func() {
				const relative = "Mötley Crüe/Theatre of Pain/01 - City Boy Blues.flac"
				path := path(heavy, relative)
				Expect(xfs.FileExists(path)).To(BeTrue())

				xfs.TraverseAny(path, func(o *xfs.AnyOptions) {
					// TODO: enforce the double callback, first without error,
					// second with readdirent/"not a directory" error
					//
					o.Callback = func(item *xfs.TraverseItem) *xfs.LocalisableError {
						GinkgoWriter.Printf("---> 🥭 ON-CALLBACK: '%v' ...\n", item.Path)

						if item.Error != nil {
							GinkgoWriter.Printf("---> 🔥 ON-CALLBACK (error): '%s' ...\n", item.Error.Inner)
						}

						return nil
					}
				})
			})
		})
	})

	Context("Path contains folders", func() {
		It("🧪 should: visit all files and directories", func() {
			const relative = "Mötley Crüe"
			path := path(heavy, relative)
			Expect(xfs.FolderExists(path)).To(BeTrue())

			xfs.TraverseAny(path, anyCallbackWithErrorCheck)
		})
	})

	Context("Path contains folders (large)", func() {
		It("🧪 should: visit all files and directories", func() {
			xfs.TraverseAny(root, anyCallbackWithErrorCheck)
		})
	})
})
