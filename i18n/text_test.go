package i18n_test

import (
	"errors"
	"strings"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ginkgo ok
	. "github.com/onsi/gomega"    //nolint:revive // gomega ok
	xi18n "github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/internal/helpers"
	"github.com/snivilised/extendio/xfs/utils"
	"golang.org/x/text/language"
)

type textTE struct {
	message           string
	path              string
	sourcePath        string
	name              string
	defaultAcceptable bool
}

const (
	expectUS = "Found graffiti on sidewalk; primary color: 'Violet'"
	expectGB = "Found graffiti on pavement; primary colour: 'Violet'"
)

func testTranslationPath(entry *textTE) string {
	// this test form required, because DescribeTable can't be used
	// due to not being able to setup state correctly, eg l10nPath
	//
	if err := xi18n.Use(func(o *xi18n.UseOptions) {
		o.Tag = language.AmericanEnglish
		o.DefaultIsAcceptable = entry.defaultAcceptable
		o.From = xi18n.LoadFrom{
			Path: entry.path,
			Sources: xi18n.TranslationFiles{
				GrafficoSourceID: xi18n.TranslationSource{
					Path: entry.sourcePath,
					Name: "test.graffico",
				},
			},
		}
	}); err != nil {
		Fail(err.Error())
	}

	return xi18n.Text(PavementGraffitiReportTemplData{
		Primary: "Violet",
	})
}

var _ = Describe("Text", Ordered, func() {
	var (
		repo                string
		l10nPath            string
		testTranslationFile xi18n.TranslationFiles
	)

	BeforeAll(func() {
		repo = helpers.Repo("../..")
		l10nPath = helpers.Path(repo, "Test/data/l10n")
		Expect(utils.FolderExists(l10nPath)).To(BeTrue())

		testTranslationFile = xi18n.TranslationFiles{
			xi18n.ExtendioSourceID: xi18n.TranslationSource{Name: "test"},
		}
	})

	BeforeEach(func() {
		xi18n.ResetTx()
	})

	Context("Default Language", func() {
		BeforeEach(func() {
			if err := xi18n.Use(func(o *xi18n.UseOptions) {
				o.Tag = xi18n.DefaultLanguage.Get()
				o.From.Sources = testTranslationFile
			}); err != nil {
				Fail(err.Error())
			}
		})

		Context("given: ThirdPartyError", func() {
			It("🧪 should: contain the third party error text", func() {
				if err := xi18n.Use(func(o *xi18n.UseOptions) {
					o.Tag = language.BritishEnglish
				}); err != nil {
					Fail(err.Error())
				}

				err := xi18n.NewThirdPartyErr(errors.New("computer says no"))
				Expect(err.Error()).To(ContainSubstring("computer says no"))
			})

			Context("Text", func() {
				Context("given: a template data instance", func() {
					It("🧪 should: evaluate translated text", func() {
						Expect(xi18n.Text(xi18n.ThirdPartyErrorTemplData{
							Error: errors.New("out of stock"),
						})).NotTo(BeNil())
					})
				})
			})
		})
	})

	Context("Foreign Language", func() {
		BeforeEach(func() {
			if err := xi18n.Use(func(o *xi18n.UseOptions) {
				o.Tag = language.AmericanEnglish
				o.From.Path = l10nPath
				o.From.Sources = testTranslationFile
			}); err != nil {
				Fail(err.Error())
			}
		})

		Context("Text", func() {
			Context("given: a template data instance", func() {
				It("🧪 should: evaluate translated text(internationalization)", func() {
					Expect(xi18n.Text(xi18n.InternationalisationTemplData{})).To(Equal("internationalization"))
				})

				It("🧪 should: evaluate translated text(localization)", func() {
					Expect(xi18n.Text(xi18n.LocalisationTemplData{})).To(Equal("localization"))
				})
			})
		})
	})

	Context("Multiple Sources", func() {
		Context("Foreign Language", func() {
			It("🧪 should: translate text with the correct localizer", func() {
				if err := xi18n.Use(func(o *xi18n.UseOptions) {
					o.Tag = language.AmericanEnglish
					o.From = xi18n.LoadFrom{
						Path: l10nPath,
						Sources: xi18n.TranslationFiles{
							xi18n.ExtendioSourceID: xi18n.TranslationSource{Name: "test"},
							GrafficoSourceID:       xi18n.TranslationSource{Name: "test.graffico"},
						},
					}
				}); err != nil {
					Fail(err.Error())
				}
				actual := xi18n.Text(PavementGraffitiReportTemplData{
					Primary: "Violet",
				})
				Expect(actual).To(Equal(expectUS))
			})
		})
	})

	When("extendio source not provided", func() {
		Context("Default Language", func() {
			It("🧪 should: create factory that contains the extendio source", func() {
				if err := xi18n.Use(func(o *xi18n.UseOptions) {
					o.Tag = language.BritishEnglish
					o.From = xi18n.LoadFrom{
						Path: l10nPath,
						Sources: xi18n.TranslationFiles{
							GrafficoSourceID: xi18n.TranslationSource{Name: "test.graffico"},
						},
					}
				}); err != nil {
					Fail(err.Error())
				}

				actual := xi18n.Text(xi18n.InternationalisationTemplData{})
				Expect(actual).To(Equal("internationalisation"))
			})
		})
	})

	Context("translation source contains path", func() {
		Context("Foreign Language", func() {
			Context("given: source path exists", func() {
				It("🧪 should: create localizer from source path", func() {
					actual := testTranslationPath(&textTE{
						sourcePath:        l10nPath,
						defaultAcceptable: true,
					})

					Expect(actual).To(Equal(expectUS))
				})
			})

			Context("given: path exists", func() {
				It("🧪 should: create localizer from path", func() {
					actual := testTranslationPath(&textTE{
						path:              l10nPath,
						defaultAcceptable: true,
					})

					Expect(actual).To(Equal(expectUS))
				})
			})

			Context("given: neither path exists", func() {
				It("🧪 should: create localizer using default language", func() {
					actual := testTranslationPath(&textTE{
						defaultAcceptable: true,
					})

					Expect(actual).To(Equal(expectGB))
				})
			})

			Context("given: neither path exists", func() {
				It("🧪 should: create localizer using default language", func() {
					defer func() {
						pe := recover()
						if err, ok := pe.(error); !ok || !strings.Contains(err.Error(),
							"could not load translations for") {
							Fail("translation file not available with exe")
						}
					}()

					_ = testTranslationPath(&textTE{
						defaultAcceptable: false,
					})

					Fail("❌ expected panic due translation file not available with exe")
				})
			})
		})
	})
})
