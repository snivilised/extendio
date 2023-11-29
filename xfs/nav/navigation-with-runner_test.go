package nav_test

import (
	"context"
	"time"

	"github.com/fortytw2/leaktest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/extendio/internal/helpers"
	"github.com/snivilised/extendio/xfs/nav"
	"github.com/snivilised/lorax/boost"

	. "github.com/snivilised/extendio/i18n"
)

var _ = Describe("NavigationWithRunner", Ordered, func() {
	var (
		root            string
		jroot           string
		fromJSONPath    string
		path            string
		now             int
		jobsChOut       boost.JobStream[nav.TraverseItemInput]
		jobsOutputChOut boost.JobOutputStream[nav.TraverseOutput]
		setOptions      func(o *nav.TraverseOptions)
	)

	BeforeAll(func() {
		root = musico()
		jroot = helpers.JoinCwd("Test", "json")
		fromJSONPath = helpers.Path(jroot, "resume-state.json")
		path = helpers.Path(root, "RETRO-WAVE")
		now = 3

		setOptions = func(o *nav.TraverseOptions) {
			o.Notify.OnBegin = begin("🛡️")
			o.Store.Subscription = nav.SubscribeAny
			o.Callback = universalCallbackNoAssert(
				"universal: Path contains folders(Prime-WithRunner)",
			)
		}
	})

	BeforeEach(func() {
		if err := Use(func(o *UseOptions) {
			o.Tag = DefaultLanguage.Get()
		}); err != nil {
			Fail(err.Error())
		}

		jobsChOut = make(boost.JobStream[nav.TraverseItemInput], DefaultJobsChSize)
		jobsOutputChOut = make(boost.JobOutputStream[nav.TraverseOutput], DefaultJobsChSize)
	})

	Context("resume and worker pool acceleration", func() {
		var (
			restorer func(o *nav.TraverseOptions, active *nav.ActiveState)
		)

		BeforeAll(func() {
			restorer = func(o *nav.TraverseOptions, active *nav.ActiveState) {
				// synthetic assignments
				//
				active.Root = helpers.Path(root, "RETRO-WAVE")
				active.NodePath = helpers.Path(root, ResumeAtTeenageColor)
				active.Listen = nav.ListenPending
				o.Store.Subscription = nav.SubscribeAny
				//
				// end of synthetic assignments
				o.Callback = universalCallbackNoAssert(
					"universal: listen pending(Resume-WithRunner)",
				)
			}
		})

		When("universal: listen pending(logged)", func() {
			It("🧪 should: resume without error",
				SpecTimeout(time.Second*5),
				func(ctxSpec SpecContext) {
					defer leaktest.Check(GinkgoT())()

					ctx, cancel := context.WithCancel(ctxSpec)
					wgan := boost.NewAnnotatedWaitGroup("🍂 traversal")
					wgan.Add(1, navigatorRoutineName)
					createWith := nav.RunnerWithResume | nav.RunnerWithPool

					result, err := nav.New().With(createWith, &nav.RunnerInfo{
						PrimeInfo: &nav.Prime{
							Path:      path,
							OptionsFn: setOptions,
						},
						ResumeInfo: &nav.Resumption{
							RestorePath: fromJSONPath,
							Restorer:    restorer,
							Strategy:    nav.ResumeStrategySpawnEn,
						},
						AccelerationInfo: &nav.Acceleration{
							WgAn:            wgan,
							RoutineName:     navigatorRoutineName,
							NoW:             now,
							JobsChOut:       jobsChOut,
							JobResultsCh:    jobsOutputChOut,
							OutputChTimeout: outputChTimeout,
						},
					}).Run(
						nav.IfWithPoolUseContext(createWith, ctx, cancel)...,
					)

					if createWith&nav.RunnerWithPool > 0 {
						wgan.Wait("👾 test-main")
					}

					Expect(err).Error().To(BeNil())
					_ = result.Session.StartedAt()
					_ = result.Session.Elapsed()
				},
			)
		})
	})

	When("Filter Applied", func() {
		It("🧪 should: only invoke sync callback for filtered items",
			SpecTimeout(time.Second*5),
			func(ctxSpec SpecContext) {
				defer leaktest.Check(GinkgoT())()

				ctx, cancel := context.WithCancel(ctxSpec)
				defer cancel()

				wgan := boost.NewAnnotatedWaitGroup("🍂 traversal")
				wgan.Add(1, navigatorRoutineName)

				filterDefs := &nav.FilterDefinitions{
					Node: nav.FilterDef{
						Type:        nav.FilterTypeGlobEn,
						Description: "flac files",
						Pattern:     "*.flac",
						Scope:       nav.ScopeLeafEn,
					},
				}

				result, err := nav.New().Primary(&nav.Prime{
					Path: path,
					OptionsFn: func(o *nav.TraverseOptions) {
						o.Notify.OnBegin = begin("🛡️")
						o.Store.Subscription = nav.SubscribeFiles
						o.Callback = universalCallbackNoAssert(
							"filtered *.flac files: WithPool",
						)
						o.Store.FilterDefs = filterDefs
					},
				}).WithPool(
					&nav.AsyncInfo{
						NavigatorRoutineName: navigatorRoutineName,
						WaitAQ:               wgan,
						JobsChanOut:          jobsChOut,
					},
				).NoW(now).Run(ctx, cancel)

				wgan.Wait("👾 test-main")

				Expect(err).Error().To(BeNil())
				_ = result.Session.StartedAt()
				_ = result.Session.Elapsed()
			},
		)
		When("using ProvidedOptions", func() {
			It("🧪 should: ProvideOptions - only invoke sync callback for filtered items",
				SpecTimeout(time.Second*5),
				func(ctxSpec SpecContext) {
					defer leaktest.Check(GinkgoT())()

					ctx, cancel := context.WithCancel(ctxSpec)
					defer cancel()

					wgan := boost.NewAnnotatedWaitGroup("🍂 traversal")
					wgan.Add(1, navigatorRoutineName)

					filterDefs := &nav.FilterDefinitions{
						Node: nav.FilterDef{
							Type:        nav.FilterTypeGlobEn,
							Description: "flac files",
							Pattern:     "*.flac",
							Scope:       nav.ScopeLeafEn,
						},
					}

					providedOptions := nav.GetDefaultOptions()
					providedOptions.Notify.OnBegin = begin("🛡️")
					providedOptions.Store.Subscription = nav.SubscribeFiles
					providedOptions.Callback = universalCallbackNoAssert(
						"filtered *.flac files: WithPool",
					)
					providedOptions.Store.FilterDefs = filterDefs

					_, err := nav.New().Primary(&nav.Prime{
						Path:            path,
						ProvidedOptions: providedOptions, // 😎
					}).WithPool(
						&nav.AsyncInfo{
							NavigatorRoutineName: navigatorRoutineName,
							WaitAQ:               wgan,
							JobsChanOut:          jobsChOut,
						},
					).NoW(now).Run(ctx, cancel)

					wgan.Wait("👾 test-main")

					Expect(err).Error().To(BeNil())
				},
			)
		})
	})
})
