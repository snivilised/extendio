package i18n_test

import (
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/text/language"

	. "github.com/snivilised/extendio/i18n"
	xi18n "github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/internal/helpers"
	"github.com/snivilised/extendio/xfs/utils"
)

type dummyCreator struct {
	invoked bool
}

func (dc *dummyCreator) create(lang *xi18n.LanguageInfo, sourceId string) *i18n.Localizer {
	dc.invoked = true
	return &i18n.Localizer{}
}

var _ = Describe("Translator", Ordered, func() {
	var (
		repo                string
		l10nPath            string
		testTranslationFile TranslationFiles
	)

	BeforeAll(func() {
		repo = helpers.Repo("../..")
		l10nPath = helpers.Path(repo, "Test/data/l10n")
		Expect(utils.FolderExists(l10nPath)).To(BeTrue())

		testTranslationFile = TranslationFiles{
			SOURCE_ID: TranslationSource{"test"},
		}
	})

	BeforeEach(func() {
		ResetTx()
	})

	Context("TxRef.IsNone", func() {
		When("not Use'd", func() {
			It("🧪 should: be true", func() {
				Expect(TxRef.IsNone()).To(BeTrue(),
					"'Use' not invoked but TxRef.IsNone indicates, its still set?",
				)
			})
		})
	})

	Context("Use", func() {
		When("requested language is available", func() {
			It("🧪 should: create Translator", func() {
				if err := Use(func(o *UseOptions) {
					o.Tag = language.AmericanEnglish
					o.From.Path = l10nPath
					o.From.Sources = testTranslationFile
				}); err != nil {
					Fail(err.Error())
				}
				Expect(TxRef.IsNone()).To(BeFalse())
				Expect(TxRef.Get().LanguageInfoRef().Get().Tag).To(Equal(language.AmericanEnglish))
			})
		})

		When("requested language is the default", func() {
			It("🧪 should: create Translator", func() {
				if err := Use(func(o *UseOptions) {
					o.Tag = language.BritishEnglish
				}); err != nil {
					Fail(err.Error())
				}
				Expect(TxRef.IsNone()).To(BeFalse())
				Expect(TxRef.Get().LanguageInfoRef().Get().Tag).To(Equal(language.BritishEnglish))
			})
		})

		Context("DefaultIsAcceptable is true", func() {
			When("requested language is NOT available", func() {
				It("🧪 should: return NOT error", func() {
					requested := language.Spanish
					if err := Use(func(o *UseOptions) {
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
					if err := Use(func(o *UseOptions) {
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
				dummy := dummyCreator{}
				if err := Use(func(o *UseOptions) {
					o.Tag = language.BritishEnglish
					o.Create = dummy.create
				}); err != nil {
					Fail(err.Error())
				}
				Expect(dummy.invoked).To(BeTrue())
				Expect(TxRef.IsNone()).To(BeFalse())
				Expect(TxRef.Get().LanguageInfoRef().Get().Tag).To(Equal(language.BritishEnglish))
			})
		})
	})

	Context("Error Checking", func() {
		Context("given: FailedToReadDirectoryContentsError", func() {
			It("🧪 should: be identifiable via query function", func() {
				reason := fmt.Errorf("file missing")
				err := NewFailedToReadDirectoryContentsError("/foo/bar/", reason)
				result := QueryFailedToReadDirectoryContentsError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: NewFailedToResumeFromFileError", func() {
			It("🧪 should: be identifiable via query function", func() {
				reason := fmt.Errorf("file missing")
				err := NewFailedToResumeFromFileError("/foo/bar/resume.json", reason)
				result := QueryFailedToResumeFromFileError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: InvalidConfigEntryError", func() {
			It("🧪 should: be identifiable via query function", func() {
				err := NewInvalidConfigEntryError("foo", "Store/Logging/Path")
				result := QueryInvalidConfigEntryError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: InvalidResumeStrategyError", func() {
			It("🧪 should: be identifiable via query function", func() {
				err := NewInvalidResumeStrategyError("foo")
				result := QueryInvalidResumeStrategyError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: MissingCallbackError", func() {
			It("🧪 should: be identifiable via query function", func() {
				err := NewMissingCallbackError()
				result := QueryMissingCallbackError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: MissingCustomFilterDefinitionError", func() {
			It("🧪 should: be identifiable via query function", func() {
				err := NewMissingCustomFilterDefinitionError(
					"Options/Store/FilterDefs/Node/Custom",
				)
				result := QueryMissingCustomFilterDefinitionError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: NotADirectoryError", func() {
			It("🧪 should: be identifiable via query function", func() {
				err := NewNotADirectoryError("/foo/bar")
				result := QueryNotADirectoryError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: SortFnFailedError", func() {
			It("🧪 should: be identifiable via query function", func() {
				err := NewSortFnFailedError()
				result := QuerySortFnFailedError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: UnknownMarshalFormatError", func() {
			It("🧪 should: be identifiable via query function", func() {
				err := NewUnknownMarshalFormatError(
					"Options/Persist/Format", "jpg",
				)
				result := QueryUnknownMarshalFormatError(err)
				Expect(result).To(BeTrue())
			})
		})
	})
})
