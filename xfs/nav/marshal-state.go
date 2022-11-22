package nav

import (
	"encoding/json"
	"os"
)

type stateMarshallerJSON struct {
	o       *TraverseOptions
	ps      *persistState
	restore PersistenceRestorer
}

func (m *stateMarshallerJSON) marshal(path string) error {
	bytes, err := json.MarshalIndent(
		m.ps,
		JSON_MARSHAL_NO_PREFIX, JSON_MARSHAL_2SPACES_INDENT,
	)

	if err == nil {
		return writeBytes(bytes, path)
	}

	return err
}

func (m *stateMarshallerJSON) unmarshal(path string) error {
	if bytes, err := os.ReadFile(path); err == nil {
		m.o = GetDefaultOptions()
		m.ps = new(persistState)
		if err = json.Unmarshal(bytes, &m.ps); err == nil {
			m.o.Store = *m.ps.Store
			m.restore(m.o, m.ps.Active)
			return nil

		} else {
			return err
		}
	} else {
		return err
	}
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
