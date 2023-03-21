package i18n_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/text/language"

	. "github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/internal/helpers"
	"github.com/snivilised/extendio/xfs/utils"
)

var _ = Describe("Translator", Ordered, func() {
	var (
		repo string

		l10nPath string
	)

	BeforeAll(func() {
		repo = helpers.Repo("../..")
		l10nPath = helpers.Path(repo, "Test/data/l10n")
		Expect(utils.FolderExists(l10nPath)).To(BeTrue())
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
				Expect(Use(func(o *UseOptions) {
					o.Tag = language.AmericanEnglish
					o.Name = "test"
					o.Path = l10nPath
				})).Error().To(BeNil())
				Expect(TxRef.IsNone()).To(BeFalse())
				Expect(TxRef.Get().LanguageInfoRef().Get().Current).To(Equal(language.AmericanEnglish))
			})
		})

		When("requested language is the default", func() {
			It("🧪 should: create Translator", func() {
				Expect(Use(func(o *UseOptions) {
					o.Tag = language.BritishEnglish
				})).Error().To(BeNil())
				Expect(TxRef.IsNone()).To(BeFalse())
				Expect(TxRef.Get().LanguageInfoRef().Get().Current).To(Equal(language.BritishEnglish))
			})
		})

		When("requested language is NOT available", func() {
			It("🧪 should: return error", func() {
				Expect(Use(func(o *UseOptions) {
					o.Tag = language.Spanish
				})).Error().NotTo(BeNil())
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
