package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/snivilised/extendio/xfs/utils"
	"golang.org/x/text/language"
)

type multiplexor struct {
}

func (mx *multiplexor) invoke(localizer *i18n.Localizer, data Localisable) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: data.Message(),
		TemplateData:   data,
	})
}

type singleLocalizer struct {
	multiplexor
	localizer *i18n.Localizer
}

func (mx *singleLocalizer) localise(data Localisable) string {
	return mx.invoke(mx.localizer, data)
}

type multipleLocalizers struct {
	multiplexor
	lookup localizerLookup
}

func (mx *multipleLocalizers) localise(data Localisable) string {
	return mx.invoke(mx.find(data.SourceId()), data)
}

func (mx *multipleLocalizers) add(info *LocalizerInfo) error {
	if _, found := mx.lookup[info.SourceId]; found {
		return NewLocalizerAlreadyExistsNativeError(info.SourceId)
	}
	mx.lookup[info.SourceId] = info.Localizer
	return nil
}

func (mx *multipleLocalizers) find(id string) *i18n.Localizer {
	if loc, found := mx.lookup[id]; found {
		return loc
	}

	panic(NewCouldNotFindLocalizerNativeError(id))
}

type translationProvider struct {
	languageInfoRef utils.RoProp[LanguageInfo]
}

func (p *translationProvider) Query(tag language.Tag) bool {
	return containsLanguage(p.languageInfoRef.Get().Supported, tag)
}

func (p *translationProvider) Create(li *LanguageInfo) *i18n.Localizer {
	// create foreign localizer for the SourceId representing
	// the dependency not supporting the requested language.
	//
	return nil
}
