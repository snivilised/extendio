package i18n

type multiplexor struct {
}

func (mx *multiplexor) invoke(localizer *Localizer, data Localisable) string {
	return localizer.MustLocalize(&LocalizeConfig{
		DefaultMessage: data.Message(),
		TemplateData:   data,
	})
}

type multiContainer struct {
	multiplexor
	localizers localizerContainer
}

func (mc *multiContainer) localise(data Localisable) string {
	return mc.invoke(mc.find(data.SourceID()), data)
}

func (mc *multiContainer) add(info *LocalizerInfo) {
	if _, found := mc.localizers[info.sourceID]; found {
		return
	}

	mc.localizers[info.sourceID] = info.Localizer
}

func (mc *multiContainer) find(id string) *Localizer {
	if loc, found := mc.localizers[id]; found {
		return loc
	}

	panic(NewCouldNotFindLocalizerNativeError(id))
}
