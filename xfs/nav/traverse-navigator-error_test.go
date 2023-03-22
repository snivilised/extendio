package nav_test

import (
	"errors"
	"fmt"
	"io/fs"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/snivilised/extendio/internal/helpers"

	. "github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("TraverseNavigator errors", Ordered, func() {
	var root string

	BeforeAll(func() {
		root = musico()
	})

	BeforeEach(func() {
		Use(func(o *UseOptions) {
			o.Tag = DefaultLanguage.Get()
		})
	})

	Context("new-navigator", func() {
		When("callback not set", func() {
			It("🧪 should: panic", func() {
				defer func() {
					_ = recover()
				}()

				session := &nav.PrimarySession{
					Path: root,
				}
				_ = session.Configure(func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🧲")
					o.Store.Subscription = nav.SubscribeAny
				})

				Fail("❌ expected panic due to missing callback")
				Expect(false)
			})
		})
	})

	Context("extend", func() {
		When("item is already extended", func() {
			It("🧪 should: panic", func() {
				defer func() {
					_ = recover()
				}()

				const relative = "RETRO-WAVE"
				path := helpers.Path(root, relative)
				session := &nav.PrimarySession{
					Path: path,
				}
				_, _ = session.Configure(func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🧲")
					o.Store.Subscription = nav.SubscribeFolders
					o.Hooks.Extend = func(navi *nav.NavigationInfo, entries *nav.DirectoryEntries) {
						navi.Item.Extension = &nav.ExtendedItem{
							Name: "dummy",
						}
						nav.DefaultExtendHookFn(navi, entries)
					}
					o.Store.DoExtend = true
					o.Callback = nav.LabelledTraverseCallback{
						Label: "test callback",
						Fn: func(item *nav.TraverseItem) error {
							return nil
						},
					}
				}).Run()

				Fail("❌ expected panic due to item already being extended")
			})
		})
	})

	Context("read error", func() {

		Context("navigator-folders", func() {
			It("🧪 should: invoke callback with error", func() {
				recording := []error{}

				const relative = "RETRO-WAVE"
				path := helpers.Path(root, relative)
				session := &nav.PrimarySession{
					Path: path,
				}
				_, _ = session.Configure(func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🧲")
					o.Store.Subscription = nav.SubscribeFolders
					o.Hooks.ReadDirectory = readDirFakeError
					o.Store.DoExtend = true
					o.Callback = nav.LabelledTraverseCallback{
						Label: "test callback",
						Fn: func(item *nav.TraverseItem) error {
							GinkgoWriter.Printf("---> 🔥 READ-ERR-CALLBACK: '%v', error: '%v'\n",
								item.Path, item.Error,
							)
							recording = append(recording, item.Error)
							return item.Error
						},
					}
				}).Run()

				Expect(len(recording)).To(Equal(2))
				Expect(recording[0]).To(BeNil())
				Expect(recording[1]).ToNot(BeNil())
			})
		})

		Context("navigator-files", func() {
			It("🧪 should: invoke callback with immediate read error", func() {
				const relative = "RETRO-WAVE"
				path := helpers.Path(root, relative)
				session := &nav.PrimarySession{
					Path: path,
				}
				_, _ = session.Configure(func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🧲")
					o.Store.Subscription = nav.SubscribeFiles
					o.Hooks.ReadDirectory = readDirFakeError
					o.Store.DoExtend = true
					o.Callback = errorCallback("(FILES):IMMEDIATE-READ-ERR", o.Store.DoExtend, false)
				}).Run()
			})

			It("🧪 should: invoke callback with error at ...", func() {
				const relative = "RETRO-WAVE"
				path := helpers.Path(root, relative)
				session := nav.PrimarySession{
					Path: path,
				}
				_, _ = session.Configure(func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🧲")
					o.Store.Subscription = nav.SubscribeFiles
					o.Hooks.ReadDirectory = readDirFakeErrorAt("Chromatics")
					o.Store.DoExtend = true
					o.Callback = errorCallback("(FILES):ERR-AT", o.Store.DoExtend, false)
				}).Run()
			})
		})
	})

	DescribeTable("given: sort generates an error",
		func(entry *errorTE) {
			defer func() {
				_ = recover()
			}()

			const relative = "RETRO-WAVE"
			path := helpers.Path(root, relative)
			session := &nav.PrimarySession{
				Path: path,
			}
			_, _ = session.Configure(func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = begin("🧲")
				o.Store.Subscription = entry.subscription
				o.Hooks.Sort = func(entries []fs.DirEntry, custom ...any) error {

					return errors.New("fake sort error")
				}
				o.Store.DoExtend = true
				o.Callback = errorCallback("SORT-ERR", o.Store.DoExtend, false)
			}).Run()

			Fail("❌ expected panic due to sort error")
		},
		func(entry *errorTE) string {
			return fmt.Sprintf("🧪 ===> ('%v') should panic", entry.message)
		},
		Entry(nil, &errorTE{naviTE{message: "universal", subscription: nav.SubscribeAny}}),
		Entry(nil, &errorTE{naviTE{message: "folders", subscription: nav.SubscribeFolders}}),
		Entry(nil, &errorTE{naviTE{message: "files", subscription: nav.SubscribeFiles}}),
	)

	Context("top level QueryStatus", func() {
		Context("given: error occurs", func() {
			It("🧪 should: halt traversal", func() {
				const relative = "RETRO-WAVE"
				path := helpers.Path(root, relative)
				session := nav.PrimarySession{
					Path: path,
				}
				_, _ = session.Configure(func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🧲")
					o.Store.Subscription = nav.SubscribeFolders
					o.Hooks.QueryStatus = func(path string) (fs.FileInfo, error) {
						return nil, errors.New("fake Lstat error")
					}
					o.Callback = errorCallback("ROOT-QUERY-STATUS", o.Store.DoExtend, true)
				}).Run()
			})
		})
	})

	Context("Extension error", func() {
		When("filter defined without setting DoExtend=true", func() {
			It("should: not panic", func() {
				filterDef := nav.FilterDef{
					Type:        nav.FilterTypeGlobEn,
					Description: "flac files",
					Pattern:     "*.flac",
					Scope:       nav.ScopeLeafEn,
				}
				const relative = "RETRO-WAVE"
				path := helpers.Path(root, relative)
				session := nav.PrimarySession{
					Path: path,
				}
				_, _ = session.Configure(func(o *nav.TraverseOptions) {
					o.Store.Subscription = nav.SubscribeAny
					o.Store.FilterDefs = &nav.FilterDefinitions{
						Node: filterDef,
					}
					o.Notify.OnBegin = begin("🧲")
					o.Callback = nav.LabelledTraverseCallback{
						Label: "test callback",
						Fn: func(item *nav.TraverseItem) error {
							GinkgoWriter.Printf("===> path:'%s'\n", item.Path)
							return nil
						},
					}
				}).Run()
			})
		})
	})
})
