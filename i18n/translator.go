package i18n

import (
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/snivilised/extendio/xfs/utils"
	"golang.org/x/text/language"
)

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
		// TODO: create a behavioural interface that denotes language created,
		// as per Dave Cheney's  recommendation
		//
		return fmt.Errorf("language '%v' not available", o.Tag)
	}

	tx = NewTranslator(li)
	TxRef = utils.NewRoProp(tx)

	if TxRef.IsNone() {
		return fmt.Errorf("failed to create translator for language '%v'", o.Tag)
	}

	return nil
}

func newLanguageInfo(o *UseOptions) *LanguageInfo {
	return &LanguageInfo{
		UseOptions: *o,
		Current:    o.Tag,
		Default:    language.BritishEnglish,
		Supported: SupportedLanguages{
			language.BritishEnglish,
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
				panic(errors.Wrapf(err, "failed to create localizer for language '%v', dependency: '%v'",
					li.Current, EXTENDIO_SOURCE_ID))
			}

			for _, forloc := range foreigners {
				if err := multi.add(&LocalizerInfo{
					SourceId:  forloc.SourceId,
					Localizer: forloc.Localizer,
				}); err != nil {
					panic(errors.Wrapf(err, "failed to create foreigner localizer for language '%v', dependency: '%v'",
						li.Current, forloc.SourceId))
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
