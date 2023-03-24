package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/samber/lo"
	"github.com/snivilised/extendio/xfs/utils"
)

// LocalizerCreatorFn represents the signature of the function can optionally
// provide to override how an i18n Localizer is created.
type LocalizerCreatorFn func(li *LanguageInfo, sourceId string) *i18n.Localizer

type TranslatorFactory struct {
	Create LocalizerCreatorFn
}

func (f *TranslatorFactory) setup(lang *LanguageInfo) {
	if lang.From.Sources == nil {
		lang.From.Sources = TranslationFiles{
			SOURCE_ID: TranslationSource{},
		}
	}

	if f.Create == nil {
		f.Create = createLocalizer
	}
}

// SingularTranslatorFactory creates Translator with the single localizer which
// represents the client's package.
type SingularTranslatorFactory struct {
	TranslatorFactory
}

func (f *SingularTranslatorFactory) New(lang *LanguageInfo) Translator {
	f.setup(lang)

	count := len(lang.From.Sources)
	if (lang.From.Sources == nil) || (count == 0) {
		panic(NoSourcesSpecifiedNativeError())
	}

	if len(lang.From.Sources) > 1 {
		panic(MultipleSourcesSpecifiedForSingularTranslatorNativeError(count))
	}
	sourceId := lo.Keys(lang.From.Sources)[0]

	liRef := utils.NewRoProp(*lang)
	native := f.Create(lang, sourceId)
	single := &singularContainer{
		localizer: native,
	}

	return &i18nTranslator{
		mx:              single,
		languageInfoRef: liRef,
	}
}

// MultiTranslatorFactory creates a translator instance from the provided
// Localizers. If the client library needs to provide localizers for itself
// and at least 1 dependency, then they should use MultiTranslatorFactory
// specify appropriate information concerning where to load the translation
// files from, otherwise SingularTranslatorFactory should be used.
//
// Note, in the case where a source client wants to provide a localizer
// for a language that one of ite dependencies does not support, then
// the translator should create the localizer based on its own default
// language, but we load the client provided translation file at the same
// name as the dependency would have created it for, then this file will
// be loaded as per usual.
type MultiTranslatorFactory struct {
	TranslatorFactory
}

func (f *MultiTranslatorFactory) New(lang *LanguageInfo) Translator {
	f.setup(lang)

	liRef := utils.NewRoProp(*lang)
	multi := &multiContainer{
		localizers: make(localizerContainer),
	}

	count := len(lang.From.Sources)
	if len(lang.From.Sources) < 2 {
		panic(InsufficientSourcesSpecifiedForMultiTranslatorNativeError(count))
	}

	for id := range lang.From.Sources {
		localizer := f.Create(lang, id)

		err := multi.add(&LocalizerInfo{
			sourceId:  id,
			Localizer: localizer,
		})

		if err != nil {
			panic(err)
		}
	}

	return &i18nTranslator{
		mx:              multi,
		languageInfoRef: liRef,
	}
}
