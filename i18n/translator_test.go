package i18n_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/text/language"

	xi18n "github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/internal/helpers"
	"github.com/snivilised/extendio/xfs/utils"
)

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
