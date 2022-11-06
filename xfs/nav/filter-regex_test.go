package nav_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/extendio/translate"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("FilterRegex", Ordered, func() {
	var root string

	BeforeAll(func() {
		root = cwd()
	})

	DescribeTable("RegexFilter",
		func(entry *filterTE) {
			recording := recordingMap{}

			navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = begin("🛡️")
				o.Subscription = entry.subscription
				o.Filters.Current = &nav.RegexFilter{
					Filter: nav.Filter{
						Name:            entry.name,
						RequiredScope:   entry.scope,
						Pattern:         entry.pattern,
						Negate:          entry.negate,
						IfNotApplicable: entry.ifNotApplicable,
					},
				}
				o.DoExtend = true
				o.Callback = func(item *nav.TraverseItem) *translate.LocalisableError {
					GinkgoWriter.Printf(
						"===> ⚗️ Regex Filter(%v) source: '%v', item-name: '%v', item-scope(fs): '%v(%v)'\n",
						o.Filters.Current.Description(), o.Filters.Current.Source(), item.Extension.Name,
						item.Extension.NodeScope, o.Filters.Current.Scope(),
					)
					Expect(o.Filters.Current.IsMatch(item)).To(BeTrue(), reason(item.Extension.Name))
					recording[item.Extension.Name] = true
					return nil
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

			navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = begin("🧲")
				o.Subscription = nav.SubscribeFolders
				o.Filters.Current = &nav.RegexFilter{
					Filter: nav.Filter{
						Name:    entry.name,
						Pattern: entry.pattern,
					},
				}
				o.Callback = func(item *nav.TraverseItem) *translate.LocalisableError {
					return nil
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
			naviTE:      naviTE{message: "regex missing pattern"},
			name:        "missing pattern test",
			expectedErr: nav.PATTERN_NOT_DEFINED_L_ERR,
		}),
		Entry(nil, &filterTE{
			naviTE:        naviTE{message: "bad regex pattern"},
			name:          "bad regex pattern test",
			pattern:       "(",
			errorContains: "Compile",
		}),
	)
})
