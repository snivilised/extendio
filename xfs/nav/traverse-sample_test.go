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

var _ = Describe("Traverse With Sample", Ordered, func() {
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

	DescribeTable("sample",
		func(entry *sampleTE) {
			path := helpers.Path(
				root,
				lo.Ternary(entry.naviTE.relative == "", "RETRO-WAVE", entry.naviTE.relative),
			)
			providedOptions := nav.GetDefaultOptions()
			providedOptions.Store.Subscription = entry.subscription
			providedOptions.Store.Sampling.SampleType = entry.sampleType
			providedOptions.Store.Sampling.SampleInReverse = entry.reverse
			providedOptions.Sampler.Custom.Each = entry.each
			providedOptions.Sampler.Custom.While = entry.while
			providedOptions.Callback = &nav.LabelledTraverseCallback{
				Label: "test universal callback",
				Fn: func(item *nav.TraverseItem) error {
					GinkgoWriter.Printf(
						"---> 🌊 SAMPLE-CALLBACK: '%v'\n", item.Path,
					)

					prohibited := fmt.Sprintf("%v, was invoked, but does not satisfy sample criteria",
						helpers.Reason(item.Extension.Name),
					)
					Expect(entry.prohibited).ToNot(ContainElement(item.Extension.Name), prohibited)

					return nil
				},
			}
			providedOptions.Notify.OnBegin = begin("🛡️")
			providedOptions.Store.Sampling.NoOf = entry.noOf

			if entry.filter != nil {
				filterDefs := &nav.FilterDefinitions{
					Node: nav.FilterDef{
						Type:        nav.FilterTypeGlobEn,
						Description: entry.filter.name,
						Pattern:     entry.filter.pattern,
						Scope:       entry.filter.scope,
					},
				}
				providedOptions.Store.FilterDefs = filterDefs
			}

			result, _ := nav.New().Primary(&nav.Prime{
				Path:            path,
				ProvidedOptions: providedOptions,
			}).Run()

			files := result.Metrics.Count(nav.MetricNoFilesInvokedEn)
			folders := result.Metrics.Count(nav.MetricNoFoldersInvokedEn)
			GinkgoWriter.Printf("===> files: '%v', folders: '%v'\n", files, folders)

			if entry.expectedNoOf.folders > 0 {
				Expect(folders).To(BeEquivalentTo(entry.expectedNoOf.folders),
					helpers.BecauseQuantity(entry.message, int(entry.expectedNoOf.folders), int(folders)),
				)
			}

			if entry.expectedNoOf.files > 0 {
				Expect(files).To(BeEquivalentTo(entry.expectedNoOf.files),
					helpers.BecauseQuantity(entry.message, int(entry.expectedNoOf.files), int(files)),
				)
			}
		},
		func(entry *sampleTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'", entry.message, entry.should)
		},

		// === universal =====================================================

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "universal(slice): first, with 2 files",
				should:       "invoke for at most 2 files per directory",
				subscription: nav.SubscribeAny,
				prohibited:   []string{"cover.night-drive.jpg"},
			},
			sampleType: nav.SampleTypeSliceEn,
			noOf: nav.EntryQuantities{
				Files: 2,
			},
			expectedNoOf: directoryQuantities{
				files: 8,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "universal(slice): first, with 2 folders",
				should:       "invoke for at most 2 folders per directory",
				subscription: nav.SubscribeAny,
				prohibited:   []string{"Electric Youth"},
			},
			sampleType: nav.SampleTypeSliceEn,
			noOf: nav.EntryQuantities{
				Folders: 2,
			},
			expectedNoOf: directoryQuantities{
				files:   11,
				folders: 6,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "universal(slice): first, with 2 files and 2 folders",
				should:       "invoke for at most 2 files and 2 folders per directory",
				subscription: nav.SubscribeAny,
				prohibited:   []string{"cover.night-drive.jpg", "Electric Youth"},
			},
			sampleType: nav.SampleTypeSliceEn,
			noOf: nav.EntryQuantities{
				Files:   2,
				Folders: 2,
			},
			expectedNoOf: directoryQuantities{
				files:   6,
				folders: 6,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "universal(filter): first, single file, first 2 folders",
				should:       "invoke for at most single file per directory",
				relative:     "edm",
				subscription: nav.SubscribeAny,
				prohibited:   []string{"02 - Swab.flac"},
			},
			filter: &filterTE{
				name:    "items with .flac suffix",
				pattern: "*.flac",
				scope:   nav.ScopeFileEn,
			},
			sampleType: nav.SampleTypeFilterEn,
			noOf: nav.EntryQuantities{
				Files:   1,
				Folders: 2,
			},
			expectedNoOf: directoryQuantities{
				files: 7,
				// folders: 7,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "universal(filter): last, last single files, last 2 folders",
				should:       "invoke for at most single file per directory",
				relative:     "edm",
				subscription: nav.SubscribeAny,
				prohibited:   []string{"01 - Dre.flac"},
			},
			filter: &filterTE{
				name:    "items with .flac suffix",
				pattern: "*.flac",
				scope:   nav.ScopeFileEn,
			},
			sampleType: nav.SampleTypeFilterEn,
			reverse:    true,
			noOf: nav.EntryQuantities{
				Files:   1,
				Folders: 2,
			},
			expectedNoOf: directoryQuantities{
				files: 8,
			},
		}),

		// === folders =======================================================

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "folders(slice): first, with 2 folders",
				should:       "invoke for at most 2 folders per directory",
				subscription: nav.SubscribeFolders,
				prohibited:   []string{"Electric Youth"},
			},
			sampleType: nav.SampleTypeSliceEn,
			noOf: nav.EntryQuantities{
				Folders: 2,
			},
			expectedNoOf: directoryQuantities{
				folders: 6,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "folders(slice): last, with last single folder",
				should:       "invoke for only last folder per directory",
				subscription: nav.SubscribeFolders,
				prohibited:   []string{"Chromatics"},
			},
			sampleType: nav.SampleTypeSliceEn,
			reverse:    true,
			noOf: nav.EntryQuantities{
				Folders: 1,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "filtered folders(filter): first, with 2 folders that start with A",
				should:       "invoke for at most 2 folders per directory",
				relative:     "edm",
				subscription: nav.SubscribeFolders,
				prohibited:   []string{"Tales Of Ephidrina"},
			},
			filter: &filterTE{
				name:    "items with that start with A",
				pattern: "A*",
				scope:   nav.ScopeFolderEn,
			},
			sampleType: nav.SampleTypeFilterEn,
			noOf: nav.EntryQuantities{
				Folders: 2,
			},
			expectedNoOf: directoryQuantities{
				// AMBIENT-TECHNO, Amorphous Androgynous, Aphex Twin
				folders: 3,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "filtered folders(filter): last, with single folder that start with A",
				should:       "invoke for at most a single folder per directory",
				relative:     "edm",
				subscription: nav.SubscribeFolders,
				prohibited:   []string{"Amorphous Androgynous"},
			},
			filter: &filterTE{
				name:    "items with that start with A",
				pattern: "A*",
				scope:   nav.ScopeAllEn,
			},
			sampleType: nav.SampleTypeFilterEn,
			reverse:    true,
			noOf: nav.EntryQuantities{
				Folders: 1,
			},
			expectedNoOf: directoryQuantities{
				folders: 2,
			},
		}),

		// === folders with files ============================================

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "folders with files(slice): first, with 2 folders",
				should:       "invoke for at most 2 folders per directory",
				subscription: nav.SubscribeFoldersWithFiles,
				prohibited:   []string{"Electric Youth"},
			},
			sampleType: nav.SampleTypeSliceEn,
			noOf: nav.EntryQuantities{
				Folders: 2,
			},
			expectedNoOf: directoryQuantities{
				folders: 6,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "folders with files(slice): last, with last single folder",
				should:       "invoke for only last folder per directory",
				subscription: nav.SubscribeFoldersWithFiles,
				prohibited:   []string{"Chromatics"},
			},
			sampleType: nav.SampleTypeSliceEn,
			reverse:    true,
			noOf: nav.EntryQuantities{
				Folders: 1,
			},
			expectedNoOf: directoryQuantities{
				folders: 3,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "filtered folders with files(filter): last, with single folder that start with A",
				should:       "invoke for at most a single folder per directory",
				relative:     "edm",
				subscription: nav.SubscribeFoldersWithFiles,
				prohibited:   []string{"Amorphous Androgynous"},
			},
			filter: &filterTE{
				name:    "items with that start with A",
				pattern: "A*",
				scope:   nav.ScopeAllEn,
			},
			sampleType: nav.SampleTypeFilterEn,
			reverse:    true,
			noOf: nav.EntryQuantities{
				Folders: 1,
			},
			expectedNoOf: directoryQuantities{
				folders: 2,
			},
		}),

		// === files =========================================================

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "files(slice): first, with 2 files",
				should:       "invoke for at most 2 files per directory",
				subscription: nav.SubscribeFiles,
				prohibited:   []string{"cover.night-drive.jpg"},
			},
			sampleType: nav.SampleTypeSliceEn,
			noOf: nav.EntryQuantities{
				Files: 2,
			},
			expectedNoOf: directoryQuantities{
				files: 8,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "files(slice): last, with last single file",
				should:       "invoke for only last file per directory",
				subscription: nav.SubscribeFiles,
				prohibited:   []string{"A1 - The Telephone Call.flac"},
			},
			sampleType: nav.SampleTypeSliceEn,
			reverse:    true,
			noOf: nav.EntryQuantities{
				Files: 1,
			},
			expectedNoOf: directoryQuantities{
				files: 4,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "filtered files(filter): first, 2 files",
				should:       "invoke for at most 2 files per directory",
				relative:     "edm/ELECTRONICA",
				subscription: nav.SubscribeFiles,
				prohibited:   []string{"03 - Mountain Goat.flac"},
			},
			filter: &filterTE{
				name:    "items with .flac suffix",
				pattern: "*.flac",
				scope:   nav.ScopeLeafEn,
			},
			sampleType: nav.SampleTypeFilterEn,
			noOf: nav.EntryQuantities{
				Files: 2,
			},
			expectedNoOf: directoryQuantities{
				files: 24,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "filtered files(filter): last, last 2 files",
				should:       "invoke for at most 2 files per directory",
				relative:     "edm",
				subscription: nav.SubscribeFiles,
				prohibited:   []string{"01 - Liquid Insects.flac"},
			},
			filter: &filterTE{
				name:    "items with .flac suffix",
				pattern: "*.flac",
				scope:   nav.ScopeAllEn,
			},
			sampleType: nav.SampleTypeFilterEn,
			reverse:    true,
			noOf: nav.EntryQuantities{
				Files: 2,
			},
			expectedNoOf: directoryQuantities{
				files: 42,
			},
		}),

		// === custom ========================================================

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "universal(custom): first, single file, 2 folders",
				should:       "invoke for at most single file per directory",
				relative:     "edm",
				subscription: nav.SubscribeAny,
				prohibited:   []string{"02 - Swab.flac"},
			},
			each: func(childItem *nav.TraverseItem) bool {
				if childItem.IsDirectory() {
					return true
				}

				return strings.HasPrefix(childItem.Extension.Name, "cover")
			},
			while: func(fi *nav.FilteredInfo) bool {
				fi.Enough.Files = fi.Counts.Files == 1
				fi.Enough.Folders = fi.Counts.Folders == 2

				return !fi.Enough.Files && !fi.Enough.Folders
			},
			sampleType: nav.SampleTypeCustomEn,
			noOf: nav.EntryQuantities{
				Files:   1,
				Folders: 2,
			},
			expectedNoOf: directoryQuantities{
				files:   7,
				folders: 14,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "filtered folders(custom): last, single folder that starts with A",
				should:       "invoke for at most a single folder per directory",
				relative:     "edm",
				subscription: nav.SubscribeFolders,
				prohibited:   []string{"Amorphous Androgynous"},
			},
			each: func(childItem *nav.TraverseItem) bool {
				return strings.HasPrefix(childItem.Extension.Name, "A")
			},
			while: func(fi *nav.FilteredInfo) bool {
				return fi.Counts.Folders < 1
			},
			sampleType: nav.SampleTypeCustomEn,
			reverse:    true,
			expectedNoOf: directoryQuantities{
				folders: 3,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "filtered files(custom): last, last 2 files",
				should:       "invoke for at most 2 files per directory",
				relative:     "edm",
				subscription: nav.SubscribeFiles,
				prohibited:   []string{"01 - Liquid Insects.flac"},
			},
			each: func(childItem *nav.TraverseItem) bool {
				return strings.HasSuffix(childItem.Extension.Name, ".flac")
			},
			while: func(fi *nav.FilteredInfo) bool {
				return fi.Counts.Files != 2
			},
			sampleType: nav.SampleTypeCustomEn,
			reverse:    true,
			expectedNoOf: directoryQuantities{
				files: 42,
			},
		}),
	)
})
