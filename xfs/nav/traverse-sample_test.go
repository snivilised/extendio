package nav_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/internal/helpers"

	"github.com/snivilised/extendio/xfs/nav"
)

/*
  ---> 🛡️ [traverse-navigator-test:BEGIN], root: '&{false /Users/plastikfan/dev/github/snivilised/extendio/Test/data/MUSICO/RETRO-WAVE}'
  ---> 🌊 UNIVERSAL//CONTAINS-FOLDERS-CALLBACK: (depth:9999) '/Users/plastikfan/dev/github/snivilised/extendio/Test/data/MUSICO/RETRO-WAVE'
  ---> 🌊 UNIVERSAL//CONTAINS-FOLDERS-CALLBACK: (depth:9999) '/Users/plastikfan/dev/github/snivilised/extendio/Test/data/MUSICO/RETRO-WAVE/Chromatics'
  ---> 🌊 UNIVERSAL//CONTAINS-FOLDERS-CALLBACK: (depth:9999) '/Users/plastikfan/dev/github/snivilised/extendio/Test/data/MUSICO/RETRO-WAVE/Chromatics/Night Drive'
  ---> 🌊 UNIVERSAL//CONTAINS-FOLDERS-CALLBACK: (depth:9999) '/Users/plastikfan/dev/github/snivilised/extendio/Test/data/MUSICO/RETRO-WAVE/Chromatics/Night Drive/A1 - The Telephone Call.flac'
  ---> 🌊 UNIVERSAL//CONTAINS-FOLDERS-CALLBACK: (depth:9999) '/Users/plastikfan/dev/github/snivilised/extendio/Test/data/MUSICO/RETRO-WAVE/Chromatics/Night Drive/A2 - Night Drive.flac'
  ---> 🌊 UNIVERSAL//CONTAINS-FOLDERS-CALLBACK: (depth:9999) '/Users/plastikfan/dev/github/snivilised/extendio/Test/data/MUSICO/RETRO-WAVE/Chromatics/Night Drive/cover.night-drive.jpg'
  ---> 🌊 UNIVERSAL//CONTAINS-FOLDERS-CALLBACK: (depth:9999) '/Users/plastikfan/dev/github/snivilised/extendio/Test/data/MUSICO/RETRO-WAVE/Chromatics/Night Drive/vinyl-info.night-drive.txt'
  ---> 🌊 UNIVERSAL//CONTAINS-FOLDERS-CALLBACK: (depth:9999) '/Users/plastikfan/dev/github/snivilised/extendio/Test/data/MUSICO/RETRO-WAVE/College'
  ---> 🌊 UNIVERSAL//CONTAINS-FOLDERS-CALLBACK: (depth:9999) '/Users/plastikfan/dev/github/snivilised/extendio/Test/data/MUSICO/RETRO-WAVE/College/Northern Council'
  ---> 🌊 UNIVERSAL//CONTAINS-FOLDERS-CALLBACK: (depth:9999) '/Users/plastikfan/dev/github/snivilised/extendio/Test/data/MUSICO/RETRO-WAVE/College/Northern Council/A1 - Incident.flac'
  ---> 🌊 UNIVERSAL//CONTAINS-FOLDERS-CALLBACK: (depth:9999) '/Users/plastikfan/dev/github/snivilised/extendio/Test/data/MUSICO/RETRO-WAVE/College/Northern Council/A2 - The Zemlya Expedition.flac'
  ---> 🌊 UNIVERSAL//CONTAINS-FOLDERS-CALLBACK: (depth:9999) '/Users/plastikfan/dev/github/snivilised/extendio/Test/data/MUSICO/RETRO-WAVE/College/Northern Council/cover.northern-council.jpg'
  ---> 🌊 UNIVERSAL//CONTAINS-FOLDERS-CALLBACK: (depth:9999) '/Users/plastikfan/dev/github/snivilised/extendio/Test/data/MUSICO/RETRO-WAVE/College/Northern Council/vinyl-info.northern-council.txt'
  ---> 🌊 UNIVERSAL//CONTAINS-FOLDERS-CALLBACK: (depth:9999) '/Users/plastikfan/dev/github/snivilised/extendio/Test/data/MUSICO/RETRO-WAVE/College/Teenage Color'
  ---> 🌊 UNIVERSAL//CONTAINS-FOLDERS-CALLBACK: (depth:9999) '/Users/plastikfan/dev/github/snivilised/extendio/Test/data/MUSICO/RETRO-WAVE/College/Teenage Color/A1 - Can You Kiss Me First.flac'
  ---> 🌊 UNIVERSAL//CONTAINS-FOLDERS-CALLBACK: (depth:9999) '/Users/plastikfan/dev/github/snivilised/extendio/Test/data/MUSICO/RETRO-WAVE/College/Teenage Color/A2 - Teenage Color.flac'
  ---> 🌊 UNIVERSAL//CONTAINS-FOLDERS-CALLBACK: (depth:9999) '/Users/plastikfan/dev/github/snivilised/extendio/Test/data/MUSICO/RETRO-WAVE/College/Teenage Color/vinyl-info.teenage-color.txt'
  ---> 🌊 UNIVERSAL//CONTAINS-FOLDERS-CALLBACK: (depth:9999) '/Users/plastikfan/dev/github/snivilised/extendio/Test/data/MUSICO/RETRO-WAVE/Electric Youth'
  ---> 🌊 UNIVERSAL//CONTAINS-FOLDERS-CALLBACK: (depth:9999) '/Users/plastikfan/dev/github/snivilised/extendio/Test/data/MUSICO/RETRO-WAVE/Electric Youth/Innerworld'
  ---> 🌊 UNIVERSAL//CONTAINS-FOLDERS-CALLBACK: (depth:9999) '/Users/plastikfan/dev/github/snivilised/extendio/Test/data/MUSICO/RETRO-WAVE/Electric Youth/Innerworld/A1 - Before Life.flac'
  ---> 🌊 UNIVERSAL//CONTAINS-FOLDERS-CALLBACK: (depth:9999) '/Users/plastikfan/dev/github/snivilised/extendio/Test/data/MUSICO/RETRO-WAVE/Electric Youth/Innerworld/A2 - Runaway.flac'
  ---> 🌊 UNIVERSAL//CONTAINS-FOLDERS-CALLBACK: (depth:9999) '/Users/plastikfan/dev/github/snivilised/extendio/Test/data/MUSICO/RETRO-WAVE/Electric Youth/Innerworld/vinyl-info.innerworld.txt'
*/

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
			path := helpers.Path(root, "RETRO-WAVE")
			providedOptions := nav.GetDefaultOptions()
			providedOptions.Store.Subscription = entry.subscription
			providedOptions.Store.DoExtend = true
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

			if entry.useLastFn {
				providedOptions.Sampler.Fn = nav.GetLastSampler(&entry.noOf)
			}

			if entry.filter != nil {
				filterDefs := &nav.FilterDefinitions{
					Children: nav.CompoundFilterDef{
						Type:        nav.FilterTypeGlobEn,
						Description: entry.filter.name,
						Pattern:     entry.filter.pattern,
						Negate:      entry.filter.negate,
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

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "universal: default (first), with 2 files",
				should:       "invoke for at most 2 files per directory",
				subscription: nav.SubscribeAny,
				prohibited:   []string{"cover.night-drive.jpg"},
			},
			noOf: nav.SampleNoOf{
				Files: 2,
			},
			expectedNoOf: directoryQuantities{
				files: 8,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "universal: default (first), with 2 folders",
				should:       "invoke for at most 2 folders per directory",
				subscription: nav.SubscribeAny,
				prohibited:   []string{"Electric Youth"},
			},
			noOf: nav.SampleNoOf{
				Folders: 2,
			},
			expectedNoOf: directoryQuantities{
				files:   11,
				folders: 6,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "universal: default (first), with 2 files and 2 folders",
				should:       "invoke for at most 2 files and 2 folders per directory",
				subscription: nav.SubscribeAny,
				prohibited:   []string{"cover.night-drive.jpg", "Electric Youth"},
			},
			noOf: nav.SampleNoOf{
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
				message:      "folders: default (first), with 2 folders",
				should:       "invoke for at most 2 folders per directory",
				subscription: nav.SubscribeFolders,
				prohibited:   []string{"Electric Youth"},
			},
			noOf: nav.SampleNoOf{
				Folders: 2,
			},
			expectedNoOf: directoryQuantities{
				folders: 6,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "folders: custom, with last single folder",
				should:       "invoke for only last folder per directory",
				subscription: nav.SubscribeFolders,
				prohibited:   []string{"Chromatics"},
			},
			useLastFn: true,
			noOf: nav.SampleNoOf{
				Folders: 1,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "folders: default (first), with 2 folders",
				should:       "invoke for at most 2 folders per directory",
				subscription: nav.SubscribeFoldersWithFiles,
				prohibited:   []string{"Electric Youth"},
			},
			noOf: nav.SampleNoOf{
				Folders: 2,
			},
			expectedNoOf: directoryQuantities{
				folders: 6,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "folders: custom, with last single folder",
				should:       "invoke for only last folder per directory",
				subscription: nav.SubscribeFoldersWithFiles,
				prohibited:   []string{"Chromatics"},
			},
			useLastFn: true,
			noOf: nav.SampleNoOf{
				Folders: 1,
			},
		}),

		// TODO: With filter for folders: first Co* folder (College)

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "files: default (first), with 2 files",
				should:       "invoke for at most 2 files per directory",
				subscription: nav.SubscribeFiles,
				prohibited:   []string{"cover.night-drive.jpg"},
			},
			noOf: nav.SampleNoOf{
				Files: 2,
			},
			expectedNoOf: directoryQuantities{
				files: 8,
			},
		}),

		Entry(nil, &sampleTE{
			naviTE: naviTE{
				message:      "files: custom, with last single file",
				should:       "invoke for only last file per directory",
				subscription: nav.SubscribeFiles,
				prohibited:   []string{"A1 - The Telephone Call.flac"},
			},
			useLastFn: true,
			noOf: nav.SampleNoOf{
				Files: 1,
			},
			expectedNoOf: directoryQuantities{
				files: 4,
			},
		}),

		// TODO: With filter for files: first cover* file
	)
})
