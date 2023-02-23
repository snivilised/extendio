package i18n_test

import (
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

	Context("TxRef.IsNone", func() {
		When("not Use'd", func() {
			It("🧪 should: be true", func() {
				Expect(TxRef.IsNone()).To(BeTrue(), "client should always call 'Use' before use")
			})
		})
	})

	Context("Use", func() {
		When("requested language is available", func() {
			It("🧪 should: create Translator", func() {
				Expect(Use(func(o *UseOptions) {
					o.Tag = language.AmericanEnglish
					o.App = "test"
					o.Path = l10nPath
				})).Error().To(BeNil())
				Expect(TxRef.IsNone()).To(BeFalse())
				Expect(TxRef.Get().LanguageInfoRef.Get().Current).To(Equal(language.AmericanEnglish))
			})
		})

		When("requested language is the default", func() {
			It("🧪 should: create Translator", func() {
				Expect(Use(func(o *UseOptions) {
					o.Tag = language.BritishEnglish
				})).Error().To(BeNil())
				Expect(TxRef.IsNone()).To(BeFalse())
				Expect(TxRef.Get().LanguageInfoRef.Get().Current).To(Equal(language.BritishEnglish))
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

	Context("Text", func() {
		Context("given: a template data instance", func() {
			It("🧪 should: evaluate translated text", func() {
				_ = Use(func(o *UseOptions) {
					o.Tag = language.BritishEnglish
				})
				Expect(Text(ThirdPartyErrorTemplData{
					Error: "out of stock",
				})).NotTo(BeNil())
			})
		})
	})
})
