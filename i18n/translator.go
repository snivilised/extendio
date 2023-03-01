package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/samber/lo"
	"github.com/snivilised/extendio/xfs/utils"
	"golang.org/x/text/language"
)

var DefaultLanguage = utils.NewRoProp(language.BritishEnglish)

var tx *Translator
var TxRef utils.RoProp[*Translator] = utils.NewRoProp(tx)

type localizerLookup map[string]*i18n.Localizer

// Use, must be called by the client before any string data
// can be translated. If the client requests the default
// language, then only the language Tag needs to be provided.
// If the requested language is not the default and therefore
// requires translation from the translation file(s), then
// the client must provide the App and Path properties indicating
// how the l18n bundle is created.
func Use(options ...UseOptionFn) error {
	o := &UseOptions{}
	for _, fo := range options {
		fo(o)
	}
	li := newLanguageInfo(o)

	if !containsLanguage(li.Supported, o.Tag) {
		return NewLanguageNotAvailableNativeError(o.Tag)
	}

	tx = NewTranslator(li)
	TxRef = utils.NewRoProp(tx)

	if TxRef.IsNone() {
		return NewFailedToCreateTranslatorNativeError(o.Tag)
	}

	return nil
}

// ResetTx, do not use, required for unit testing only and is not considered
// part of the public api and may be removed without corresponding version
// number change.
func ResetTx() {
	// having to do this smells a bit, but required so unit tests can
	// remain isolated (this is why package globals are bad).
	tx = nil
	TxRef = utils.NewRoProp(tx)
}

func newLanguageInfo(o *UseOptions) *LanguageInfo {
	return &LanguageInfo{
		UseOptions: *o,
		Current:    o.Tag,
		Default:    DefaultLanguage.Get(),
		Supported: SupportedLanguages{
			DefaultLanguage.Get(),
			language.AmericanEnglish,
		},
	}
}

// Text is the function to use to obtain a string created from
// registered Localizers. The data parameter must be a go template
// defining the input parameters and the translatable message content.
func Text(data Localisable) string {
	return tx.localise(data)
}

// Translator provides the translation implementation used by the
// Text function
type Translator struct {
	mx              localizerMultiplexor
	LanguageInfoRef utils.RoProp[LanguageInfo]
}

// since extendio is not trying to provide foreign translations for any
// of its dependencies, we only need create a localizer for this module
// only (extendio). If we do need to provide these additional translations,
// then set _USE_MULTI to true and then provide additional localizers
// as indicated with the add method.
const _USE_MULTI = false

// NewTranslator creates a translator instance from the provided
// Localizers. If no foreign localizers are provided, then the
// Translator will be created with the single localizer which represents
// the client's package. If foreign localizers are present, then
// these are added as registered localizers.
func NewTranslator(li *LanguageInfo, foreigners ...*LocalizerInfo) *Translator {
	liRef := utils.NewRoProp(*li)

	factory := LocalizerFactory{
		provider: &translationProvider{
			languageInfoRef: liRef,
		},
	}
	// The native localizer represents the one that is used for this
	// module's translations requirements.
	//
	native := factory.New(li)

	mx := lo.TernaryF(_USE_MULTI,
		func() localizerMultiplexor {
			multi := multipleLocalizers{}

			if err := multi.add(&LocalizerInfo{
				SourceId:  EXTENDIO_SOURCE_ID,
				Localizer: native,
			}); err != nil {
				panic(NewFailedToCreateLocalizerNativeError(li.Current, EXTENDIO_SOURCE_ID))
			}

			for _, forloc := range foreigners {
				if err := multi.add(&LocalizerInfo{
					SourceId:  forloc.SourceId,
					Localizer: forloc.Localizer,
				}); err != nil {
					panic(NewFailedToCreateLocalizerNativeError(li.Current, EXTENDIO_SOURCE_ID))
				}
			}

			return &multi
		},
		func() localizerMultiplexor {
			return &singleLocalizer{
				localizer: native,
			}
		},
	)

	return &Translator{
		mx:              mx,
		LanguageInfoRef: liRef,
	}
}

func (t *Translator) localise(data Localisable) string {
	return t.mx.localise(data)
}

func containsLanguage(languages SupportedLanguages, tag language.Tag) bool {
	return lo.ContainsBy(languages, func(t language.Tag) bool {
		return t == tag
	})
}
