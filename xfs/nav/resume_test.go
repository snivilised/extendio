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
			Skip("NOT-READY")
			restore := func(o *nav.TraverseOptions, as *nav.ActiveState) {
				GinkgoWriter.Printf("===> 🐚 restoring ...\n")
			}
			info := &nav.NewResumerInfo{
				RestorePath: fromJsonPath,
				Restorer:    restore,
				Strategy:    nav.ResumeStrategyFastwardEn,
			}
			result, err := nav.Resume(info)
			Expect(err).To(BeNil())
			Expect(result).ToNot(BeNil())
		})
	})
})
