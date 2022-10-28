package nav_test

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	. "github.com/snivilised/extendio/translate"
	"github.com/snivilised/extendio/xfs/nav"
)

func normalise(p string) string {
	return strings.ReplaceAll(p, "/", string(filepath.Separator))
}

func reason(item *nav.TraverseItem) string {
	return fmt.Sprintf("❌ for item named: '%v'", item.Extension.Name)
}

func begin(em string) nav.BeginHandler {
	return func(root string) {
		GinkgoWriter.Printf("---> %v [traverse-navigator-test:BEGIN], root: '%v'\n", em, root)
	}
}

func path(parent, relative string) string {
	segments := strings.Split(relative, "/")
	return filepath.Join(append([]string{parent}, segments...)...)
}

type naviTE struct {
	message       string
	relative      string
	extended      bool
	once          bool
	visit         bool
	caseSensitive bool
	subscription  nav.TraverseSubscription
	callback      nav.TraverseCallback
}

type skipTE struct {
	naviTE
	skip    string
	exclude string
}

func universalCallback(item *nav.TraverseItem) *LocalisableError {
	GinkgoWriter.Printf("---> 🌊 ON-NAVIGATOR-CALLBACK: '%v'\n", item.Path)
	Expect(item.Extension).To(BeNil(), fmt.Sprintf("❌ %v", item.Path))

	return nil
}

func universalCallbackEx(item *nav.TraverseItem) *LocalisableError {
	GinkgoWriter.Printf("---> 🌊 ON-NAVIGATOR-CALLBACK-EX: '%v'\n", item.Path)
	Expect(item.Extension).NotTo(BeNil(), fmt.Sprintf("❌ %v", item.Path))

	return nil
}

func foldersCallback(item *nav.TraverseItem) *LocalisableError {
	GinkgoWriter.Printf("---> ☀️ ON-NAVIGATOR-CALLBACK: '%v'\n", item.Path)
	Expect(item.Info.IsDir()).To(BeTrue())
	Expect(item.Extension).To(BeNil(), fmt.Sprintf("❌ %v", item.Path))

	return nil
}

