package nav_test

import (
	"errors"
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

type errorTE struct {
	naviTE
}

func readDirFakeError(dirname string) ([]fs.DirEntry, error) {

	entries := []fs.DirEntry{}
	err := errors.New("fake read error")
	return entries, err
}

func readDirFakeErrorAt(name string) func(dirname string) ([]fs.DirEntry, error) {

	return func(dirname string) ([]fs.DirEntry, error) {
		if strings.HasSuffix(dirname, name) {
			return readDirFakeError(dirname)
		}

		return nav.ReadEntries(dirname)
	}
}

func errorCallback(name string, extended bool, hasError bool) nav.TraverseCallback {

	ex := lo.Ternary(extended, "-EX", "")
	return func(item *nav.TraverseItem) *LocalisableError {
		GinkgoWriter.Printf("---> 🔥 %v-CALLBACK%v: '%v'\n", name, ex, item.Path)

		if extended {
			Expect(item.Extension).NotTo(BeNil(), fmt.Sprintf("❌ %v", item.Path))
		} else {
			Expect(item.Extension).To(BeNil(), fmt.Sprintf("❌ %v", item.Path))
		}
		if hasError {
			Expect(item.Error).ToNot(BeNil())
		}
		return item.Error
	}
}

var _ = Describe("TraverseNavigator errors", Ordered, func() {
	var root string

	BeforeAll(func() {
		if current, err := os.Getwd(); err == nil {
			parent, _ := filepath.Split(current)
			grand := filepath.Dir(parent)
			great := filepath.Dir(grand)
			root = filepath.Join(great, "Test", "data", "MUSICO")
		}
	})

	Context("new-navigator", func() {
		When("callback not set", func() {
			It("🧪 should: panic", func() {
				defer func() {
					_ = recover()
				}()
				_ = root

				nav.NewNavigator(func(o *nav.TraverseOptions) {
					o.Subscription = nav.SubscribeAny
					o.Notify.OnBegin = begin("🧲")
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

				navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🧲")
					o.Subscription = nav.SubscribeFolders
					o.Hooks.Extend = func(navi *nav.NavigationParams, descendants []fs.DirEntry) {
						navi.Item.Extension = &nav.ExtendedItem{
							Name: "dummy",
						}
						nav.DefaultExtendHookFn(navi, descendants)
					}
					o.DoExtend = true
					o.Callback = func(item *nav.TraverseItem) *LocalisableError {
						return nil
					}
				})
				const relative = "RETRO-WAVE"
				path := path(root, relative)
				_ = navigator.Walk(path)

				Fail("❌ expected panic due to item already being extended")
			})
		})
	})

	Context("read error", func() {

		Context("navigator-folders", func() {
			It("🧪 should: invoke callback with error", func() {
				recording := []error{}
				navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🧲")
					o.Subscription = nav.SubscribeFolders
					o.Hooks.ReadDirectory = readDirFakeError
					o.DoExtend = true
					o.Callback = func(item *nav.TraverseItem) *LocalisableError {
						GinkgoWriter.Printf("---> 🔥 READ-ERR-CALLBACK: '%v', error: '%v'\n",
							item.Path, item.Error,
						)
						recording = append(recording, item.Error)
						return item.Error
					}
				})
				const relative = "RETRO-WAVE"
				path := path(root, relative)
				_ = navigator.Walk(path)

				Expect(len(recording)).To(Equal(2))
				Expect(recording[0]).To(BeNil())
				Expect(recording[1]).ToNot(BeNil())
			})
		})

		Context("navigator-files", func() {
			It("🧪 should: invoke callback with immediate read error", func() {
				navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🧲")
					o.Subscription = nav.SubscribeFiles
					o.Hooks.ReadDirectory = readDirFakeError
					o.DoExtend = true
					o.Callback = errorCallback("(FILES):IMMEDIATE-READ-ERR", o.DoExtend, false)
				})
				const relative = "RETRO-WAVE"
				path := path(root, relative)
				_ = navigator.Walk(path)
			})

			It("🧪 should: invoke callback with error at ...", func() {
				navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🧲")
					o.Subscription = nav.SubscribeFiles
					o.Hooks.ReadDirectory = readDirFakeErrorAt("Chromatics")
					o.DoExtend = true
					o.Callback = errorCallback("(FILES):ERR-AT", o.DoExtend, false)
				})
				const relative = "RETRO-WAVE"
				path := path(root, relative)
				_ = navigator.Walk(path)
			})
		})
	})

	DescribeTable("given: sort generates an error",
		func(entry *errorTE) {
			defer func() {
				_ = recover()
			}()

			navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = begin("🧲")
				o.Subscription = entry.subscription
				o.Hooks.Sort = func(entries []fs.DirEntry, custom ...any) error {

					return errors.New("fake sort error")
				}
				o.DoExtend = true
				o.Callback = errorCallback("SORT-ERR", o.DoExtend, false)
			})
			const relative = "RETRO-WAVE"
			path := path(root, relative)
			_ = navigator.Walk(path)

			Fail("❌ expected panic due to sort error")
		},
		func(entry *errorTE) string {
			return fmt.Sprintf("🧪 ===> ('%v') should panic", entry.message)
		},
		Entry(nil, &errorTE{naviTE{message: "universal", subscription: nav.SubscribeAny}}),
		Entry(nil, &errorTE{naviTE{message: "folders", subscription: nav.SubscribeFolders}}),
		Entry(nil, &errorTE{naviTE{message: "files", subscription: nav.SubscribeFiles}}),
	)

	DescribeTable("given: root is not a folder",
		func(entry *errorTE) {
			navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = begin("🧲")
				o.Subscription = entry.subscription
				o.DoExtend = true
				o.Callback = errorCallback("ROOT-NOT-FOLDER-ERR", o.DoExtend, true)
			})
			const relative = "RETRO-WAVE/Electric Youth/Innerworld/A2 - Runaway.flac"
			path := path(root, relative)
			_ = navigator.Walk(path)
		},

		func(entry *errorTE) string {
			return fmt.Sprintf("🧪 ===> ('%v') should return error", entry.message)
		},
		Entry(nil, &errorTE{naviTE{message: "universal", subscription: nav.SubscribeAny}}),
		Entry(nil, &errorTE{naviTE{message: "folders", subscription: nav.SubscribeFolders}}),
		Entry(nil, &errorTE{naviTE{message: "files", subscription: nav.SubscribeFiles}}),
	)

	Context("top level QueryStatus", func() {
		Context("given: error occurs", func() {
			It("🧪 should: halt traversal", func() {
				navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🧲")
					o.Subscription = nav.SubscribeFolders
					o.Hooks.QueryStatus = func(path string) (fs.FileInfo, error) {
						return nil, errors.New("fake Lstat error")
					}
					o.Callback = errorCallback("ROOT-QUERY-STATUS", o.DoExtend, true)
				})
				const relative = "RETRO-WAVE"
				path := path(root, relative)
				_ = navigator.Walk(path)
			})
		})
	})
})
