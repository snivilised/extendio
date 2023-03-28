package i18n_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	xi18n "github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/internal/helpers"
	"github.com/snivilised/extendio/xfs/utils"
	"golang.org/x/text/language"
)

var _ = Describe("SingularTranslatorFactory", Ordered, func() {
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
		factory = &xi18n.SingularTranslatorFactory{}
	})

	Context("given: dependency supports requested language", func() {
		var local *xi18n.LanguageInfo

		BeforeEach(func() {
			from = xi18n.LoadFrom{
				Path: l10nPath,
				Sources: xi18n.TranslationFiles{
					xi18n.SOURCE_ID: xi18n.TranslationSource{"test"},
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

			Context("given: dependency supports requested language", func() {
				var local *xi18n.LanguageInfo

				BeforeEach(func() {
					from = xi18n.LoadFrom{
						Path: l10nPath,
						Sources: xi18n.TranslationFiles{
							xi18n.SOURCE_ID: xi18n.TranslationSource{"test"},
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
							from = xi18n.LoadFrom{
								Path: l10nPath,
								Sources: xi18n.TranslationFiles{
									xi18n.SOURCE_ID: xi18n.TranslationSource{"test.graffico"},
								},
							}
							local = xi18n.NewLanguageInfo(&xi18n.UseOptions{
								Tag:  language.AmericanEnglish,
								From: from,
							})
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
						factory = &xi18n.SingularTranslatorFactory{
							AbstractTranslatorFactory: xi18n.AbstractTranslatorFactory{
								Create: dummyLocalizer,
							},
						}
						local = xi18n.NewLanguageInfo(&xi18n.UseOptions{
							Tag:    language.AmericanEnglish,
							From:   from,
							Custom: any("foo-bar"),
						})

						translator := factory.New(local)
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
									xi18n.SOURCE_ID: xi18n.TranslationSource{"scooby-doo"},
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
						It("🧪 should: return original text", func() {
							// Since the singular translator does not need to perform a lookup
							// of the localizer (since there is only 1), it doesn't need to use
							// the message's source id, so that source id is irrelevant and
							// doesn't matter if it's incorrect.
							local = xi18n.NewLanguageInfo(&xi18n.UseOptions{
								Tag:  language.AmericanEnglish,
								From: from,
							})
							translator := factory.New(local)
							_ = xi18n.UseTx(translator)
							data := WrongSourceIdTemplData{}
							text := xi18n.Text(data)
							Expect(text).To(Equal(data.Message().Other))
						})
					})

					When("no sources specified have been specified", func() {
						It("🧪 should: panic", func() {
							defer func() {
								pe := recover()
								if err, ok := pe.(error); !ok || !strings.Contains(err.Error(),
									"no sources specified") {
									Fail(fmt.Sprintf(
										"Incorrect error reported, when: no sources specified sources have been specified 💥(%v)",
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
	})
})
