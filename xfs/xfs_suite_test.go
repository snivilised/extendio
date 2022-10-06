package xfs_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestXfs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Xfs Suite")
}
