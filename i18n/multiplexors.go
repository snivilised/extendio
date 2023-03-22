package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
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
	if _, found := mx.lookup[info.sourceId]; found {
		return NewLocalizerAlreadyExistsNativeError(info.sourceId)
	}
	mx.lookup[info.sourceId] = info.Localizer
	return nil
}

func (mx *multipleLocalizers) find(id string) *i18n.Localizer {
	if loc, found := mx.lookup[id]; found {
		return loc
	}

	panic(NewCouldNotFindLocalizerNativeError(id))
}
