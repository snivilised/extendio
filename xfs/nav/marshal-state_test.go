package nav_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/extendio/internal/helpers"

	. "github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("MarshalOptions", Ordered, func() {
	var (
		root       string
		jroot      string
		toJsonPath string
		filterDefs nav.FilterDefinitions
	)

	BeforeAll(func() {
		root = musico()
		jroot = helpers.JoinCwd("Test", "json")
		toJsonPath = helpers.Path(jroot, "test-state-marshal.json")

		filterDefs = nav.FilterDefinitions{
			Node: nav.FilterDef{
				Type:            nav.FilterTypeGlobEn,
				Description:     "items with .flac suffix",
				Pattern:         "*.flac",
				Scope:           nav.ScopeLeafEn,
				Negate:          false,
				IfNotApplicable: nav.TriStateBoolTrueEn,
			},
			Children: nav.CompoundFilterDef{
				Type:        nav.FilterTypeRegexEn,
				Description: "jpg files",
				Pattern:     "\\.jpg$",
				Negate:      true,
			},
		}
	})

	BeforeEach(func() {
		ResetTx()
		_ = Use(func(o *UseOptions) {
			o.Tag = DefaultLanguage.Get()
		})
	})

	Context("Marshal", func() {
		Context("given: correct config", func() {
			It("🧪 should: write options in JSON", func() {
				path := helpers.Path(root, "RETRO-WAVE")
				session := &nav.PrimarySession{
					Path: path,
				}

				_, _ = session.Configure(func(o *nav.TraverseOptions) {
					o.Store.Subscription = nav.SubscribeAny
					o.Store.DoExtend = true
					o.Store.FilterDefs = &filterDefs
					o.Callback = nav.LabelledTraverseCallback{
						Label: "test marshal state callback",
						Fn: func(item *nav.TraverseItem) error {
							return nil
						},
					}
				}).Run()

				err := session.Save(toJsonPath)
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

				path := helpers.Path(root, entry.relative)
				session := &nav.PrimarySession{
					Path: path,
				}

				_, _ = session.Configure(func(o *nav.TraverseOptions) {
					o.Persist.Format = entry.format
					o.Store.DoExtend = true
					o.Store.FilterDefs = &filterDefs
					o.Callback = nav.LabelledTraverseCallback{
						Label: "test marshal state callback",
						Fn: func(item *nav.TraverseItem) error {
							return nil
						},
					}
				}).Run()

				_ = session.Save(toJsonPath)

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
})
