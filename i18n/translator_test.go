package i18n_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/text/language"

	. "github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/internal/helpers"
	"github.com/snivilised/extendio/xfs/utils"
)

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
				Use(func(o *UseOptions) {
					o.Tag = language.AmericanEnglish
					o.From.Path = l10nPath
					o.From.Sources = testTranslationFile
				})
				Expect(TxRef.IsNone()).To(BeFalse())
				Expect(TxRef.Get().LanguageInfoRef().Get().Tag).To(Equal(language.AmericanEnglish))
			})
		})

		When("requested language is the default", func() {
			It("🧪 should: create Translator", func() {
				Use(func(o *UseOptions) {
					o.Tag = language.BritishEnglish
				})
				Expect(TxRef.IsNone()).To(BeFalse())
				Expect(TxRef.Get().LanguageInfoRef().Get().Tag).To(Equal(language.BritishEnglish))
			})
		})

		When("requested language is NOT available", func() {
			It("🧪 should: return error", func() {
				defer func() {
					pe := recover()
					if err, ok := pe.(error); !ok || !strings.Contains(err.Error(), "not available") {
						Fail("FAILED")
					}
				}()

				requested := language.Spanish
				Use(func(o *UseOptions) {
					o.Tag = requested
				})
				Fail(fmt.Sprintf("❌ expected panic due to request language: '%v' not being supported", requested))
			})
		})
	})

	Context("Error Checking", func() {
		Context("given: FailedToReadDirectoryContentsError", func() {
			It("🧪 should: be identifiable via query function", func() {
				reason := fmt.Errorf("file missing")
				var err error = NewFailedToReadDirectoryContentsError("/foo/bar/", reason)
				result := QueryFailedToReadDirectoryContentsError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: NewFailedToResumeFromFileError", func() {
			It("🧪 should: be identifiable via query function", func() {
				reason := fmt.Errorf("file missing")
				var err error = NewFailedToResumeFromFileError("/foo/bar/resume.json", reason)
				result := QueryFailedToResumeFromFileError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: InvalidConfigEntryError", func() {
			It("🧪 should: be identifiable via query function", func() {
				var err error = NewInvalidConfigEntryError("foo", "Store/Logging/Path")
				result := QueryInvalidConfigEntryError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: InvalidResumeStrategyError", func() {
			It("🧪 should: be identifiable via query function", func() {
				var err error = NewInvalidResumeStrategyError("foo")
				result := QueryInvalidResumeStrategyError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: MissingCallbackError", func() {
			It("🧪 should: be identifiable via query function", func() {
				var err error = NewMissingCallbackError()
				result := QueryMissingCallbackError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: MissingCustomFilterDefinitionError", func() {
			It("🧪 should: be identifiable via query function", func() {
				var err error = NewMissingCustomFilterDefinitionError(
					"Options/Store/FilterDefs/Node/Custom",
				)
				result := QueryMissingCustomFilterDefinitionError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: NotADirectoryError", func() {
			It("🧪 should: be identifiable via query function", func() {
				var err error = NewNotADirectoryError("/foo/bar")
				result := QueryNotADirectoryError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: SortFnFailedError", func() {
			It("🧪 should: be identifiable via query function", func() {
				var err error = NewSortFnFailedError()
				result := QuerySortFnFailedError(err)
				Expect(result).To(BeTrue())
			})
		})

		Context("given: UnknownMarshalFormatError", func() {
			It("🧪 should: be identifiable via query function", func() {
				var err error = NewUnknownMarshalFormatError(
					"Options/Persist/Format", "jpg",
				)
				result := QueryUnknownMarshalFormatError(err)
				Expect(result).To(BeTrue())
			})
		})
	})
})
