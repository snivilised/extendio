package nav_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/snivilised/extendio/translate"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("TraverseNavigatorScope", Ordered, func() {
	var root string

	BeforeAll(func() {
		root = origin()
	})

	DescribeTable("scope",
		func(entry *scopeTE) {
			recording := recordingScopeMap{}

			scopeRecorder := nav.LabelledTraverseCallback{
				Label: "test callback",
				Fn: func(item *nav.TraverseItem) *LocalisableError {
					_, found := recording[item.Extension.Name]

					if !found {
						recording[item.Extension.Name] = item.Extension.NodeScope
					}
					return entry.callback.Fn(item)
				},
			}

			navigator := nav.NavigatorFactory{}.Construct(func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = begin("🛡️")
				o.Store.Subscription = entry.subscription
				o.Store.DoExtend = true
				o.Callback = scopeRecorder
			})

			path := path(root, entry.relative)
			_ = navigator.Walk(path)

			for name, expected := range entry.expectedScopes {
				actual := recording[name]
				Expect(actual).To(Equal(expected), reason(name))
			}
		},
		func(entry *scopeTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.message)
		},

		// === universal =====================================================

		Entry(nil, &scopeTE{
			naviTE: naviTE{
				message:      "universal: Path is leaf",
				relative:     "RETRO-WAVE/Chromatics/Night Drive",
				subscription: nav.SubscribeAny,
				callback:     universalScopeCallback("LEAF-PATH"),
			},
			expectedScopes: recordingScopeMap{
				"Night Drive":                  nav.ScopeRootEn | nav.ScopeLeafEn,
				"A1 - The Telephone Call.flac": nav.ScopeLeafEn,
			},
		}),
		Entry(nil, &scopeTE{
			naviTE: naviTE{
				message:      "universal: Path contains folders",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
				callback:     universalScopeCallback("CONTAINS-FOLDERS"),
			},
			expectedScopes: recordingScopeMap{
				"RETRO-WAVE":                   nav.ScopeRootEn,
				"Night Drive":                  nav.ScopeLeafEn,
				"A1 - The Telephone Call.flac": nav.ScopeLeafEn,
			},
		}),

		Entry(nil, &scopeTE{
			naviTE: naviTE{
				message:      "universal: Path contains folders (Top & Leaf)",
				relative:     "RETRO-WAVE/Chromatics",
				subscription: nav.SubscribeAny,
				callback:     universalScopeCallback("CONTAINS-FOLDERS"),
			},
			expectedScopes: recordingScopeMap{
				"Night Drive": nav.ScopeTopEn | nav.ScopeLeafEn,
			},
		}),

		// === folders =======================================================

		Entry(nil, &scopeTE{
			naviTE: naviTE{
				message:      "folders: Path is leaf",
				relative:     "RETRO-WAVE/Chromatics/Night Drive",
				subscription: nav.SubscribeFolders,
				callback:     foldersScopeCallback("LEAF-PATH"),
			},
			expectedScopes: recordingScopeMap{
				"Night Drive": nav.ScopeRootEn | nav.ScopeLeafEn,
			},
		}),
		Entry(nil, &scopeTE{
			naviTE: naviTE{
				message:      "folders: Path contains folders",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFolders,
				callback:     foldersScopeCallback("CONTAINS-FOLDERS"),
			},
			expectedScopes: recordingScopeMap{
				"RETRO-WAVE":  nav.ScopeRootEn,
				"Chromatics":  nav.ScopeTopEn,
				"Night Drive": nav.ScopeLeafEn,
			},
		}),

		// === files =========================================================

		Entry(nil, &scopeTE{
			naviTE: naviTE{
				message:      "files: Path contains non-leaf files",
				relative:     "bass",
				subscription: nav.SubscribeFiles,
				callback:     filesScopeCallback("CONTAINS-FOLDERS"),
			},
			expectedScopes: recordingScopeMap{
				"segments.bass.infex.txt": nav.ScopeLeafEn,
			},
		}),
	)
})
