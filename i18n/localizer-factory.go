package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/snivilised/extendio/internal/lo"
	"github.com/snivilised/extendio/xfs/utils"
	"golang.org/x/text/language"
)

func createLocalizer(lang *LanguageInfo, sourceID string) (*Localizer, error) {
	bundle := i18n.NewBundle(lang.Tag)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	if lang.Tag != lang.Default {
		txSource := lang.From.Sources[sourceID]
		path := resolveBundlePath(lang, txSource)
		_, err := bundle.LoadMessageFile(path)

		if (err != nil) && (!lang.DefaultIsAcceptable) {
			return nil, NewCouldNotLoadTranslationsNativeError(lang.Tag, path, err)
		}
	}

	supported := lo.Map(lang.Supported, func(t language.Tag, _ int) string {
		return t.String()
	})

	return i18n.NewLocalizer(bundle, supported...), nil
}

func resolveBundlePath(lang *LanguageInfo, txSource TranslationSource) string {
	filename := lo.TernaryF(txSource.Name == "",
		func() string {
			return fmt.Sprintf("active.%v.json", lang.Tag)
		},
		func() string {
			return fmt.Sprintf("%v.active.%v.json", txSource.Name, lang.Tag)
		},
	)

	path := lo.Ternary(txSource.Path != "" && utils.FolderExists(txSource.Path),
		txSource.Path,
		lang.From.Path,
	)

	directory := lo.TernaryF(path != "" && utils.FolderExists(path),
		func() string {
			resolved, _ := filepath.Abs(path)
			return resolved
		},
		func() string {
			exe, _ := os.Executable()
			return filepath.Dir(exe)
		},
	)

	return filepath.Join(directory, filename)
}
