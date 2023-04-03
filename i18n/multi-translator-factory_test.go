package i18n_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/text/language"

	xi18n "github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/internal/helpers"
	"github.com/snivilised/extendio/xfs/utils"
)

const (
	expectUS = "Found graffiti on sidewalk; primary color: 'Violet'"
	expectGB = "Found graffiti on pavement; primary colour: 'Violet'"
)

// This is an example of how a client library should implement
// their Use function. In this particular case, the addition
// source is already present, so the result of AppendSources
// will contain the source id of the duplicated dependency, ie
// "test.graffico", which is simply ignored as it should be.
func clientUse(options ...xi18n.UseOptionFn) error {
	o := append(options, func(uo *xi18n.UseOptions) {
		_ = uo.From.AppendSources(&xi18n.TranslationFiles{
			GRAFFICO_SOURCE_ID: xi18n.TranslationSource{
				Name: "test.graffico",
			},
		})
	})
	return xi18n.Use(o...)
}

var _ = Describe("MultiTranslatorFactory", Ordered, func() {
	var (
		repo     string
		l10nPath string
		factory  xi18n.TranslatorFactory
		from     xi18n.LoadFrom
	)

	BeforeAll(func() {
		repo = helpers.Repo("../..")
		l10nPath = helpers.Path(repo, "Test/data/l10n")
		Expect(utils.FolderExists(l10nPath)).To(BeTrue())
	})

	BeforeEach(func() {
		factory = &xi18n.MultiTranslatorFactory{}
	})

	Context("given: dependency supports requested language", func() {
		var local *xi18n.LanguageInfo

		BeforeEach(func() {
			from = xi18n.LoadFrom{
				Path: l10nPath,
				Sources: xi18n.TranslationFiles{
					xi18n.SOURCE_ID:    xi18n.TranslationSource{Name: "test"},
					GRAFFICO_SOURCE_ID: xi18n.TranslationSource{Name: "test.graffico"},
				},
			}
		})

		Context("Default Language", func() {
			BeforeEach(func() {
				local = xi18n.NewLanguageInfo(&xi18n.UseOptions{
					Tag:  language.BritishEnglish,
					From: from,
				})
			})

			It("🧪 should: create Translator containing localizer", func() {
				translator := factory.New(local)
				Expect(translator).ToNot(BeNil())
				Expect(xi18n.UseTx(translator)).Error().To(BeNil())
			})

			Context("Text", func() {
				It("🧪 should: translate text with the correct localizer", func() {
					translator := factory.New(local)
					_ = xi18n.UseTx(translator)
					expect := expectGB
					actual := xi18n.Text(PavementGraffitiReportTemplData{
						Primary: "Violet",
					})
					Expect(actual).To(Equal(expect))
				})
			})

			When("extendio source not provided", func() {
				It("🧪 should: create factory that contains the extendio source", func() {
					from = xi18n.LoadFrom{
						Path: l10nPath,
						Sources: xi18n.TranslationFiles{
							GRAFFICO_SOURCE_ID: xi18n.TranslationSource{Name: "test.graffico"},
						},
					}

					local = xi18n.NewLanguageInfo(&xi18n.UseOptions{
						Tag:  language.BritishEnglish,
						From: from,
					})
					translator := factory.New(local)
					Expect(translator).ToNot(BeNil())
					Expect(xi18n.UseTx(translator)).Error().To(BeNil())
				})
			})

			When("duplicate sources for the same dependency", func() {
				It("🧪 should: return duplicate", func() {
					from = xi18n.LoadFrom{
						Path: l10nPath,
						Sources: xi18n.TranslationFiles{
							GRAFFICO_SOURCE_ID: xi18n.TranslationSource{
								Name: "test.graffico",
							},
						},
					}
					duplicates := from.AppendSources(&xi18n.TranslationFiles{
						GRAFFICO_SOURCE_ID: xi18n.TranslationSource{
							Name: "test.graffico",
						},
					})
					Expect(duplicates[0]).To(Equal("test.graffico"))
				})
			})

			When("client Use function uses a dependency already registered", func() {
				It("🧪 should: ignore the duplicate", func() {
					err := clientUse(func(uo *xi18n.UseOptions) {
						uo.From = xi18n.LoadFrom{
							Path: l10nPath,
							Sources: xi18n.TranslationFiles{
								GRAFFICO_SOURCE_ID: xi18n.TranslationSource{
									Name: "test.graffico",
								},
							},
						}
					})
					Expect(err).Error().To(BeNil())
				})
			})
		})

		Context("Foreign Language", func() {
			BeforeEach(func() {
				local = xi18n.NewLanguageInfo(&xi18n.UseOptions{
					Tag:  language.AmericanEnglish,
					From: from,
				})
			})

			Context("Text", func() {
				It("🧪 should: translate text with the correct localizer", func() {
					translator := factory.New(local)
					_ = xi18n.UseTx(translator)
					actual := xi18n.Text(PavementGraffitiReportTemplData{
						Primary: "Violet",
					})
					Expect(actual).To(Equal(expectUS))
				})
			})
		})

		When("custom function provided", func() {
			It("🧪 should: use custom localizer creator", func() {
				dummy := helpers.DummyCreator{}
				factory = &xi18n.MultiTranslatorFactory{
					AbstractTranslatorFactory: xi18n.AbstractTranslatorFactory{
						Create: dummy.Create,
					},
				}
				local = xi18n.NewLanguageInfo(&xi18n.UseOptions{
					Tag:    language.AmericanEnglish,
					From:   from,
					Custom: any("foo-bar"),
				})

				translator := factory.New(local)
				Expect(dummy.Invoked).To(BeTrue())
				Expect(translator).ToNot(BeNil())
			})
		})

		Context("Error Scenarios", func() {
			When("invalid translation source specified", func() {
				It("🧪 should: panic", func() {
					defer func() {
						pe := recover()
						if err, ok := pe.(error); !ok || !strings.Contains(err.Error(),
							"could not load translations") {
							Fail(fmt.Sprintf(
								"Incorrect error reported, when: invalid translation source specified 💥(%v)",
								err.Error()),
							)
						}
					}()

					from = xi18n.LoadFrom{
						Path: l10nPath,
						Sources: xi18n.TranslationFiles{
							xi18n.SOURCE_ID:    xi18n.TranslationSource{Name: "scooby-doo"},
							GRAFFICO_SOURCE_ID: xi18n.TranslationSource{Name: "test.graffico"},
						},
					}
					local := xi18n.NewLanguageInfo(&xi18n.UseOptions{
						Tag:  language.AmericanEnglish,
						From: from,
					})
					_ = factory.New(local)
					Fail("❌ expected panic due to invalid path: 'scooby-doo.active.en-US.json'")
				})
			})

			When("message defined with non-existent source id", func() {
				It("🧪 should: panic", func() {
					defer func() {
						pe := recover()
						if err, ok := pe.(error); !ok || !strings.Contains(err.Error(),
							"could not find localizer for source") {
							Fail(fmt.Sprintf(
								"Incorrect error reported, when: message defined with non-existent source id 💥(%v)",
								err.Error()),
							)
						}
					}()

					local = xi18n.NewLanguageInfo(&xi18n.UseOptions{
						Tag:  language.AmericanEnglish,
						From: from,
					})
					translator := factory.New(local)
					_ = xi18n.UseTx(translator)
					_ = xi18n.Text(WrongSourceIdTemplData{})
					Fail("❌ expected panic due to invalid path: 'FOO-BAR'")
				})
			})

			When("insufficient number of sources have been specified", func() {
				It("🧪 should: panic", func() {
					defer func() {
						pe := recover()
						if err, ok := pe.(error); !ok || !strings.Contains(err.Error(),
							"insufficient sources") {
							Fail(fmt.Sprintf(
								"Incorrect error reported, when: insufficient number of sources have been specified 💥(%v)",
								err.Error()),
							)
						}
					}()

					from = xi18n.LoadFrom{
						Path:    l10nPath,
						Sources: xi18n.TranslationFiles{},
					}
					local = xi18n.NewLanguageInfo(&xi18n.UseOptions{
						Tag:  language.AmericanEnglish,
						From: from,
					})

					translator := factory.New(local)
					_ = xi18n.UseTx(translator)
					_ = xi18n.Text(WrongSourceIdTemplData{})
					Fail("❌ expected panic due to invalid path: 'FOO-BAR'")
				})
			})
		})
	})
})

// NB: sha1 created manually with sha1sum command
// eg: sha1sum <text-file-containing-content-to-hash>
