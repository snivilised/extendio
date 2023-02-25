package nav_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/snivilised/extendio/internal/helpers"
	. "github.com/snivilised/extendio/translate"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("TraverseNavigatorSort", Ordered, func() {
	var root string

	BeforeAll(func() {
		root = musico()
	})

	DescribeTable("sort",
		func(entry *sortTE) {
			recording := recordingOrderMap{}
			counter := 0

			recorder := nav.LabelledTraverseCallback{
				Label: "test recorder callback",
				Fn: func(item *nav.TraverseItem) *LocalisableError {
					_, found := recording[item.Extension.Name]

					if !found {
						recording[item.Extension.Name] = counter
					}
					counter++
					return entry.callback.Fn(item)
				},
			}

			path := helpers.Path(root, entry.relative)
			session := &nav.PrimarySession{
				Path: path,
			}
			_ = session.Configure(func(o *nav.TraverseOptions) {
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
				o.Store.DoExtend = true
				o.Callback = recorder
			}).Run()

			sequence := -1
			for _, n := range entry.expectedOrder {
				Expect(recording[n] > sequence).To(BeTrue(), helpers.Reason(n))
				sequence = recording[n]
			}
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
			order:         nav.DirectoryEntryOrderFoldersFirstEn,
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
			order:         nav.DirectoryEntryOrderFoldersFirstEn,
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
			order:         nav.DirectoryEntryOrderFilesFirstEn,
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
			order: nav.DirectoryEntryOrderFoldersFirstEn,
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
			order: nav.DirectoryEntryOrderFilesFirstEn,
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
