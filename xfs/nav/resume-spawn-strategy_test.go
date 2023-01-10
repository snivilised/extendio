package nav_test

import (
	"fmt"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"

	. "github.com/snivilised/extendio/translate"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("ResumeSpawnStrategy", Ordered, func() {
	var (
		root         string
		jroot        string
		fromJsonPath string
		prohibited   map[string]string
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
			recording := recordingMap{}
			profile, ok := profiles[entry.profile]
			if !ok {
				Fail(fmt.Sprintf("bad test, missing profile for '%v'", entry.profile))
			}

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

				if !profile.filtered {
					o.Store.FilterDefs = nil
				}
				//
				// end of synthetic assignments

				o.Notify.OnBegin = begin("🧲")
				GinkgoWriter.Printf("===> 🐚 restoring ...\n")

				o.Callback = nav.LabelledTraverseCallback{
					Label: "unit test callback for resume",
					Fn: func(item *nav.TraverseItem) *LocalisableError {
						depth := lo.TernaryF(o.Store.DoExtend,
							func() uint { return item.Extension.Depth },
							func() uint { return 9999 },
						)
						GinkgoWriter.Printf(
							"---> 🍤 SPAWN: (depth:%v) '%v'\n", depth, item.Path,
						)

						segments := strings.Split(item.Path, string(filepath.Separator))
						last := segments[len(segments)-1]
						if _, found := prohibited[last]; found {
							Fail(fmt.Sprintf("item: '%v' should have been fast forwarded over", item.Path))
						}
						recording[item.Extension.Name] = len(item.Children)
						return nil
					},
				}
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
		// A note about the profiles: since the expected results for the spawn strategy is
		// exactly the same as that for the fastward strategy, the profiles can be shared
		// between the 2 and therefore do not need to be re-defined.
		//

		// === Listening (uni/folder/file) (pend/active)
		//

		Entry(nil, &spawnTE{
			naviTE: naviTE{
				message:      "universal(spawn): listen pending",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
			},
			active: activeTE{
				resumeAtPath: RESUME_AT_TEENAGE_COLOR,
				listenState:  nav.ListenPending,
			},
			resumeAt: START_AT_ELECTRIC_YOUTH,
			profile:  "-> universal(pending): unfiltered",
		}),

		// ...

		Entry(nil, &spawnTE{
			naviTE: naviTE{
				message:      "files(fastward): listen pending",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFiles,
			},
			active: activeTE{
				resumeAtPath: RESUME_AT_CAN_YOU_KISS_ME_FIRST,
				listenState:  nav.ListenPending,
			},
			resumeAt: START_AT_BEFORE_LIFE,
			profile:  "-> files(pending): unfiltered",
		}),
	)
})
