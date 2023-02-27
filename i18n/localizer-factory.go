package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/samber/lo"
	"golang.org/x/text/language"
)

type LocalizerFactory struct {
	provider CreateLocaliser
}

func (f LocalizerFactory) New(li *LanguageInfo) *i18n.Localizer {

	return lo.TernaryF(f.provider.Query(li.Current),
		func() *i18n.Localizer {
			return f.create(li)
		},
		func() *i18n.Localizer {
			return f.provider.Create(li)
		},
	)
}

func (f LocalizerFactory) create(li *LanguageInfo) *i18n.Localizer {
	bundle := i18n.NewBundle(li.Current)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	if li.Current != li.Default {
		// TODO: consider changing this, too rigid, active may not be liked by client,
		// they should be able to control the name of the translation file.
		//
		filename := fmt.Sprintf("%v.active.%v.json", li.Name, li.Current)
		resolved, _ := filepath.Abs(li.Path)

		directory := lo.TernaryF(li.Path != "",
			func() string {
				return resolved
			},
			func() string {
				exe, _ := os.Executable()
				return filepath.Dir(exe)
			},
		)
		path := filepath.Join(directory, filename)
		_, err := bundle.LoadMessageFile(path)

		if err != nil {
			// Since, translations failed to load, we will never be in a situation where
			// this error message is able to be generated in translated form, so
			// we are forced to generate an error message in the default language.
			//
			panic(NewCouldNotLoadTranslationsNativeError(li.Current, path, err))
		}
	}

	supported := lo.Map(li.Supported, func(t language.Tag, _ int) string {
		return t.String()
	})

	return i18n.NewLocalizer(bundle, supported...)
}
