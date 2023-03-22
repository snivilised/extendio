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

func createLocalizer(lang *LanguageInfo, sourceId string) *i18n.Localizer {
	bundle := i18n.NewBundle(lang.Tag)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	if lang.Tag != lang.Default {
		name := lang.From.Sources[sourceId].Name
		path := resolveBundlePath(lang, name)
		_, err := bundle.LoadMessageFile(path)

		if err != nil {
			panic(NewCouldNotLoadTranslationsNativeError(lang.Tag, path, err))
		}
	}

	supported := lo.Map(lang.Supported, func(t language.Tag, _ int) string {
		return t.String()
	})

	return i18n.NewLocalizer(bundle, supported...)
}

func resolveBundlePath(lang *LanguageInfo, dependencyName string) string {
	filename := lo.TernaryF(dependencyName == "",
		func() string {
			return fmt.Sprintf("active.%v.json", lang.Tag)
		},
		func() string {
			return fmt.Sprintf("%v.active.%v.json", dependencyName, lang.Tag)
		},
	)

	resolved, _ := filepath.Abs(lang.From.Path)

	directory := lo.TernaryF(lang.From.Path != "",
		func() string {
			return resolved
		},
		func() string {
			exe, _ := os.Executable()
			return filepath.Dir(exe)
		},
	)
	return filepath.Join(directory, filename)
}
