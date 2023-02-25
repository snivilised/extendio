package nav_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	"github.com/snivilised/extendio/internal/helpers"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("TraverseNavigatorSkip", Ordered, func() {
	var root string

	BeforeAll(func() {
		root = musico()
	})

	When("folder is skipped", func() {
		Context("folder navigator", func() {
			It("🧪 should: not invoke skipped folder descendants", func() {
				path := helpers.Path(root, "RETRO-WAVE")
				session := &nav.PrimarySession{
					Path: path,
				}
				_ = session.Configure(func(o *nav.TraverseOptions) {
					o.Store.Subscription = nav.SubscribeFolders
					o.Store.DoExtend = true
					o.Callback = skipFolderCallback("College", "Northern Council")
					o.Notify.OnBegin = begin("🛡️")
				}).Run()
			})
		})

		Context("universal navigator", func() {
			It("🧪 should: not invoke skipped folder descendants", func() {
				path := helpers.Path(root, "RETRO-WAVE")
				session := &nav.PrimarySession{
					Path: path,
				}
				session.Configure(func(o *nav.TraverseOptions) {
					o.Store.Subscription = nav.SubscribeAny
					o.Store.DoExtend = true
					o.Callback = skipFolderCallback("College", "Northern Council")
					o.Notify.OnBegin = begin("🛡️")
				}).Run()
			})
		})
	})

	DescribeTable("skip",
		func(entry *skipTE) {
			path := helpers.Path(root, "RETRO-WAVE")
			session := &nav.PrimarySession{
				Path: path,
			}
			_ = session.Configure(func(o *nav.TraverseOptions) {
				o.Store.Subscription = entry.subscription
				o.Callback = skipFolderCallback("College", "Northern Council")
				o.Notify.OnBegin = begin("🛡️")
			}).Run()
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
