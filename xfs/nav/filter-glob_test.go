package nav_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"

	"github.com/snivilised/extendio/translate"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("FilterGlob", Ordered, func() {
	var root string

	BeforeAll(func() {
		root = origin()
	})

	DescribeTable("GlobFilter",
		func(entry *filterTE) {
			recording := recordingMap{}
			filterDefs := &nav.FilterDefinitions{
				Node: nav.FilterDef{
					Type:            nav.FilterTypeGlobEn,
					Description:     entry.name,
					Pattern:         entry.pattern,
					Scope:           entry.scope,
					Negate:          entry.negate,
					IfNotApplicable: entry.ifNotApplicable,
				},
			}
			var filter nav.TraverseFilter

			path := path(root, entry.relative)
			session := nav.PrimarySession{
				Path: path,
			}
			_ = session.Configure(func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = func(state *nav.NavigationState) {
					GinkgoWriter.Printf(
						"---> 🛡️ [traverse-navigator-test:BEGIN], root: '%v'\n", state.Root,
					)
					filter = state.Filters.Node
				}

				o.Store.Subscription = entry.subscription
				o.Store.FilterDefs = filterDefs
				o.Store.DoExtend = true
				o.Callback = nav.LabelledTraverseCallback{
					Label: "test glob filter callback",
					Fn: func(item *nav.TraverseItem) *translate.LocalisableError {
						GinkgoWriter.Printf(
							"===> 💠 Glob Filter(%v) source: '%v', item-name: '%v', item-scope(fs): '%v(%v)'\n",
							filter.Description(),
							filter.Source(),
							item.Extension.Name,
							item.Extension.NodeScope,
							filter.Scope(),
						)
						if lo.Contains(entry.mandatory, item.Extension.Name) {
							Expect(item).Should(MatchCurrentGlobFilter(filter))
						}

						recording[item.Extension.Name] = len(item.Children)
						return nil
					},
				}
			}).Run()

			if entry.mandatory != nil {
				for _, name := range entry.mandatory {
					_, found := recording[name]
					Expect(found).To(BeTrue(), reason(name))
				}
			}

			if entry.prohibited != nil {
				for _, name := range entry.prohibited {
					_, found := recording[name]
					Expect(found).To(BeFalse(), reason(name))
				}
			}
		},
		func(entry *filterTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.message)
		},

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "universal(any scope): glob filter",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
			},
			name:    "items with '.flac' suffix",
			pattern: "*.flac",
			scope:   nav.ScopeAllEn,
		}),
		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "universal(any scope): glob filter (negate)",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
			},
			name:    "items without .flac suffix",
			pattern: "*.flac",
			scope:   nav.ScopeAllEn,
			negate:  true,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "universal(undefined scope): glob filter",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
			},
			name:    "items with '.flac' suffix",
			pattern: "*.flac",
		}),

		// === ifNotApplicable ===============================================

		Entry(nil, &filterTE{ // THIS MAYBE AN INCORRECT TEST, because ifNotApplicable=true, so re-check mandatory/prohibited
			naviTE: naviTE{
				message:      "universal(any scope): glob filter (ifNotApplicable=true)",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
				mandatory:    []string{"A1 - Can You Kiss Me First.flac"},
			},
			name:            "items with '.flac' suffix",
			pattern:         "*.flac",
			scope:           nav.ScopeLeafEn,
			ifNotApplicable: true,
		}),
		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "universal(leaf scope): glob filter (ifNotApplicable=false)",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
				mandatory:    []string{"A1 - Can You Kiss Me First.flac"},
				prohibited:   []string{"vinyl-info.teenage-color"},
			},
			name:            "items with '.flac' suffix",
			pattern:         "*.flac",
			scope:           nav.ScopeLeafEn,
			ifNotApplicable: false,
		}),
	)

	DescribeTable("Filter Children (glob)",
		func(entry *filterTE) {
			recording := recordingMap{}
			filterDefs := &nav.FilterDefinitions{
				Children: nav.CompoundFilterDef{
					Type:        nav.FilterTypeGlobEn,
					Description: entry.name,
					Pattern:     entry.pattern,
					Negate:      entry.negate,
				},
			}
			var filter nav.CompoundTraverseFilter

			path := path(root, entry.relative)
			session := nav.PrimarySession{
				Path: path,
			}
			_ = session.Configure(func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = func(state *nav.NavigationState) {
					GinkgoWriter.Printf(
						"---> 🛡️ [traverse-navigator-test:BEGIN], root: '%v'\n", state.Root,
					)
					filter = state.Filters.Children
				}
				o.Store.Subscription = entry.subscription
				o.Store.FilterDefs = filterDefs
				o.Store.DoExtend = true
				o.Callback = nav.LabelledTraverseCallback{
					Label: "test glob filter callback",
					Fn: func(item *nav.TraverseItem) *translate.LocalisableError {
						actualNoChildren := len(item.Children)
						GinkgoWriter.Printf(
							"===> 💠 Glob Filter(%v, children: %v) source: '%v', item-name: '%v', item-scope: '%v'\n",
							filter.Description(),
							actualNoChildren,
							filter.Source(),
							item.Extension.Name,
							item.Extension.NodeScope,
						)

						recording[item.Extension.Name] = len(item.Children)
						return nil
					},
				}
			}).Run()

			if entry.mandatory != nil {
				for _, name := range entry.mandatory {
					_, found := recording[name]
					Expect(found).To(BeTrue(), reason(name))
				}
			}

			if entry.prohibited != nil {
				for _, name := range entry.prohibited {
					_, found := recording[name]
					Expect(found).To(BeFalse(), reason(name))
				}
			}
			for n, actualNoChildren := range entry.expectedNoOf.children {
				Expect(recording[n]).To(Equal(actualNoChildren), reason(n))
			}
		},
		func(entry *filterTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.message)
		},
		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "folder(with files): glob filter",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFoldersWithFiles,
				expectedNoOf: expectedNo{
					children: map[string]int{
						"Night Drive":      2,
						"Northern Council": 2,
						"Teenage Color":    2,
						"Innerworld":       2,
					},
				},
			},
			name:    "items with '.flac' suffix",
			pattern: "*.flac",
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "folder(with files): glob filter (negate)",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFoldersWithFiles,
				expectedNoOf: expectedNo{
					children: map[string]int{
						"Night Drive":      3,
						"Northern Council": 3,
						"Teenage Color":    2,
						"Innerworld":       2,
					},
				},
			},
			name:    "items without '.txt' suffix",
			pattern: "*.txt",
			negate:  true,
		}),
	)
})
