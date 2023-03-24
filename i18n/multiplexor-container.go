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

type singularContainer struct {
	multiplexor
	localizer *i18n.Localizer
}

func (sc *singularContainer) localise(data Localisable) string {
	return sc.invoke(sc.localizer, data)
}

type multiContainer struct {
	multiplexor
	localizers localizerContainer
}

func (mc *multiContainer) localise(data Localisable) string {
	return mc.invoke(mc.find(data.SourceId()), data)
}

func (mc *multiContainer) add(info *LocalizerInfo) error {
	if _, found := mc.localizers[info.sourceId]; found {
		return NewLocalizerAlreadyExistsNativeError(info.sourceId)
	}
	mc.localizers[info.sourceId] = info.Localizer
	return nil
}

func (mc *multiContainer) find(id string) *i18n.Localizer {
	if loc, found := mc.localizers[id]; found {
		return loc
	}

	panic(NewCouldNotFindLocalizerNativeError(id))
}
