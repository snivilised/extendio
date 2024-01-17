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

var _ = Describe("Filter Extended glob", Ordered, func() {
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

	DescribeTable("Filter Extended glob",
		func(entry *filterTE) {
			recording := make(recordingMap)
			filterDefs := &nav.FilterDefinitions{
				Node: nav.FilterDef{
					Type:            nav.FilterTypeExtendedGlobEn,
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
				o.Callback = &nav.LabelledTraverseCallback{
					Label: "test extended glob filter callback",
					Fn: func(item *nav.TraverseItem) error {
						indicator := lo.Ternary(item.IsDir(), "📁", "💠")
						GinkgoWriter.Printf(
							"===> %v extended glob Filter(%v) source: '%v', item-name: '%v', item-scope(fs): '%v(%v)'\n",
							indicator,
							filter.Description(),
							filter.Source(),
							item.Extension.Name,
							item.Extension.NodeScope,
							filter.Scope(),
						)
						if lo.Contains(entry.mandatory, item.Extension.Name) {
							Expect(item).Should(MatchCurrentExtendedFilter(filter))
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

		// === universal =====================================================

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "universal(any scope): extended glob filter",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: nav.SubscribeAny,
				expectedNoOf: directoryQuantities{
					files:   16,
					folders: 5,
				},
				prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			name:    "items with 'flac' suffix",
			pattern: "*|flac",
			scope:   nav.ScopeAllEn,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "universal(any scope): extended glob filter, with dot extension",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: nav.SubscribeAny,
				expectedNoOf: directoryQuantities{
					files:   16,
					folders: 5,
				},
				prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			name:    "items with 'flac' suffix",
			pattern: "*|.flac",
			scope:   nav.ScopeAllEn,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "universal(any scope): extended glob filter, with multiple extensions",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: nav.SubscribeAny,
				expectedNoOf: directoryQuantities{
					files:   19,
					folders: 5,
				},
				mandatory:  []string{"front.jpg"},
				prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			name:    "items with 'flac' suffix",
			pattern: "*|flac,jpg",
			scope:   nav.ScopeAllEn,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "universal(any scope): extended glob filter, without extension",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: nav.SubscribeAny,
				expectedNoOf: directoryQuantities{
					files:   3,
					folders: 5,
				},
				mandatory:  []string{"cover-clutching-at-straws-jpg"},
				prohibited: []string{"01 - Hotel Hobbies.flac"},
			},
			name:    "items with 'flac' suffix",
			pattern: "*|",
			scope:   nav.ScopeAllEn,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "universal(file scope): extended glob filter (negate)",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: nav.SubscribeAny,
				expectedNoOf: directoryQuantities{
					files:   7,
					folders: 5,
				},
				prohibited: []string{"01 - Hotel Hobbies.flac"},
			},
			name:    "files without .flac suffix",
			pattern: "*|flac",
			scope:   nav.ScopeFileEn,
			negate:  true,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "universal(undefined scope): extended glob filter",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: nav.SubscribeAny,
				expectedNoOf: directoryQuantities{
					files:   16,
					folders: 5,
				},
				prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			name:    "items with '.flac' suffix",
			pattern: "*|flac",
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "universal(any scope): extended glob filter, any extension",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: nav.SubscribeAny,
				expectedNoOf: directoryQuantities{
					files:   4,
					folders: 1,
				},
				mandatory:  []string{"cover-clutching-at-straws-jpg"},
				prohibited: []string{"01 - Hotel Hobbies.flac"},
			},
			name:    "starts with c, any extension",
			pattern: "c*|*",
			scope:   nav.ScopeAllEn,
		}),

		// === folders =======================================================

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "folders(any scope): extended glob filter",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: nav.SubscribeFolders,
				expectedNoOf: directoryQuantities{
					files:   0,
					folders: 2,
				},
				mandatory:  []string{"Marillion"},
				prohibited: []string{"Fugazi"},
			},
			name:    "folders starting with M",
			pattern: "M*|",
			scope:   nav.ScopeFolderEn,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "folders(folder scope): extended glob filter (negate)",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: nav.SubscribeFolders,
				expectedNoOf: directoryQuantities{
					files:   0,
					folders: 3,
				},
				mandatory:  []string{"Fugazi"},
				prohibited: []string{"Marillion"},
			},
			name:    "folders NOT starting with M",
			pattern: "M*|",
			scope:   nav.ScopeFolderEn,
			negate:  true,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "universal(undefined scope): extended glob filter",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: nav.SubscribeFolders,
				expectedNoOf: directoryQuantities{
					files:   0,
					folders: 2,
				},
				mandatory:  []string{"Marillion"},
				prohibited: []string{"Fugazi"},
			},
			name:    "folders starting with M",
			pattern: "M*|",
		}),

		// === files =========================================================

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "files(file scope): extended glob filter",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: nav.SubscribeFiles,
				expectedNoOf: directoryQuantities{
					files:   16,
					folders: 0,
				},
				mandatory:  []string{"01 - Hotel Hobbies.flac"},
				prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			name:    "items with 'flac' suffix",
			pattern: "*|flac",
			scope:   nav.ScopeFileEn,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "files(any scope): extended glob filter, with dot extension",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: nav.SubscribeFiles,
				expectedNoOf: directoryQuantities{
					files:   16,
					folders: 0,
				},
				mandatory:  []string{"01 - Hotel Hobbies.flac"},
				prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			name:    "items with 'flac' suffix",
			pattern: "*|.flac",
			scope:   nav.ScopeFileEn,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "files(file scope): extended glob filter, with multiple extensions",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: nav.SubscribeFiles,
				expectedNoOf: directoryQuantities{
					files:   19,
					folders: 0,
				},
				mandatory:  []string{"front.jpg"},
				prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			name:    "items with 'flac' suffix",
			pattern: "*|flac,jpg",
			scope:   nav.ScopeFileEn,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "file(file scope): extended glob filter, without extension",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: nav.SubscribeFiles,
				expectedNoOf: directoryQuantities{
					files:   3,
					folders: 0,
				},
				mandatory:  []string{"cover-clutching-at-straws-jpg"},
				prohibited: []string{"01 - Hotel Hobbies.flac"},
			},
			name:    "items with 'flac' suffix",
			pattern: "*|",
			scope:   nav.ScopeFileEn,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "file(file scope): extended glob filter (negate)",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: nav.SubscribeFiles,
				expectedNoOf: directoryQuantities{
					files:   7,
					folders: 0,
				},
				mandatory:  []string{"cover-clutching-at-straws-jpg"},
				prohibited: []string{"01 - Hotel Hobbies.flac"},
			},
			name:    "files without .flac suffix",
			pattern: "*|flac",
			scope:   nav.ScopeFileEn,
			negate:  true,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "file(undefined scope): extended glob filter",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: nav.SubscribeFiles,
				expectedNoOf: directoryQuantities{
					files:   16,
					folders: 0,
				},
				mandatory:  []string{"01 - Hotel Hobbies.flac"},
				prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			name:    "items with '.flac' suffix",
			pattern: "*|flac",
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "file(any scope): extended glob filter, any extension",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: nav.SubscribeFiles,
				expectedNoOf: directoryQuantities{
					files:   4,
					folders: 0,
				},
				mandatory:  []string{"cover-clutching-at-straws-jpg"},
				prohibited: []string{"01 - Hotel Hobbies.flac"},
			},
			name:    "starts with c, any extension",
			pattern: "c*|*",
			scope:   nav.ScopeAllEn,
		}),

		// === ifNotApplicable ===============================================

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "universal(leaf scope): extended glob filter (ifNotApplicable=true)",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: nav.SubscribeAny,
				expectedNoOf: directoryQuantities{
					files:   16,
					folders: 5,
				},
				mandatory:  []string{"Marillion"},
				prohibited: []string{"cover-clutching-at-straws-jpg"},
			},
			name:            "leaf items with 'flac' suffix",
			pattern:         "*|flac",
			scope:           nav.ScopeLeafEn,
			ifNotApplicable: nav.TriStateBoolTrueEn,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "universal(leaf scope): extended glob filter (ifNotApplicable=false)",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: nav.SubscribeAny,
				expectedNoOf: directoryQuantities{
					files:   16,
					folders: 4,
				},
				prohibited: []string{"Marillion"},
			},
			name:            "items with '.flac' suffix",
			pattern:         "*|flac",
			scope:           nav.ScopeLeafEn,
			ifNotApplicable: nav.TriStateBoolFalseEn,
		}),

		// === with-exclusion ================================================

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "universal(any scope): extended glob filter with exclusion",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: nav.SubscribeFiles,
				expectedNoOf: directoryQuantities{
					files:   12,
					folders: 0,
				},
				prohibited: []string{"01 - Hotel Hobbies.flac"},
			},
			name:    "files starting with 0, except 01 items and flac suffix",
			pattern: "0*/*01*|flac",
			scope:   nav.ScopeFileEn,
		}),
	)

	DescribeTable("Filter Children (extended glob)",
		func(entry *filterTE) {
			recording := make(recordingMap)
			filterDefs := &nav.FilterDefinitions{
				Children: nav.CompoundFilterDef{
					Type:        nav.FilterTypeExtendedGlobEn,
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
				o.Callback = &nav.LabelledTraverseCallback{
					Label: "test extended glob filter callback",
					Fn: func(item *nav.TraverseItem) error {
						actualNoChildren := len(item.Children)
						indicator := lo.Ternary(item.IsDir(), "📁", "💠")
						GinkgoWriter.Printf(
							"===> %v Compound extended glob Filter(%v, children: %v) source: '%v', item-name: '%v', item-scope: '%v', depth: '%v'\n",
							indicator,
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
				message:      "folder(with files): extended glob filter",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: nav.SubscribeFoldersWithFiles,
				expectedNoOf: directoryQuantities{
					files:   0,
					folders: 5,
					children: map[string]int{
						"Clutching At Straws":       4,
						"Fugazi":                    4,
						"Misplaced Childhood":       4,
						"Script for a Jesters Tear": 4,
					},
				},
			},
			name:    "items with 'flac' suffix",
			pattern: "*|flac",
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "folder(with files): extended glob filter (negate)",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: nav.SubscribeFoldersWithFiles,
				expectedNoOf: directoryQuantities{
					files:   0,
					folders: 5,
					children: map[string]int{
						"Clutching At Straws":       7,
						"Fugazi":                    5,
						"Misplaced Childhood":       5,
						"Script for a Jesters Tear": 5,
					},
				},
			},
			name:    "items without '.txt' suffix",
			pattern: "*|txt",
			negate:  true,
		}),

		// === with-exclusion ================================================

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "folder(with files): extended glob filter with exclusion",
				relative:     "rock/PROGRESSIVE-ROCK/Marillion",
				subscription: nav.SubscribeFoldersWithFiles,
				expectedNoOf: directoryQuantities{
					files:   0,
					folders: 5,
					children: map[string]int{
						"Clutching At Straws":       3,
						"Fugazi":                    3,
						"Misplaced Childhood":       3,
						"Script for a Jesters Tear": 3,
					},
				},
			},
			name:    "files starting with 0, except 01 items and flac suffix",
			pattern: "0*/*01*|flac",
		}),
	)
})
