package i18n

import "github.com/nicksnyder/go-i18n/v2/i18n"

// ❌ Internationalisation

// Internationalisation
type InternationalisationTemplData struct {
	ExtendioTemplData
}

func (td InternationalisationTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "internationalisation.general.extendio",
		Description: "Internationalisation",
		Other:       "internationalisation",
	}
}

// ❌ Localisation

// Internationalisation
type LocalisationTemplData struct {
	ExtendioTemplData
}

func (td LocalisationTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "localisation.general.extendio",
		Description: "Localisation",
		Other:       "localisation",
	}
}
