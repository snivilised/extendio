package nav_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	. "github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/internal/helpers"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("TraverseNavigatorSkip", Ordered, func() {
	var root string

	BeforeAll(func() {
		root = musico()
	})

	BeforeEach(func() {
		if err := Use(func(o *UseOptions) {
			o.Tag = DefaultLanguage.Get()
		}); err != nil {
			Fail(err.Error())
		}
	})

	DescribeTable("skip",
		func(entry *skipTE) {
			path := helpers.Path(root, "RETRO-WAVE")
			optionFn := func(o *nav.TraverseOptions) {
				o.Store.Subscription = entry.subscription
				o.Callback = lo.Ternary(entry.all,
					skipAllFolderCallback(entry.skipAt, entry.prohibit),
					skipDirFolderCallback(entry.skipAt, entry.prohibit),
				)
				o.Notify.OnBegin = begin("🛡️")
			}

			result, _ := nav.New().Primary(&nav.Prime{
				Path:      path,
				OptionsFn: optionFn,
			}).Run()

			_ = result.Session.StartedAt()
			_ = result.Session.Elapsed()

			files := result.Metrics.Count(nav.MetricNoFilesInvokedEn)
			folders := result.Metrics.Count(nav.MetricNoFoldersInvokedEn)
			GinkgoWriter.Printf("===> files: '%v', folders: '%v'\n", files, folders)

			if entry.folderCount > 0 {
				Expect(folders).To(BeEquivalentTo(entry.folderCount))
			}

			if entry.fileCount > 0 {
				Expect(files).To(BeEquivalentTo(entry.fileCount))
			}
		},
		func(entry *skipTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.message)
		},
		Entry(nil, &skipTE{
			naviTE: naviTE{
				message:      "universal: SkipAll (skipAt:folder)",
				subscription: nav.SubscribeAny,
			},
			skipAt:      "College",
			prohibit:    "Northern Council",
			all:         true,
			fileCount:   4,
			folderCount: 4,
		}),
		Entry(nil, &skipTE{
			naviTE: naviTE{
				message:      "universal: SkipDir (skipAt:folder)",
				subscription: nav.SubscribeAny,
			},
			skipAt:      "College",
			prohibit:    "Northern Council",
			fileCount:   4,
			folderCount: 4,
		}),
		Entry(nil, &skipTE{
			naviTE: naviTE{
				message:      "universal: SkipAll (skipAt:file)",
				subscription: nav.SubscribeAny,
			},
			skipAt:      "A1 - The Telephone Call.flac",
			prohibit:    "A2 - Night Drive.flac",
			all:         true,
			fileCount:   1,
			folderCount: 3,
		}),
		Entry(nil, &skipTE{
			naviTE: naviTE{
				message:      "universal: SkipDir (skipAt:file)",
				subscription: nav.SubscribeAny,
			},
			skipAt:      "A1 - The Telephone Call.flac",
			prohibit:    "A2 - Night Drive.flac",
			fileCount:   11,
			folderCount: 8,
		}),
		Entry(nil, &skipTE{
			naviTE: naviTE{
				message:      "folders: SkipAll (skipAt:folder)",
				subscription: nav.SubscribeFolders,
			},
			skipAt:      "College",
			prohibit:    "Northern Council",
			all:         true,
			folderCount: 4,
		}),
		Entry(nil, &skipTE{
			naviTE: naviTE{
				message:      "folders: SkipDir (skipAt:folder)",
				subscription: nav.SubscribeFolders,
			},
			skipAt:      "Northern Council",
			prohibit:    "Teenage Color",
			folderCount: 7,
		}),
		Entry(nil, &skipTE{
			naviTE: naviTE{
				message:      "files: SkipAll (skipAt:file)",
				subscription: nav.SubscribeFiles,
			},
			skipAt:    "A1 - The Telephone Call.flac",
			prohibit:  "A2 - Night Drive.flac",
			all:       true,
			fileCount: 1,
		}),
		Entry(nil, &skipTE{
			naviTE: naviTE{
				message:      "files: SkipDir (skipAt:file)",
				subscription: nav.SubscribeFiles,
			},
			skipAt:    "A1 - The Telephone Call.flac",
			prohibit:  "A2 - Night Drive.flac",
			fileCount: 11,
		}),
	)
})
