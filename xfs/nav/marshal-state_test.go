package nav_test

import (
	"fmt"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/extendio/translate"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("MarshalOptions", Ordered, func() {
	var (
		o            *nav.TraverseOptions
		root         string
		jroot        string
		fromJsonPath string
		toJsonPath   string

		filterDefs nav.FilterDefinitions
	)

	BeforeAll(func() {
		root = origin()
		jroot = joinCwd("Test", "json")
		fromJsonPath = strings.Join([]string{jroot, "persisted-state.json"}, string(filepath.Separator))
		toJsonPath = strings.Join([]string{jroot, "test-state-marshal.json"}, string(filepath.Separator))
		filterDefs = nav.FilterDefinitions{
			Current: nav.FilterDef{
				Type:            nav.FilterTypeGlobEn,
				Description:     "items with .flac suffix",
				Source:          "*.flac",
				Scope:           nav.ScopeLeafEn,
				Negate:          false,
				IfNotApplicable: true,
			},
			Children: nav.CompoundFilterDef{
				Type:        nav.FilterTypeRegexEn,
				Description: "jpg files",
				Source:      "\\.jpg$",
				Negate:      true,
			},
		}
	})

	BeforeEach(func() {
		o = nav.GetDefaultOptions()
	})

	Context("Marshal", func() {
		Context("given: correct config", func() {
			It("🧪 should: write options in JSON", func() {
				navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
					o.Store.DoExtend = true
					o.Store.FilterDefs = &filterDefs
					o.Callback = nav.LabelledTraverseCallback{
						Label: "test marshal state callback",
						Fn: func(item *nav.TraverseItem) *translate.LocalisableError {
							return nil
						},
					}
				})
				path := path(root, "RETRO-WAVE/Chromatics/Night Drive")
				_ = navigator.Walk(path)

				err := navigator.Save(toJsonPath)
				Expect(err).To(BeNil())
			})
		})

		DescribeTable("marshall error",
			func(entry *marshalTE) {
				defer func() {
					pe := recover()
					if err, ok := pe.(error); ok {
						Expect(
							strings.Contains(err.Error(), entry.errorContains)).To(BeTrue(),
							err.Error(),
						)
					}
				}()

				navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
					o.Persist.Format = entry.format
					o.Store.DoExtend = true
					o.Store.FilterDefs = &filterDefs
					o.Callback = nav.LabelledTraverseCallback{
						Label: "test marshal state callback",
						Fn: func(item *nav.TraverseItem) *translate.LocalisableError {
							return nil
						},
					}
				})
				path := path(root, entry.relative)
				_ = navigator.Walk(path)
				_ = navigator.Save(toJsonPath)

				Fail(fmt.Sprintf("❌ expected panic due to %v", entry.errorContains))
			},
			func(entry *marshalTE) string {
				return fmt.Sprintf("🧪 ===> given: '%v', should panic", entry.message)
			},
			Entry(nil, &marshalTE{
				naviTE: naviTE{
					message:  "unknown marshal format",
					relative: "RETRO-WAVE/Chromatics/Night Drive",
				},
				errorContains: "unknown marshal format",
			}),
		)
	})

	Context("Unmarshal", func() {
		Context("given: correct config", func() {
			It("🧪 should: write options in JSON", func() {
				Skip("needs resume function")
				restore := func(o *nav.TraverseOptions) {
					GinkgoWriter.Printf("---> marshaller ...\n")
				}
				restore(o)
				_ = fromJsonPath

				// err := o.UnmarshalDefunct(fromJsonPath)
				// Expect(err).To(BeNil())
			})
		})
	})
})
