package nav_test

import (
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("ResumeFastward", Ordered, func() {

	var (
		root         string
		jroot        string
		fromJsonPath string
	)

	BeforeAll(func() {
		root = origin()
		jroot = joinCwd("Test", "json")
		fromJsonPath = strings.Join(
			[]string{jroot, "fastward-state.json"},
			string(filepath.Separator),
		)
	})

	Context("resume from default listen state", func() {
		It("should: ---", func() {
			Skip("NOT-READY")
			restore := func(o *nav.TraverseOptions, active *nav.ActiveState) {

				relative := "RETRO-WAVE"
				active.Root = path(root, relative)
				active.Listen = nav.ListenDefault

				o.Notify.OnBegin = begin("🛡️")
				GinkgoWriter.Printf("===> 🐚 restoring ...\n")
			}
			info := nav.NewResumerInfo{
				Path:     fromJsonPath,
				Restore:  restore,
				Strategy: nav.ResumeStrategyFastwardEn,
			}
			resumer, err := nav.NewResumer(info)
			Expect(err).To(BeNil())

			resumer.Walk()
		})
	})

	Context("resume from active listen state", func() {
		It("should: ---", func() {
			Skip("NOT-READY")
			restore := func(o *nav.TraverseOptions, active *nav.ActiveState) {

				relative := "RETRO-WAVE"
				active.Root = path(root, relative)
				active.Listen = nav.ListenActive

				// o.Listen = nav.ListenOptions{
				// 	Start: ,
				// }
				o.Notify.OnBegin = begin("🛡️")
				GinkgoWriter.Printf("===> 🐚 restoring ...\n")
			}
			info := nav.NewResumerInfo{
				Path:     fromJsonPath,
				Restore:  restore,
				Strategy: nav.ResumeStrategyFastwardEn,
			}
			resumer, err := nav.NewResumer(info)
			Expect(err).To(BeNil())

			resumer.Walk()
		})
	})
})
