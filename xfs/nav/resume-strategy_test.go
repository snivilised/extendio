package nav_test

import (
	"fmt"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"

	. "github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/internal/helpers"
	"github.com/snivilised/extendio/xfs/nav"
)

type strategyTheme struct {
	label string
}

type strategyInvokeInfo struct {
	files   uint
	folders uint
}

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
	strategies = []nav.ResumeStrategyEnum{
		nav.ResumeStrategyFastwardEn,
		nav.ResumeStrategySpawnEn,
	}
	themes = map[nav.ResumeStrategyEnum]strategyTheme{
		nav.ResumeStrategyFastwardEn: {label: "FASTWARD"},
		nav.ResumeStrategySpawnEn:    {label: "SPAWN"},
	}

	profiles = map[string]resumeTestProfile{
		// === Listening (uni/folder/file) (pend/active)

		"-> universal(pending): unfiltered": {
			filtered:   false,
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

var _ = Describe("Resume", Ordered, func() {

	var (
		root         string
		jroot        string
		fromJsonPath string
	)

	BeforeAll(func() {
		root = musico()
		jroot = helpers.JoinCwd("Test", "json")
		fromJsonPath = strings.Join(
			[]string{jroot, "resume-state.json"},
			string(filepath.Separator),
		)
	})

	BeforeEach(func() {
		_ = Use(func(o *UseOptions) {
			o.Tag = DefaultLanguage.Get()
		})
	})

	DescribeTable("resume",
		func(entry *resumeTE) {

			invocations := map[nav.ResumeStrategyEnum]*strategyInvokeInfo{}

			for _, strategyEn := range strategies {
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
					active.Root = helpers.Path(root, entry.relative)
					active.NodePath = helpers.Path(root, entry.active.resumeAt)
					active.Listen = entry.active.listenState
					o.Store.Subscription = entry.subscription

					if !profile.filtered {
						o.Store.FilterDefs = nil
					}
					//
					// end of synthetic assignments

					if strategyEn == nav.ResumeStrategyFastwardEn {
						o.Notify.OnBegin = func(state *nav.NavigationState) {
							panic("begin handler should not be invoked because begin notification muted")
						}
					}
					GinkgoWriter.Printf("===> 🐚 restoring ...\n")

					once := nav.LabelledTraverseCallback{
						Label: "test once callback",
						Fn: func(item *nav.TraverseItem) error {
							_, found := recording[item.Extension.Name]

							Expect(found).To(BeFalse(), fmt.Sprintf("once only invoke failure -> %v",
								helpers.Reason(item.Extension.Name)))

							recording[item.Extension.Name] = len(item.Children)
							return nil
						},
					}

					o.Callback = nav.LabelledTraverseCallback{
						Label: "unit test callback for resume",
						Fn: func(item *nav.TraverseItem) error {
							depth := lo.TernaryF(o.Store.DoExtend,
								func() int { return item.Extension.Depth },
								func() int { return 9999 },
							)
							GinkgoWriter.Printf(
								"---> ⏩ %v: (depth:%v) '%v'\n", themes[strategyEn].label, depth, item.Path,
							)

							if strategyEn == nav.ResumeStrategyFastwardEn {
								segments := strings.Split(item.Path, string(filepath.Separator))
								last := segments[len(segments)-1]
								if _, found := prohibited[last]; found {
									Fail(fmt.Sprintf("item: '%v' should have been fast forwarded over", item.Path))
								}
							}
							return once.Fn(item)
						},
					}
					if strategyEn == nav.ResumeStrategyFastwardEn {
						if entry.clientListenAt != "" {
							o.Store.ListenDefs.StartAt = &nav.FilterDef{
								Type:        nav.FilterTypeGlobEn,
								Description: fmt.Sprintf("Start Listening At: %v", entry.clientListenAt),
								Pattern:     entry.clientListenAt,
							}
						}
					}
					if entry.log {
						o.Store.Logging = logo()
					}
				}

				session := &nav.ResumeSession{
					Path:     fromJsonPath,
					Strategy: strategyEn,
				}
				result, _ := session.Configure(restore).Run()

				if profile.mandatory != nil {
					for _, name := range profile.mandatory {
						_, found := recording[name]
						Expect(found).To(BeTrue(),
							fmt.Sprintf("mandatory item failure -> %v", helpers.Reason(name)),
						)
					}
				}

				invocations[strategyEn] = &strategyInvokeInfo{
					files:   (*result.Metrics)[nav.MetricNoFilesEn].Count,
					folders: (*result.Metrics)[nav.MetricNoFoldersEn].Count,
				}
			}

			for _, strategyEn := range strategies {
				fmt.Printf("💡💡 invocations(%v) - files:%v, folders:%v\n",
					themes[strategyEn].label,
					invocations[strategyEn].files,
					invocations[strategyEn].folders,
				)
			}

			if len(strategies) == 2 {
				Expect(invocations[nav.ResumeStrategyFastwardEn].files).To(
					Equal(invocations[nav.ResumeStrategySpawnEn].files),
					"Both strategies should invoke same no of files")

				Expect(invocations[nav.ResumeStrategyFastwardEn].folders).To(
					Equal(invocations[nav.ResumeStrategySpawnEn].folders),
					"Both strategies should invoke same no of folders")
			}
		},
		func(entry *resumeTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.message)
		},

		// === Listening (uni/folder/file) (pend/active)
		//
		// for the active cases, it doesn't really matter what the resumeAt is set
		// to, because the listener is already in the active listening state. But resumeAt
		// still has to be set because that is what would happen in the real world.
		//

		Entry(nil, &resumeTE{
			naviTE: naviTE{
				message:      "universal: listen pending(logged)",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
			},
			active: activeTE{
				resumeAt:    RESUME_AT_TEENAGE_COLOR,
				listenState: nav.ListenPending,
			},
			clientListenAt: START_AT_ELECTRIC_YOUTH,
			profile:        "-> universal(pending): unfiltered",
			log:            true,
		}),

		Entry(nil, &resumeTE{
			naviTE: naviTE{
				message:      "universal: listen active",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
			},
			active: activeTE{
				resumeAt:    RESUME_AT_TEENAGE_COLOR,
				listenState: nav.ListenActive,
			},
			// For these scenarios (START_AT_CLIENT_ALREADY_ACTIVE), since
			// listening is already active, the value of resumeAt is irrelevant,
			// because the client was already listening in the previous session,
			// which is reflected by the state being active. So in essence, the client
			// listen value is a historical event, so the value defined here is a moot
			// point.
			//
			clientListenAt: START_AT_CLIENT_ALREADY_ACTIVE,
			profile:        "-> universal(active): unfiltered",
		}),

		Entry(nil, &resumeTE{
			naviTE: naviTE{
				message:      "folders: listen pending",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFolders,
			},
			active: activeTE{
				resumeAt:    RESUME_AT_TEENAGE_COLOR,
				listenState: nav.ListenPending,
			},
			clientListenAt: START_AT_ELECTRIC_YOUTH,
			profile:        "-> folders(pending): unfiltered",
		}),

		Entry(nil, &resumeTE{
			naviTE: naviTE{
				message:      "folders: listen active",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFolders,
			},
			active: activeTE{
				resumeAt:    RESUME_AT_TEENAGE_COLOR,
				listenState: nav.ListenActive,
			},
			clientListenAt: START_AT_CLIENT_ALREADY_ACTIVE,
			profile:        "-> folders(active): unfiltered",
		}),

		Entry(nil, &resumeTE{
			naviTE: naviTE{
				message:      "files: listen pending",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFiles,
			},
			active: activeTE{
				resumeAt:    RESUME_AT_CAN_YOU_KISS_ME_FIRST,
				listenState: nav.ListenPending,
			},
			clientListenAt: START_AT_BEFORE_LIFE,
			profile:        "-> files(pending): unfiltered",
		}),

		Entry(nil, &resumeTE{
			naviTE: naviTE{
				message:      "files: listen active",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFiles,
			},
			active: activeTE{
				resumeAt:    RESUME_AT_CAN_YOU_KISS_ME_FIRST,
				listenState: nav.ListenActive,
			},
			clientListenAt: START_AT_CLIENT_ALREADY_ACTIVE,
			profile:        "-> files(active): unfiltered",
		}),

		// === Filtering (uni/folder/file)

		Entry(nil, &resumeTE{
			naviTE: naviTE{
				message:      "universal: listen not active/deaf",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
			},
			active: activeTE{
				resumeAt:    RESUME_AT_TEENAGE_COLOR,
				listenState: nav.ListenDeaf,
			},
			profile: "-> universal: filtered",
		}),

		Entry(nil, &resumeTE{
			naviTE: naviTE{
				message:      "folders: listen not active/deaf",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFolders,
			},
			active: activeTE{
				resumeAt:    RESUME_AT_TEENAGE_COLOR,
				listenState: nav.ListenDeaf,
			},
			profile: "-> folders: filtered",
		}),

		Entry(nil, &resumeTE{
			naviTE: naviTE{
				message:      "files: listen not active/deaf",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFiles,
			},
			active: activeTE{
				resumeAt:    RESUME_AT_CAN_YOU_KISS_ME_FIRST,
				listenState: nav.ListenDeaf,
			},
			profile: "-> files: filtered",
		}),

		// === Listening and filtering (uni/folder/file)

		Entry(nil, &resumeTE{
			naviTE: naviTE{
				message:      "universal: listen pending and filtered",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
			},
			active: activeTE{
				resumeAt:    RESUME_AT_TEENAGE_COLOR,
				listenState: nav.ListenPending,
			},
			clientListenAt: START_AT_ELECTRIC_YOUTH,
			profile:        "-> universal: listen pending and filtered",
		}),

		Entry(nil, &resumeTE{
			naviTE: naviTE{
				message:      "universal: listen active and filtered",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeAny,
			},
			active: activeTE{
				resumeAt:    RESUME_AT_TEENAGE_COLOR,
				listenState: nav.ListenActive,
			},
			clientListenAt: START_AT_ELECTRIC_YOUTH,
			profile:        "-> universal: filtered",
		}),

		Entry(nil, &resumeTE{
			naviTE: naviTE{
				message:      "folders: listen pending and filtered",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFolders,
			},
			active: activeTE{
				resumeAt:    RESUME_AT_TEENAGE_COLOR,
				listenState: nav.ListenPending,
			},
			clientListenAt: START_AT_ELECTRIC_YOUTH,
			profile:        "-> folders: listen pending and filtered",
		}),

		Entry(nil, &resumeTE{
			naviTE: naviTE{
				message:      "folders: listen active and filtered",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFolders,
			},
			active: activeTE{
				resumeAt:    RESUME_AT_TEENAGE_COLOR,
				listenState: nav.ListenActive,
			},
			clientListenAt: START_AT_ELECTRIC_YOUTH,
			profile:        "-> folders: filtered",
		}),

		Entry(nil, &resumeTE{
			naviTE: naviTE{
				message:      "files: listen pending and filtered",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFiles,
			},
			active: activeTE{
				resumeAt:    RESUME_AT_CAN_YOU_KISS_ME_FIRST,
				listenState: nav.ListenPending,
			},
			clientListenAt: START_AT_BEFORE_LIFE,
			profile:        "-> files: listen pending and filtered",
		}),

		Entry(nil, &resumeTE{
			naviTE: naviTE{
				message:      "files: listen active and filtered",
				relative:     "RETRO-WAVE",
				subscription: nav.SubscribeFiles,
			},
			active: activeTE{
				resumeAt:    RESUME_AT_CAN_YOU_KISS_ME_FIRST,
				listenState: nav.ListenActive,
			},
			clientListenAt: START_AT_BEFORE_LIFE,
			profile:        "-> files: filtered",
		}),
	)
})
