package nav_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"           //nolint:revive // ginkgo ok
	. "github.com/onsi/gomega"              //nolint:revive // gomega ok
	. "github.com/snivilised/extendio/i18n" //nolint:revive // i18n ok
	"github.com/snivilised/extendio/internal/helpers"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("TraverseNavigatorScope", Ordered, func() {
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

	DescribeTable("scope",
		func(entry *scopeTE) {
			recording := make(recordingScopeMap)

			scopeRecorder := &nav.LabelledTraverseCallback{
				Label: "test callback",
				Fn: func(item *nav.TraverseItem) error {
					_, found := recording[item.Extension.Name]

					if !found {
						recording[item.Extension.Name] = item.Extension.NodeScope
					}
					return entry.callback.Fn(item)
				},
			}

			path := helpers.Path(root, entry.relative)
			optionFn := func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = begin("🛡️")
				o.Store.Subscription = entry.subscription
				o.Callback = scopeRecorder
			}

			result, _ := nav.New().Primary(&nav.Prime{
				Path:      path,
				OptionsFn: optionFn,
			}).Run()

			for name, expected := range entry.expectedScopes {
				actual := recording[name]
				Expect(actual).To(Equal(expected), helpers.Reason(name))
			}

			_ = result.Session.StartedAt()
			_ = result.Session.Elapsed()
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
				"Night Drive":                  nav.ScopeRootEn | nav.ScopeLeafEn | nav.ScopeFolderEn,
				"A1 - The Telephone Call.flac": nav.ScopeLeafEn | nav.ScopeFileEn,
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
				"RETRO-WAVE":                   nav.ScopeRootEn | nav.ScopeFolderEn,
				"Night Drive":                  nav.ScopeLeafEn | nav.ScopeFolderEn,
				"A1 - The Telephone Call.flac": nav.ScopeLeafEn | nav.ScopeFileEn,
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
				"Night Drive": nav.ScopeTopEn | nav.ScopeLeafEn | nav.ScopeFolderEn,
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
				"Night Drive": nav.ScopeRootEn | nav.ScopeLeafEn | nav.ScopeFolderEn,
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
				"RETRO-WAVE":  nav.ScopeRootEn | nav.ScopeFolderEn,
				"Chromatics":  nav.ScopeTopEn | nav.ScopeFolderEn,
				"Night Drive": nav.ScopeLeafEn | nav.ScopeFolderEn,
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
				"segments.bass.infex.txt": nav.ScopeLeafEn | nav.ScopeFileEn,
			},
		}),
	)
})
