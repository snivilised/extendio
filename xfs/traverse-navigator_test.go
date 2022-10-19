package xfs_test

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/snivilised/extendio/xfs"
)

func genericCallback(item *xfs.TraverseItem) *xfs.LocalisableError {
	GinkgoWriter.Printf("---> 🌊 ON-NAVIGATOR-CALLBACK: '%v'\n", item.Path)

	return nil
}

func foldersCallback(item *xfs.TraverseItem) *xfs.LocalisableError {
	GinkgoWriter.Printf("---> 🌊 ON-NAVIGATOR-CALLBACK: '%v'\n", item.Path)
	Expect(item.Info.IsDir()).To(BeTrue())

	return nil
}

func filesCallback(item *xfs.TraverseItem) *xfs.LocalisableError {
	GinkgoWriter.Printf("---> 🌊 ON-NAVIGATOR-CALLBACK: '%v'\n", item.Path)
	Expect(item.Info.IsDir()).To(BeFalse())

	return nil
}

var _ = Describe("TraverseNavigator", Ordered, func() {
	var root string

	BeforeAll(func() {
		if current, err := os.Getwd(); err == nil {
			parent, _ := filepath.Split(current)
			root = filepath.Join(parent, "Test", "data", "MUSICO")
		}
	})

	Context("Path exists", func() {
		DescribeTable("Navigator",
			func(message, relative string, subscription xfs.TraverseSubscription, callback xfs.TraverseCallback) {
				path := path(root, relative)
				navigator := xfs.NewNavigator(func(o *xfs.TraverseOptions) {
					o.Callback = callback
					o.Subscription = subscription
				})
				_ = navigator.Walk(path)

			},
			func(message, relative string, subscription xfs.TraverseSubscription, callback xfs.TraverseCallback) string {
				return fmt.Sprintf("🧪 ===> '%v'", message)
			},
			Entry(nil, "universal: Path is leaf",
				"RETRO-WAVE/Chromatics/Night Drive", xfs.SubscribeAny, genericCallback,
			),
			Entry(nil, "universal: Path contains folders",
				"RETRO-WAVE", xfs.SubscribeAny, genericCallback,
			),
			Entry(nil, "universal: Path contains folders (large)",
				"", xfs.SubscribeAny, genericCallback,
			),

			Entry(nil, "folders: Path is leaf",
				"RETRO-WAVE/Chromatics/Night Drive",
				xfs.SubscribeFolders, foldersCallback,
			),
			Entry(nil, "folders: Path contains folders",
				"RETRO-WAVE", xfs.SubscribeFolders, foldersCallback,
			),
			Entry(nil, "folders: Path contains folders (large)",
				"", xfs.SubscribeFolders, foldersCallback,
			),

			Entry(nil, "files: Path is leaf",
				"RETRO-WAVE/Chromatics/Night Drive", xfs.SubscribeFiles, filesCallback,
			),
			Entry(nil, "files: Path contains folders",
				"RETRO-WAVE", xfs.SubscribeFiles, filesCallback,
			),
			Entry(nil, "files: Path contains folders (large)",
				"", xfs.SubscribeFiles, filesCallback,
			),
		)
	})
})
