package nav_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ginkgo ok
	. "github.com/onsi/gomega"    //nolint:revive // gomega ok
	"github.com/snivilised/extendio/internal/helpers"
	"github.com/snivilised/extendio/xfs/nav"

	. "github.com/snivilised/extendio/i18n" //nolint:revive // i18n ok
)

var _ = Describe("TraverseNavigatorSort", Ordered, func() {
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

	DescribeTable("sort",
		func(entry *sortTE) {
			recording := make(recordingOrderMap)
			counter := 0

			recorder := &nav.LabelledTraverseCallback{
				Label: "test recorder callback",
				Fn: func(item *nav.TraverseItem) error {
					_, found := recording[item.Extension.Name]

					if !found {
						recording[item.Extension.Name] = counter
					}
					counter++
					return entry.callback.Fn(item)
				},
			}

			path := helpers.Path(root, entry.relative)
			optionFn := func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = begin("🛡️")
				o.Store.Subscription = entry.subscription
				o.Store.FilterDefs = &nav.FilterDefinitions{
					Node: nav.FilterDef{
						Type:            nav.FilterTypeGlobEn,
						Description:     entry.name,
						Pattern:         entry.pattern,
						Scope:           entry.scope,
						Negate:          entry.negate,
						IfNotApplicable: entry.ifNotApplicable,
					},
				}
				o.Store.Behaviours.Sort.DirectoryEntryOrder = entry.order
				o.Callback = recorder
			}

			result, _ := nav.New().Primary(&nav.Prime{
				Path:      path,
				OptionsFn: optionFn,
			}).Run()

			sequence := -1

			for _, n := range entry.expectedOrder {
				Expect(recording[n] > sequence).To(BeTrue(), helpers.Reason(n))
				sequence = recording[n]
			}

			_ = result.Session.StartedAt()
			_ = result.Session.Elapsed()
		},
		func(entry *sortTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.message)
		},

		// === universal =====================================================

		Entry(nil, &sortTE{
			filterTE: filterTE{
				naviTE: naviTE{
					message:      "universal: Path contains folders",
					relative:     "",
					subscription: nav.SubscribeAny,
					callback:     universalSortCallback("CONTAINS-FOLDERS"),
				},
				name:    "items containing 'POP'",
				pattern: "*POP*",
				scope:   nav.ScopeAllEn,
			},
			order:         nav.DirectoryContentsOrderFoldersFirstEn,
			expectedOrder: []string{"DREAM-POP", "ELECTRONIC-POP", "POP"},
		}),

		Entry(nil, &sortTE{
			filterTE: filterTE{
				naviTE: naviTE{
					message:      "universal: folders before files",
					relative:     "bass",
					subscription: nav.SubscribeAny,
					callback:     universalDepthCallback("FOLDERS-FIRST", 2),
				},
				name:    "any",
				pattern: "*",
				scope:   nav.ScopeAllEn,
			},
			order:         nav.DirectoryContentsOrderFoldersFirstEn,
			expectedOrder: []string{"DUB", "DUBSTEP", "segments.bass.infex.txt"},
		}),

		Entry(nil, &sortTE{
			filterTE: filterTE{
				naviTE: naviTE{
					message:      "universal: files before folders",
					relative:     "bass",
					subscription: nav.SubscribeAny,
					callback:     universalDepthCallback("FILES-FIRST", 2),
				},
				name:    "any",
				pattern: "*",
				scope:   nav.ScopeAllEn,
			},
			order:         nav.DirectoryContentsOrderFilesFirstEn,
			expectedOrder: []string{"segments.bass.infex.txt", "DUB", "DUBSTEP"},
		}),

		// === folders =======================================================

		Entry(nil, &sortTE{
			filterTE: filterTE{
				naviTE: naviTE{
					message:      "folders: Path contains folders",
					relative:     "rock/metal",
					subscription: nav.SubscribeFolders,
					callback:     foldersSortCallback("CONTAINS-FOLDERS"),
				},
				name:    "items containing 'METAL'",
				pattern: "*METAL*",
				scope:   nav.ScopeAllEn,
			},
			order: nav.DirectoryContentsOrderFoldersFirstEn,
			expectedOrder: []string{
				"HEAVY-METAL",
				"THRASH-METAL",
				"HARD-METAL",
			},
		}),

		// === files =========================================================

		Entry(nil, &sortTE{
			filterTE: filterTE{
				naviTE: naviTE{
					message:      "files: Path contains folders",
					relative:     "rock/metal/dark",
					subscription: nav.SubscribeFiles,
					callback:     filesSortCallback("CONTAINS-FOLDERS"),
				},
				name:    "first track items with '.flac' suffix",
				pattern: "01*.flac",
				scope:   nav.ScopeLeafEn,
			},
			order: nav.DirectoryContentsOrderFilesFirstEn,
			expectedOrder: []string{
				"01 - Neon Knights.flac",
				"01 - Turn Up The Night.flac",
				"01 - The Ides of March.flac",
				"01 - Where Eagles Dare.flac",
				"01 - Wake Up Dead.flac",
				"01 - Holy Wars...The Punishment Due.flac",
			},
		}),
	)
})
