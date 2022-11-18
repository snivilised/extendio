package nav_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestNav(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nav Suite")
}
