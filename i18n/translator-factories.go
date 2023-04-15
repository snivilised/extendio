package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/samber/lo"
	"github.com/snivilised/extendio/xfs/utils"
)

// LocalizerCreatorFn represents the signature of the function can optionally
// provide to override how an i18n Localizer is created.
type LocalizerCreatorFn func(li *LanguageInfo, sourceID string) (*i18n.Localizer, error)

type AbstractTranslatorFactory struct {
	Create LocalizerCreatorFn
}

func (f *AbstractTranslatorFactory) setup(lang *LanguageInfo) {
	verifyLanguage(lang)

	if f.Create == nil {
		f.Create = createLocalizer
	}
}

// singularTranslatorFactory creates Translator with the single localizer which
// represents the client's package.
type singularTranslatorFactory struct {
	AbstractTranslatorFactory
}

func (f *singularTranslatorFactory) New(lang *LanguageInfo) Translator {
	f.setup(lang)

	sourceID := lo.Keys(lang.From.Sources)[0]

	liRef := utils.NewRoProp(*lang)
	native, err := f.Create(lang, sourceID)

	if err != nil {
		panic(err)
	}

	single := &singularContainer{
		localizer: native,
	}

	return &i18nTranslator{
		mx:              single,
		languageInfoRef: liRef,
	}
}

// multiTranslatorFactory creates a translator instance from the provided
// Localizers. If the client library needs to provide localizers for itself
// and at least 1 dependency, then they should use multiTranslatorFactory
// specify appropriate information concerning where to load the translation
// files from, otherwise SingularTranslatorFactory should be used.
//
// Note, in the case where a source client wants to provide a localizer
// for a language that one of ite dependencies does not support, then
// the translator should create the localizer based on its own default
// language, but we load the client provided translation file at the same
// name as the dependency would have created it for, then this file will
// be loaded as per usual.
type multiTranslatorFactory struct {
	AbstractTranslatorFactory
}

func (f *multiTranslatorFactory) New(lang *LanguageInfo) Translator {
	f.setup(lang)

	liRef := utils.NewRoProp(*lang)
	multi := &multiContainer{
		localizers: make(localizerContainer),
	}

	for id := range lang.From.Sources {
		localizer, err := f.Create(lang, id)

		if err != nil {
			panic(err)
		}

		multi.add(&LocalizerInfo{
			sourceID:  id,
			Localizer: localizer,
		})
	}

	return &i18nTranslator{
		mx:              multi,
		languageInfoRef: liRef,
	}
}
