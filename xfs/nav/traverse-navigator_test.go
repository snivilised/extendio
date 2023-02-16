package nav_test

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	. "github.com/snivilised/extendio/translate"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("TraverseNavigator(logged)", Ordered, func() {
	var root string

	BeforeAll(func() {
		root = origin()
	})

	Context("Navigator", func() {
		DescribeTable("Ensure Callback Invoked Once",
			func(entry *naviTE) {
				recording := recordingMap{}
				visited := []string{}

				once := nav.LabelledTraverseCallback{
					Label: "test once decorator",
					Fn: func(item *nav.TraverseItem) *LocalisableError {
						_, found := recording[item.Path]
						Expect(found).To(BeFalse())
						recording[item.Path] = len(item.Children)

						return entry.callback.Fn(item)
					},
				}

				visitor := nav.LabelledTraverseCallback{
					Fn: func(item *nav.TraverseItem) *LocalisableError {
						// just kept to enable visitor specific debug activity
						//
						return once.Fn(item)
					},
				}
				callback := lo.Ternary(entry.once, once, lo.Ternary(entry.visit, visitor, entry.callback))

				path := path(root, entry.relative)
				session := &nav.PrimarySession{
					Path: path,
				}

				// TODO: check that the metric counts from the result are as expected
				//
				_ = session.Configure(func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🛡️")
					o.Store.Subscription = entry.subscription
					o.Store.Behaviours.Sort.IsCaseSensitive = entry.caseSensitive
					o.Store.DoExtend = entry.extended
					o.Callback = callback
				}).Run()

				if entry.visit {
					_ = filepath.WalkDir(path, func(path string, de fs.DirEntry, err error) error {
						if subscribes(entry.subscription, de) {
							visited = append(visited, path)
						}
						return nil
					})
				}

				if entry.visit {
					every := lo.EveryBy(visited, func(p string) bool {
						_, found := recording[p]
						return found
					})
					Expect(every).To(BeTrue())
				}
			},
			func(entry *naviTE) string {
				return fmt.Sprintf("🧪 ===> given: '%v'", entry.message)
			},

			// === universal =====================================================

			Entry(nil, &naviTE{
				message:      "universal: Path is leaf",
				relative:     "RETRO-WAVE/Chromatics/Night Drive",
				extended:     IsExtended,
				subscription: nav.SubscribeAny,
				callback:     universalCallback("LEAF-PATH", IsExtended),
			}),
			Entry(nil, &naviTE{
				message:      "universal: Path contains folders",
				relative:     "RETRO-WAVE",
				extended:     NotExtended,
				subscription: nav.SubscribeAny,
				callback:     universalCallback("CONTAINS-FOLDERS", NotExtended),
			}),
			Entry(nil, &naviTE{
				message:      "universal: Path contains folders",
				relative:     "RETRO-WAVE",
				extended:     NotExtended,
				visit:        true,
				subscription: nav.SubscribeAny,
				callback:     universalCallback("VISIT-CONTAINS-FOLDERS", NotExtended),
			}),
			Entry(nil, &naviTE{
				message:      "universal: Path contains folders (large)",
				relative:     "",
				extended:     NotExtended,
				subscription: nav.SubscribeAny,
				callback:     universalCallback("CONTAINS-FOLDERS (large)", NotExtended),
			}),
			Entry(nil, &naviTE{
				message:      "universal: Path contains folders (large, ensure single invoke)",
				relative:     "",
				extended:     NotExtended,
				once:         true,
				subscription: nav.SubscribeAny,
				callback:     universalCallback("CONTAINS-FOLDERS (large, ensure single invoke)", NotExtended),
			}),

			// === folders =======================================================

			Entry(nil, &naviTE{
				message:      "folders: Path is leaf",
				relative:     "RETRO-WAVE/Chromatics/Night Drive",
				extended:     NotExtended,
				subscription: nav.SubscribeFolders,
				callback:     foldersCallback("LEAF-PATH", NotExtended),
			}),
			Entry(nil, &naviTE{
				message:      "folders: Path contains folders",
				relative:     "RETRO-WAVE",
				extended:     IsExtended,
				subscription: nav.SubscribeFolders,
				callback:     foldersCallback("CONTAINS-FOLDERS ", IsExtended),
			}),
			Entry(nil, &naviTE{
				message:      "folders: Path contains folders (check all invoked)",
				relative:     "RETRO-WAVE",
				extended:     IsExtended,
				visit:        true,
				subscription: nav.SubscribeFolders,
				callback:     foldersCallback("CONTAINS-FOLDERS (check all invoked)", IsExtended),
			}),
			Entry(nil, &naviTE{
				message:      "folders: Path contains folders (large)",
				relative:     "",
				extended:     NotExtended,
				subscription: nav.SubscribeFolders,
				callback:     foldersCallback("CONTAINS-FOLDERS (large)", NotExtended),
			}),
			Entry(nil, &naviTE{
				message:      "folders: Path contains folders (large, ensure single invoke)",
				relative:     "",
				extended:     NotExtended,
				once:         true,
				subscription: nav.SubscribeFolders,
				callback:     foldersCallback("CONTAINS-FOLDERS (large, ensure single invoke)", NotExtended),
			}),

			Entry(nil, &naviTE{
				message:       "folders: case sensitive sort",
				relative:      "rock/metal",
				extended:      NotExtended,
				subscription:  nav.SubscribeFolders,
				caseSensitive: true,
				callback:      foldersCaseSensitiveCallback("rock/metal/HARD-METAL", "rock/metal/dark"),
			}),

			// === files =========================================================

			Entry(nil, &naviTE{
				message:      "files: Path is leaf",
				relative:     "RETRO-WAVE/Chromatics/Night Drive",
				extended:     NotExtended,
				subscription: nav.SubscribeFiles,
				callback:     filesCallback("LEAF-PATH", NotExtended),
			}),
			Entry(nil, &naviTE{
				message:      "files: Path contains folders",
				relative:     "RETRO-WAVE",
				extended:     NotExtended,
				subscription: nav.SubscribeFiles,
				callback:     filesCallback("CONTAINS-FOLDERS", NotExtended),
			}),
			Entry(nil, &naviTE{
				message:      "files: Path contains folders",
				relative:     "RETRO-WAVE",
				extended:     NotExtended,
				visit:        true,
				subscription: nav.SubscribeFiles,
				callback:     filesCallback("VISIT-CONTAINS-FOLDERS", NotExtended),
			}),
			Entry(nil, &naviTE{
				message:      "files: Path contains folders (large)",
				relative:     "",
				extended:     IsExtended,
				subscription: nav.SubscribeFiles,
				callback:     filesCallback("CONTAINS-FOLDERS (large)", IsExtended),
			}),
			Entry(nil, &naviTE{
				message:      "files: Path contains folders (large, ensure single invoke)",
				relative:     "",
				extended:     IsExtended,
				once:         true,
				subscription: nav.SubscribeFiles,
				callback:     filesCallback("CONTAINS-FOLDERS (large, ensure single invoke)", IsExtended),
			}),
		)
	})

	DescribeTable("Folders With Files",
		func(entry *naviTE) {
			recording := recordingMap{}
			visited := []string{}

			once := nav.LabelledTraverseCallback{
				Label: "test once callback",
				Fn: func(item *nav.TraverseItem) *LocalisableError {
					_, found := recording[item.Extension.Name]
					Expect(found).To(BeFalse())
					recording[item.Extension.Name] = len(item.Children)

					return entry.callback.Fn(item)
				},
			}

			path := path(root, entry.relative)
			session := nav.PrimarySession{
				Path: path,
			}
			_ = session.Configure(func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = begin("🛡️")
				o.Store.Subscription = entry.subscription
				o.Store.Behaviours.Sort.IsCaseSensitive = entry.caseSensitive
				o.Store.DoExtend = entry.extended
				o.Callback = once
			}).Run()

			if entry.visit {
				_ = filepath.WalkDir(path, func(path string, de fs.DirEntry, err error) error {
					if subscribes(entry.subscription, de) {
						visited = append(visited, path)
					}
					return nil
				})
			}

			if entry.visit {
				every := lo.EveryBy(visited, func(p string) bool {

					segments := strings.Split(p, string(filepath.Separator))
					name, err := lo.Last(segments)

					if err == nil {
						_, found := recording[name]
						return found
					}
					return false
				})
				Expect(every).To(BeTrue())
			}

			for n, actualNoChildren := range entry.expectedNoChildren {
				Expect(recording[n]).To(Equal(actualNoChildren), reason(n))
			}
		},
		func(entry *naviTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.message)
		},
		// === folders (with files) ==========================================

		Entry(nil, &naviTE{
			message:      "folders(with files): Path is leaf",
			relative:     "RETRO-WAVE/Chromatics/Night Drive",
			extended:     IsExtended,
			subscription: nav.SubscribeFoldersWithFiles,
			callback:     foldersCallback("LEAF-PATH", IsExtended),
			expectedNoChildren: map[string]int{
				"Night Drive": 4,
			},
		}),

		Entry(nil, &naviTE{
			message:      "folders(with files): Path contains folders (check all invoked)",
			relative:     "RETRO-WAVE",
			extended:     IsExtended,
			visit:        true,
			subscription: nav.SubscribeFoldersWithFiles,
			callback:     foldersCallback("CONTAINS-FOLDERS (check all invoked)", IsExtended),
			expectedNoChildren: map[string]int{
				"Night Drive":      4,
				"Northern Council": 4,
				"Teenage Color":    3,
				"Innerworld":       3,
			},
		}),
	)
})
