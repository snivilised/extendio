package nav_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"

	"github.com/snivilised/extendio/internal/helpers"

	. "github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("FilterGlob", Ordered, func() {
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

	DescribeTable("GlobFilter",
		func(entry *filterTE) {
			recording := make(recordingMap)
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

			path := helpers.Path(root, entry.relative)
			optionFn := func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = func(state *nav.NavigationState) {
					GinkgoWriter.Printf(
						"---> 🛡️ [traverse-navigator-test:BEGIN], root: '%v'\n", state.Root,
					)
					filter = state.Filters.Node
				}

				o.Store.Subscription = entry.subscription
				o.Store.FilterDefs = filterDefs
				o.Store.DoExtend = true
				o.Callback = &nav.LabelledTraverseCallback{
					Label: "test glob filter callback",
					Fn: func(item *nav.TraverseItem) error {
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
			}
			result, err := nav.New().Primary(&nav.Prime{
				Path:      path,
				OptionsFn: optionFn,
			}).Run()

			if entry.mandatory != nil {
				for _, name := range entry.mandatory {
					_, found := recording[name]
					Expect(found).To(BeTrue(), helpers.Reason(name))
				}
			}

			if entry.prohibited != nil {
				for _, name := range entry.prohibited {
					_, found := recording[name]
					Expect(found).To(BeFalse(), helpers.Reason(name))
				}
			}

			Expect(err).Error().To(BeNil())

			Expect(result.Metrics.Count(nav.MetricNoFilesInvokedEn)).To(Equal(entry.expectedNoOf.files),
				helpers.BecauseQuantity("Incorrect no of files",
					int(entry.expectedNoOf.files),
					int(result.Metrics.Count(nav.MetricNoFilesInvokedEn)),
				),
			)

			Expect(result.Metrics.Count(nav.MetricNoFoldersInvokedEn)).To(Equal(entry.expectedNoOf.folders),
				helpers.BecauseQuantity("Incorrect no of folders",
					int(entry.expectedNoOf.folders),
					int(result.Metrics.Count(nav.MetricNoFoldersInvokedEn)),
				),
			)

			sum := lo.Sum(lo.Values(entry.expectedNoOf.children))

			Expect(result.Metrics.Count(nav.MetricNoChildFilesFoundEn)).To(Equal(uint(sum)),
				helpers.BecauseQuantity("Incorrect total no of child files",
					sum,
					int(result.Metrics.Count(nav.MetricNoChildFilesFoundEn)),
				),
			)
		},
		func(entry *filterTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.message)
		},

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "universal(any scope): glob filter",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
				expectedNoOf: directoryQuantities{
					files:   8,
					folders: 0,
				},
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
				expectedNoOf: directoryQuantities{
					files:   6,
					folders: 8,
				},
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
				expectedNoOf: directoryQuantities{
					files:   8,
					folders: 0,
				},
			},
			name:    "items with '.flac' suffix",
			pattern: "*.flac",
		}),

		// === ifNotApplicable ===============================================

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "universal(any scope): glob filter (ifNotApplicable=true)",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
				expectedNoOf: directoryQuantities{
					files:   8,
					folders: 4,
				},
				mandatory: []string{"A1 - Can You Kiss Me First.flac"},
			},
			name:            "items with '.flac' suffix",
			pattern:         "*.flac",
			scope:           nav.ScopeLeafEn,
			ifNotApplicable: nav.TriStateBoolTrueEn,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "universal(leaf scope): glob filter (ifNotApplicable=false)",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
				expectedNoOf: directoryQuantities{
					files:   8,
					folders: 0,
				},
				mandatory:  []string{"A1 - Can You Kiss Me First.flac"},
				prohibited: []string{"vinyl-info.teenage-color"},
			},
			name:            "items with '.flac' suffix",
			pattern:         "*.flac",
			scope:           nav.ScopeLeafEn,
			ifNotApplicable: nav.TriStateBoolFalseEn,
		}),
	)

	DescribeTable("Filter Children (glob)",
		func(entry *filterTE) {
			recording := make(recordingMap)
			filterDefs := &nav.FilterDefinitions{
				Children: nav.CompoundFilterDef{
					Type:        nav.FilterTypeGlobEn,
					Description: entry.name,
					Pattern:     entry.pattern,
					Negate:      entry.negate,
				},
			}
			var filter nav.CompoundTraverseFilter

			path := helpers.Path(root, entry.relative)
			optionFn := func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = func(state *nav.NavigationState) {
					GinkgoWriter.Printf(
						"---> 🛡️ [traverse-navigator-test:BEGIN], root: '%v'\n", state.Root,
					)
					filter = state.Filters.Children
				}
				o.Store.Subscription = entry.subscription
				o.Store.FilterDefs = filterDefs
				o.Store.DoExtend = true
				o.Callback = &nav.LabelledTraverseCallback{
					Label: "test glob filter callback",
					Fn: func(item *nav.TraverseItem) error {
						actualNoChildren := len(item.Children)
						GinkgoWriter.Printf(
							"===> 💠 Compound Glob Filter(%v, children: %v) source: '%v', item-name: '%v', item-scope: '%v', depth: '%v'\n",
							filter.Description(),
							actualNoChildren,
							filter.Source(),
							item.Extension.Name,
							item.Extension.NodeScope,
							item.Extension.Depth,
						)

						recording[item.Extension.Name] = len(item.Children)
						return nil
					},
				}
			}

			result, _ := nav.New().Primary(&nav.Prime{
				Path:      path,
				OptionsFn: optionFn,
			}).Run()

			if entry.mandatory != nil {
				for _, name := range entry.mandatory {
					_, found := recording[name]
					Expect(found).To(BeTrue(), helpers.Reason(name))
				}
			}

			if entry.prohibited != nil {
				for _, name := range entry.prohibited {
					_, found := recording[name]
					Expect(found).To(BeFalse(), helpers.Reason(name))
				}
			}
			for n, actualNoChildren := range entry.expectedNoOf.children {
				expected := recording[n]

				Expect(expected).To(Equal(actualNoChildren),
					helpers.BecauseQuantity("Incorrect no of children",
						expected,
						actualNoChildren,
					),
				)
			}

			Expect(result.Metrics.Count(nav.MetricNoFilesInvokedEn)).To(Equal(entry.expectedNoOf.files),
				helpers.BecauseQuantity("Incorrect no of files",
					int(entry.expectedNoOf.files),
					int(result.Metrics.Count(nav.MetricNoFilesInvokedEn)),
				),
			)

			Expect(result.Metrics.Count(nav.MetricNoFoldersInvokedEn)).To(Equal(entry.expectedNoOf.folders),
				helpers.BecauseQuantity("Incorrect no of folders",
					int(entry.expectedNoOf.folders),
					int(result.Metrics.Count(nav.MetricNoFoldersInvokedEn)),
				),
			)
		},
		func(entry *filterTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.message)
		},

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "folder(with files): glob filter",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFoldersWithFiles,
				expectedNoOf: directoryQuantities{
					files:   0,
					folders: 8,
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
				expectedNoOf: directoryQuantities{
					files:   0,
					folders: 8,
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
