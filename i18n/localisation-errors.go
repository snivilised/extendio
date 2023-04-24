package i18n

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

// NB: These errors occur prior to or during the process of creating a localizer
// which by definition means translated content can't be served to the client using
// the requested locale and therefore have to be displayed untranslated.

// ❌ Could Not Find Localizer

// NewFailedToCreateLocalizerNativeError creates an untranslated error to
// indicate the Translator already contains a localizer for the source
// specified. (internal error)
func NewCouldNotFindLocalizerNativeError(sourceID string) error {
	return fmt.Errorf(
		"i18n: could not find localizer for source: '%v'", sourceID,
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

// ❌ Failed To Create Translator

// NewFailedToCreateTranslatorNativeError creates an untranslated error to
// indicate failure to create a Translator instance
func NewFailedToCreateTranslatorNativeError(tag language.Tag) error {
	return fmt.Errorf(
		"i18n: failed to create translator for language '%v'", tag,
	)
}
