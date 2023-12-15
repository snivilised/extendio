package nav_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"

	. "github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/internal/helpers"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("FilterPoly", Ordered, func() {
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

	Context("FilterScopeBiEnum", func() {
		Context("Clear", func() {
			It("should: clear bit", func() {
				scope := nav.ScopeFolderEn | nav.ScopeLeafEn
				scope.Clear(nav.ScopeFolderEn)
				Expect(scope).To(Equal(nav.ScopeLeafEn))
			})
		})

		Context("Set", func() {
			It("should: set bit", func() {
				scope := nav.ScopeLeafEn
				scope.Set(nav.ScopeFolderEn)
				Expect(scope).To(Equal(nav.ScopeFolderEn | nav.ScopeLeafEn))
			})
		})
	})

	DescribeTable("PolyFilter",
		func(entry *polyTE) {
			recording := make(recordingMap)
			filterDefs := &nav.FilterDefinitions{
				Node: nav.FilterDef{
					Type: nav.FilterTypePolyEn,
					Poly: &nav.PolyFilterDef{
						File:   entry.file,
						Folder: entry.folder,
					},
				},
			}
			var filter nav.TraverseFilter

			path := helpers.Path(root, entry.relative)
			optionFn := func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = func(state *nav.NavigationState) {
					GinkgoWriter.Printf(
						"---> 🛡️ [traverse-navigator-test:BEGIN], root: '%v'\n", state.Root,
					)
					filter = state.Filters.Node
				}

				o.Store.Subscription = entry.subscription
				o.Store.FilterDefs = filterDefs
				o.Callback = &nav.LabelledTraverseCallback{
					Label: "test poly filter callback",
					Fn: func(item *nav.TraverseItem) error {
						GinkgoWriter.Printf(
							"===> ⚗️ Poly Filter(%v) source: '%v', item-name: '%v', item-scope(fs): '%v(%v)'\n",
							filter.Description(),
							filter.Source(),
							item.Extension.Name,
							item.Extension.NodeScope,
							filter.Scope(),
						)
						if lo.Contains(entry.mandatory, item.Extension.Name) {
							Expect(item).Should(MatchCurrentRegexFilter(filter))
						}

						recording[item.Extension.Name] = len(item.Children)
						return nil
					},
				}
			}
			result, err := nav.New().Primary(&nav.Prime{
				Path:      path,
				OptionsFn: optionFn,
			}).Run()

			if entry.mandatory != nil {
				for _, name := range entry.mandatory {
					_, found := recording[name]
					Expect(found).To(BeTrue(), helpers.Reason(name))
				}
			}

			if entry.prohibited != nil {
				for _, name := range entry.prohibited {
					_, found := recording[name]
					Expect(found).To(BeFalse(), helpers.Reason(name))
				}
			}

			Expect(err).Error().To(BeNil())
			Expect(result.Metrics.Count(nav.MetricNoFilesInvokedEn)).To(Equal(entry.expectedNoOf.files),
				helpers.BecauseQuantity("Incorrect no of files",
					int(entry.expectedNoOf.files),
					int(result.Metrics.Count(nav.MetricNoFilesInvokedEn)),
				),
			)

			Expect(result.Metrics.Count(nav.MetricNoFoldersInvokedEn)).To(Equal(entry.expectedNoOf.folders),
				helpers.BecauseQuantity("Incorrect no of folders",
					int(entry.expectedNoOf.folders),
					int(result.Metrics.Count(nav.MetricNoFoldersInvokedEn)),
				),
			)
		},
		func(entry *polyTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.message)
		},

		Entry(nil, &polyTE{
			naviTE: naviTE{
				message:      "poly - files:regex; folders:glob",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
				expectedNoOf: directoryQuantities{
					// file is 2 not 3 because *i* is case sensitive so Innerworld is not a match
					// The next(not this one) regex test case, fixes this because folder regex has better
					// control over case sensitivity
					files:   2,
					folders: 8,
				},
			},
			file: nav.FilterDef{
				Type:        nav.FilterTypeRegexEn,
				Description: "files: starts with vinyl",
				Pattern:     "^vinyl",
				Scope:       nav.ScopeFileEn,
			},
			folder: nav.FilterDef{
				Type:        nav.FilterTypeGlobEn,
				Description: "folders: contains i (case sensitive)",
				Pattern:     "*i*",
				Scope:       nav.ScopeFolderEn | nav.ScopeLeafEn,
			},
		}),

		Entry(nil, &polyTE{
			naviTE: naviTE{
				message:      "poly - files:regex; folders:regex",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
				expectedNoOf: directoryQuantities{
					files:   3,
					folders: 8,
				},
			},
			file: nav.FilterDef{
				Type:        nav.FilterTypeRegexEn,
				Description: "files: starts with vinyl",
				Pattern:     "^vinyl",
				Scope:       nav.ScopeFileEn,
			},
			folder: nav.FilterDef{
				Type:        nav.FilterTypeRegexEn,
				Description: "folders: contains i (case insensitive)",
				Pattern:     "[iI]",
				Scope:       nav.ScopeFolderEn | nav.ScopeLeafEn,
			},
		}),

		// For the poly filter, the file/folder scopes must be set correctly, but because
		// they can be set automatically, the client is not forced to set them. This test
		// checks that when the file/folder scopes are not set, then poly filtering still works
		// properly.
		Entry(nil, &polyTE{
			naviTE: naviTE{
				message:      "poly(scopes omitted) - files:regex; folders:regex",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
				expectedNoOf: directoryQuantities{
					files:   3,
					folders: 8,
				},
			},
			file: nav.FilterDef{
				Type:        nav.FilterTypeRegexEn,
				Description: "files: starts with vinyl",
				Pattern:     "^vinyl",
				// file scope omitted
			},
			folder: nav.FilterDef{
				Type:        nav.FilterTypeRegexEn,
				Description: "folders: contains i (case insensitive)",
				Pattern:     "[iI]",
				Scope:       nav.ScopeLeafEn, // folder scope omitted
			},
		}),

		Entry(nil, &polyTE{
			naviTE: naviTE{
				message:      "poly(subscribe:files)",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFiles,
				expectedNoOf: directoryQuantities{
					files:   3,
					folders: 0,
				},
			},
			file: nav.FilterDef{
				Type:        nav.FilterTypeRegexEn,
				Description: "files: starts with vinyl",
				Pattern:     "^vinyl",
			},
			folder: nav.FilterDef{
				Type:        nav.FilterTypeRegexEn,
				Description: "folders: contains i",
				Pattern:     "[iI]",
				Scope:       nav.ScopeLeafEn,
			},
		}),
	)

	DescribeTable("Panic: PolyFilter",
		func(entry *polyTE) {
			defer func() {
				pe := recover()
				if err, ok := pe.(error); !ok || !strings.Contains(err.Error(),
					"invalid subscription type for poly filter") {
					Fail("incorrect panic")
				}
			}()

			fileDef := nav.FilterDef{
				Type:        nav.FilterTypeRegexEn,
				Description: "files: starts with vinyl",
				Pattern:     "^vinyl",
			}
			folderDef := nav.FilterDef{
				Type:        nav.FilterTypeRegexEn,
				Description: "folders: contains i",
				Pattern:     "[iI]",
				Scope:       nav.ScopeLeafEn,
			}

			filterDefs := &nav.FilterDefinitions{
				Node: nav.FilterDef{
					Type: nav.FilterTypePolyEn,
					Poly: &nav.PolyFilterDef{
						File:   fileDef,
						Folder: folderDef,
					},
				},
			}

			path := helpers.Path(root, entry.relative)
			optionFn := func(o *nav.TraverseOptions) {
				o.Store.Subscription = entry.subscription
				o.Store.FilterDefs = filterDefs
				o.Callback = &nav.LabelledTraverseCallback{
					Label: "(panic): test poly filter callback",
					Fn: func(item *nav.TraverseItem) error {
						return nil
					},
				}
			}
			_, _ = nav.New().Primary(&nav.Prime{
				Path:      path,
				OptionsFn: optionFn,
			}).Run()

			Fail("❌ expected panic due to invalid subscription type for poly filter")
		},
		func(entry *polyTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.message)
		},

		// Poly filtering is not valid when files are not subscribed to
		// a warning is issued in the logs to indicate that this
		// scenario is invalid
		Entry(nil, &polyTE{
			naviTE: naviTE{
				message:      "poly(subscribe:folders/invalid)",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFolders,
			},
		}),

		Entry(nil, &polyTE{
			naviTE: naviTE{
				message:      "poly(subscribe:folders/invalid)",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFoldersWithFiles,
			},
		}),
	)
})
