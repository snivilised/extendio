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
		root = cwd()
	})

	DescribeTable("Listener",
		func(entry *listenTE) {
			navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = begin("🛡️")
				o.Subscription = entry.subscription
				o.Behaviours.Listen.InclusiveStart = entry.incStart
				o.Behaviours.Listen.InclusiveStop = entry.incStop
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
				o.DoExtend = entry.extended
				o.Callback = func(item *nav.TraverseItem) *LocalisableError {
					GinkgoWriter.Printf("---> 🔊 LISTENING-CALLBACK: name: '%v'\n",
						item.Extension.Name,
					)

					Expect(lo.Contains(entry.prohibited, item.Extension.Name)).To(
						BeFalse(), reason(item.Extension.Name),
					)
					Expect(lo.Contains(entry.mandatory, item.Extension.Name)).To(
						BeTrue(), reason(item.Extension.Name),
					)

					entry.mandatory = lo.Reject(entry.mandatory, func(s string, _ int) bool {
						return s == item.Extension.Name
					})
					return nil
				}
			})
			path := path(root, entry.relative)
			navigator.Walk(path)

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
			navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = begin("🛡️")
				o.Subscription = nav.SubscribeFolders
				o.Listen.Stop = &nav.ListenBy{
					Name: "Name: DREAM-POP",
					Fn: func(item *nav.TraverseItem) bool {
						return item.Extension.Name == "DREAM-POP"
					},
				}
				o.Notify.OnStop = func(description string) {
					GinkgoWriter.Printf("===> ⛔ Stop Listening: '%v'\n", description)
				}
				o.DoExtend = true
				o.Callback = foldersCallback("EARLY-EXIT-😴", o.DoExtend)
			})
			path := path(root, "")
			navigator.Walk(path)
		})

		It("should: exit early (files)", func() {
			navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
				o.Notify.OnBegin = begin("🛡️")
				o.Subscription = nav.SubscribeFiles
				o.Listen.Stop = &nav.ListenBy{
					Name: "Name(contains): Captain",
					Fn: func(item *nav.TraverseItem) bool {
						return strings.Contains(item.Extension.Name, "Captain")
					},
				}
				o.Notify.OnStop = func(description string) {
					GinkgoWriter.Printf("===> ⛔ Stop Listening: '%v'\n", description)
				}
				o.DoExtend = true
				o.Callback = filesCallback("EARLY-EXIT-😴", o.DoExtend)
			})
			path := path(root, "")
			navigator.Walk(path)
		})
	})

	Context("folders", func() {
		Context("given: filter and listen both active", func() {
			It("🧪 should: apply filter within the listen range", func() {
				navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🛡️")
					o.Subscription = nav.SubscribeFolders
					o.Filters.Current = &nav.RegexFilter{
						Filter: nav.Filter{
							Name:          "Contains 'o'",
							RequiredScope: nav.ScopeAllEn,
							Pattern:       "(i?)o",
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
					o.DoExtend = true
					o.Callback = func(item *nav.TraverseItem) *translate.LocalisableError {
						GinkgoWriter.Printf("---> 🔊 LISTENING-CALLBACK: name: '%v'\n",
							item.Extension.Name,
						)
						GinkgoWriter.Printf(
							"===> ⚗️ Regex Filter(%v) source: '%v', item-name: '%v', item-scope(fs): '%v(%v)'\n",
							o.Filters.Current.Description(), o.Filters.Current.Source(), item.Extension.Name,
							item.Extension.NodeScope, o.Filters.Current.Scope(),
						)
						Expect(o.Filters.Current.IsMatch(item)).To(BeTrue(), reason(item.Extension.Name))
						return nil
					}
				})
				path := path(root, "edm/ELECTRONICA")
				_ = navigator.Walk(path)
			})
		})
	})
})
