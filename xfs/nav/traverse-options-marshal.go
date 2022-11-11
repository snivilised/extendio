package nav

import (
	"encoding/json"
	"os"
)

type PersistenceFormatEnum uint

const (
	PersistInUndefinedEn PersistenceFormatEnum = iota
	PersistInJSONEn
)

const (
	NO_JSON_MARSHAL_PREFIX        = ""
	TWO_SPACE_JSON_MARSHAL_INDENT = "  "
)

type PersistCurrentFilterDef struct {
	Object      string
	Description string
	Source      string
	Scope       FilterScopeEnum
}

type PersistChildrenFilterDef struct {
	Object      string
	Description string
	Source      string
}

type PersistFilters struct {
	Current  *PersistCurrentFilterDef
	Children *PersistChildrenFilterDef
}

// AdaptedPersistOptions represents that part of the TraverseOptions that can't be serialised
// as is and therefore need a persistence friendly representation. Typically, function
// parameters and interfaces can't be representing in persisted form. The client is expected
// to restore other properties manually like any custom Hooks via the RestoreFN on
// OptionMarshallerJSON
type AdaptedPersistOptions struct {
	Filters PersistFilters
}

type TraverseOptionsAsJSON struct {
	Subscription   TraverseSubscription
	DoExtend       bool
	WithMetrics    bool
	Behaviours     NavigationBehaviours
	Persist        PersistOptions
	AdaptedOptions AdaptedPersistOptions
}

type OptionMarshaller interface {
	marshal(path string) error
	unmarshal(path string) error
}

func (o *TraverseOptions) Marshal(path string) error {
	marshaller := newOptionMarshaller(o)
	return marshaller.marshal(path)
}

func (o *TraverseOptions) Unmarshal(path string) error {
	marshaller := newOptionMarshaller(o)
	return marshaller.unmarshal(path)
}

type OptionsRestorer func(o *TraverseOptions)

type OptionMarshallerJSON struct {
	Options     *TraverseOptions
	OptionsJSON *TraverseOptionsAsJSON
}

func (m *OptionMarshallerJSON) marshal(path string) error {
	jo := m.toJSON()
	bytes, err := json.MarshalIndent(
		jo,
		NO_JSON_MARSHAL_PREFIX, TWO_SPACE_JSON_MARSHAL_INDENT,
	)

	if err == nil {
		return writeBytes(bytes, path)
	}

	return err
}

func (m *OptionMarshallerJSON) unmarshal(path string) error {

	if bytes, err := os.ReadFile(path); err != nil {
		return err
	} else {
		m.OptionsJSON = new(TraverseOptionsAsJSON)
		if err = json.Unmarshal(bytes, m.OptionsJSON); err != nil {
			return err
		} else {
			m.restore()
			return nil
		}
	}
}

func (m *OptionMarshallerJSON) restore() {
	defOptions := GetDefaultOptions()
	m.Options.Hooks = defOptions.Hooks
	m.Options.Notify = defOptions.Notify

	m.Options.Persist.Restorer(m.Options)
	m.Options.Subscription = m.OptionsJSON.Subscription
	m.Options.DoExtend = m.OptionsJSON.DoExtend
	m.Options.WithMetrics = m.OptionsJSON.WithMetrics
	m.Options.Behaviours = m.OptionsJSON.Behaviours
	m.Options.Persist = m.OptionsJSON.Persist
}

func (m *OptionMarshallerJSON) toJSON() *TraverseOptionsAsJSON {
	result := &TraverseOptionsAsJSON{
		Subscription: m.Options.Subscription,
		DoExtend:     m.Options.DoExtend,
		WithMetrics:  m.Options.WithMetrics,
		Behaviours:   m.Options.Behaviours,
		Persist:      m.Options.Persist,
	}

	if m.Options.Filters.Current != nil {
		result.AdaptedOptions.Filters.Current = &PersistCurrentFilterDef{
			Description: m.Options.Filters.Current.Description(),
			Source:      m.Options.Filters.Current.Source(),
			Scope:       m.Options.Filters.Current.Scope(),
		}
	}

	if m.Options.Filters.Children != nil {
		result.AdaptedOptions.Filters.Children = &PersistChildrenFilterDef{
			Description: m.Options.Filters.Children.Description(),
			Source:      m.Options.Filters.Children.Source(),
		}
	}

	return result
}

func writeBytes(bytes []byte, path string) error {
	if file, err := os.Create(path); err == nil {
		defer file.Close()

		_, e := file.Write(bytes)
		return e
	} else {
		return nil
	}
}
