package xfs_test

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/snivilised/extendio/xfs"
)

func universalCallback(item *xfs.TraverseItem) *xfs.LocalisableError {
	GinkgoWriter.Printf("---> 🌊 ON-NAVIGATOR-CALLBACK: '%v'\n", item.Path)
	Expect(item.Extension).To(BeNil(), fmt.Sprintf("❌ %v", item.Path))

	return nil
}

func universalCallbackEx(item *xfs.TraverseItem) *xfs.LocalisableError {
	GinkgoWriter.Printf("---> 🌊 ON-NAVIGATOR-CALLBACK-EX: '%v'\n", item.Path)
	Expect(item.Extension).NotTo(BeNil(), fmt.Sprintf("❌ %v", item.Path))

	return nil
}

func foldersCallback(item *xfs.TraverseItem) *xfs.LocalisableError {
	GinkgoWriter.Printf("---> ☀️ ON-NAVIGATOR-CALLBACK: '%v'\n", item.Path)
	Expect(item.Info.IsDir()).To(BeTrue())
	Expect(item.Extension).To(BeNil(), fmt.Sprintf("❌ %v", item.Path))

	return nil
}

func foldersCallbackEx(item *xfs.TraverseItem) *xfs.LocalisableError {
	GinkgoWriter.Printf("---> ☀️ ON-NAVIGATOR-CALLBACK-EX: '%v'\n", item.Path)
	Expect(item.Info.IsDir()).To(BeTrue())
	Expect(item.Extension).NotTo(BeNil(), fmt.Sprintf("❌ %v", item.Path))

	return nil
}

func filesCallback(item *xfs.TraverseItem) *xfs.LocalisableError {
	GinkgoWriter.Printf("---> 🌙 ON-NAVIGATOR-CALLBACK: '%v'\n", item.Path)
	Expect(item.Info.IsDir()).To(BeFalse())
	Expect(item.Extension).To(BeNil(), fmt.Sprintf("❌ %v", item.Path))

	return nil
}

func filesCallbackEx(item *xfs.TraverseItem) *xfs.LocalisableError {
	GinkgoWriter.Printf("---> 🌙 ON-NAVIGATOR-CALLBACK-EX: '%v'\n", item.Path)
	Expect(item.Info.IsDir()).To(BeFalse())
	Expect(item.Extension).NotTo(BeNil(), fmt.Sprintf("❌ %v", item.Path))
	return nil
}

var _ = Describe("TraverseNavigator", Ordered, func() {
	var root string
	const IsExtended = true
	const NotExtended = false

	BeforeAll(func() {
		if current, err := os.Getwd(); err == nil {
			parent, _ := filepath.Split(current)
			root = filepath.Join(parent, "Test", "data", "MUSICO")
		}
	})

	Context("Path exists", func() {
		DescribeTable("Navigator",
			func(message, relative string, extended bool,
				subscription xfs.TraverseSubscription, callback xfs.TraverseCallback) {

				path := path(root, relative)
				navigator := xfs.NewNavigator(func(o *xfs.TraverseOptions) {
					o.Callback = callback
					o.Subscription = subscription
					o.DoExtend = extended
				})

				_ = navigator.Walk(path)
			},
			func(message, relative string, extended bool,
				subscription xfs.TraverseSubscription, callback xfs.TraverseCallback) string {

				return fmt.Sprintf("🧪 ===> '%v'", message)
			},
			Entry(nil, "universal: Path is leaf",
				"RETRO-WAVE/Chromatics/Night Drive", IsExtended, xfs.SubscribeAny, universalCallbackEx,
			),
			Entry(nil, "universal: Path contains folders",
				"RETRO-WAVE", NotExtended, xfs.SubscribeAny, universalCallback,
			),
			Entry(nil, "universal: Path contains folders (large)",
				"", NotExtended, xfs.SubscribeAny, universalCallback,
			),

			Entry(nil, "folders: Path is leaf",
				"RETRO-WAVE/Chromatics/Night Drive",
				NotExtended, xfs.SubscribeFolders, foldersCallback,
			),
			Entry(nil, "folders: Path contains folders",
				"RETRO-WAVE", IsExtended, xfs.SubscribeFolders, foldersCallbackEx,
			),
			Entry(nil, "folders: Path contains folders (large)",
				"", NotExtended, xfs.SubscribeFolders, foldersCallback,
			),

			Entry(nil, "files: Path is leaf",
				"RETRO-WAVE/Chromatics/Night Drive", NotExtended, xfs.SubscribeFiles, filesCallback,
			),
			Entry(nil, "files: Path contains folders",
				"RETRO-WAVE", NotExtended, xfs.SubscribeFiles, filesCallback,
			),
			Entry(nil, "files: Path contains folders (large)",
				"", IsExtended, xfs.SubscribeFiles, filesCallbackEx,
			),
		)
	})
})
