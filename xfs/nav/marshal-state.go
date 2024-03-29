package nav

import (
	"encoding/json"
	"os"

	"github.com/snivilised/extendio/i18n"
)

type stateMarshallerJSON struct {
	o       *TraverseOptions
	ps      *persistState
	restore PersistenceRestorer
}

func (m *stateMarshallerJSON) marshal(path string) error {
	bytes, err := json.MarshalIndent(
		m.ps,
		JSONMarshallNoPrefix, JSONMarshall2SpacesIndent,
	)

	if err == nil {
		return writeBytes(bytes, path)
	}

	return err
}

func (m *stateMarshallerJSON) unmarshal(path string) error {
	bytes, err := os.ReadFile(path)

	if err == nil {
		m.o = GetDefaultOptions()
		m.ps = new(persistState)

		err = json.Unmarshal(bytes, &m.ps)

		if err == nil {
			m.o.Store = *m.ps.Store
			m.restore(m.o, m.ps.Active)
			m.o.afterUserOptions()
			m.validate()
		}
	}

	return err
}

func (m *stateMarshallerJSON) validate() {
	if m.o.Callback.Fn == nil {
		panic(i18n.NewMissingCallbackError())
	}
}

func writeBytes(bytes []byte, path string) error {
	if file, err := os.Create(path); err == nil {
		defer file.Close()

		_, e := file.Write(bytes)

		return e
	}

	return nil
}
