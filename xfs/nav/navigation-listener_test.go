package nav_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	"github.com/snivilised/extendio/translate"
	. "github.com/snivilised/extendio/translate"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("Listener", Ordered, func() {
	var root string

	BeforeAll(func() {
		root = origin()
	})

	DescribeTable("Listener",
		func(entry *listenTE) {
			path := path(root, entry.relative)
			session := &nav.PrimarySession{
				Path: path,
			}
			_ = session.Configure(func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = begin("🛡️")
				o.Store.Subscription = entry.subscription
				o.Store.Behaviours.Listen.InclusiveStart = entry.incStart
				o.Store.Behaviours.Listen.InclusiveStop = entry.incStop
				o.Listen.Start = entry.start
				o.Listen.Stop = entry.stop
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
					Fn: func(item *nav.TraverseItem) *LocalisableError {
						GinkgoWriter.Printf("---> 🔊 LISTENING-CALLBACK: name: '%v'\n",
							item.Extension.Name,
						)

						prohibited := fmt.Sprintf("%v, was invoked, but should NOT have been", reason(item.Extension.Name))
						Expect(lo.Contains(entry.prohibited, item.Extension.Name)).To(
							BeFalse(), prohibited,
						)
						mandatory := fmt.Sprintf("%v, was not invoked, but should have been", reason(item.Extension.Name))
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
			// navigator.Walk(path)

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
			start: &nav.ListenBy{
				Name: "Name: Night Drive",
				Fn: func(item *nav.TraverseItem) bool {
					return item.Extension.Name == "Night Drive"
				},
			},
			stop: &nav.ListenBy{
				Name: "Name: Electric Youth",
				Fn: func(item *nav.TraverseItem) bool {
					return item.Extension.Name == "Electric Youth"
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
			start: &nav.ListenBy{
				Name: "Name: Night Drive",
				Fn: func(item *nav.TraverseItem) bool {
					return item.Extension.Name == "Night Drive"
				},
			},
			stop: &nav.ListenBy{
				Name: "Name: Electric Youth",
				Fn: func(item *nav.TraverseItem) bool {
					return item.Extension.Name == "Electric Youth"
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
			start: &nav.ListenBy{
				Name: "Name: Night Drive",
				Fn: func(item *nav.TraverseItem) bool {
					return item.Extension.Name == "Night Drive"
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
			stop: &nav.ListenBy{
				Name: "Name: Electric Youth",
				Fn: func(item *nav.TraverseItem) bool {
					return item.Extension.Name == "Electric Youth"
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
			stop: &nav.ListenBy{
				Name: "Name: Night Drive",
				Fn: func(item *nav.TraverseItem) bool {
					return item.Extension.Name == "Night Drive"
				},
			},
			incStart: true,
			incStop:  false,
		}),
	)

	Context("given: Early Exit", func() {
		It("should: exit early (folders)", func() {
			path := path(root, "")
			session := &nav.PrimarySession{
				Path: path,
			}
			session.Configure(func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = begin("🛡️")
				o.Store.Subscription = nav.SubscribeFolders
				o.Listen.Stop = &nav.ListenBy{
					Name: "Name: DREAM-POP",
					Fn: func(item *nav.TraverseItem) bool {
						return item.Extension.Name == "DREAM-POP"
					},
				}
				o.Notify.OnStop = func(description string) {
					GinkgoWriter.Printf("===> ⛔ Stop Listening: '%v'\n", description)
				}
				o.Store.DoExtend = true
				o.Callback = foldersCallback("EARLY-EXIT-😴", o.Store.DoExtend)
			}).Run()
			// navigator.Walk(path)
		})

		It("should: exit early (files)", func() {
			path := path(root, "")
			session := &nav.PrimarySession{
				Path: path,
			}
			session.Configure(func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = begin("🛡️")
				o.Store.Subscription = nav.SubscribeFiles
				o.Listen.Stop = &nav.ListenBy{
					Name: "Name(contains): Captain",
					Fn: func(item *nav.TraverseItem) bool {
						return strings.Contains(item.Extension.Name, "Captain")
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
				path := path(root, "edm/ELECTRONICA")
				session := &nav.PrimarySession{
					Path: path,
				}
				session.Configure(func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🛡️")
					o.Store.Subscription = nav.SubscribeFolders
					o.Store.FilterDefs = &nav.FilterDefinitions{
						Node: nav.FilterDef{
							Type:        nav.FilterTypeRegexEn,
							Description: "Contains 'o'",
							Scope:       nav.ScopeAllEn,
							Source:      "(i?)o",
						},
					}
					o.Listen.Start = &nav.ListenBy{
						Name: "Name: Orbital",
						Fn: func(item *nav.TraverseItem) bool {
							return item.Extension.Name == "Orbital"
						},
					}
					o.Listen.Stop = &nav.ListenBy{
						Name: "Name: Underworld",
						Fn: func(item *nav.TraverseItem) bool {
							return item.Extension.Name == "Underworld"
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
						Fn: func(item *nav.TraverseItem) *translate.LocalisableError {
							GinkgoWriter.Printf("---> 🔊 LISTENING-CALLBACK: name: '%v'\n",
								item.Extension.Name,
							)
							GinkgoWriter.Printf(
								"===> ⚗️ Regex Filter(%v) source: '%v', item-name: '%v', item-scope(fs): '%v(%v)'\n",
								o.Store.FilterDefs.Node.Description, o.Store.FilterDefs.Node.Source, item.Extension.Name,
								item.Extension.NodeScope, o.Store.FilterDefs.Node.Scope,
							)
							Expect(item.Extension.Name).To(MatchRegexp(o.Store.FilterDefs.Node.Source), reason(item.Extension.Name))
							return nil
						},
					}
					o.Store.Logging = logo()
				}).Run()
			})
		})
	})
})
