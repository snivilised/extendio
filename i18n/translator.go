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

func applyLanguage(lang *LanguageInfo) {
	// since extendio is not trying to provide foreign translations for any
	// of its dependencies, we only need create a localizer for this module
	// (extendio). If we do need to provide these additional translations,
	// then use NewMultiTranslator instead and then provide additional
	// localizers.
	//
	factory := SingularTranslatorFactory{}

	tx = factory.New(lang)
	TxRef = utils.NewRoProp(tx)
}

// Deprecated: ResetTx, do not use, required for unit testing only and is
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

// Deprecated: UseTx, do not use, required for unit testing only and is
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
