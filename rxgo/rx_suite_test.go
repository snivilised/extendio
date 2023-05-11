//nolint:gocritic // foo bar
package rxgo_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRx(t *testing.T) {
	// defer goleak.VerifyNone(t)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rx Suite")
}
