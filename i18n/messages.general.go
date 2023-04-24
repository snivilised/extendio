package i18n

// 🧊 Internationalisation

// Internationalisation
type InternationalisationTemplData struct {
	ExtendioTemplData
}

func (td InternationalisationTemplData) Message() *Message {
	return &Message{
		ID:          "internationalisation.general.extendio",
		Description: "Internationalisation",
		Other:       "internationalisation",
	}
}

// 🧊 Localisation

// Internationalisation
type LocalisationTemplData struct {
	ExtendioTemplData
}

func (td LocalisationTemplData) Message() *Message {
	return &Message{
		ID:          "localisation.general.extendio",
		Description: "Localisation",
		Other:       "localisation",
	}
}
