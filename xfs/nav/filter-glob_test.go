package nav_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/extendio/translate"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("FilterGlob", Ordered, func() {
	var root string

	BeforeAll(func() {
		root = cwd()
	})

	DescribeTable("GlobFilter",
		func(entry *filterTE) {
			recording := recordingMap{}

			navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = begin("🛡️")
				o.Subscription = entry.subscription
				o.Filters.Current = &nav.GlobFilter{
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
						"===> 💠 Glob Filter(%v) source: '%v', item-name: '%v', item-scope(fs): '%v(%v)'\n",
						o.Filters.Current.Description(), o.Filters.Current.Source(),
						item.Extension.Name, item.Extension.NodeScope, o.Filters.Current.Scope(),
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

		// === ifNotApplicable ===============================================

		Entry(nil, &filterTE{
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
				message:      "universal(any scope): glob filter (ifNotApplicable=false)",
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
})
