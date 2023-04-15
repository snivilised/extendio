package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// ExtendioSourceID the id that represents this module. If client want
// to provides translations for languages that extendio does not, then
// the localizer the create created for this purpose should use this
// SourceID. So whenever the Text function is used on templates defined
// inside this module, the translation process is directed to use the
// correct i18n.Localizer (identified by the SourceID). The Source is
// statically defined for all templates defined in extendio.
const ExtendioSourceID = "github.com/snivilised/extendio"

type Localisable interface {
	Message() *i18n.Message
	SourceID() string
}

type SupportedLanguages []language.Tag

// LoadFrom denotes where to load the translation file from
type LoadFrom struct {
	// Path denoting where to load language file from, defaults to exe location
	//
	Path string

	// Sources are the translation files that need to be loaded. They represent
	// the client app/library its dependencies.
	//
	// The source id would typically be the name of a package that is the source
	// of string messages that are to be translated. Actually, we could use
	// the top level url of the package by convention, as that is unique.
	// So extendio would use "github.com/snivilised/extendio" but clients
	// are free to use whatever naming scheme they want to use for their own
	// dependencies.
	//
	Sources TranslationFiles
}

// AppendSources is designed to be used by dependencies that need to add
// their own dependencies to the Sources that are to be used to create the
// translator.
func (lf *LoadFrom) AppendSources(appendFiles *TranslationFiles) []string {
	duplicates := []string{}

	for id, source := range *appendFiles {
		if _, found := lf.Sources[id]; found {
			duplicates = append(duplicates, source.Name)
		} else {
			lf.Sources[id] = source
		}
	}

	if len(duplicates) > 0 {
		return duplicates
	}

	return nil
}

type TranslationSource struct {
	// Name of dependency's translation file
	Name string
}

// TranslationFiles maps a source id to a TranslationSource
type TranslationFiles map[string]TranslationSource

// UseOptions the options provided to the Use function
type UseOptions struct {
	// Tag sets the language to use
	//
	Tag language.Tag

	// From denotes where to load the translation file from
	//
	From LoadFrom

	// DefaultIsAcceptable controls whether an error is returned if the
	// request language is not available. By default DefaultIsAcceptable
	// is true so that the application continues in the default language
	// even if the requested language is not available.
	//
	DefaultIsAcceptable bool

	// Create allows the client to  override the default function to create
	// the i18n Localizer(s) (1 per language).
	//
	Create LocalizerCreatorFn

	// Custom set-able by the client for what ever purpose is required.
	//
	Custom any
}

// LanguageInfo information pertaining to setting language. Auto detection
// is not supported. Any executable that supports i18n, should perform
// auto detection and then invoke Use, with the detected language tag

type LanguageInfo struct {
	UseOptions

	// Default language reflects the base language. If all else fails, messages will
	// be in this language. It is fixed at BritishEnglish reflecting the language this
	// package is written in.
	//
	Default language.Tag

	// Supported indicates the list of languages for which translations are available.
	//
	Supported SupportedLanguages
}

// UseOptionFn functional options function required by Use.
type UseOptionFn func(*UseOptions)

type localizerMultiplexor interface {
	localise(data Localisable) string
}

// LocalizerInfo
type LocalizerInfo struct {
	// Localizer by default created internally, but can be overridden by
	// the client if they provide a create function to the Translator Factory
	// (either SingularTranslatorFactory or MultiTranslatorFactory)
	//
	Localizer *i18n.Localizer

	sourceID string
}

// TranslatorFactory
type TranslatorFactory interface {
	New(lang *LanguageInfo) Translator
}
