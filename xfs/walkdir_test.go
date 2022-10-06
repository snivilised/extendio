package xfs_test

import (
	"path/filepath"

	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/snivilised/extendio/xfs"
)

// new function to create will be Traverse
// TRAP CTRL-C: https://nathanleclaire.com/blog/2014/08/24/handling-ctrl-c-interrupt-signal-in-golang-programs/

var _ = Describe("WalkDir", func() {

	Context("MUSICO", func() {
		It("🧪 should: walk", func() {
			if current, err := os.Getwd(); err == nil {
				parent, _ := filepath.Split(current)
				root := filepath.Join(parent, "Test", "data", "MUSICO")

				GinkgoWriter.Printf("---> 🔰 ROOT-PATH: '%v' ...\n", root)
				xfs.Walkover(root)
			}
			Expect(true)
		})
	})
})
