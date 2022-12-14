package nav_test

import (
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/snivilised/extendio/translate"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("NavigatorResume", Ordered, func() {
	var (
		root         string
		jroot        string
		fromJsonPath string
	)

	BeforeAll(func() {
		root = origin()
		jroot = joinCwd("Test", "json")
		fromJsonPath = strings.Join(
			[]string{jroot, "resume-state.json"},
			string(filepath.Separator),
		)
	})

	Context("client not previously using listener", func() {
		It("should: blah", func() {
			Skip("NOT-READY")
			restore := func(o *nav.TraverseOptions, active *nav.ActiveState) {
				GinkgoWriter.Printf("---> marshaller ...\n")

				active.Root = path(root, "RETRO-WAVE")
				active.NodePath = path(root, "RETRO-WAVE/Electric Youth")
				o.Notify.OnBegin = begin("🛡️")
				// subscribe-any
				o.Callback = nav.LabelledTraverseCallback{
					Label: "test resume callback",
					Fn: func(item *nav.TraverseItem) *LocalisableError {
						return nil
					},
				}
			}

			resumeInfo := &nav.NewResumerInfo{
				RestorePath: fromJsonPath,
				Restorer:    restore,
				Strategy:    nav.ResumeStrategyFastwardEn,
			}
			Expect(resumeInfo).ToNot(BeNil())

			_, _ = nav.Resume(resumeInfo)
		})
	})
})
