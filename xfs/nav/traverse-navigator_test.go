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
	"github.com/snivilised/extendio/translate"
	. "github.com/snivilised/extendio/translate"
	"github.com/snivilised/extendio/xfs/nav"
)

func normalise(p string) string {
	return strings.ReplaceAll(p, "/", string(filepath.Separator))
}

func xname(item *nav.TraverseItem) string {
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

type recordingMap map[string]bool
type recordingScopeMap map[string]nav.FilterScopeEnum

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

type listenTE struct {
	naviTE
	start      nav.Listener
	stop       nav.Listener
	incStart   bool
	incStop    bool
	mute       bool
	mandatory  []string
	prohibited []string
}

type filterTE struct {
	naviTE
	name          string
	pattern       string
	scope         nav.FilterScopeEnum
	negate        bool
	expectedErr   error
	errorContains string
}

type scopeTE struct {
	naviTE
	expectedScopes recordingScopeMap
}

func universalCallback(name string, extended bool) nav.TraverseCallback {

	ex := lo.Ternary(extended, "-EX", "")
	return func(item *nav.TraverseItem) *LocalisableError {
		GinkgoWriter.Printf("---> 🌊 %v-CALLBACK%v: '%v'\n", name, ex, item.Path)

		if extended {
			Expect(item.Extension).NotTo(BeNil(), fmt.Sprintf("❌ %v", item.Path))
		} else {
			Expect(item.Extension).To(BeNil(), fmt.Sprintf("❌ %v", item.Path))
		}
		return nil
	}
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

func foldersCallback(name string, extended bool) nav.TraverseCallback {

	ex := lo.Ternary(extended, "-EX", "")
	return func(item *nav.TraverseItem) *LocalisableError {
		GinkgoWriter.Printf("---> ☀️ FOLDERS:%v-CALLBACK%v: '%v'\n", name, ex, item.Path)
		Expect(item.Info.IsDir()).To(BeTrue())

		if extended {
			Expect(item.Extension).NotTo(BeNil(), fmt.Sprintf("❌ %v", item.Path))
		} else {
			Expect(item.Extension).To(BeNil(), fmt.Sprintf("❌ %v", item.Path))
		}
		return nil
	}
}

func universalScopeCallback(name string) nav.TraverseCallback {

	return func(item *nav.TraverseItem) *LocalisableError {
		Expect(item.Extension).NotTo(BeNil(), fmt.Sprintf("❌ %v", item.Extension.Name))
		GinkgoWriter.Printf("---> 🌠 UNIVERSAL:%v-CALLBACK-EX item-scope: (%v) '%v'\n",
			name, item.Extension.NodeScope, item.Extension.Name,
		)
		return nil
	}
}

func foldersScopeCallback(name string) nav.TraverseCallback {

	return func(item *nav.TraverseItem) *LocalisableError {
		Expect(item.Extension).NotTo(BeNil(), fmt.Sprintf("❌ %v", item.Extension.Name))
		GinkgoWriter.Printf("---> 🌟 FOLDERS:%v-CALLBACK-EX item-scope: (%v) '%v'\n",
			name, item.Extension.NodeScope, item.Extension.Name,
		)
		Expect(item.Info.IsDir()).To(BeTrue())
		return nil
	}
}

func filesScopeCallback(name string) nav.TraverseCallback {

	return func(item *nav.TraverseItem) *LocalisableError {
		Expect(item.Extension).NotTo(BeNil(), fmt.Sprintf("❌ %v", item.Extension.Name))
		GinkgoWriter.Printf("---> 🌬️ FILES:%v-CALLBACK-EX item-scope: (%v) '%v'\n",
			name, item.Extension.NodeScope, item.Extension.Name,
		)
		Expect(item.Info.IsDir()).To(BeFalse())
		return nil
	}
}

func filesCallback(name string, extended bool) nav.TraverseCallback {

	ex := lo.Ternary(extended, "-EX", "")
	return func(item *nav.TraverseItem) *LocalisableError {
		GinkgoWriter.Printf("---> 🌙 FILES:%v-CALLBACK%v: '%v'\n", name, ex, item.Path)
		Expect(item.Info.IsDir()).To(BeFalse())

		if extended {
			Expect(item.Extension).NotTo(BeNil(), fmt.Sprintf("❌ %v", item.Path))
		}
		return nil
	}
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

	Context("Navigator", func() {
		DescribeTable("Ensure Callback Invoked Once",
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
					o.Notify.OnBegin = begin("🛡️")
					o.Subscription = entry.subscription
					o.Behaviours.Sort.IsCaseSensitive = entry.caseSensitive
					o.DoExtend = entry.extended
					o.Callback = callback
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
				return fmt.Sprintf("🧪 ===> given: '%v'", entry.message)
			},

			// === universal =====================================================

			Entry(nil, &naviTE{
				message:      "universal: Path is leaf",
				relative:     "RETRO-WAVE/Chromatics/Night Drive",
				extended:     IsExtended,
				subscription: nav.SubscribeAny,
				callback:     universalCallback("(UNIVERSAL):LEAF-PATH", IsExtended),
			}),
			Entry(nil, &naviTE{
				message:      "universal: Path contains folders",
				relative:     "RETRO-WAVE",
				extended:     NotExtended,
				subscription: nav.SubscribeAny,
				callback:     universalCallback("(UNIVERSAL):CONTAINS-FOLDERS", NotExtended),
			}),
			Entry(nil, &naviTE{
				message:      "universal: Path contains folders",
				relative:     "RETRO-WAVE",
				extended:     NotExtended,
				visit:        true,
				subscription: nav.SubscribeAny,
				callback:     universalCallback("(UNIVERSAL):VISIT-CONTAINS-FOLDERS", NotExtended),
			}),
			Entry(nil, &naviTE{
				message:      "universal: Path contains folders (large)",
				relative:     "",
				extended:     NotExtended,
				subscription: nav.SubscribeAny,
				callback:     universalCallback("(UNIVERSAL):CONTAINS-FOLDERS (large)", NotExtended),
			}),
			Entry(nil, &naviTE{
				message:      "universal: Path contains folders (large, ensure single invoke)",
				relative:     "",
				extended:     NotExtended,
				once:         true,
				subscription: nav.SubscribeAny,
				callback:     universalCallback("(UNIVERSAL):CONTAINS-FOLDERS (large, ensure single invoke)", NotExtended),
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

		When("folder is skipped", func() {
			Context("folder navigator", func() {
				It("🧪 should: not invoke skipped folder descendants", func() {
					navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
						o.Subscription = nav.SubscribeFolders
						o.DoExtend = true
						o.Callback = skipFolderCallback("College", "Northern Council")
						o.Notify.OnBegin = begin("🛡️")
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
						o.Notify.OnBegin = begin("🛡️")
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
					o.Notify.OnBegin = begin("🛡️")
				})
				path := path(root, "RETRO-WAVE")
				navigator.Walk(path)
			},
			func(entry *skipTE) string {
				return fmt.Sprintf("🧪 ===> given: '%v'", entry.message)
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
						o.Notify.OnBegin = begin("🛡️")
						o.Subscription = nav.SubscribeAny
						o.Behaviours.SubPath.KeepTrailingSep = true
						o.DoExtend = true
						o.Callback = func(item *nav.TraverseItem) *LocalisableError {
							if expected, ok := expectations[item.Extension.Name]; ok {
								Expect(item.Extension.SubPath).To(Equal(expected), xname(item))
								GinkgoWriter.Printf("---> 🧩 SUB-PATH-CALLBACK(with): '%v', name: '%v', scope: '%v'\n",
									item.Extension.SubPath, item.Extension.Name, item.Extension.NodeScope,
								)
							}

							return nil
						}
					})
					path := path(root, "RETRO-WAVE")
					navigator.Walk(path)
				})

				When("using RootItemSubPath", func() {
					It("should: calculate subpath WITH trailing separator", func() {

						expectations := map[string]string{
							"edm":                         "",
							"_segments.def.infex.txt":     normalise("/_segments.def.infex.txt"),
							"Orbital 2 (The Brown Album)": normalise("/ELECTRONICA/Orbital/Orbital 2 (The Brown Album)"),
							"03 - Lush 3-1.flac":          normalise("/ELECTRONICA/Orbital/Orbital 2 (The Brown Album)/03 - Lush 3-1.flac"),
						}
						navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
							o.Notify.OnBegin = begin("🛡️")
							o.Subscription = nav.SubscribeAny
							o.Hooks.FolderSubPath = nav.RootItemSubPath
							o.Hooks.FileSubPath = nav.RootItemSubPath
							o.Behaviours.SubPath.KeepTrailingSep = true
							o.DoExtend = true
							o.Callback = func(item *nav.TraverseItem) *LocalisableError {
								if expected, ok := expectations[item.Extension.Name]; ok {
									Expect(item.Extension.SubPath).To(Equal(expected), xname(item))
									GinkgoWriter.Printf("---> 🧩🧩 SUB-PATH-CALLBACK(with): '%v', name: '%v', scope: '%v'\n",
										item.Extension.SubPath, item.Extension.Name, item.Extension.NodeScope,
									)
								}

								return nil
							}
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
						"Electric Youth":        "",
						"Innerworld":            normalise("/Electric Youth"),
						"A1 - Before Life.flac": normalise("/Electric Youth/Innerworld"),
					}
					navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
						o.Notify.OnBegin = begin("🛡️")
						o.Behaviours.SubPath.KeepTrailingSep = false
						o.Subscription = nav.SubscribeAny
						o.DoExtend = true
						o.Callback = func(item *nav.TraverseItem) *LocalisableError {
							if expected, ok := expectations[item.Extension.Name]; ok {
								Expect(item.Extension.SubPath).To(Equal(expected), xname(item))
								GinkgoWriter.Printf("---> 🧩 SUB-PATH-CALLBACK(without): '%v', name: '%v', scope: '%v'\n",
									item.Extension.SubPath, item.Extension.Name, item.Extension.NodeScope,
								)
							}

							return nil
						}
					})
					path := path(root, "RETRO-WAVE")
					navigator.Walk(path)
				})
			})
		})

		DescribeTable("listening",
			func(entry *listenTE) {
				navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🛡️")
					o.Subscription = entry.subscription
					o.Behaviours.Listen.InclusiveStart = entry.incStart
					o.Behaviours.Listen.InclusiveStop = entry.incStop
					o.Listen.Start = entry.start
					o.Listen.Stop = entry.stop
					if !entry.mute {
						o.Notify.OnStart = func(description string) {
							GinkgoWriter.Printf("===> 🎶 Start Listening: '%v'\n", description)
						}
						o.Notify.OnStop = func(description string) {
							GinkgoWriter.Printf("===> ⛔ Stop Listening: '%v'\n", description)
						}
					}
					o.DoExtend = entry.extended
					o.Callback = func(item *nav.TraverseItem) *LocalisableError {
						GinkgoWriter.Printf("---> 🔊 LISTENING-CALLBACK: name: '%v'\n",
							item.Extension.Name,
						)

						Expect(lo.Contains(entry.prohibited, item.Extension.Name)).To(BeFalse(), xname(item))
						Expect(lo.Contains(entry.mandatory, item.Extension.Name)).To(BeTrue(), xname(item))

						entry.mandatory = lo.Reject(entry.mandatory, func(s string, _ int) bool {
							return s == item.Extension.Name
						})
						return nil
					}
				})
				path := path(root, entry.relative)
				navigator.Walk(path)

				reason := fmt.Sprintf("❌ remaining: '%v'", strings.Join(entry.mandatory, ", "))
				Expect(len(entry.mandatory)).To(Equal(0), reason)
			},
			func(entry *listenTE) string {
				return fmt.Sprintf("🧪 ===> given: '%v'", entry.message)
			},

			Entry(nil, &listenTE{
				naviTE: naviTE{
					message:      "listening, start and stop (folders, inc:default)",
					relative:     "RETRO-WAVE",
					extended:     true,
					subscription: nav.SubscribeFolders,
				},
				start: &nav.ListenBy{
					Name: "Name: Night Drive",
					Fn: func(item *nav.TraverseItem) bool {
						return item.Extension.Name == "Night Drive"
					},
				},
				stop: &nav.ListenBy{
					Name: "Name: Electric Youth",
					Fn: func(item *nav.TraverseItem) bool {
						return item.Extension.Name == "Electric Youth"
					},
				},
				incStart:   true,
				incStop:    false,
				mandatory:  []string{"Night Drive", "College", "Northern Council", "Teenage Color"},
				prohibited: []string{"RETRO-WAVE", "Chromatics", "Electric Youth", "Innerworld"},
			}),
			Entry(nil, &listenTE{
				naviTE: naviTE{
					message:      "listening, start and stop (folders, excl:start, inc:stop, mute)",
					relative:     "RETRO-WAVE",
					extended:     true,
					subscription: nav.SubscribeFolders,
				},
				start: &nav.ListenBy{
					Name: "Name: Night Drive",
					Fn: func(item *nav.TraverseItem) bool {
						return item.Extension.Name == "Night Drive"
					},
				},
				stop: &nav.ListenBy{
					Name: "Name: Electric Youth",
					Fn: func(item *nav.TraverseItem) bool {
						return item.Extension.Name == "Electric Youth"
					},
				},
				incStart:  false,
				incStop:   true,
				mute:      true,
				mandatory: []string{"College", "Northern Council", "Teenage Color", "Electric Youth"},
				prohibited: []string{"Night Drive", "RETRO-WAVE", "Chromatics",
					"Innerworld",
				},
			}),
			Entry(nil, &listenTE{
				naviTE: naviTE{
					message:      "listening, start only (folders, inc:default)",
					relative:     "RETRO-WAVE",
					extended:     true,
					subscription: nav.SubscribeFolders,
				},
				start: &nav.ListenBy{
					Name: "Name: Night Drive",
					Fn: func(item *nav.TraverseItem) bool {
						return item.Extension.Name == "Night Drive"
					},
				},
				incStart: true,
				incStop:  false,
				mandatory: []string{"Night Drive", "College", "Northern Council", "Teenage Color",
					"Electric Youth", "Innerworld",
				},
				prohibited: []string{"RETRO-WAVE", "Chromatics"},
			}),
			Entry(nil, &listenTE{
				naviTE: naviTE{
					message:      "listening, stop only (folders, inc:default)",
					relative:     "RETRO-WAVE",
					extended:     true,
					subscription: nav.SubscribeFolders,
				},
				stop: &nav.ListenBy{
					Name: "Name: Electric Youth",
					Fn: func(item *nav.TraverseItem) bool {
						return item.Extension.Name == "Electric Youth"
					},
				},
				incStart: true,
				incStop:  false,
				mandatory: []string{"RETRO-WAVE", "Chromatics", "Night Drive", "College",
					"Northern Council", "Teenage Color",
				},
				prohibited: []string{"Electric Youth", "Innerworld"},
			}),
			Entry(nil, &listenTE{
				naviTE: naviTE{
					message:      "listening, stop only (folders, inc:default)",
					relative:     "RETRO-WAVE",
					extended:     true,
					subscription: nav.SubscribeFolders,
				},
				stop: &nav.ListenBy{
					Name: "Name: Night Drive",
					Fn: func(item *nav.TraverseItem) bool {
						return item.Extension.Name == "Night Drive"
					},
				},
				incStart:  true,
				incStop:   false,
				mandatory: []string{"RETRO-WAVE", "Chromatics"},
				prohibited: []string{"Night Drive", "College", "Northern Council",
					"Teenage Color", "Electric Youth", "Innerworld",
				},
			}),
		)

		Context("given: Early Exit", func() {
			It("should: exit early (folders)", func() {
				navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🛡️")
					o.Subscription = nav.SubscribeFolders
					o.Listen.Stop = &nav.ListenBy{
						Name: "Name: DREAM-POP",
						Fn: func(item *nav.TraverseItem) bool {
							return item.Extension.Name == "DREAM-POP"
						},
					}
					o.Notify.OnStop = func(description string) {
						GinkgoWriter.Printf("===> ⛔ Stop Listening: '%v'\n", description)
					}
					o.DoExtend = true
					o.Callback = foldersCallback("EARLY-EXIT-😴", o.DoExtend)
				})
				path := path(root, "")
				navigator.Walk(path)
			})

			It("should: exit early (files)", func() {
				navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🛡️")
					o.Subscription = nav.SubscribeFiles
					o.Listen.Stop = &nav.ListenBy{
						Name: "Name(contains): Captain",
						Fn: func(item *nav.TraverseItem) bool {
							return strings.Contains(item.Extension.Name, "Captain")
						},
					}
					o.Notify.OnStop = func(description string) {
						GinkgoWriter.Printf("===> ⛔ Stop Listening: '%v'\n", description)
					}
					o.DoExtend = true
					o.Callback = filesCallback("EARLY-EXIT-😴", o.DoExtend)
				})
				path := path(root, "")
				navigator.Walk(path)
			})
		})

		DescribeTable("RegexFilter",
			func(entry *filterTE) {
				navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🛡️")
					o.Subscription = entry.subscription
					o.Filter = &nav.RegexFilter{
						Filter: nav.Filter{
							Name:          entry.name,
							RequiredScope: entry.scope,
							Pattern:       entry.pattern,
							Negate:        entry.negate,
						},
					}
					o.DoExtend = true
					o.Callback = func(item *nav.TraverseItem) *translate.LocalisableError {
						GinkgoWriter.Printf("===> ⚗️ Regex Filter(%v) source: '%v', item-name: '%v', item-scope(fs): '%v(%v)'\n",
							o.Filter.Description(), o.Filter.Source(), item.Extension.Name, item.Extension.NodeScope, o.Filter.Scope(),
						)
						Expect(o.Filter.IsMatch(item)).To(BeTrue(), xname(item))
						return nil
					}
				})
				path := path(root, entry.relative)
				_ = navigator.Walk(path)
			},
			func(entry *filterTE) string {
				return fmt.Sprintf("🧪 ===> given: '%v'", entry.message)
			},

			// !!! -> if the folder doesn't pass the filter, should that folder be skipped?
			// make this a behaviour option

			Entry(nil, &filterTE{
				naviTE: naviTE{
					message:      "files(any scope): regex filter",
					relative:     "RETRO-WAVE",
					subscription: nav.SubscribeFiles,
				},
				name:    "items that start with 'vinyl'",
				pattern: "^vinyl",
				scope:   nav.AllScopesEn,
			}),
			Entry(nil, &filterTE{
				naviTE: naviTE{
					message:      "files(any scope): regex filter (negate)",
					relative:     "RETRO-WAVE",
					subscription: nav.SubscribeFiles,
				},
				name:    "items that don't start with 'vinyl'",
				pattern: "^vinyl",
				scope:   nav.AllScopesEn,
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
			Entry(nil, &filterTE{
				naviTE: naviTE{
					message:      "folders(any scope): regex filter",
					relative:     "RETRO-WAVE",
					subscription: nav.SubscribeFolders,
				},
				name:    "items that start with 'C'",
				pattern: "^C",
				scope:   nav.AllScopesEn,
			}),
			Entry(nil, &filterTE{
				naviTE: naviTE{
					message:      "folders(any scope): regex filter (negate)",
					relative:     "RETRO-WAVE",
					subscription: nav.SubscribeFolders,
				},
				name:    "items that don't start with 'C'",
				pattern: "^C",
				scope:   nav.AllScopesEn,
				negate:  true,
			}),

			// Entry(nil, &filterTE{
			// 	// THIS TEST NOT YET FINALISED. WHEN A SCOPE IS NOT APPLICABLE, NODE IS STILL VISITED
			// 	naviTE: naviTE{
			// 		message:      "folders(top): regex filter",
			// 		relative:     "PROGRESSIVE-HOUSE",
			// 		subscription: nav.SubscribeFolders,
			// 	},
			// 	name:    "top items that contain 'a'",
			// 	pattern: "a",
			// 	scope:   nav.TopScopeEn,
			// }),
		)

		Context("folders", func() {
			Context("given: filter and listen both active", func() {
				It("🧪 should: apply filter within the listen range", func() {
					navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
						o.Notify.OnBegin = begin("🛡️")
						o.Subscription = nav.SubscribeFolders
						o.Filter = &nav.RegexFilter{
							Filter: nav.Filter{
								Name:          "Contains 'o'",
								RequiredScope: nav.AllScopesEn,
								Pattern:       "(i?)o",
							},
						}
						o.Listen.Start = &nav.ListenBy{
							Name: "Name: Orbital",
							Fn: func(item *nav.TraverseItem) bool {
								return item.Extension.Name == "Orbital"
							},
						}
						o.Listen.Stop = &nav.ListenBy{
							Name: "Name: Underworld",
							Fn: func(item *nav.TraverseItem) bool {
								return item.Extension.Name == "Underworld"
							},
						}
						o.Notify.OnStart = func(description string) {
							GinkgoWriter.Printf("===> 🎶 Start Listening: '%v'\n", description)
						}
						o.Notify.OnStop = func(description string) {
							GinkgoWriter.Printf("===> ⛔ Stop Listening: '%v'\n", description)
						}
						o.DoExtend = true
						o.Callback = func(item *nav.TraverseItem) *translate.LocalisableError {
							GinkgoWriter.Printf("---> 🔊 LISTENING-CALLBACK: name: '%v'\n",
								item.Extension.Name,
							)
							GinkgoWriter.Printf("===> ⚗️ Regex Filter(%v) source: '%v', item-name: '%v', item-scope(fs): '%v(%v)'\n",
								o.Filter.Description(), o.Filter.Source(), item.Extension.Name, item.Extension.NodeScope, o.Filter.Scope(),
							)
							Expect(o.Filter.IsMatch(item)).To(BeTrue(), xname(item))
							return nil
						}
					})
					path := path(root, "edm/ELECTRONICA")
					_ = navigator.Walk(path)
				})
			})
		})

		DescribeTable("GlobFilter",
			func(entry *filterTE) {
				navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🛡️")
					o.Subscription = entry.subscription
					o.Filter = &nav.GlobFilter{
						Filter: nav.Filter{
							Name:          entry.name,
							RequiredScope: entry.scope,
							Pattern:       entry.pattern,
							Negate:        entry.negate,
						},
					}
					o.DoExtend = true
					o.Callback = func(item *nav.TraverseItem) *translate.LocalisableError {
						GinkgoWriter.Printf("===> 💠 Glob Filter(%v) source: '%v', item-name: '%v', item-scope(fs): '%v(%v)'\n",
							o.Filter.Description(), o.Filter.Source(), item.Extension.Name, item.Extension.NodeScope, o.Filter.Scope(),
						)
						Expect(o.Filter.IsMatch(item)).To(BeTrue(), xname(item))
						return nil
					}
				})
				path := path(root, entry.relative)
				_ = navigator.Walk(path)
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
				name:    "items with .flac suffix",
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
		)
	})
})
