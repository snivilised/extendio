package i18n_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/snivilised/extendio/internal/helpers"
	"github.com/snivilised/extendio/xfs/utils"
	"golang.org/x/text/language"

	. "github.com/snivilised/extendio/i18n"
)

var _ = Describe("Text", Ordered, func() {
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

	Context("native", func() {
		BeforeEach(func() {
			Use(func(o *UseOptions) {
				o.Tag = DefaultLanguage.Get()
				o.From.Sources = testTranslationFile
			})
		})

		Context("given: ThirdPartyError", func() {
			It("🧪 should: contain the third party error text", func() {
				Use(func(o *UseOptions) {
					o.Tag = language.BritishEnglish
				})
				err := NewThirdPartyErr(errors.New("computer says no"))
				Expect(err.Error()).To(ContainSubstring("computer says no"))
			})
		})

		Context("Text", func() {
			Context("given: a template data instance", func() {
				It("🧪 should: evaluate translated text", func() {
					Expect(Text(ThirdPartyErrorTemplData{
						Error: errors.New("out of stock"),
					})).NotTo(BeNil())
				})
			})
		})
	})

	Context("foreign", func() {
		BeforeEach(func() {
			Use(func(o *UseOptions) {
				o.Tag = language.AmericanEnglish
				o.From.Path = l10nPath
				o.From.Sources = testTranslationFile
			})
		})

		Context("Text", func() {
			Context("given: a template data instance", func() {
				It("🧪 should: evaluate translated text(internationalization)", func() {
					Expect(Text(InternationalisationTemplData{})).To(Equal("internationalization"))
				})

				It("🧪 should: evaluate translated text(localization)", func() {
					Expect(Text(LocalisationTemplData{})).To(Equal("localization"))
				})
			})
		})
	})
})
