package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/samber/lo"
	"github.com/snivilised/extendio/xfs/utils"
	"golang.org/x/text/language"
)

type Translator interface {
	Localise(data Localisable) string
	LanguageInfoRef() utils.RoProp[LanguageInfo]
}

var DefaultLanguage = utils.NewRoProp(language.BritishEnglish)
var tx Translator
var TxRef utils.RoProp[Translator] = utils.NewRoProp(tx)

type localizerContainer map[string]*i18n.Localizer

// Use, must be called by the client before any string data
// can be translated. If the client requests the default
// language, then only the language Tag needs to be provided.
// If the requested language is not the default and therefore
// requires translation from the translation file(s), then
// the client must provide the App and Path properties indicating
// how the l18n bundle is created.
// If the client just wishes to use the Default language, then Use
// can even be called without specifying the Tag and in this case
// the default language will be used. The client MUST call Use
// before using any functionality in this package.
func Use(options ...UseOptionFn) error {
	var err error
	o := &UseOptions{}
	o.DefaultIsAcceptable = true
	o.Tag = DefaultLanguage.Get()

	for _, fo := range options {
		fo(o)
	}
	lang := NewLanguageInfo(o)

	if !ContainsLanguage(lang.Supported, o.Tag) {
		if o.DefaultIsAcceptable {
			o.Tag = DefaultLanguage.Get()
			lang.Tag = o.Tag
		} else {
			err = NewFailedToCreateTranslatorNativeError(o.Tag)
		}
	}

	if err == nil {
		applyLanguage(lang)
	}

	return err
}

func verifyLanguage(lang *LanguageInfo) {
	if lang.From.Sources == nil {
		lang.From.Sources = make(TranslationFiles)
	}

	// By adding in the source for extendio, we relieve the client from having
	// to do this. After-all, it should be taken as read that if the client is
	// using extendio then the translations for extendio should be loaded,
	// otherwise extendio will not be able to convey these translations to the
	// client. The client app has to make sure that when their app is deployed,
	// the translations file(s) for extendio are named as 'extendio', as you
	// can see below, that that is the name assigned to the app name of the
	// source. There is little value in making this customisable as this would
	// just lead to confusion. If the client really wants to control the name
	// of the translation file for extendio, they can provide an override
	// 'Create' function on UseOptions.
	//
	if _, found := lang.From.Sources[SOURCE_ID]; !found {
		lang.From.Sources[SOURCE_ID] = TranslationSource{Name: "extendio"}
	}
}

func applyLanguage(lang *LanguageInfo) {
	verifyLanguage(lang)
	factory := lo.TernaryF(len(lang.From.Sources) > 1,
		func() TranslatorFactory {
			return &MultiTranslatorFactory{
				AbstractTranslatorFactory: AbstractTranslatorFactory{
					Create: lang.Create,
				},
			}
		},
		func() TranslatorFactory {
			return &SingularTranslatorFactory{
				AbstractTranslatorFactory: AbstractTranslatorFactory{
					Create: lang.Create,
				},
			}
		},
	)

	tx = factory.New(lang)
	TxRef = utils.NewRoProp(tx)
}

// ResetTx, do not use, required for unit testing only and is
// not considered part of the public api and may be removed without
// corresponding version number change.
func ResetTx() {
	// having to do this smells a bit, but required so unit tests can
	// remain isolated (this is why package globals are bad, but sometimes
	// unavoidable). This is all because we want to be able to call the Text
	// function easily. If we defined the Text function on an object, then that
	// would require passing that state around in many places, making the code
	// much more brittle and cumbersome to maintain.
	//
	tx = nil
	TxRef = utils.NewRoProp(tx)
}

// UseTx, do not use, required for unit testing only and is
// not considered part of the public api and may be removed without
// corresponding version number change.
func UseTx(with Translator, setters ...UseOptionFn) error {
	o := &UseOptions{}
	for _, fo := range setters {
		fo(o)
	}

	tx = with
	TxRef = utils.NewRoProp(tx)

	if TxRef.IsNone() {
		return NewFailedToCreateTranslatorNativeError(o.Tag)
	}

	return nil
}

// NewLanguageInfo gets a new instance of Language info from the use options
// specified. This is specific to extendio. Client applications should
// provide their own version that reflects their own defaults.
func NewLanguageInfo(o *UseOptions) *LanguageInfo {
	return &LanguageInfo{
		UseOptions: *o,
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
	return tx.Localise(data)
}

// i18nTranslator provides the translation implementation used by the
// Text function
type i18nTranslator struct {
	mx              localizerMultiplexor
	languageInfoRef utils.RoProp[LanguageInfo]
}

func (t *i18nTranslator) LanguageInfoRef() utils.RoProp[LanguageInfo] {
	return t.languageInfoRef
}

func (t *i18nTranslator) Localise(data Localisable) string {
	return t.mx.localise(data)
}

func ContainsLanguage(languages SupportedLanguages, tag language.Tag) bool {

	return lo.ContainsBy(languages, func(t language.Tag) bool {
		return t == tag
	})
}
