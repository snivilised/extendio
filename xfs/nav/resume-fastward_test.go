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

const (
	NOTHING                         = ""
	RESUME_AT_TEENAGE_COLOR         = "RETRO-WAVE/College/Teenage Color"
	RESUME_AT_CAN_YOU_KISS_ME_FIRST = "RETRO-WAVE/College/Teenage Color/A1 - Can You Kiss Me First.flac"
	START_AT_ELECTRIC_YOUTH         = "Electric Youth"
	START_AT_BEFORE_LIFE            = "A1 - Before Life.flac"
	START_AT_CLIENT_ALREADY_ACTIVE  = "this value doesn't matter"
)

var (
	prohibited = map[string]string{
		"RETRO-WAVE":                      NOTHING,
		"Chromatics":                      NOTHING,
		"Night Drive":                     NOTHING,
		"A1 - The Telephone Call.flac":    NOTHING,
		"A2 - Night Drive.flac":           NOTHING,
		"cover.night-drive.jpg":           NOTHING,
		"vinyl-info.night-drive.txt":      NOTHING,
		"College":                         NOTHING,
		"Northern Council":                NOTHING,
		"A1 - Incident.flac":              NOTHING,
		"A2 - The Zemlya Expedition.flac": NOTHING,
		"cover.northern-council.jpg":      NOTHING,
		"vinyl-info.northern-council.txt": NOTHING,
	}
	filteredListenFlacs = []string{
		"A1 - Before Life.flac",
		"A2 - Runaway.flac",
	}
	filteredFlacs = []string{
		"A1 - Can You Kiss Me First.flac",
		"A2 - Teenage Color.flac",
		"A1 - Before Life.flac",
		"A2 - Runaway.flac",
	}
	textFiles = []string{
		"vinyl-info.teenage-color.txt",
		"vinyl-info.innerworld.txt",
	}

	profiles = map[string]resumeTestProfile{
		// === Listening (uni/folder/file) (pend/active)

		"-> universal(pending): unfiltered": {
			filtered: false,
			// state here?
			prohibited: prohibited,
			mandatory: append(append([]string{
				"Electric Youth",
				"Innerworld",
			}, filteredListenFlacs...), "vinyl-info.innerworld.txt"),
		},

		"-> universal(active): unfiltered": {
			filtered:   false,
			prohibited: prohibited,
			mandatory: append(append([]string{
				"Electric Youth",
				"Innerworld",
			}, filteredFlacs...), textFiles...),
		},

		"-> folders(pending): unfiltered": {
			filtered:   false,
			prohibited: prohibited,
			mandatory: []string{
				"Electric Youth",
				"Innerworld",
			},
		},

		"-> folders(active): unfiltered": {
			filtered:   false,
			prohibited: prohibited,
			mandatory: []string{
				"Teenage Color",
				"Electric Youth",
				"Innerworld",
			},
		},

		"-> files(pending): unfiltered": {
			filtered:   false,
			prohibited: prohibited,
			mandatory: []string{
				"A1 - Before Life.flac",
				"A2 - Runaway.flac",
				"vinyl-info.innerworld.txt",
			},
		},

		"-> files(active): unfiltered": {
			filtered:   false,
			prohibited: prohibited,
			mandatory:  append(filteredFlacs, textFiles...),
		},

		// === Filtering (uni/folder/file)

		"-> universal: filtered": {
			filtered:   true,
			prohibited: prohibited,
			mandatory: append([]string{
				"Electric Youth",
			}, filteredFlacs...),
		},

		"-> folders: filtered": {
			filtered:   true,
			prohibited: prohibited,
			mandatory: []string{
				"Electric Youth",
			},
		},

		"-> files: filtered": {
			filtered:   true,
			prohibited: prohibited,
			mandatory:  filteredFlacs,
		},

		// Listening and filtering (uni/folder/file)

		"-> universal: listen pending and filtered": {
			filtered:   true,
			prohibited: prohibited,
			mandatory: append([]string{
				"Electric Youth"}, filteredListenFlacs...),
		},

		"-> folders: listen pending and filtered": {
			filtered:   true,
			prohibited: prohibited,
			mandatory: []string{
				"Electric Youth",
			},
		},

		"-> files: listen pending and filtered": {
			filtered:   true,
			prohibited: prohibited,
			mandatory:  filteredListenFlacs,
		},
	}
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
			[]string{jroot, "resume-state.json"},
			string(filepath.Separator),
		)
	})

	DescribeTable("fastward",
		func(entry *fastwardTE) {
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

				o.Notify.OnBegin = func(state *nav.NavigationState) {
					panic("begin handler should not be invoked because begin notification muted")
				}
				GinkgoWriter.Printf("===> 🐚 restoring ...\n")

				o.Callback = nav.LabelledTraverseCallback{
					Label: "unit test callback for resume",
					Fn: func(item *nav.TraverseItem) *LocalisableError {
						depth := lo.TernaryF(o.Store.DoExtend,
							func() uint { return item.Extension.Depth },
							func() uint { return 9999 },
						)
						GinkgoWriter.Printf(
							"---> ⏩ FASTWARD: (depth:%v) '%v'\n", depth, item.Path,
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
				if entry.listenStart != "" {
					o.Listen = nav.ListenOptions{
						Start: &nav.ListenBy{
							Fn: func(item *nav.TraverseItem) bool {
								return item.Extension.Name == entry.listenStart
							},
						},
					}
				}
			}
			info := &nav.NewResumerInfo{
				RestorePath: fromJsonPath,
				Restorer:    restore,
				Strategy:    nav.ResumeStrategyFastwardEn,
			}
			result, err := nav.Resume(info)
			Expect(err).To(BeNil())
			Expect(result).ToNot(BeNil())

			if profile.mandatory != nil {
				for _, name := range profile.mandatory {
					_, found := recording[name]
					Expect(found).To(BeTrue(), fmt.Sprintf("mandatory item failure -> %v", reason(name)))
				}
			}
		},
		func(entry *fastwardTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.message)
		},

		// === Listening (uni/folder/file) (pend/active)
		//
		// for the active cases, it doesn't really matter what the listenStart is set
		// to, because the listener is already in the active listening state. But listenStart
		// still has to be set because that is what would happen in the real world.
		//

		Entry(nil, &fastwardTE{
			naviTE: naviTE{
				message:      "universal(fastward): listen pending",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
			},
			active: activeTE{
				resumeAtPath: RESUME_AT_TEENAGE_COLOR,
				listenState:  nav.ListenPending,
			},
			listenStart: START_AT_ELECTRIC_YOUTH,
			profile:     "-> universal(pending): unfiltered",
		}),

		Entry(nil, &fastwardTE{
			naviTE: naviTE{
				message:      "universal(fastward): listen active",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
			},
			active: activeTE{
				resumeAtPath: RESUME_AT_TEENAGE_COLOR,
				listenState:  nav.ListenActive,
			},
			// For these scenarios (START_AT_CLIENT_ALREADY_ACTIVE), since
			// listening is already active, the value of listenStart is irrelevant,
			// because the client was already listening in the previous session,
			// which is reflected by the state being active. So in essence, the client
			// listen value is a historical event, so the value defined here is a moot
			// point.
			//
			listenStart: START_AT_CLIENT_ALREADY_ACTIVE,
			profile:     "-> universal(active): unfiltered",
		}),

		Entry(nil, &fastwardTE{
			naviTE: naviTE{
				message:      "folders(fastward): listen pending",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFolders,
			},
			active: activeTE{
				resumeAtPath: RESUME_AT_TEENAGE_COLOR,
				listenState:  nav.ListenPending,
			},
			listenStart: START_AT_ELECTRIC_YOUTH,
			profile:     "-> folders(pending): unfiltered",
		}),

		Entry(nil, &fastwardTE{
			naviTE: naviTE{
				message:      "folders(fastward): listen active",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFolders,
			},
			active: activeTE{
				resumeAtPath: RESUME_AT_TEENAGE_COLOR,
				listenState:  nav.ListenActive,
			},
			listenStart: START_AT_CLIENT_ALREADY_ACTIVE,
			profile:     "-> folders(active): unfiltered",
		}),

		Entry(nil, &fastwardTE{
			naviTE: naviTE{
				message:      "files(fastward): listen pending",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFiles,
			},
			active: activeTE{
				resumeAtPath: RESUME_AT_CAN_YOU_KISS_ME_FIRST,
				listenState:  nav.ListenPending,
			},
			listenStart: START_AT_BEFORE_LIFE,
			profile:     "-> files(pending): unfiltered",
		}),

		Entry(nil, &fastwardTE{
			naviTE: naviTE{
				message:      "files(fastward): listen active",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFiles,
			},
			active: activeTE{
				resumeAtPath: RESUME_AT_CAN_YOU_KISS_ME_FIRST,
				listenState:  nav.ListenActive,
			},
			listenStart: START_AT_CLIENT_ALREADY_ACTIVE,
			profile:     "-> files(active): unfiltered",
		}),

		// === Filtering (uni/folder/file)

		Entry(nil, &fastwardTE{
			naviTE: naviTE{
				message:      "universal(fastward): listen not active/deaf",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
			},
			active: activeTE{
				resumeAtPath: RESUME_AT_TEENAGE_COLOR,
				listenState:  nav.ListenDeaf,
			},
			profile: "-> universal: filtered",
		}),

		Entry(nil, &fastwardTE{
			naviTE: naviTE{
				message:      "folders(fastward): listen not active/deaf",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFolders,
			},
			active: activeTE{
				resumeAtPath: RESUME_AT_TEENAGE_COLOR,
				listenState:  nav.ListenDeaf,
			},
			profile: "-> folders: filtered",
		}),

		Entry(nil, &fastwardTE{
			naviTE: naviTE{
				message:      "files(fastward): listen not active/deaf",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFiles,
			},
			active: activeTE{
				resumeAtPath: RESUME_AT_CAN_YOU_KISS_ME_FIRST,
				listenState:  nav.ListenDeaf,
			},
			profile: "-> files: filtered",
		}),

		// Listening and filtering (uni/folder/file)

		Entry(nil, &fastwardTE{
			naviTE: naviTE{
				message:      "universal(fastward): listen pending and filtered",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
			},
			active: activeTE{
				resumeAtPath: RESUME_AT_TEENAGE_COLOR,
				listenState:  nav.ListenPending,
			},
			listenStart: START_AT_ELECTRIC_YOUTH,
			profile:     "-> universal: listen pending and filtered",
		}),

		Entry(nil, &fastwardTE{
			naviTE: naviTE{
				message:      "universal(fastward): listen active and filtered",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
			},
			active: activeTE{
				resumeAtPath: RESUME_AT_TEENAGE_COLOR,
				listenState:  nav.ListenActive,
			},
			listenStart: START_AT_ELECTRIC_YOUTH,
			profile:     "-> universal: filtered",
		}),

		Entry(nil, &fastwardTE{
			naviTE: naviTE{
				message:      "folders(fastward): listen pending and filtered",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFolders,
			},
			active: activeTE{
				resumeAtPath: RESUME_AT_TEENAGE_COLOR,
				listenState:  nav.ListenPending,
			},
			listenStart: START_AT_ELECTRIC_YOUTH,
			profile:     "-> folders: listen pending and filtered",
		}),

		Entry(nil, &fastwardTE{
			naviTE: naviTE{
				message:      "folders(fastward): listen active and filtered",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFolders,
			},
			active: activeTE{
				resumeAtPath: RESUME_AT_TEENAGE_COLOR,
				listenState:  nav.ListenActive,
			},
			listenStart: START_AT_ELECTRIC_YOUTH,
			profile:     "-> folders: filtered",
		}),

		Entry(nil, &fastwardTE{
			naviTE: naviTE{
				message:      "files(fastward): listen pending and filtered",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFiles,
			},
			active: activeTE{
				resumeAtPath: RESUME_AT_CAN_YOU_KISS_ME_FIRST,
				listenState:  nav.ListenPending,
			},
			listenStart: START_AT_BEFORE_LIFE,
			profile:     "-> files: listen pending and filtered",
		}),

		Entry(nil, &fastwardTE{
			naviTE: naviTE{
				message:      "files(fastward): listen active and filtered",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFiles,
			},
			active: activeTE{
				resumeAtPath: RESUME_AT_CAN_YOU_KISS_ME_FIRST,
				listenState:  nav.ListenActive,
			},
			listenStart: START_AT_BEFORE_LIFE,
			profile:     "-> files: filtered",
		}),
	)
})
