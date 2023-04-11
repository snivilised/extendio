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

var _ = Describe("Translator", Ordered, func() {
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
			xi18n.SOURCE_ID: xi18n.TranslationSource{"test"},
		}
	})

	BeforeEach(func() {
		xi18n.ResetTx()
	})

	Context("xi18n.TxRef.IsNone", func() {
		When("not Use'd", func() {
			It("🧪 should: be true", func() {
				Expect(xi18n.TxRef.IsNone()).To(BeTrue(),
					"'Use' not invoked but xi18n.TxRef.IsNone indicates, its still set?",
				)
			})
		})
	})

	Context("LoadFrom", func() {
		When("duplicate sources for the same dependency", func() {
			It("🧪 should: return duplicate", func() {
				from := xi18n.LoadFrom{
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

	Context("Use", func() {
		When("requested language is available", func() {
			It("🧪 should: create Translator", func() {
				if err := xi18n.Use(func(o *xi18n.UseOptions) {
					o.Tag = language.AmericanEnglish
					o.From.Path = l10nPath
					o.From.Sources = testTranslationFile
				}); err != nil {
					Fail(err.Error())
				}
				Expect(xi18n.TxRef.IsNone()).To(BeFalse())
				Expect(xi18n.TxRef.Get().LanguageInfoRef().Get().Tag).To(Equal(language.AmericanEnglish))
			})
		})

		When("requested language is the default", func() {
			It("🧪 should: create Translator", func() {
				if err := xi18n.Use(func(o *xi18n.UseOptions) {
					o.Tag = language.BritishEnglish
				}); err != nil {
					Fail(err.Error())
				}
				Expect(xi18n.TxRef.IsNone()).To(BeFalse())
				Expect(xi18n.TxRef.Get().LanguageInfoRef().Get().Tag).To(Equal(language.BritishEnglish))
			})
		})

		Context("DefaultIsAcceptable is true", func() {
			When("requested language is NOT available", func() {
				It("🧪 should: return NOT error", func() {
					requested := language.Spanish
					if err := xi18n.Use(func(o *xi18n.UseOptions) {
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
					if err := xi18n.Use(func(o *xi18n.UseOptions) {
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
				if err := xi18n.Use(func(o *xi18n.UseOptions) {
					o.Tag = language.BritishEnglish
					o.Create = dummy.Create
				}); err != nil {
					Fail(err.Error())
				}
				Expect(dummy.Invoked).To(BeTrue())
				Expect(xi18n.TxRef.IsNone()).To(BeFalse())
				Expect(xi18n.TxRef.Get().LanguageInfoRef().Get().Tag).To(Equal(language.BritishEnglish))
			})
		})

		When("extendio source not provided", func() {
			It("🧪 should: create Translator", func() {
				from := xi18n.LoadFrom{
					Path: l10nPath,
					Sources: xi18n.TranslationFiles{
						GRAFFICO_SOURCE_ID: xi18n.TranslationSource{Name: "test.graffico"},
					},
				}

				if err := xi18n.Use(func(o *xi18n.UseOptions) {
					o.Tag = language.BritishEnglish
					o.From = from
				}); err != nil {
					Fail(err.Error())
				}
				Expect(xi18n.TxRef.IsNone()).To(BeFalse())
				Expect(xi18n.TxRef.Get().LanguageInfoRef().Get().Tag).To(Equal(language.BritishEnglish))
			})
		})
	})

	Context("Error Scenarios", func() {
		Context("given: default is acceptable", func() {
			When("invalid translation source specified", func() {
				It("🧪 should: NOT return error", func() {
					if err := xi18n.Use(func(o *xi18n.UseOptions) {
						o.Tag = language.EuropeanSpanish
						o.From = xi18n.LoadFrom{
							Path: l10nPath,
							Sources: xi18n.TranslationFiles{
								xi18n.SOURCE_ID:    xi18n.TranslationSource{Name: "scooby-doo"},
								GRAFFICO_SOURCE_ID: xi18n.TranslationSource{Name: "test.graffico"},
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

						_ = xi18n.Use(func(o *xi18n.UseOptions) {
							o.DefaultIsAcceptable = false
							o.Tag = language.AmericanEnglish
							o.From = xi18n.LoadFrom{
								Sources: xi18n.TranslationFiles{
									GRAFFICO_SOURCE_ID: xi18n.TranslationSource{Name: "test.graffico"},
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
						if err := xi18n.Use(func(o *xi18n.UseOptions) {
							o.DefaultIsAcceptable = false
							o.Tag = language.AmericanEnglish
							o.From = xi18n.LoadFrom{
								Path: l10nPath,
								Sources: xi18n.TranslationFiles{
									xi18n.SOURCE_ID:    xi18n.TranslationSource{Name: "scooby-doo"},
									GRAFFICO_SOURCE_ID: xi18n.TranslationSource{Name: "test.graffico"},
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
						if err := xi18n.Use(func(o *xi18n.UseOptions) {
							o.DefaultIsAcceptable = false
							o.Tag = language.EuropeanSpanish
							o.From = xi18n.LoadFrom{
								Path: l10nPath,
								Sources: xi18n.TranslationFiles{
									xi18n.SOURCE_ID:    xi18n.TranslationSource{Name: "scooby-doo"},
									GRAFFICO_SOURCE_ID: xi18n.TranslationSource{Name: "test.graffico"},
								},
							}
						}); err == nil {
							Fail(err.Error())
						}
					})
				})
			})

			XContext("could not load translations for", func() {

			})
		})

		When("message defined with non-existent source id", func() {
			Context("single source", func() {
				It("🧪 should: return default string", func() {
					// the singular translator does not check source ids, so
					// won't report an error if the source id of the message
					// does not correspond to the localizer. This is not currently
					// deemed to be problem worth actioning. (May change in future
					// though)
					if err := xi18n.Use(func(o *xi18n.UseOptions) {
						o.Tag = language.AmericanEnglish
						o.From.Path = l10nPath
						o.From.Sources = xi18n.TranslationFiles{
							xi18n.SOURCE_ID: xi18n.TranslationSource{Name: "test"},
						}
					}); err != nil {
						Fail(err.Error())
					}
					actual := xi18n.Text(WrongSourceIdTemplData{})
					Expect(actual).To(Equal("Message with wrong id"))
				})
			})

			Context("multi sources", func() {
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
					if err := xi18n.Use(func(o *xi18n.UseOptions) {
						o.Tag = language.AmericanEnglish
						o.From.Path = l10nPath
						o.From.Sources = xi18n.TranslationFiles{
							xi18n.SOURCE_ID:    xi18n.TranslationSource{Name: "test"},
							GRAFFICO_SOURCE_ID: xi18n.TranslationSource{Name: "test.graffico"},
						}
					}); err != nil {
						Fail(err.Error())
					}
					_ = xi18n.Text(WrongSourceIdTemplData{})
					Fail("❌ expected panic due to invalid path: 'FOO-BAR'")
				})
			})
		})
	})

	Context("Error Checking", func() {
		Context("given: FailedToReadDirectoryContentsError", func() {
			It("🧪 should: be identifiable via query function", func() {
				reason := fmt.Errorf("file missing")
				err := xi18n.NewFailedToReadDirectoryContentsError("/foo/bar/", reason)
				result := xi18n.QueryFailedToReadDirectoryContentsError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: NewFailedToResumeFromFileError", func() {
			It("🧪 should: be identifiable via query function", func() {
				reason := fmt.Errorf("file missing")
				err := xi18n.NewFailedToResumeFromFileError("/foo/bar/resume.json", reason)
				result := xi18n.QueryFailedToResumeFromFileError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: InvalidConfigEntryError", func() {
			It("🧪 should: be identifiable via query function", func() {
				err := xi18n.NewInvalidConfigEntryError("foo", "Store/Logging/Path")
				result := xi18n.QueryInvalidConfigEntryError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: InvalidResumeStrategyError", func() {
			It("🧪 should: be identifiable via query function", func() {
				err := xi18n.NewInvalidResumeStrategyError("foo")
				result := xi18n.QueryInvalidResumeStrategyError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: MissingCallbackError", func() {
			It("🧪 should: be identifiable via query function", func() {
				err := xi18n.NewMissingCallbackError()
				result := xi18n.QueryMissingCallbackError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: MissingCustomFilterDefinitionError", func() {
			It("🧪 should: be identifiable via query function", func() {
				err := xi18n.NewMissingCustomFilterDefinitionError(
					"Options/Store/FilterDefs/Node/Custom",
				)
				result := xi18n.QueryMissingCustomFilterDefinitionError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: NotADirectoryError", func() {
			It("🧪 should: be identifiable via query function", func() {
				err := xi18n.NewNotADirectoryError("/foo/bar")
				result := xi18n.QueryNotADirectoryError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: SortFnFailedError", func() {
			It("🧪 should: be identifiable via query function", func() {
				err := xi18n.NewSortFnFailedError()
				result := xi18n.QuerySortFnFailedError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: UnknownMarshalFormatError", func() {
			It("🧪 should: be identifiable via query function", func() {
				err := xi18n.NewUnknownMarshalFormatError(
					"Options/Persist/Format", "jpg",
				)
				result := xi18n.QueryUnknownMarshalFormatError(err)
				Expect(result).To(BeTrue())
			})
		})
	})
})
