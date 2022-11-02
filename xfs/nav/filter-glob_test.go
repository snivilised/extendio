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
				o.Filter = &nav.GlobFilter{
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
					GinkgoWriter.Printf("===> 💠 Glob Filter(%v) source: '%v', item-name: '%v', item-scope(fs): '%v(%v)'\n",
						o.Filter.Description(), o.Filter.Source(), item.Extension.Name, item.Extension.NodeScope, o.Filter.Scope(),
					)
					Expect(o.Filter.IsMatch(item)).To(BeTrue(), named(item.Extension.Name))
					recording[item.Extension.Name] = true
					return nil
				}
			})
			path := path(root, entry.relative)
			_ = navigator.Walk(path)

			if entry.mandatory != nil {
				for _, name := range entry.mandatory {
					_, found := recording[name]
					Expect(found).To(BeTrue(), named(name))
				}
			}

			if entry.prohibited != nil {
				for _, name := range entry.prohibited {
					_, found := recording[name]
					Expect(found).To(BeFalse(), named(name))
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
			scope:   nav.AllScopesEn,
		}),
		Entry(nil, &filterTE{
			naviTE: naviTE{
				message:      "universal(any scope): glob filter (negate)",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
			},
			name:    "items without .flac suffix",
			pattern: "*.flac",
			scope:   nav.AllScopesEn,
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
			scope:           nav.LeafScopeEn,
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
			scope:           nav.LeafScopeEn,
			ifNotApplicable: false,
		}),
	)
})