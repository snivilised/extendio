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
		root = cwd()
	})

	DescribeTable("scope",
		func(entry *scopeTE) {
			recording := recordingScopeMap{}

			scopeRecorder := func(item *nav.TraverseItem) *LocalisableError {
				_, found := recording[item.Extension.Name]

				if !found {
					recording[item.Extension.Name] = item.Extension.NodeScope
				}
				return entry.callback(item)
			}

			navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = begin("🛡️")
				o.Subscription = entry.subscription
				o.DoExtend = true
				o.Callback = scopeRecorder
			})

			path := path(root, entry.relative)
			_ = navigator.Walk(path)

			for p, expected := range entry.expectedScopes {
				actual := recording[p]
				Expect(expected).To(Equal(actual))
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
				"Night Drive":                  nav.TopScopeEn | nav.LeafScopeEn,
				"A1 - The Telephone Call.flac": nav.LeafScopeEn,
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
				"RETRO-WAVE":                   nav.TopScopeEn,
				"Night Drive":                  nav.LeafScopeEn,
				"A1 - The Telephone Call.flac": nav.LeafScopeEn,
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
				"Night Drive": nav.TopScopeEn | nav.LeafScopeEn,
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
				"RETRO-WAVE":  nav.TopScopeEn,
				"Chromatics":  nav.IntermediateScopeEn,
				"Night Drive": nav.LeafScopeEn,
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
				"segments.bass.infex.txt": nav.LeafScopeEn,
			},
		}),
	)
})