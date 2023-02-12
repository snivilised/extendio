package nav_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("TraverseNavigatorSkip", Ordered, func() {
	var root string

	BeforeAll(func() {
		root = origin()
	})

	When("folder is skipped", func() {
		Context("folder navigator", func() {
			It("🧪 should: not invoke skipped folder descendants", func() {
				navigator := nav.NavigatorFactory{}.Construct(func(o *nav.TraverseOptions) {
					o.Store.Subscription = nav.SubscribeFolders
					o.Store.DoExtend = true
					o.Callback = skipFolderCallback("College", "Northern Council")
					o.Notify.OnBegin = begin("🛡️")
				})
				path := path(root, "RETRO-WAVE")
				navigator.Walk(path)
			})
		})

		Context("universal navigator", func() {
			It("🧪 should: not invoke skipped folder descendants", func() {
				navigator := nav.NavigatorFactory{}.Construct(func(o *nav.TraverseOptions) {
					o.Store.Subscription = nav.SubscribeAny
					o.Store.DoExtend = true
					o.Callback = skipFolderCallback("College", "Northern Council")
					o.Notify.OnBegin = begin("🛡️")
				})
				path := path(root, "RETRO-WAVE")
				navigator.Walk(path)
			})
		})
	})

	DescribeTable("skip",
		func(entry *skipTE) {
			navigator := nav.NavigatorFactory{}.Construct(func(o *nav.TraverseOptions) {
				o.Store.Subscription = entry.subscription
				o.Callback = skipFolderCallback("College", "Northern Council")
				o.Notify.OnBegin = begin("🛡️")
			})
			path := path(root, "RETRO-WAVE")
			navigator.Walk(path)
		},
		func(entry *skipTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.message)
		},
		Entry(nil, &skipTE{
			naviTE: naviTE{
				message:      "universal: skip",
				subscription: nav.SubscribeAny,
			},
			skip:    "College",
			exclude: "Northern Council",
		}),
		Entry(nil, &skipTE{
			naviTE: naviTE{
				message:      "folders: skip",
				subscription: nav.SubscribeFolders,
			},
			skip:    "College",
			exclude: "Northern Council",
		}),
	)
})
