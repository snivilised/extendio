package nav_test

import (
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("Resume", Ordered, func() {

	var (
		jroot        string
		fromJsonPath string
	)

	BeforeAll(func() {
		jroot = joinCwd("Test", "json")
		fromJsonPath = strings.Join(
			[]string{jroot, "persisted-state.json"},
			string(filepath.Separator),
		)
	})

	Context("given: existing persisted state", func() {
		It("should: create resumer", func() {
			restore := func(o *nav.TraverseOptions) {
				GinkgoWriter.Printf("===> 🐚 restoring ...\n")
			}
			info := nav.NewResumerInfo{
				Path:     fromJsonPath,
				Restore:  restore,
				Strategy: nav.ResumeStrategyFastwardEn,
			}
			_, err := nav.NewResumer(info)
			Expect(err).To(BeNil())
		})
	})
})
