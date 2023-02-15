package nav_test

import (
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/extendio/translate"
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
			info := &nav.ResumerInfo{
				RestorePath: fromJsonPath,
				Restorer:    restore,
				Strategy:    nav.ResumeStrategyFastwardEn,
			}
			result, err := nav.ResumeLegacy(info)
			Expect(err).To(BeNil())
			Expect(result).ToNot(BeNil())
		})
	})

	Context("temp: spawn dead code fix lint errors", func() {
		It("should: prevent lint errors", func() {
			Skip("NOT-READY")
			restore := func(o *nav.TraverseOptions, as *nav.ActiveState) {
				GinkgoWriter.Printf("===> 🐚 restoring ...\n")

				o.Callback = nav.LabelledTraverseCallback{
					Label: "test spawn callback",
					Fn: func(item *nav.TraverseItem) *translate.LocalisableError {
						return nil
					},
				}
			}
			info := &nav.ResumerInfo{
				RestorePath: fromJsonPath,
				Restorer:    restore,
				Strategy:    nav.ResumeStrategyFastwardEn,
			}
			// panic: callback is not set
			result, err := nav.ResumeLegacy(info)
			Expect(err).To(BeNil())
			Expect(result).ToNot(BeNil())
		})
	})
})
