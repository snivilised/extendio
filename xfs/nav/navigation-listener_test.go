package nav_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	"github.com/snivilised/extendio/internal/helpers"

	. "github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("Listener", Ordered, func() {
	var root string

	BeforeAll(func() {
		_ = Use(func(o *UseOptions) {
			o.Tag = DefaultLanguage.Get()
		})

		root = musico()
	})

	DescribeTable("Listener",
		func(entry *listenTE) {
			path := helpers.Path(root, entry.relative)
			session := &nav.PrimarySession{
				Path: path,
			}
			_, _ = session.Configure(func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = begin("🛡️")
				o.Store.Subscription = entry.subscription
				o.Store.Behaviours.Listen.InclusiveStart = entry.incStart
				o.Store.Behaviours.Listen.InclusiveStop = entry.incStop
				if entry.listenDefs != nil {
					o.Store.ListenDefs = *entry.listenDefs
				}
				if !entry.mute {
					o.Notify.OnStart = func(description string) {
						GinkgoWriter.Printf("===> 🎶 Start Listening: '%v'\n", description)
					}
					o.Notify.OnStop = func(description string) {
						GinkgoWriter.Printf("===> ⛔ Stop Listening: '%v'\n", description)
					}
				}
				o.Store.DoExtend = entry.extended
				o.Callback = nav.LabelledTraverseCallback{
					Label: "test listener callback",
					Fn: func(item *nav.TraverseItem) error {
						GinkgoWriter.Printf("---> 🔊 LISTENING-CALLBACK: name: '%v'\n",
							item.Extension.Name,
						)

						prohibited := fmt.Sprintf("%v, was invoked, but should NOT have been",
							helpers.Reason(item.Extension.Name),
						)
						Expect(lo.Contains(entry.prohibited, item.Extension.Name)).To(
							BeFalse(), prohibited,
						)

						mandatory := fmt.Sprintf("%v, was not invoked, but should have been",
							helpers.Reason(item.Extension.Name),
						)
						Expect(lo.Contains(entry.mandatory, item.Extension.Name)).To(
							BeTrue(), mandatory,
						)

						entry.mandatory = lo.Reject(entry.mandatory, func(s string, _ int) bool {
							return s == item.Extension.Name
						})
						return nil
					},
				}
			}).Run()

			reason := fmt.Sprintf("❌ remaining: '%v'", strings.Join(entry.mandatory, ", "))
			Expect(len(entry.mandatory)).To(Equal(0), reason)
		},
		func(entry *listenTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.message)
		},

		// === folders =======================================================

		Entry(nil, &listenTE{
			naviTE: naviTE{
				message:      "listening, start and stop (folders, inc:default)",
				relative:     "RETRO-WAVE",
				extended:     true,
				subscription: nav.SubscribeFolders,
				mandatory:    []string{"Night Drive", "College", "Northern Council", "Teenage Color"},
				prohibited:   []string{"RETRO-WAVE", "Chromatics", "Electric Youth", "Innerworld"},
			},
			listenDefs: &nav.ListenDefinitions{
				StartAt: &nav.FilterDef{
					Type:        nav.FilterTypeGlobEn,
					Description: "Start Listening At: Night Drive",
					Pattern:     "Night Drive",
				},
				StopAt: &nav.FilterDef{
					Type:        nav.FilterTypeGlobEn,
					Description: "Stop Listening At: Electric Youth",
					Pattern:     "Electric Youth",
				},
			},
			incStart: true,
			incStop:  false,
		}),

		Entry(nil, &listenTE{
			naviTE: naviTE{
				message:      "listening, start and stop (folders, excl:start, inc:stop, mute)",
				relative:     "RETRO-WAVE",
				extended:     true,
				subscription: nav.SubscribeFolders,
				mandatory:    []string{"College", "Northern Council", "Teenage Color", "Electric Youth"},
				prohibited: []string{"Night Drive", "RETRO-WAVE", "Chromatics",
					"Innerworld",
				},
			},
			listenDefs: &nav.ListenDefinitions{
				StartAt: &nav.FilterDef{
					Type:        nav.FilterTypeRegexEn,
					Description: "Start Listening At: Night Drive",
					Pattern:     "Night Drive",
				},
				StopAt: &nav.FilterDef{
					Type:        nav.FilterTypeGlobEn,
					Description: "Stop Listening At: Electric Youth",
					Pattern:     "Electric Youth",
				},
			},
			incStart: false,
			incStop:  true,
			mute:     true,
		}),

		Entry(nil, &listenTE{
			naviTE: naviTE{
				message:      "listening, start only (folders, inc:default)",
				relative:     "RETRO-WAVE",
				extended:     true,
				subscription: nav.SubscribeFolders,
				mandatory: []string{"Night Drive", "College", "Northern Council", "Teenage Color",
					"Electric Youth", "Innerworld",
				},
				prohibited: []string{"RETRO-WAVE", "Chromatics"},
			},
			listenDefs: &nav.ListenDefinitions{
				StartAt: &nav.FilterDef{
					Type:        nav.FilterTypeRegexEn,
					Description: "Start Listening At: Night Drive",
					Pattern:     "Night Drive",
				},
			},
			incStart: true,
			incStop:  false,
		}),

		Entry(nil, &listenTE{
			naviTE: naviTE{
				message:      "listening, stop only (folders, inc:default)",
				relative:     "RETRO-WAVE",
				extended:     true,
				subscription: nav.SubscribeFolders,
				mandatory: []string{"RETRO-WAVE", "Chromatics", "Night Drive", "College",
					"Northern Council", "Teenage Color",
				},
				prohibited: []string{"Electric Youth", "Innerworld"},
			},

			listenDefs: &nav.ListenDefinitions{
				StopAt: &nav.FilterDef{
					Type:        nav.FilterTypeGlobEn,
					Description: "Stop Listening At: Electric Youth",
					Pattern:     "Electric Youth",
				},
			},
			incStart: true,
			incStop:  false,
		}),

		Entry(nil, &listenTE{
			naviTE: naviTE{
				message:      "listening, stop only (folders, inc:default)",
				relative:     "RETRO-WAVE",
				extended:     true,
				subscription: nav.SubscribeFolders,
				mandatory:    []string{"RETRO-WAVE", "Chromatics"},
				prohibited: []string{"Night Drive", "College", "Northern Council",
					"Teenage Color", "Electric Youth", "Innerworld",
				},
			},
			listenDefs: &nav.ListenDefinitions{
				StopAt: &nav.FilterDef{
					Type:        nav.FilterTypeGlobEn,
					Description: "Stop Listening At: Night Drive",
					Pattern:     "Night Drive",
				},
			},
			incStart: true,
			incStop:  false,
		}),
	)

	Context("given: Early Exit", func() {
		It("should: exit early (folders)", func() {
			path := helpers.Path(root, "")
			session := &nav.PrimarySession{
				Path: path,
			}
			_, _ = session.Configure(func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = begin("🛡️")
				o.Store.Subscription = nav.SubscribeFolders
				o.Store.ListenDefs = nav.ListenDefinitions{
					StopAt: &nav.FilterDef{
						Type:        nav.FilterTypeGlobEn,
						Description: "Stop Listening At: DREAM-POP",
						Pattern:     "DREAM-POP",
					},
				}

				o.Notify.OnStop = func(description string) {
					GinkgoWriter.Printf("===> ⛔ Stop Listening: '%v'\n", description)
				}
				o.Store.DoExtend = true
				o.Callback = foldersCallback("EARLY-EXIT-😴", o.Store.DoExtend)
			}).Run()
		})

		It("should: exit early (files)", func() {
			path := helpers.Path(root, "")
			session := &nav.PrimarySession{
				Path: path,
			}
			_, _ = session.Configure(func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = begin("🛡️")
				o.Store.Subscription = nav.SubscribeFiles
				o.Store.ListenDefs = nav.ListenDefinitions{
					StopAt: &nav.FilterDef{
						Type:        nav.FilterTypeGlobEn,
						Description: "Stop Listening At: Item containing Captain",
						Pattern:     "*Captain*",
					},
				}

				o.Notify.OnStop = func(description string) {
					GinkgoWriter.Printf("===> ⛔ Stop Listening: '%v'\n", description)
				}
				o.Store.DoExtend = true
				o.Callback = filesCallback("EARLY-EXIT-😴", o.Store.DoExtend)
			}).Run()
		})
	})

	Context("folders", func() {
		Context("given: filter and listen both active", func() {
			It("🧪 should: apply filter within the listen range", func() {
				path := helpers.Path(root, "edm/ELECTRONICA")
				session := &nav.PrimarySession{
					Path: path,
				}
				result, _ := session.Configure(func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🛡️")
					o.Store.Subscription = nav.SubscribeFolders
					o.Store.FilterDefs = &nav.FilterDefinitions{
						Node: nav.FilterDef{
							Type:        nav.FilterTypeRegexEn,
							Description: "Contains 'o'",
							Scope:       nav.ScopeAllEn,
							Pattern:     "(i?)o",
						},
					}

					o.Store.ListenDefs = nav.ListenDefinitions{
						StartAt: &nav.FilterDef{
							Type: nav.FilterTypeCustomEn,
							Custom: &helpers.CustomFilter{
								Value: "Orbital",
								Name:  "Start Listening At: Orbital",
							},
						},
						StopAt: &nav.FilterDef{
							Type: nav.FilterTypeCustomEn,
							Custom: &helpers.CustomFilter{
								Value: "Underworld",
								Name:  "Stop Listening At: Underworld",
							},
						},
					}

					o.Notify.OnStart = func(description string) {
						GinkgoWriter.Printf("===> 🎶 Start Listening: '%v'\n", description)
					}
					o.Notify.OnStop = func(description string) {
						GinkgoWriter.Printf("===> ⛔ Stop Listening: '%v'\n", description)
					}
					o.Store.DoExtend = true
					o.Callback = nav.LabelledTraverseCallback{
						Label: "Listener Test Callback",
						Fn: func(item *nav.TraverseItem) error {
							GinkgoWriter.Printf("---> 🔊 LISTENING-CALLBACK: name: '%v'\n",
								item.Extension.Name,
							)
							GinkgoWriter.Printf(
								"===> ⚗️ Regex Filter(%v) source: '%v', item-name: '%v', item-scope(fs): '%v(%v)'\n",
								o.Store.FilterDefs.Node.Description,
								o.Store.FilterDefs.Node.Pattern,
								item.Extension.Name,
								item.Extension.NodeScope,
								o.Store.FilterDefs.Node.Scope,
							)
							Expect(item.Extension.Name).To(MatchRegexp(o.Store.FilterDefs.Node.Pattern),
								helpers.Reason(item.Extension.Name),
							)
							return nil
						},
					}
					o.Store.Logging = logo()
				}).Run()

				files := (*result.Metrics)[nav.MetricNoFilesEn].Count
				folders := (*result.Metrics)[nav.MetricNoFoldersEn].Count

				GinkgoWriter.Printf("---> 🍕🍕 Metrics, files:'%v', folders:'%v'\n",
					files, folders,
				)
			})
		})
	})
})
