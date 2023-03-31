package i18n

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

// NB: These errors occur prior to or during the process of creating a localizer
// which by definition means translated content can't be served to the client using
// the requested locale and therefore have to be displayed untranslated.

// ❌ CouldNotFindLocalizer

// NewFailedToCreateLocalizerNativeError creates an untranslated error to
// indicate the Translator already contains a localizer for the source
// specified. (internal error)
func NewCouldNotFindLocalizerNativeError(sourceId string) error {
	return fmt.Errorf(
		"i18n: could not find localizer for source: '%v'", sourceId,
	)
}

// ❌ Could Not Load Translations

// NewCouldNotLoadTranslationsNativeError creates an untranslated error to
// indicate translations file could not be loaded
func NewCouldNotLoadTranslationsNativeError(tag language.Tag, path string, reason error) error {
	return errors.Wrapf(
		reason, "i18n: could not load translations for '%v', from: '%v'", tag, path,
	)
}

// ❌ FailedToCreateLocalizer

// NewFailedToCreateLocalizerNativeError creates an untranslated error to
// indicate failure to create a localizer instance
func NewFailedToCreateLocalizerNativeError(tag language.Tag, sourceId string) error {
	return fmt.Errorf(
		"i18n: failed to create localizer for Language '%v', dependency: '%v'", tag, sourceId,
	)
}

// ❌ FailedToCreateTranslator

// NewFailedToCreateTranslatorNativeError creates an untranslated error to
// indicate failure to create a Translator instance
func NewFailedToCreateTranslatorNativeError(tag language.Tag) error {
	return fmt.Errorf(
		"i18n: failed to create translator for language '%v'", tag,
	)
}

// ❌ LanguageNotAvailable

// NewLanguageNotAvailableNativeError creates an untranslated error to indicate
// the requested language s not available
func NewLanguageNotAvailableNativeError(tag language.Tag) error {
	return fmt.Errorf(
		"i18n: language '%v' not available", tag,
	)
}

// ❌ LocalizerAlreadyExists

// NewFailedToCreateLocalizerNativeError creates an untranslated error to
// indicate the Translator already contains a localizer for the source
// specified. (internal error)

func NewLocalizerAlreadyExistsNativeError(sourceId string) error {
	return fmt.Errorf(
		"i18n: localizer already exists for source: '%v'", sourceId,
	)
}

// ❌ MultipleSourcesSpecifiedForSingularTranslator
func MultipleSourcesSpecifiedForSingularTranslatorNativeError(count int) error {
	return fmt.Errorf(
		"i18n: multiple sources (%v) have been specified for SingularTranslator", count,
	)
}

// ❌ InsufficientSourcesSpecifiedForSingularTranslator
func InsufficientSourcesSpecifiedForMultiTranslatorNativeError(count int) error {
	return fmt.Errorf(
		"i18n: insufficient sources (%v) have been specified for MultiTranslator", count,
	)
}
