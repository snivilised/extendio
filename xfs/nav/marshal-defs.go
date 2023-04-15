package nav

type PersistenceFormatEnum uint

const (
	PersistInUndefinedEn PersistenceFormatEnum = iota
	PersistInJSONEn
)

const (
	JSONMarshallNoPrefix      = ""
	JSONMarshall2SpacesIndent = "  "
)

type PersistFilterDef struct {
	Description string
	Source      string
	Scope       FilterScopeBiEnum
}

type PersistCompoundFilterDef struct {
	Description string
	Source      string
}

type PersistFilters struct {
	Node     *PersistFilterDef
	Children *PersistCompoundFilterDef
}

type PersistenceRestorer func(o *TraverseOptions, active *ActiveState)

type ActiveState struct {
	Root     string
	Listen   ListeningState
	NodePath string
	Depth    int
	Metrics  *MetricCollection
}

type persistState struct {
	Store  *OptionsStore
	Active *ActiveState
}

type stateMarshaller interface {
	marshal(path string) error
	unmarshal(path string) error
}
