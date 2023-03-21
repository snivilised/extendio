package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// EXTENDIO_SOURCE_ID the id that represents this module. If client want
// to provides translations for languages that extendio does not, then
// the localizer the create created for this purpose should use this
// SourceId. So whenever the Text function is used on templates defined
// inside this module, the translation process is directed to use the
// correct i18n.Localizer (identified by the SourceId). The Source is
// statically defined for all templates defined in extendio.
const EXTENDIO_SOURCE_ID = "github.com/snivilised/extendio"

type Localisable interface {
	Message() *i18n.Message
	SourceId() string
}

type SupportedLanguages []language.Tag

// UseOptions the options provided to the Use function
type UseOptions struct {
	// Tag sets the language to use
	//
	Tag language.Tag

	// Name, name which forms part of the language filename, typically
	// the name of the client app
	//
	Name string

	// Path denoting where to load language file from, defaults to exe location
	//
	Path string
}

// LanguageInfo information pertaining to setting language. Auto detection
// is not supported. Any executable that supports i18n, should perform
// auto detection and then invoke Use, with the detected language tag

type LanguageInfo struct {
	UseOptions

	// This would typically be the name of a package that is the source
	// of string messages that are to be translated. Actually, we could use
	// the top level url of the package by convention, as that is unique.
	// So extendio would use "github.com/snivilised/extendio"
	//
	SourceId string

	// Default language reflects the base language. If all else fails, messages will
	// be in this language. It is fixed at BritishEnglish reflecting the language this
	// package is written in.
	//
	Default language.Tag

	// Current reflects the language currently in force. Will by default be the detected
	// language. Client can change this with the UseTag function.
	//
	Current language.Tag

	// Supported indicates the list of languages for which translations are available.
	//
	Supported SupportedLanguages
}

// UseOptionFn functional options function required by Use.
type UseOptionFn func(*UseOptions)

// LocalizerProvider
type LocalizerProvider interface {
	Query(tag language.Tag) bool
	Create(li *LanguageInfo) *i18n.Localizer
}

type localizerMultiplexor interface {
	localise(data Localisable) string
}

// LocalizerInfo
type LocalizerInfo struct {
	// SourceId identifiers the module to provide translations for
	SourceId string

	// Localizer created by the client
	Localizer *i18n.Localizer
}
