package xfs_test

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	"github.com/snivilised/extendio/xfs"
)

func normalise(p string) string {
	return strings.ReplaceAll(p, "/", string(filepath.Separator))
}

func reason(item *xfs.TraverseItem) string {
	return fmt.Sprintf("❌ for item named: '%v'", item.Extension.Name)
}

func begin(em string) xfs.BeginHandler {
	return func(root string) {
		GinkgoWriter.Printf("---> %v [traverse-navigator-test:BEGIN], root: '%v'\n", em, root)
	}
}

type naviTE struct {
	message       string
	relative      string
	extended      bool
	once          bool
	visit         bool
	caseSensitive bool
	subscription  xfs.TraverseSubscription
	callback      xfs.TraverseCallback
}

type skipTE struct {
	naviTE
	skip    string
	exclude string
}

func universalCallback(item *xfs.TraverseItem) *xfs.LocalisableError {
	GinkgoWriter.Printf("---> 🌊 ON-NAVIGATOR-CALLBACK: '%v'\n", item.Path)
	Expect(item.Extension).To(BeNil(), fmt.Sprintf("❌ %v", item.Path))

	return nil
}

func universalCallbackEx(item *xfs.TraverseItem) *xfs.LocalisableError {
	GinkgoWriter.Printf("---> 🌊 ON-NAVIGATOR-CALLBACK-EX: '%v'\n", item.Path)
	Expect(item.Extension).NotTo(BeNil(), fmt.Sprintf("❌ %v", item.Path))

	return nil
}

func foldersCallback(item *xfs.TraverseItem) *xfs.LocalisableError {
	GinkgoWriter.Printf("---> ☀️ ON-NAVIGATOR-CALLBACK: '%v'\n", item.Path)
	Expect(item.Info.IsDir()).To(BeTrue())
	Expect(item.Extension).To(BeNil(), fmt.Sprintf("❌ %v", item.Path))

	return nil
}

func foldersCaseSensitiveCallback(first, second string) xfs.TraverseCallback {
	recording := recordingMap{}

	return func(item *xfs.TraverseItem) *xfs.LocalisableError {
		recording[item.Path] = true

		GinkgoWriter.Printf("---> ☀️ CASE-SENSITIVE-CALLBACK: '%v'\n", item.Path)
		Expect(item.Info.IsDir()).To(BeTrue())

		if strings.HasSuffix(item.Path, second) {
			GinkgoWriter.Printf("---> 💧 FIRST: '%v', 💧 SECOND: '%v'\n", first, second)

			paths := lo.Keys(recording)
			_, found := lo.Find(paths, func(s string) bool {
				return strings.HasSuffix(s, first)
			})

			Expect(found).To(BeTrue())
		}

		return nil
	}
}

func foldersCallbackEx(item *xfs.TraverseItem) *xfs.LocalisableError {
	GinkgoWriter.Printf("---> ☀️ ON-NAVIGATOR-CALLBACK-EX: '%v'\n", item.Path)
	Expect(item.Info.IsDir()).To(BeTrue())
	Expect(item.Extension).NotTo(BeNil(), fmt.Sprintf("❌ %v", item.Path))

	return nil
}

func filesCallback(item *xfs.TraverseItem) *xfs.LocalisableError {
	GinkgoWriter.Printf("---> 🌙 ON-NAVIGATOR-CALLBACK: '%v'\n", item.Path)
	Expect(item.Info.IsDir()).To(BeFalse())
	Expect(item.Extension).To(BeNil(), fmt.Sprintf("❌ %v", item.Path))

	return nil
}

func filesCallbackEx(item *xfs.TraverseItem) *xfs.LocalisableError {
	GinkgoWriter.Printf("---> 🌙 ON-NAVIGATOR-CALLBACK-EX: '%v'\n", item.Path)
	Expect(item.Info.IsDir()).To(BeFalse())
	Expect(item.Extension).NotTo(BeNil(), fmt.Sprintf("❌ %v", item.Path))
	return nil
}

func skipFolderCallback(skip, exclude string) xfs.TraverseCallback {

	return func(item *xfs.TraverseItem) *xfs.LocalisableError {
		GinkgoWriter.Printf("---> ♻️ ON-NAVIGATOR-SKIP-CALLBACK(skip:%v): '%v'\n", skip, item.Path)

		Expect(strings.HasSuffix(item.Path, exclude)).To(BeFalse())

		return lo.Ternary(strings.HasSuffix(item.Path, skip),
			&xfs.LocalisableError{Inner: fs.SkipDir}, nil,
		)
	}
}

