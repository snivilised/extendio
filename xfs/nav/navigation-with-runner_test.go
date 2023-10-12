package nav_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/extendio/internal/helpers"
	"github.com/snivilised/extendio/xfs/nav"
	"github.com/snivilised/lorax/boost"

	. "github.com/snivilised/extendio/i18n"
)

var _ = Describe("NavigationWithRunner", Ordered, func() {
	var (
		root         string
		jroot        string
		fromJSONPath string
	)

	BeforeAll(func() {
		root = musico()
		jroot = helpers.JoinCwd("Test", "json")
		fromJSONPath = helpers.Path(jroot, "resume-state.json")
	})

	BeforeEach(func() {
		if err := Use(func(o *UseOptions) {
			o.Tag = DefaultLanguage.Get()
		}); err != nil {
			Fail(err.Error())
		}
	})

	Context("resume and worker pool acceleration", func() {
		When("universal: listen pending(logged)", func() {
			It("🧪 should: ...", SpecTimeout(time.Second*5), func(ctxSpec SpecContext) {
				ctx, cancel := context.WithCancel(ctxSpec)
				path := helpers.Path(root, "RETRO-WAVE")
				restorer := func(o *nav.TraverseOptions, active *nav.ActiveState) {
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
						NotExtended,
					)
				}

				wgan := boost.NewAnnotatedWaitGroup("🍂 traversal")
				wgan.Add(1, navigatorRoutineName)
				createWith := nav.RunnerWithResume | nav.RunnerWithPool
				now := 3
				JobsChOut := make(boost.JobStream[nav.TraverseItemInput], DefaultJobsChSize)
				jobsOutputChOut := make(boost.JobOutputStream[nav.TraverseOutput], DefaultJobsChSize)

				result, err := nav.New().With(createWith, &nav.RunnerInfo{
					PrimeInfo: &nav.Prime{
						Path: path,
						OptionsFn: func(o *nav.TraverseOptions) {
							o.Notify.OnBegin = begin("🛡️")
							o.Store.Subscription = nav.SubscribeAny
							o.Callback = universalCallbackNoAssert(
								"universal: Path contains folders(Prime-WithRunner)",
								NotExtended,
							)
						},
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
						JobsChOut:       JobsChOut,
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
			})
		})
	})

	When("Filter Applied", func() {
		It("🧪 should: only invoke sync callback for filtered items", func(ctxSpec SpecContext) {
			ctx, cancel := context.WithCancel(ctxSpec)
			defer cancel()

			path := helpers.Path(root, "RETRO-WAVE")

			wgan := boost.NewAnnotatedWaitGroup("🍂 traversal")
			wgan.Add(1, navigatorRoutineName)
			now := 3
			jobsChOut := make(boost.JobStream[nav.TraverseItemInput], DefaultJobsChSize)

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
						NotExtended,
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
		})
	})
})
