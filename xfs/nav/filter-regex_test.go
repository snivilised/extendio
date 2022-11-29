package nav_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"

	"github.com/snivilised/extendio/translate"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("FilterRegex", Ordered, func() {
	var root string

	BeforeAll(func() {
		root = origin()
	})

	DescribeTable("RegexFilter",
		func(entry *filterTE) {
			recording := recordingMap{}
			filterDefs := &nav.FilterDefinitions{
				Current: nav.FilterDef{
					Type:            nav.FilterTypeRegexEn,
					Description:     entry.name,
					Source:          entry.pattern,
					Scope:           entry.scope,
					Negate:          entry.negate,
					IfNotApplicable: entry.ifNotApplicable,
				},
			}
			var filter nav.TraverseFilter

			navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = func(state *nav.NavigationState) {
					GinkgoWriter.Printf(
						"---> 🛡️ [traverse-navigator-test:BEGIN], root: '%v'\n", state.Root,
					)
					filter = state.Filters.Current
				}

				o.Store.Subscription = entry.subscription
				o.Store.FilterDefs = filterDefs
				o.Store.DoExtend = true
				o.Callback = nav.LabelledTraverseCallback{
					Label: "test regex filter callback",
					Fn: func(item *nav.TraverseItem) *translate.LocalisableError {
						GinkgoWriter.Printf(
							"===> ⚗️ Regex Filter(%v) source: '%v', item-name: '%v', item-scope(fs): '%v(%v)'\n",
							filter.Description(),
							filter.Source(),
							item.Extension.Name,
							item.Extension.NodeScope,
							filter.Scope(),
						)
						if lo.Contains(entry.mandatory, item.Extension.Name) {
							Expect(item).Should(MatchCurrentRegexFilter(filter))
						}

						recording[item.Extension.Name] = len(item.Children)
						return nil
					},
				}
			})
			path := path(root, entry.relative)
			_ = navigator.Walk(path)

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

		// === files =========================================================

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "files(any scope): regex filter",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFiles,
			},
			name:    "items that start with 'vinyl'",
			pattern: "^vinyl",
			scope:   nav.ScopeAllEn,
		}),
		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "files(any scope): regex filter (negate)",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFiles,
			},
			name:    "items that don't start with 'vinyl'",
			pattern: "^vinyl",
			scope:   nav.ScopeAllEn,
			negate:  true,
		}),
		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "files(default to any scope): regex filter",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFiles,
			},
			name:    "items that start with 'vinyl'",
			pattern: "^vinyl",
		}),

		// === folders =======================================================

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "folders(any scope): regex filter",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFolders,
			},
			name:    "items that start with 'C'",
			pattern: "^C",
			scope:   nav.ScopeAllEn,
		}),
		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "folders(any scope): regex filter (negate)",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFolders,
			},
			name:    "items that don't start with 'C'",
			pattern: "^C",
			scope:   nav.ScopeAllEn,
			negate:  true,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "folders(undefined scope): regex filter",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFolders,
			},
			name:    "items that start with 'C'",
			pattern: "^C",
		}),

		// === ifNotApplicable ===============================================

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "folders(top): regex filter (ifNotApplicable=true)",
				relative:     "PROGRESSIVE-HOUSE",
				subscription: nav.SubscribeFolders,
				mandatory:    []string{"PROGRESSIVE-HOUSE"},
			},
			name:            "top items that contain 'HOUSE'",
			pattern:         "HOUSE",
			scope:           nav.ScopeTopEn,
			ifNotApplicable: true,
		}),
		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "folders(top): regex filter (ifNotApplicable=false)",
				relative:     "",
				subscription: nav.SubscribeFolders,
				mandatory:    []string{"PROGRESSIVE-HOUSE"},
				prohibited:   []string{"Blue Amazon", "The Javelin"},
			},
			name:            "top items that contain 'HOUSE'",
			pattern:         "HOUSE",
			scope:           nav.ScopeTopEn,
			ifNotApplicable: false,
		}),
	)

	DescribeTable("Filter Children (regex)",
		func(entry *filterTE) {
			recording := recordingMap{}
			filterDefs := &nav.FilterDefinitions{
				Children: nav.CompoundFilterDef{
					Type:        nav.FilterTypeRegexEn,
					Description: entry.name,
					Source:      entry.pattern,
					Negate:      entry.negate,
				},
			}
			var filter nav.CompoundTraverseFilter

			navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = func(state *nav.NavigationState) {
					GinkgoWriter.Printf(
						"---> 🛡️ [traverse-navigator-test:BEGIN], root: '%v'\n", state.Root,
					)
					filter = state.Filters.Compound
				}
				o.Store.Subscription = entry.subscription
				o.Store.FilterDefs = filterDefs
				o.Store.DoExtend = true
				o.Callback = nav.LabelledTraverseCallback{
					Label: "test regex filter callback",
					Fn: func(item *nav.TraverseItem) *translate.LocalisableError {
						actualNoChildren := len(item.Children)
						GinkgoWriter.Printf(
							"===> 💠 Regex Filter(%v, children: %v) source: '%v', item-name: '%v', item-scope: '%v'\n",
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
			})
			path := path(root, entry.relative)
			_ = navigator.Walk(path)

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
			for n, actualNoChildren := range entry.expectedNoChildren {
				Expect(recording[n]).To(Equal(actualNoChildren), reason(n))
			}
		},
		func(entry *filterTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.message)
		},
		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "folder(with files): regex filter",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFoldersWithFiles,
				expectedNoChildren: map[string]int{
					"Night Drive":      2,
					"Northern Council": 2,
					"Teenage Color":    2,
					"Innerworld":       2,
				},
			},
			name:    "items with '.flac' suffix",
			pattern: "\\.flac$",
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "folder(with files): regex filter (negate)",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFoldersWithFiles,
				expectedNoChildren: map[string]int{
					"Night Drive":      3,
					"Northern Council": 3,
					"Teenage Color":    2,
					"Innerworld":       2,
				},
			},
			name:    "items without '.txt' suffix",
			pattern: "\\.txt$",
			negate:  true,
		}),
	)

	DescribeTable("RegexFilter (error)",
		func(entry *filterTE) {
			defer func() {
				pe := recover()
				if entry.errorContains != "" {
					if err, ok := pe.(error); ok {
						Expect(strings.Contains(err.Error(), entry.errorContains)).To(BeTrue())
					}
				} else {
					Expect(pe).To(Equal(entry.expectedErr))
				}
			}()

			filterDefs := &nav.FilterDefinitions{
				Current: nav.FilterDef{
					Type:        nav.FilterTypeRegexEn,
					Description: entry.name,
					Source:      entry.pattern,
				},
			}

			navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = begin("🧲")
				o.Store.Subscription = nav.SubscribeFolders
				o.Store.FilterDefs = filterDefs
				o.Callback = nav.LabelledTraverseCallback{
					Label: "test regex filter callback",
					Fn: func(item *nav.TraverseItem) *translate.LocalisableError {
						return nil
					},
				}
			})
			const relative = "RETRO-WAVE"
			path := path(root, relative)
			_ = navigator.Walk(path)

			Fail(fmt.Sprintf("❌ expected panic due to '%v'", entry.name))
		},
		func(entry *filterTE) string {
			return fmt.Sprintf("🧪 ===> '%v'", entry.message)
		},

		Entry(nil, &filterTE{
			naviTE:        naviTE{message: "bad regex pattern"},
			name:          "bad regex pattern test",
			pattern:       "(",
			errorContains: "Compile",
		}),
	)

	DescribeTable("CompoundRegexFilter (error)",
		func(entry *filterTE) {
			defer func() {
				pe := recover()
				if entry.errorContains != "" {
					if err, ok := pe.(error); ok {
						Expect(strings.Contains(err.Error(), entry.errorContains)).To(BeTrue())
					}
				} else {
					Expect(pe).To(Equal(entry.expectedErr))
				}
			}()

			filterDefs := &nav.FilterDefinitions{
				Children: nav.CompoundFilterDef{
					Type:        nav.FilterTypeRegexEn,
					Description: entry.name,
					Source:      entry.pattern,
					Negate:      entry.negate,
				},
			}

			navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = begin("🧲")
				o.Store.Subscription = nav.SubscribeFoldersWithFiles
				o.Store.FilterDefs = filterDefs
				o.Store.DoExtend = true
				o.Callback = nav.LabelledTraverseCallback{
					Label: "test regex filter callback",
					Fn: func(item *nav.TraverseItem) *translate.LocalisableError {
						return nil
					},
				}
			})
			const relative = "RETRO-WAVE"
			path := path(root, relative)
			_ = navigator.Walk(path)

			Fail(fmt.Sprintf("❌ expected panic due to '%v'", entry.name))
		},
		func(entry *filterTE) string {
			return fmt.Sprintf("🧪 ===> '%v'", entry.message)
		},

		Entry(nil, &filterTE{
			naviTE:        naviTE{message: "bad regex pattern"},
			name:          "bad regex pattern test",
			pattern:       "(",
			errorContains: "Compile",
		}),
	)
})
