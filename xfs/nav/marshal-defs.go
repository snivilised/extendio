package nav

type PersistenceFormatEnum uint

const (
	PersistInUndefinedEn PersistenceFormatEnum = iota
	PersistInJSONEn
)

const (
	JSON_MARSHAL_NO_PREFIX      = ""
	JSON_MARSHAL_2SPACES_INDENT = "  "
)

type PersistFilterDef struct {
	Description string
	Source      string
	Scope       FilterScopeEnum
}

type PersistCompoundFilterDef struct {
	Description string
	Source      string
}

type PersistFilters struct {
	Current  *PersistFilterDef
	Children *PersistCompoundFilterDef
}

type PersistenceRestorer func(o *TraverseOptions)

type activeState struct {
	Root   string
	Listen ListeningState
}

type persistState struct {
	Store  *OptionsStore
	Active *activeState
}

type stateMarshaller interface {
	marshal(path string) error
	unmarshal(path string) error
}