func subscribes(subscription xfs.TraverseSubscription, de fs.DirEntry) bool {

	any := (subscription == xfs.SubscribeAny)
	files := (subscription == xfs.SubscribeFiles) && (!de.IsDir())
	folders := (subscription == xfs.SubscribeFolders) && (de.IsDir())

	return any || files || folders
}

type recordingMap map[string]bool

var _ = Describe("TraverseNavigator", Ordered, func() {
	var root string
	const IsExtended = true
	const NotExtended = false

	BeforeAll(func() {
		if current, err := os.Getwd(); err == nil {
			parent, _ := filepath.Split(current)
			root = filepath.Join(parent, "Test", "data", "MUSICO")
		}
	})

	Context("Path exists", func() {
		DescribeTable("Navigator",
			func(entry *naviTE) {
				recording := recordingMap{}
				visited := []string{}

				once := func(item *xfs.TraverseItem) *xfs.LocalisableError {
					_, found := recording[item.Path]
					Expect(found).To(BeFalse())
					recording[item.Path] = true

					return entry.callback(item)
				}

				visitor := func(item *xfs.TraverseItem) *xfs.LocalisableError {
					// just kept to enable visitor specific debug activity
					//
					return once(item)
				}
				callback := lo.Ternary(entry.once, once, lo.Ternary(entry.visit, visitor, entry.callback))

				path := path(root, entry.relative)
				navigator := xfs.NewNavigator(func(o *xfs.TraverseOptions) {
					o.Callback = callback
					o.Subscription = entry.subscription
					o.DoExtend = entry.extended
					o.IsCaseSensitive = entry.caseSensitive
					o.OnBegin = begin("🛡️")
				})

				if entry.visit {
					_ = filepath.WalkDir(path, func(path string, de fs.DirEntry, err error) error {
						if subscribes(entry.subscription, de) {
							visited = append(visited, path)
						}
						return nil
					})
				}

				_ = navigator.Walk(path)

				if entry.visit {
					every := lo.EveryBy(visited, func(p string) bool {
						_, found := recording[p]
						return found
					})
					Expect(every).To(BeTrue())
				}
			},
			func(entry *naviTE) string {
				return fmt.Sprintf("🧪 ===> '%v'", entry.message)
			},

			// === universal =====================================================

			Entry(nil, &naviTE{
				message:      "universal: Path is leaf",
				relative:     "RETRO-WAVE/Chromatics/Night Drive",
				extended:     IsExtended,
				subscription: xfs.SubscribeAny,
				callback:     universalCallbackEx,
			}),
			Entry(nil, &naviTE{
				message:      "universal: Path contains folders",
				relative:     "RETRO-WAVE",
				extended:     NotExtended,
				subscription: xfs.SubscribeAny,
				callback:     universalCallback,
			}),
			Entry(nil, &naviTE{
				message:      "universal: Path contains folders",
				relative:     "RETRO-WAVE",
				extended:     NotExtended,
				visit:        true,
				subscription: xfs.SubscribeAny,
				callback:     universalCallback,
			}),
			Entry(nil, &naviTE{
				message:      "universal: Path contains folders (large)",
				relative:     "",
				extended:     NotExtended,
				subscription: xfs.SubscribeAny,
				callback:     universalCallback,
			}),
			Entry(nil, &naviTE{
				message:      "universal: Path contains folders (large, ensure single invoke)",
				relative:     "",
				extended:     NotExtended,
				once:         true,
				subscription: xfs.SubscribeAny,
				callback:     universalCallback,
			}),

			// === folders =======================================================

			Entry(nil, &naviTE{
				message:      "folders: Path is leaf",
				relative:     "RETRO-WAVE/Chromatics/Night Drive",
				extended:     NotExtended,
				subscription: xfs.SubscribeFolders,
				callback:     foldersCallback,
			}),
			Entry(nil, &naviTE{
				message:      "folders: Path contains folders",
				relative:     "RETRO-WAVE",
				extended:     IsExtended,
				subscription: xfs.SubscribeFolders,
				callback:     foldersCallbackEx,
			}),
			Entry(nil, &naviTE{
				message:      "folders: Path contains folders (check all invoked)",
				relative:     "RETRO-WAVE",
				extended:     IsExtended,
				visit:        true,
				subscription: xfs.SubscribeFolders,
				callback:     foldersCallbackEx,
			}),
			Entry(nil, &naviTE{
				message:      "folders: Path contains folders (large)",
				relative:     "",
				extended:     NotExtended,
				subscription: xfs.SubscribeFolders,
				callback:     foldersCallback,
			}),
			Entry(nil, &naviTE{
				message:      "folders: Path contains folders (large, ensure single invoke)",
				relative:     "",
				extended:     NotExtended,
				once:         true,
				subscription: xfs.SubscribeFolders,
				callback:     foldersCallback,
			}),

			Entry(nil, &naviTE{
				message:       "folders: case sensitive sort",
				relative:      "rock/metal",
				extended:      NotExtended,
				subscription:  xfs.SubscribeFolders,
				caseSensitive: true,
				callback:      foldersCaseSensitiveCallback("rock/metal/HARD-METAL", "rock/metal/dark"),
			}),

			// === files =========================================================

			Entry(nil, &naviTE{
				message:      "files: Path is leaf",
				relative:     "RETRO-WAVE/Chromatics/Night Drive",
				extended:     NotExtended,
				subscription: xfs.SubscribeFiles,
				callback:     filesCallback,
			}),
			Entry(nil, &naviTE{
				message:      "files: Path contains folders",
				relative:     "RETRO-WAVE",
				extended:     NotExtended,
				subscription: xfs.SubscribeFiles,
				callback:     filesCallback,
			}),
			Entry(nil, &naviTE{
				message:      "files: Path contains folders",
				relative:     "RETRO-WAVE",
				extended:     NotExtended,
				visit:        true,
				subscription: xfs.SubscribeFiles,
				callback:     filesCallback,
			}),
			Entry(nil, &naviTE{
				message:      "files: Path contains folders (large)",
				relative:     "",
				extended:     IsExtended,
				subscription: xfs.SubscribeFiles,
				callback:     filesCallbackEx,
			}),
			Entry(nil, &naviTE{
				message:      "files: Path contains folders (large, ensure single invoke)",
				relative:     "",
				extended:     IsExtended,
				once:         true,
				subscription: xfs.SubscribeFiles,
				callback:     filesCallbackEx,
			}),
		)

		When("folder is skipped", func() {
			Context("folder navigator", func() {
				It("🧪 should: not invoke skipped folder descendants", func() {
					navigator := xfs.NewNavigator(func(o *xfs.TraverseOptions) {
						o.Subscription = xfs.SubscribeFolders
						o.DoExtend = true
						o.Callback = skipFolderCallback("College", "Northern Council")
						o.OnBegin = begin("🛡️")
					})
					path := path(root, "RETRO-WAVE")
					navigator.Walk(path)
				})
			})

			Context("universal navigator", func() {
				It("🧪 should: not invoke skipped folder descendants", func() {
					navigator := xfs.NewNavigator(func(o *xfs.TraverseOptions) {
						o.Subscription = xfs.SubscribeAny
						o.DoExtend = true
						o.Callback = skipFolderCallback("College", "Northern Council")
						o.OnBegin = begin("🛡️")
					})
					path := path(root, "RETRO-WAVE")
					navigator.Walk(path)
				})
			})
		})

		DescribeTable("TraverseNavigator",
			func(entry *skipTE) {
				navigator := xfs.NewNavigator(func(o *xfs.TraverseOptions) {
					o.Subscription = entry.subscription
					o.Callback = skipFolderCallback("College", "Northern Council")
					o.OnBegin = begin("🛡️")
				})
				path := path(root, "RETRO-WAVE")
				navigator.Walk(path)
			},
			func(entry *skipTE) string {
				return fmt.Sprintf("🧪 ===> '%v'", entry.message)
			},
			Entry(nil, &skipTE{
				naviTE: naviTE{
					message:      "universal: skip",
					subscription: xfs.SubscribeAny,
				},
				skip:    "College",
				exclude: "Northern Council",
			}),
			Entry(nil, &skipTE{
				naviTE: naviTE{
					message:      "folders: skip",
					subscription: xfs.SubscribeFolders,
				},
				skip:    "College",
				exclude: "Northern Council",
			}),
		)

		Context("sub-path", func() {
			When("KeepTrailingSep set to true", func() {
				It("should: calculate subpath WITH trailing separator", func() {

					expectations := map[string]string{
						"RETRO-WAVE":                   "",
						"Chromatics":                   normalise("/"),
						"Night Drive":                  normalise("/Chromatics/"),
						"A1 - The Telephone Call.flac": normalise("/Chromatics/Night Drive/"),
					}
					navigator := xfs.NewNavigator(func(o *xfs.TraverseOptions) {
						o.Subscription = xfs.SubscribeAny
						o.DoExtend = true
						o.Callback = func(item *xfs.TraverseItem) *xfs.LocalisableError {
							if expected, ok := expectations[item.Extension.Name]; ok {
								Expect(item.Extension.SubPath).To(Equal(expected), reason(item))
								GinkgoWriter.Printf("---> 🧩 SUB-PATH-CALLBACK(with): '%v', name: '%v', scope: '%v'\n",
									item.Extension.SubPath, item.Extension.Name, item.Extension.NodeScope,
								)
							}

							return nil
						}
						o.Behaviours.SubPath.KeepTrailingSep = true
						o.OnBegin = begin("🛡️")
					})
					path := path(root, "RETRO-WAVE")
					navigator.Walk(path)
				})

				When("using RootItemSubPath", func() {
					It("should: calculate subpath WITH trailing separator", func() {

						expectations := map[string]string{
							"edm":                         "",
							"_segments.def.infex.txt":     "/_segments.def.infex.txt",
							"Orbital 2 (The Brown Album)": normalise("/ELECTRONICA/Orbital/Orbital 2 (The Brown Album)"),
							"03 - Lush 3-1.flac":          normalise("/ELECTRONICA/Orbital/Orbital 2 (The Brown Album)/03 - Lush 3-1.flac"),
						}
						navigator := xfs.NewNavigator(func(o *xfs.TraverseOptions) {
							o.Subscription = xfs.SubscribeAny
							o.DoExtend = true
							o.Callback = func(item *xfs.TraverseItem) *xfs.LocalisableError {
								if expected, ok := expectations[item.Extension.Name]; ok {
									Expect(item.Extension.SubPath).To(Equal(expected), reason(item))
									GinkgoWriter.Printf("---> 🧩🧩 SUB-PATH-CALLBACK(with): '%v', name: '%v', scope: '%v'\n",
										item.Extension.SubPath, item.Extension.Name, item.Extension.NodeScope,
									)
								}

								return nil
							}
							o.Hooks.FolderSubPath = xfs.RootItemSubPath
							o.Hooks.FileSubPath = xfs.RootItemSubPath
							o.Behaviours.SubPath.KeepTrailingSep = true
							o.OnBegin = begin("🛡️")
						})
						path := path(root, "edm")
						navigator.Walk(path)
					})
				})
			})

			When("KeepTrailingSep set to false", func() {
				It("should: calculate subpath WITHOUT trailing separator", func() {
					expectations := map[string]string{
						"RETRO-WAVE":            "",
						"Electric Youth":        normalise(""),
						"Innerworld":            normalise("/Electric Youth"),
						"A1 - Before Life.flac": normalise("/Electric Youth/Innerworld"),
					}
					navigator := xfs.NewNavigator(func(o *xfs.TraverseOptions) {
						o.Subscription = xfs.SubscribeAny
						o.DoExtend = true
						o.Callback = func(item *xfs.TraverseItem) *xfs.LocalisableError {
							if expected, ok := expectations[item.Extension.Name]; ok {
								Expect(item.Extension.SubPath).To(Equal(expected), reason(item))
								GinkgoWriter.Printf("---> 🧩 SUB-PATH-CALLBACK(without): '%v', name: '%v', scope: '%v'\n",
									item.Extension.SubPath, item.Extension.Name, item.Extension.NodeScope,
								)
							}

							return nil
						}
						o.Behaviours.SubPath.KeepTrailingSep = false
						o.OnBegin = begin("🛡️")
					})
					path := path(root, "RETRO-WAVE")
					navigator.Walk(path)
				})
			})
		})
	})
})
