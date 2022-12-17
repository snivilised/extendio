package nav_test

import (
	"fmt"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// this should be a dot import:
	_ "github.com/snivilised/extendio/translate"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("ResumeSpawnStrategy", Ordered, func() {
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

	DescribeTable("spawn",
		func(entry *spawnTE) {
			_ = root
			_ = jroot
			_ = fromJsonPath
			Expect(true)

			restore := func(o *nav.TraverseOptions, active *nav.ActiveState) {
				// synthetic assignments: The client should not perform these
				// types of assignments. Only being done here for testing purposes
				// to avoid the need to create many restore files
				// (eg resume-state.json) for different test cases.
				//
				active.Root = path(root, entry.relative)
				active.NodePath = path(root, entry.active.resumeAtPath)
				active.Listen = entry.active.listenState
				o.Store.Subscription = entry.subscription
				//
				// end of synthetic assignments
			}

			info := &nav.NewResumerInfo{
				RestorePath: fromJsonPath,
				Restorer:    restore,
				Strategy:    nav.ResumeStrategySpawnEn,
			}
			result, err := nav.Resume(info)
			Expect(err).To(BeNil())
			Expect(result).ToNot(BeNil())

		},
		func(entry *spawnTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.message)
		},
	)
})
