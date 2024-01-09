package matchers

import (
	"fmt"

	"github.com/onsi/gomega/types"
	"github.com/snivilised/extendio/xfs/storage"
)

type PathExistsMatcher struct {
	vfs interface{}
}

type AsDirectory string
type AsFile string

func ExistInFS(fs interface{}) types.GomegaMatcher {
	return &PathExistsMatcher{
		vfs: fs,
	}
}

func (m *PathExistsMatcher) Match(actual interface{}) (bool, error) {
	vfs, fileSystemOK := m.vfs.(storage.VirtualFS)
	if !fileSystemOK {
		return false, fmt.Errorf("❌ matcher expected a VirtualFS instance (%T)", vfs)
	}

	if actualPath, dirOK := actual.(AsDirectory); dirOK {
		return vfs.DirectoryExists(string(actualPath)), nil
	}

	if actualPath, fileOK := actual.(AsFile); fileOK {
		return vfs.FileExists(string(actualPath)), nil
	}

	return false, fmt.Errorf("❌ matcher expected an AsDirectory or AsFile instance (%T)", actual)
}

func (m *PathExistsMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("🔥 Expected\n\t%v\npath to exist", actual)
}

func (m *PathExistsMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("🔥 Expected\n\t%v\npath NOT to exist\n", actual)
}
