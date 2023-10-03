package nav_test

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	. "github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/internal/helpers"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("TraverseNavigator(logged)", Ordered, func() {
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

	Context("Navigator", func() {
		DescribeTable("Ensure Callback Invoked Once",
			func(entry *naviTE) {
				recording := make(recordingMap)
				visited := []string{}

				once := nav.LabelledTraverseCallback{
					Label: "test once decorator",
					Fn: func(item *nav.TraverseItem) error {
						_, found := recording[item.Path]
						Expect(found).To(BeFalse())
						recording[item.Path] = len(item.Children)

						return entry.callback.Fn(item)
					},
				}

				visitor := nav.LabelledTraverseCallback{
					Fn: func(item *nav.TraverseItem) error {
						return once.Fn(item)
					},
				}
				callback := lo.Ternary(entry.once, once, lo.Ternary(entry.visit, visitor, entry.callback))

				path := helpers.Path(root, entry.relative)
				optionFn := func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🛡️")
					o.Store.Subscription = entry.subscription
					o.Store.Behaviours.Sort.IsCaseSensitive = entry.caseSensitive
					o.Store.DoExtend = entry.extended
					o.Callback = callback
				}

				result, _ := nav.New().Primary(&nav.Prime{
					Path:      path,
					OptionsFn: optionFn,
				}).Run()

				// _ = session.StartedAt()
				// _ = session.Elapsed()

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

				Expect(result.Metrics.Count(nav.MetricNoFilesInvokedEn)).To(Equal(entry.expectedNoOf.files),
					"Incorrect no of files")
				Expect(result.Metrics.Count(nav.MetricNoFoldersInvokedEn)).To(Equal(entry.expectedNoOf.folders),
					"Incorrect no of folders")
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
				expectedNoOf: expectedNo{
					files:   4,
					folders: 1,
				},
			}),
			Entry(nil, &naviTE{
				message:      "universal: Path contains folders",
				relative:     "RETRO-WAVE",
				extended:     NotExtended,
				subscription: nav.SubscribeAny,
				callback:     universalCallback("CONTAINS-FOLDERS", NotExtended),
				expectedNoOf: expectedNo{
					files:   14,
					folders: 8,
				},
			}),
			Entry(nil, &naviTE{
				message:      "universal: Path contains folders",
				relative:     "RETRO-WAVE",
				extended:     NotExtended,
				visit:        true,
				subscription: nav.SubscribeAny,
				callback:     universalCallback("VISIT-CONTAINS-FOLDERS", NotExtended),
				expectedNoOf: expectedNo{
					files:   14,
					folders: 8,
				},
			}),
			Entry(nil, &naviTE{
				message:      "universal: Path contains folders (large)",
				relative:     "",
				extended:     NotExtended,
				subscription: nav.SubscribeAny,
				callback:     universalCallback("CONTAINS-FOLDERS (large)", NotExtended),
				expectedNoOf: expectedNo{
					files:   656,
					folders: 178,
				},
			}),
			Entry(nil, &naviTE{
				message:      "universal: Path contains folders (large, ensure single invoke)",
				relative:     "",
				extended:     NotExtended,
				once:         true,
				subscription: nav.SubscribeAny,
				callback:     universalCallback("CONTAINS-FOLDERS (large, ensure single invoke)", NotExtended),
				expectedNoOf: expectedNo{
					files:   656,
					folders: 178,
				},
			}),

			// === folders =======================================================

			Entry(nil, &naviTE{
				message:      "folders: Path is leaf",
				relative:     "RETRO-WAVE/Chromatics/Night Drive",
				extended:     NotExtended,
				subscription: nav.SubscribeFolders,
				callback:     foldersCallback("LEAF-PATH", NotExtended),
				expectedNoOf: expectedNo{
					files:   0,
					folders: 1,
				},
			}),
			Entry(nil, &naviTE{
				message:      "folders: Path contains folders",
				relative:     "RETRO-WAVE",
				extended:     IsExtended,
				subscription: nav.SubscribeFolders,
				callback:     foldersCallback("CONTAINS-FOLDERS ", IsExtended),
				expectedNoOf: expectedNo{
					files:   0,
					folders: 8,
				},
			}),
			Entry(nil, &naviTE{
				message:      "folders: Path contains folders (check all invoked)",
				relative:     "RETRO-WAVE",
				extended:     IsExtended,
				visit:        true,
				subscription: nav.SubscribeFolders,
				callback:     foldersCallback("CONTAINS-FOLDERS (check all invoked)", IsExtended),
				expectedNoOf: expectedNo{
					files:   0,
					folders: 8,
				},
			}),
			Entry(nil, &naviTE{
				message:      "folders: Path contains folders (large)",
				relative:     "",
				extended:     NotExtended,
				subscription: nav.SubscribeFolders,
				callback:     foldersCallback("CONTAINS-FOLDERS (large)", NotExtended),
				expectedNoOf: expectedNo{
					files:   0,
					folders: 178,
				},
			}),
			Entry(nil, &naviTE{
				message:      "folders: Path contains folders (large, ensure single invoke)",
				relative:     "",
				extended:     NotExtended,
				once:         true,
				subscription: nav.SubscribeFolders,
				callback:     foldersCallback("CONTAINS-FOLDERS (large, ensure single invoke)", NotExtended),
				expectedNoOf: expectedNo{
					files:   0,
					folders: 178,
				},
			}),
			Entry(nil, &naviTE{
				message:       "folders: case sensitive sort",
				relative:      "rock/metal",
				extended:      NotExtended,
				subscription:  nav.SubscribeFolders,
				caseSensitive: true,
				callback:      foldersCaseSensitiveCallback("rock/metal/HARD-METAL", "rock/metal/dark"),
				expectedNoOf: expectedNo{
					files:   0,
					folders: 41,
				},
			}),

			// === files =========================================================

			Entry(nil, &naviTE{
				message:      "files: Path is leaf",
				relative:     "RETRO-WAVE/Chromatics/Night Drive",
				extended:     NotExtended,
				subscription: nav.SubscribeFiles,
				callback:     filesCallback("LEAF-PATH", NotExtended),
				expectedNoOf: expectedNo{
					files:   4,
					folders: 0,
				},
			}),
			Entry(nil, &naviTE{
				message:      "files: Path contains folders",
				relative:     "RETRO-WAVE",
				extended:     NotExtended,
				subscription: nav.SubscribeFiles,
				callback:     filesCallback("CONTAINS-FOLDERS", NotExtended),
				expectedNoOf: expectedNo{
					files:   14,
					folders: 0,
				},
			}),
			Entry(nil, &naviTE{
				message:      "files: Path contains folders",
				relative:     "RETRO-WAVE",
				extended:     NotExtended,
				visit:        true,
				subscription: nav.SubscribeFiles,
				callback:     filesCallback("VISIT-CONTAINS-FOLDERS", NotExtended),
				expectedNoOf: expectedNo{
					files:   14,
					folders: 0,
				},
			}),
			Entry(nil, &naviTE{
				message:      "files: Path contains folders (large)",
				relative:     "",
				extended:     IsExtended,
				subscription: nav.SubscribeFiles,
				callback:     filesCallback("CONTAINS-FOLDERS (large)", IsExtended),
				expectedNoOf: expectedNo{
					files:   656,
					folders: 0,
				},
			}),
			Entry(nil, &naviTE{
				message:      "files: Path contains folders (large, ensure single invoke)",
				relative:     "",
				extended:     IsExtended,
				once:         true,
				subscription: nav.SubscribeFiles,
				callback:     filesCallback("CONTAINS-FOLDERS (large, ensure single invoke)", IsExtended),
				expectedNoOf: expectedNo{
					files:   656,
					folders: 0,
				},
			}),
		)
	})

	DescribeTable("Folders With Files",
		func(entry *naviTE) {
			recording := make(recordingMap)
			visited := []string{}

			once := nav.LabelledTraverseCallback{
				Label: "test once callback",
				Fn: func(item *nav.TraverseItem) error {
					_, found := recording[item.Extension.Name]
					Expect(found).To(BeFalse())
					recording[item.Extension.Name] = len(item.Children)

					return entry.callback.Fn(item)
				},
			}

			path := helpers.Path(root, entry.relative)
			optionFn := func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = begin("🛡️")
				o.Store.Subscription = entry.subscription
				o.Store.Behaviours.Sort.IsCaseSensitive = entry.caseSensitive
				o.Store.DoExtend = entry.extended
				o.Callback = once
			}
			result, _ := nav.New().Primary(&nav.Prime{
				Path:      path,
				OptionsFn: optionFn,
			}).Run()

			if entry.visit {
				_ = filepath.WalkDir(path, func(path string, de fs.DirEntry, err error) error {
					if subscribes(entry.subscription, de) {
						visited = append(visited, path)
					}
					return nil
				})

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

			for n, actualNoChildren := range entry.expectedNoOf.children {
				Expect(recording[n]).To(Equal(actualNoChildren), helpers.Reason(
					fmt.Sprintf("folder: '%v'", n)),
				)
			}

			Expect(result.Metrics.Count(nav.MetricNoFilesInvokedEn)).To(Equal(entry.expectedNoOf.files),
				"Incorrect no of files")
			Expect(result.Metrics.Count(nav.MetricNoFoldersInvokedEn)).To(Equal(entry.expectedNoOf.folders),
				"Incorrect no of folders")

			sum := lo.Sum(lo.Values(entry.expectedNoOf.children))
			Expect(result.Metrics.Count(nav.MetricNoChildFilesFoundEn)).To(Equal(uint(sum)),
				helpers.Reason("Incorrect total no of child files"),
			)
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
			expectedNoOf: expectedNo{
				files:   0,
				folders: 1,
				children: map[string]int{
					"Night Drive": 4,
				},
			},
		}),

		Entry(nil, &naviTE{
			message:      "folders(with files): Path contains folders (check all invoked)",
			relative:     "RETRO-WAVE",
			extended:     IsExtended,
			visit:        true,
			subscription: nav.SubscribeFoldersWithFiles,
			expectedNoOf: expectedNo{
				files:   0,
				folders: 8,
				children: map[string]int{
					"Night Drive":      4,
					"Northern Council": 4,
					"Teenage Color":    3,
					"Innerworld":       3,
				},
			},
			callback: foldersCallback("CONTAINS-FOLDERS (check all invoked)", IsExtended),
		}),
	)
})
