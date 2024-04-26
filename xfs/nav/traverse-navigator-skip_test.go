package nav_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ginkgo ok
	. "github.com/onsi/gomega"    //nolint:revive // gomega ok
	"github.com/samber/lo"
	. "github.com/snivilised/extendio/i18n" //nolint:revive // i18n ok
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

			if entry.expectedNoOf.folders > 0 {
				Expect(folders).To(BeEquivalentTo(entry.expectedNoOf.folders))
			}

			if entry.expectedNoOf.files > 0 {
				Expect(files).To(BeEquivalentTo(entry.expectedNoOf.files))
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
			skipAt:   "College",
			prohibit: "Northern Council",
			all:      true,
			expectedNoOf: directoryQuantities{
				files:   4,
				folders: 4,
			},
		}),
		Entry(nil, &skipTE{
			naviTE: naviTE{
				message:      "universal: SkipDir (skipAt:folder)",
				subscription: nav.SubscribeAny,
			},
			skipAt:   "College",
			prohibit: "Northern Council",
			expectedNoOf: directoryQuantities{
				files:   4,
				folders: 4,
			},
		}),
		Entry(nil, &skipTE{
			naviTE: naviTE{
				message:      "universal: SkipAll (skipAt:file)",
				subscription: nav.SubscribeAny,
			},
			skipAt:   "A1 - The Telephone Call.flac",
			prohibit: "A2 - Night Drive.flac",
			all:      true,
			expectedNoOf: directoryQuantities{
				files:   1,
				folders: 3,
			},
		}),
		Entry(nil, &skipTE{
			naviTE: naviTE{
				message:      "universal: SkipDir (skipAt:file)",
				subscription: nav.SubscribeAny,
			},
			skipAt:   "A1 - The Telephone Call.flac",
			prohibit: "A2 - Night Drive.flac",
			expectedNoOf: directoryQuantities{
				files:   11,
				folders: 8,
			},
		}),
		Entry(nil, &skipTE{
			naviTE: naviTE{
				message:      "folders: SkipAll (skipAt:folder)",
				subscription: nav.SubscribeFolders,
			},
			skipAt:   "College",
			prohibit: "Northern Council",
			all:      true,
			expectedNoOf: directoryQuantities{
				folders: 4,
			},
		}),
		Entry(nil, &skipTE{
			naviTE: naviTE{
				message:      "folders: SkipDir (skipAt:folder)",
				subscription: nav.SubscribeFolders,
			},
			skipAt:   "Northern Council",
			prohibit: "Teenage Color",
			expectedNoOf: directoryQuantities{
				folders: 7,
			},
		}),
		Entry(nil, &skipTE{
			naviTE: naviTE{
				message:      "files: SkipAll (skipAt:file)",
				subscription: nav.SubscribeFiles,
			},
			skipAt:   "A1 - The Telephone Call.flac",
			prohibit: "A2 - Night Drive.flac",
			all:      true,
			expectedNoOf: directoryQuantities{
				files: 1,
			},
		}),
		Entry(nil, &skipTE{
			naviTE: naviTE{
				message:      "files: SkipDir (skipAt:file)",
				subscription: nav.SubscribeFiles,
			},
			skipAt:   "A1 - The Telephone Call.flac",
			prohibit: "A2 - Night Drive.flac",
			expectedNoOf: directoryQuantities{
				files: 11,
			},
		}),
	)
})