func foldersCaseSensitiveCallback(first, second string) nav.TraverseCallback {
	recording := recordingMap{}

	return func(item *nav.TraverseItem) *LocalisableError {
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

func foldersCallbackEx(item *nav.TraverseItem) *LocalisableError {
	GinkgoWriter.Printf("---> ☀️ ON-NAVIGATOR-CALLBACK-EX: '%v'\n", item.Path)
	Expect(item.Info.IsDir()).To(BeTrue())
	Expect(item.Extension).NotTo(BeNil(), fmt.Sprintf("❌ %v", item.Path))

	return nil
}

func filesCallback(item *nav.TraverseItem) *LocalisableError {
	GinkgoWriter.Printf("---> 🌙 ON-NAVIGATOR-CALLBACK: '%v'\n", item.Path)
	Expect(item.Info.IsDir()).To(BeFalse())
	Expect(item.Extension).To(BeNil(), fmt.Sprintf("❌ %v", item.Path))

	return nil
}

func filesCallbackEx(item *nav.TraverseItem) *LocalisableError {
	GinkgoWriter.Printf("---> 🌙 ON-NAVIGATOR-CALLBACK-EX: '%v'\n", item.Path)
	Expect(item.Info.IsDir()).To(BeFalse())
	Expect(item.Extension).NotTo(BeNil(), fmt.Sprintf("❌ %v", item.Path))
	return nil
}

func skipFolderCallback(skip, exclude string) nav.TraverseCallback {

	return func(item *nav.TraverseItem) *LocalisableError {
		GinkgoWriter.Printf("---> ♻️ ON-NAVIGATOR-SKIP-CALLBACK(skip:%v): '%v'\n", skip, item.Path)

		Expect(strings.HasSuffix(item.Path, exclude)).To(BeFalse())

		return lo.Ternary(strings.HasSuffix(item.Path, skip),
			&LocalisableError{Inner: fs.SkipDir}, nil,
		)
	}
}

func subscribes(subscription nav.TraverseSubscription, de fs.DirEntry) bool {

	any := (subscription == nav.SubscribeAny)
	files := (subscription == nav.SubscribeFiles) && (!de.IsDir())
	folders := (subscription == nav.SubscribeFolders) && (de.IsDir())

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
			grand := filepath.Dir(parent)
			great := filepath.Dir(grand)
			root = filepath.Join(great, "Test", "data", "MUSICO")
		}
	})

	Context("Path exists", func() {
		DescribeTable("Navigator",
			func(entry *naviTE) {
				recording := recordingMap{}
				visited := []string{}

				once := func(item *nav.TraverseItem) *LocalisableError {
					_, found := recording[item.Path]
					Expect(found).To(BeFalse())
					recording[item.Path] = true

					return entry.callback(item)
				}

				visitor := func(item *nav.TraverseItem) *LocalisableError {
					// just kept to enable visitor specific debug activity
					//
					return once(item)
				}
				callback := lo.Ternary(entry.once, once, lo.Ternary(entry.visit, visitor, entry.callback))

				path := path(root, entry.relative)
				navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
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
				subscription: nav.SubscribeAny,
				callback:     universalCallbackEx,
			}),
			Entry(nil, &naviTE{
				message:      "universal: Path contains folders",
				relative:     "RETRO-WAVE",
				extended:     NotExtended,
				subscription: nav.SubscribeAny,
				callback:     universalCallback,
			}),
			Entry(nil, &naviTE{
				message:      "universal: Path contains folders",
				relative:     "RETRO-WAVE",
				extended:     NotExtended,
				visit:        true,
				subscription: nav.SubscribeAny,
				callback:     universalCallback,
			}),
			Entry(nil, &naviTE{
				message:      "universal: Path contains folders (large)",
				relative:     "",
				extended:     NotExtended,
				subscription: nav.SubscribeAny,
				callback:     universalCallback,
			}),
			Entry(nil, &naviTE{
				message:      "universal: Path contains folders (large, ensure single invoke)",
				relative:     "",
				extended:     NotExtended,
				once:         true,
				subscription: nav.SubscribeAny,
				callback:     universalCallback,
			}),

			// === folders =======================================================

			Entry(nil, &naviTE{
				message:      "folders: Path is leaf",
				relative:     "RETRO-WAVE/Chromatics/Night Drive",
				extended:     NotExtended,
				subscription: nav.SubscribeFolders,
				callback:     foldersCallback,
			}),
			Entry(nil, &naviTE{
				message:      "folders: Path contains folders",
				relative:     "RETRO-WAVE",
				extended:     IsExtended,
				subscription: nav.SubscribeFolders,
				callback:     foldersCallbackEx,
			}),
			Entry(nil, &naviTE{
				message:      "folders: Path contains folders (check all invoked)",
				relative:     "RETRO-WAVE",
				extended:     IsExtended,
				visit:        true,
				subscription: nav.SubscribeFolders,
				callback:     foldersCallbackEx,
			}),
			Entry(nil, &naviTE{
				message:      "folders: Path contains folders (large)",
				relative:     "",
				extended:     NotExtended,
				subscription: nav.SubscribeFolders,
				callback:     foldersCallback,
			}),
			Entry(nil, &naviTE{
				message:      "folders: Path contains folders (large, ensure single invoke)",
				relative:     "",
				extended:     NotExtended,
				once:         true,
				subscription: nav.SubscribeFolders,
				callback:     foldersCallback,
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
				callback:     filesCallback,
			}),
			Entry(nil, &naviTE{
				message:      "files: Path contains folders",
				relative:     "RETRO-WAVE",
				extended:     NotExtended,
				subscription: nav.SubscribeFiles,
				callback:     filesCallback,
			}),
			Entry(nil, &naviTE{
				message:      "files: Path contains folders",
				relative:     "RETRO-WAVE",
				extended:     NotExtended,
				visit:        true,
				subscription: nav.SubscribeFiles,
				callback:     filesCallback,
			}),
			Entry(nil, &naviTE{
				message:      "files: Path contains folders (large)",
				relative:     "",
				extended:     IsExtended,
				subscription: nav.SubscribeFiles,
				callback:     filesCallbackEx,
			}),
			Entry(nil, &naviTE{
				message:      "files: Path contains folders (large, ensure single invoke)",
				relative:     "",
				extended:     IsExtended,
				once:         true,
				subscription: nav.SubscribeFiles,
				callback:     filesCallbackEx,
			}),
		)

		When("folder is skipped", func() {
			Context("folder navigator", func() {
				It("🧪 should: not invoke skipped folder descendants", func() {
					navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
						o.Subscription = nav.SubscribeFolders
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
					navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
						o.Subscription = nav.SubscribeAny
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
				navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
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
					subscription: nav.SubscribeAny,
				},
				skip:    "College",
				exclude: "Northern Council",
			}),
			Entry(nil, &skipTE{
				naviTE: naviTE{
					message:      "folders: skip",
					subscription: nav.SubscribeFolders,
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
					navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
						o.Subscription = nav.SubscribeAny
						o.DoExtend = true
						o.Callback = func(item *nav.TraverseItem) *LocalisableError {
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
						navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
							o.Subscription = nav.SubscribeAny
							o.DoExtend = true
							o.Callback = func(item *nav.TraverseItem) *LocalisableError {
								if expected, ok := expectations[item.Extension.Name]; ok {
									Expect(item.Extension.SubPath).To(Equal(expected), reason(item))
									GinkgoWriter.Printf("---> 🧩🧩 SUB-PATH-CALLBACK(with): '%v', name: '%v', scope: '%v'\n",
										item.Extension.SubPath, item.Extension.Name, item.Extension.NodeScope,
									)
								}

								return nil
							}
							o.Hooks.FolderSubPath = nav.RootItemSubPath
							o.Hooks.FileSubPath = nav.RootItemSubPath
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
					navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
						o.Subscription = nav.SubscribeAny
						o.DoExtend = true
						o.Callback = func(item *nav.TraverseItem) *LocalisableError {
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
