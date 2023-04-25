package i18n_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/text/language"

	"github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/internal/helpers"
	"github.com/snivilised/extendio/xfs/utils"
)

var _ = Describe("Translator", Ordered, func() {
	var (
		repo                string
		l10nPath            string
		testTranslationFile i18n.TranslationFiles
	)

	BeforeAll(func() {
		repo = helpers.Repo("../..")
		l10nPath = helpers.Path(repo, "Test/data/l10n")
		Expect(utils.FolderExists(l10nPath)).To(BeTrue())

		testTranslationFile = i18n.TranslationFiles{
			i18n.ExtendioSourceID: i18n.TranslationSource{Name: "test"},
		}
	})

	BeforeEach(func() {
		i18n.ResetTx()
	})

	Context("xi18n.TxRef.IsNone", func() {
		When("not Use'd", func() {
			It("🧪 should: be true", func() {
				Expect(i18n.TxRef.IsNone()).To(BeTrue(),
					"'Use' not invoked but xi18n.TxRef.IsNone indicates, its still set?",
				)
			})
		})
	})

	Context("Use", func() {
		When("requested language is available", func() {
			It("🧪 should: create Translator", func() {
				if err := i18n.Use(func(o *i18n.UseOptions) {
					o.Tag = language.AmericanEnglish
					o.From.Path = l10nPath
					o.From.Sources = testTranslationFile
				}); err != nil {
					Fail(err.Error())
				}
				Expect(i18n.TxRef.IsNone()).To(BeFalse())
				Expect(i18n.TxRef.Get().LanguageInfoRef().Get().Tag).To(Equal(language.AmericanEnglish))
			})
		})

		When("requested language is the default", func() {
			It("🧪 should: create Translator", func() {
				if err := i18n.Use(func(o *i18n.UseOptions) {
					o.Tag = language.BritishEnglish
				}); err != nil {
					Fail(err.Error())
				}
				Expect(i18n.TxRef.IsNone()).To(BeFalse())
				Expect(i18n.TxRef.Get().LanguageInfoRef().Get().Tag).To(Equal(language.BritishEnglish))
			})
		})

		Context("DefaultIsAcceptable is true", func() {
			When("requested language is NOT available", func() {
				It("🧪 should: return NOT error", func() {
					requested := language.Spanish
					if err := i18n.Use(func(o *i18n.UseOptions) {
						o.Tag = requested
					}); err != nil {
						Fail(
							fmt.Sprintf(
								"❌ request language: '%v' not supported, but default should be acceptable",
								requested,
							),
						)
					}
				})
			})
		})

		Context("DefaultIsAcceptable is false", func() {
			When("requested language is NOT available", func() {
				It("🧪 should: return error", func() {
					requested := language.Spanish
					if err := i18n.Use(func(o *i18n.UseOptions) {
						o.DefaultIsAcceptable = false
						o.Tag = requested
					}); err == nil {
						Fail(
							fmt.Sprintf("❌ expected error due to request language: '%v' not being supported",
								requested,
							),
						)
					}
				})
			})
		})

		When("client provides Create function", func() {
			It("🧪 should: create the localizer with the override", func() {
				dummy := helpers.DummyCreator{}
				if err := i18n.Use(func(o *i18n.UseOptions) {
					o.Tag = language.BritishEnglish
					o.Create = dummy.Create
				}); err != nil {
					Fail(err.Error())
				}
				Expect(dummy.Invoked).To(BeTrue())
				Expect(i18n.TxRef.IsNone()).To(BeFalse())
				Expect(i18n.TxRef.Get().LanguageInfoRef().Get().Tag).To(Equal(language.BritishEnglish))
			})
		})

		When("extendio source not provided", func() {
			It("🧪 should: create Translator", func() {
				from := i18n.LoadFrom{
					Path: l10nPath,
					Sources: i18n.TranslationFiles{
						GrafficoSourceID: i18n.TranslationSource{Name: "test.graffico"},
					},
				}

				if err := i18n.Use(func(o *i18n.UseOptions) {
					o.Tag = language.BritishEnglish
					o.From = from
				}); err != nil {
					Fail(err.Error())
				}
				Expect(i18n.TxRef.IsNone()).To(BeFalse())
				Expect(i18n.TxRef.Get().LanguageInfoRef().Get().Tag).To(Equal(language.BritishEnglish))
			})
		})

		Context("negotiate translators", func() {
			When("new source already present", func() {
				It("🧪 should: ignore subsequent use source", func() {
					if err := i18n.Use(func(o *i18n.UseOptions) {
						o.Tag = language.AmericanEnglish
						o.From.Path = l10nPath
						o.From.Sources = i18n.TranslationFiles{
							i18n.ExtendioSourceID: i18n.TranslationSource{Name: "test"},
						}
					}); err != nil {
						Fail(err.Error())
					}

					if err := i18n.Use(func(o *i18n.UseOptions) {
						o.Tag = language.AmericanEnglish
						o.From.Path = "/foo-bar"
						o.From.Sources = i18n.TranslationFiles{
							i18n.ExtendioSourceID: i18n.TranslationSource{Name: "test"},
						}
					}); err != nil {
						Fail(err.Error())
					}
				})
			})

			When("new source already NOT already present", func() {
				It("🧪 should: add new source", func() {
					if err := i18n.Use(func(o *i18n.UseOptions) {
						o.Tag = language.AmericanEnglish
						o.From.Path = l10nPath
						o.From.Sources = i18n.TranslationFiles{
							i18n.ExtendioSourceID: i18n.TranslationSource{Name: "test"},
						}
					}); err != nil {
						Fail(err.Error())
					}

					if err := i18n.Use(func(o *i18n.UseOptions) {
						o.Tag = language.AmericanEnglish
						o.From.Path = l10nPath
						o.From.Sources = i18n.TranslationFiles{
							GrafficoSourceID: i18n.TranslationSource{Name: "test.graffico"},
						}
					}); err != nil {
						Fail(err.Error())
					}

					actual := i18n.Text(PavementGraffitiReportTemplData{
						Primary: "Violet",
					})
					Expect(actual).To(Equal(expectUS))
				})
			})
		})
	})

	Context("Error Scenarios", func() {
		Context("given: default is acceptable", func() {
			When("invalid translation source specified", func() {
				It("🧪 should: NOT return error", func() {
					if err := i18n.Use(func(o *i18n.UseOptions) {
						o.Tag = language.EuropeanSpanish
						o.From = i18n.LoadFrom{
							Path: l10nPath,
							Sources: i18n.TranslationFiles{
								i18n.ExtendioSourceID: i18n.TranslationSource{Name: "scooby-doo"},
								GrafficoSourceID:      i18n.TranslationSource{Name: "test.graffico"},
							},
						}
					}); err != nil {
						Fail(err.Error())
					}
				})

				Context("Path not specified", func() {
					It("🧪 should: try and load from exe path", func() {
						defer func() {
							pe := recover()
							if err, ok := pe.(error); !ok || !strings.Contains(err.Error(),
								"could not load translations for") {
								Fail(fmt.Sprintf(
									"Incorrect error reported, when: message defined with non-existent source id 💥(%v)",
									err.Error()),
								)
							}
						}()

						_ = i18n.Use(func(o *i18n.UseOptions) {
							o.DefaultIsAcceptable = false
							o.Tag = language.AmericanEnglish
							o.From = i18n.LoadFrom{
								Sources: i18n.TranslationFiles{
									GrafficoSourceID: i18n.TranslationSource{Name: "test.graffico"},
								},
							}
						})

						Fail("❌ expected panic due to invalid path: 'FOO-BAR'")
					})
				})
			})
		})

		Context("given: default is NOT acceptable", func() {
			Context("requested supported", func() {
				When("invalid translation source specified", func() {
					It("🧪 should: return error", func() {
						defer func() {
							pe := recover()
							if err, ok := pe.(error); !ok || !strings.Contains(err.Error(),
								"could not load translations for") {
								Fail(fmt.Sprintf(
									"Incorrect error reported, when: message defined with non-existent source id 💥(%v)",
									err.Error()),
								)
							}
						}()
						if err := i18n.Use(func(o *i18n.UseOptions) {
							o.DefaultIsAcceptable = false
							o.Tag = language.AmericanEnglish
							o.From = i18n.LoadFrom{
								Path: l10nPath,
								Sources: i18n.TranslationFiles{
									i18n.ExtendioSourceID: i18n.TranslationSource{Name: "scooby-doo"},
									GrafficoSourceID:      i18n.TranslationSource{Name: "test.graffico"},
								},
							}
						}); err == nil {
							Fail(err.Error())
						}
					})
				})
			})

			Context("requested NOT supported", func() {
				When("invalid translation source specified", func() {
					It("🧪 should: return error", func() {
						if err := i18n.Use(func(o *i18n.UseOptions) {
							o.DefaultIsAcceptable = false
							o.Tag = language.EuropeanSpanish
							o.From = i18n.LoadFrom{
								Path: l10nPath,
								Sources: i18n.TranslationFiles{
									i18n.ExtendioSourceID: i18n.TranslationSource{Name: "scooby-doo"},
									GrafficoSourceID:      i18n.TranslationSource{Name: "test.graffico"},
								},
							}
						}); err == nil {
							Fail(err.Error())
						}
					})
				})
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
				if err := i18n.Use(func(o *i18n.UseOptions) {
					o.Tag = language.AmericanEnglish
					o.From.Path = l10nPath
					o.From.Sources = i18n.TranslationFiles{
						i18n.ExtendioSourceID: i18n.TranslationSource{Name: "test"},
						GrafficoSourceID:      i18n.TranslationSource{Name: "test.graffico"},
					}
				}); err != nil {
					Fail(err.Error())
				}
				_ = i18n.Text(WrongSourceIDTemplData{})
				Fail("❌ expected panic due to invalid path: 'FOO-BAR'")
			})
		})
	})

	Context("Error Checking", func() {
		Context("given: FailedToReadDirectoryContentsError", func() {
			It("🧪 should: be identifiable via query function", func() {
				reason := fmt.Errorf("file missing")
				err := i18n.NewFailedToReadDirectoryContentsError("/foo/bar/", reason)
				result := i18n.QueryFailedToReadDirectoryContentsError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: NewFailedToResumeFromFileError", func() {
			It("🧪 should: be identifiable via query function", func() {
				reason := fmt.Errorf("file missing")
				err := i18n.NewFailedToResumeFromFileError("/foo/bar/resume.json", reason)
				result := i18n.QueryFailedToResumeFromFileError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: InvalidConfigEntryError", func() {
			It("🧪 should: be identifiable via query function", func() {
				err := i18n.NewInvalidConfigEntryError("foo", "Store/Logging/Path")
				result := i18n.QueryInvalidConfigEntryError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: InvalidResumeStrategyError", func() {
			It("🧪 should: be identifiable via query function", func() {
				err := i18n.NewInvalidResumeStrategyError("foo")
				result := i18n.QueryInvalidResumeStrategyError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: MissingCallbackError", func() {
			It("🧪 should: be identifiable via query function", func() {
				err := i18n.NewMissingCallbackError()
				result := i18n.QueryMissingCallbackError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: MissingCustomFilterDefinitionError", func() {
			It("🧪 should: be identifiable via query function", func() {
				err := i18n.NewMissingCustomFilterDefinitionError(
					"Options/Store/FilterDefs/Node/Custom",
				)
				result := i18n.QueryMissingCustomFilterDefinitionError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: NotADirectoryError", func() {
			It("🧪 should: be identifiable via query function", func() {
				err := i18n.NewNotADirectoryError("/foo/bar")
				result := i18n.QueryNotADirectoryError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: SortFnFailedError", func() {
			It("🧪 should: be identifiable via query function", func() {
				err := i18n.NewSortFnFailedError()
				result := i18n.QuerySortFnFailedError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: UnknownMarshalFormatError", func() {
			It("🧪 should: be identifiable via query function", func() {
				err := i18n.NewUnknownMarshalFormatError(
					"Options/Persist/Format", "jpg",
				)
				result := i18n.QueryUnknownMarshalFormatError(err)
				Expect(result).To(BeTrue())
			})
		})
	})
})
