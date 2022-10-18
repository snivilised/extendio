package xfs_test

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/snivilised/extendio/xfs"
)

var _ = Describe("TraverseNavigator", Ordered, func() {
	var root, heavy string

	BeforeAll(func() {
		if current, err := os.Getwd(); err == nil {
			parent, _ := filepath.Split(current)
			root = filepath.Join(parent, "Test", "data", "MUSICO")
			heavy = filepath.Join(root, "rock", "metal", "dark", "HEAVY-METAL")
		}
	})

	It("should: do nothing", func() {
		Expect(true)
	})

	Context("Create navigators", func() {
		It("🧪 should: ", func() {
			subs := []xfs.TraverseSubscription{xfs.SubscribeAny, xfs.SubscribeFolders, xfs.SubscribeFiles}

			for _, subscriber := range subs {

				navigator := xfs.NewNavigator(func(o *xfs.TraverseOptions) {
					o.Callback = func(item *xfs.TraverseItem) *xfs.LocalisableError {
						fmt.Printf("---> 🍧 ON-NAVIGATOR-CALLBACK: '%v' ...\n", item.Path)
						return nil
					}
					o.Subscription = subscriber
				})
				_ = navigator.Walk(heavy)
			}
		})
	})

	Context("universal", func() {
		It("🧪 should: walk all directories and folders", func() {
			navigator := xfs.NewNavigator(func(o *xfs.TraverseOptions) {
				o.Callback = func(item *xfs.TraverseItem) *xfs.LocalisableError {
					GinkgoWriter.Printf("---> 🍧 ON-NAVIGATOR-CALLBACK: '%v' ...\n", item.Path)
					return nil
				}
				o.Subscription = xfs.SubscribeAny
			})
			_ = navigator.Walk(heavy)
		})
	})
})
