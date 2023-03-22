package i18n_test

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const (
	GRAFFICO_SOURCE_ID = "github.com/snivilised/graffico"
)

type GrafficoData struct{}

func (td GrafficoData) SourceId() string {
	return GRAFFICO_SOURCE_ID
}

// 🧊 Pavement Graffiti Report

// PavementGraffitiReportTemplData
type PavementGraffitiReportTemplData struct {
	GrafficoData
	Primary string
}

func (td PavementGraffitiReportTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "pavement-graffiti-report.graffico.unit-test",
		Description: "Report of graffiti found on a pavement",
		Other:       "Found graffiti on pavement; primary colour: '{{.Primary}}'",
	}
}

// ☢️ Wrong Source Id

// WrongSourceIdTemplData
type WrongSourceIdTemplData struct {
	GrafficoData
}

func (td WrongSourceIdTemplData) SourceId() string {
	return "FOO-BAR"
}

func (td WrongSourceIdTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "wrong-source-id.graffico.unit-test",
		Description: "Incorrect Source Id for which doesn't match the one n the localizer",
		Other:       "Message with wrong id",
	}
}
