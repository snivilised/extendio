package nav_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/internal/helpers"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("TraverseNavigatorCascade", Ordered, func() {
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

	DescribeTable("cascade",
		func(entry *cascadeTE) {
			path := helpers.Path(root, entry.relative)
			optionFn := func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = begin("🛡️")
				o.Store.Subscription = entry.subscription
				o.Callback = entry.callback
				o.Store.Behaviours.Cascade.NoRecurse = entry.noRecurse
				o.Store.Behaviours.Cascade.Depth = entry.depth
			}

			result, err := nav.New().Primary(&nav.Prime{
				Path:      path,
				OptionsFn: optionFn,
			}).Run()
			_ = result

			Expect(err).Error().To(BeNil())

			Expect(result.Metrics.Count(nav.MetricNoFilesInvokedEn)).To(Equal(entry.expectedNoOf.files),
				"Incorrect no of files")
			Expect(result.Metrics.Count(nav.MetricNoFoldersInvokedEn)).To(Equal(entry.expectedNoOf.folders),
				"Incorrect no of folders")
		},
		func(entry *cascadeTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'", entry.message, entry.should)
		},

		// === universal =====================================================

		Entry(nil, &cascadeTE{
			naviTE: naviTE{
				message:      "universal: Path contains folders only, no-recurse",
				should:       "traverse single level",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
				callback:     universalScopeCallback("CONTAINS-FOLDERS"),
				expectedNoOf: directoryQuantities{
					files:   0,
					folders: 4,
				},
			},
			noRecurse: true,
		}),

		Entry(nil, &cascadeTE{
			naviTE: naviTE{
				message:      "universal: Path contains files only, no-recurse",
				should:       "traverse single level (containing files)",
				relative:     "RETRO-WAVE/Chromatics/Night Drive",
				subscription: nav.SubscribeAny,
				callback:     universalScopeCallback("CONTAINS-FILES"),
				expectedNoOf: directoryQuantities{
					files:   4,
					folders: 1,
				},
			},
			noRecurse: true,
		}),

		Entry(nil, &cascadeTE{
			naviTE: naviTE{
				message:      "universal: Path contains folders only, depth=1",
				should:       "traverse single level",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
				callback:     universalScopeCallback("CONTAINS-FOLDERS"),
				expectedNoOf: directoryQuantities{
					files:   0,
					folders: 4,
				},
			},
			depth: 1,
		}),

		Entry(nil, &cascadeTE{
			naviTE: naviTE{
				message:      "universal: Path contains folders only, depth=2",
				should:       "traverse 2 levels",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
				callback:     universalScopeCallback("CONTAINS-FOLDERS"),
				expectedNoOf: directoryQuantities{
					files:   0,
					folders: 8,
				},
			},
			depth: 2,
		}),

		Entry(nil, &cascadeTE{
			naviTE: naviTE{
				message:      "universal: Path contains folders only, depth=3",
				should:       "traverse 3 levels (containing files)",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
				callback:     universalScopeCallback("CONTAINS-FOLDERS"),
				expectedNoOf: directoryQuantities{
					files:   14,
					folders: 8,
				},
			},
			depth: 3,
		}),

		// === folders =======================================================

		Entry(nil, &cascadeTE{
			naviTE: naviTE{
				message:      "universal: Path contains folders only, no-recurse",
				should:       "traverse single level",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFolders,
				callback:     foldersCallback("CONTAINS-FILES"),
				expectedNoOf: directoryQuantities{
					files:   0,
					folders: 4,
				},
			},
			noRecurse: true,
		}),

		Entry(nil, &cascadeTE{
			naviTE: naviTE{
				message:      "universal: Path contains files only, no-recurse",
				should:       "traverse single level (containing files)",
				relative:     "RETRO-WAVE/Chromatics/Night Drive",
				subscription: nav.SubscribeFolders,
				callback:     universalScopeCallback("LEAF-PATH"),
				expectedNoOf: directoryQuantities{
					files:   0,
					folders: 1,
				},
			},
			noRecurse: true,
		}),

		Entry(nil, &cascadeTE{
			naviTE: naviTE{
				message:      "universal: Path contains folders only, depth=1",
				should:       "traverse single level",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFolders,
				callback:     universalScopeCallback("CONTAINS-FOLDERS"),
				expectedNoOf: directoryQuantities{
					files:   0,
					folders: 4,
				},
			},
			depth: 1,
		}),

		Entry(nil, &cascadeTE{
			naviTE: naviTE{
				message:      "universal: Path contains folders only, depth=2",
				should:       "traverse 2 levels",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFolders,
				callback:     universalScopeCallback("CONTAINS-FOLDERS"),
				expectedNoOf: directoryQuantities{
					files:   0,
					folders: 8,
				},
			},
			depth: 2,
		}),

		Entry(nil, &cascadeTE{
			naviTE: naviTE{
				message:      "universal: Path contains folders only, depth=3",
				should:       "traverse 3 levels (containing files)",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFolders,
				callback:     universalScopeCallback("CONTAINS-FOLDERS"),
				expectedNoOf: directoryQuantities{
					files:   0,
					folders: 8,
				},
			},
			depth: 3,
		}),

		// === files =========================================================

		Entry(nil, &cascadeTE{
			naviTE: naviTE{
				message:      "file: Path contains folders only, no-recurse",
				should:       "traverse single level",
				relative:     "RETRO-WAVE/Chromatics/Night Drive",
				subscription: nav.SubscribeFiles,
				callback:     filesCallback("FILE"),
				expectedNoOf: directoryQuantities{
					files:   4,
					folders: 0,
				},
			},
			noRecurse: true,
		}),

		Entry(nil, &cascadeTE{
			naviTE: naviTE{
				message:      "file: Path contains folders only, depth=1",
				should:       "traverse single level",
				relative:     "RETRO-WAVE/Chromatics/Night Drive",
				subscription: nav.SubscribeFiles,
				callback:     filesCallback("FILE"),
				expectedNoOf: directoryQuantities{
					files:   4,
					folders: 0,
				},
			},
			depth: 1,
		}),
	)
})
