package nav_test

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fortytw2/leaktest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"

	"github.com/snivilised/extendio/internal/helpers"
	"github.com/snivilised/extendio/xfs/nav"
	"github.com/snivilised/lorax/async"

	. "github.com/snivilised/extendio/i18n"
)

var (
	navigatorRoutineName = async.GoRoutineName("✨ observable-navigator")
)

type (
	asyncResumeTE struct {
		Strategy nav.ResumeStrategyEnum
		Listen   nav.ListeningState
	}

	operatorFunc func(r nav.NavigationRunner) nav.NavigationRunner

	asyncTE struct {
		given    string
		should   string
		operator operatorFunc
		resume   *asyncResumeTE
	}

	asyncOkTE struct {
		asyncTE
	}

	asyncErrorTE struct {
		asyncTE
		fragment string
	}
)

const (
	// we use a large job queue size to prevent blocking as these unit
	// tests don't have a consumer
	DefaultJobsChSize    = 50
	DefaultOutputsChSize = 50
)

type Consumer[O any] struct {
	waitAQ      async.AnnotatedWgAQ
	RoutineName async.GoRoutineName
	OutputChIn  async.OutputStream[O]
	Count       int
}

func StartConsumer[O any](
	ctx context.Context,
	waitAQ async.AnnotatedWgAQ,
	outputChIn async.OutputStream[O],
) *Consumer[O] {
	consumer := &Consumer[O]{
		waitAQ:      waitAQ,
		RoutineName: async.GoRoutineName("💠 consumer"),
		OutputChIn:  outputChIn,
	}

	waitAQ.Add(1, consumer.RoutineName)
	go consumer.run(ctx)

	return consumer
}

func (c *Consumer[O]) run(ctx context.Context) {
	defer func() {
		c.waitAQ.Done(c.RoutineName)
		fmt.Printf("<<<< 💠 consumer.run - finished (QUIT). 💠💠💠 \n")
	}()
	fmt.Printf("<<<< 💠 consumer.run ...(ctx:%+v)\n", ctx)

	for running := true; running; {
		select {
		case <-ctx.Done():
			running = false

			fmt.Println("<<<< 💠 consumer.run - done received 💔💔💔")

		case result, ok := <-c.OutputChIn:
			if ok {
				c.Count++
				fmt.Printf("<<<< 💠 consumer.run - new result arrived(#%v): '%+v' \n",
					c.Count, result.Payload,
				)
			} else {
				running = false
				fmt.Printf("<<<< 💠 consumer.run - no more results available (running: %+v)\n", running)
			}
		}
	}
}

func getSession(entry *asyncTE, root, path, resumeJSONPath string) nav.TraverseSession {
	getOptions := func(o *nav.TraverseOptions) {
		o.Store.Subscription = nav.SubscribeFolders
		o.Store.DoExtend = true
		o.Callback = asyncCallback("async primary session")
		o.Notify.OnBegin = begin("🛡️")
	}

	return lo.TernaryF(entry.resume == nil,
		func() nav.TraverseSession {
			return &nav.PrimarySession{
				Path:     path,
				OptionFn: getOptions,
			}
		},
		func() nav.TraverseSession {
			restorer := func(o *nav.TraverseOptions, active *nav.ActiveState) {
				// synthetic assignments: The client should not perform these
				// types of assignments.
				//
				active.Root = helpers.Path(root, "RETRO-WAVE")
				active.NodePath = helpers.Path(root, ResumeAtTeenageColor)
				active.Listen = entry.resume.Listen
				o.Store.Subscription = nav.SubscribeAny
				//
				// end of synthetic assignments
				o.Callback = asyncCallback(fmt.Sprintf("%v/%v", entry.given, entry.should))
			}
			return &nav.ResumeSession{
				Path:     resumeJSONPath,
				Restorer: restorer,
				Strategy: entry.resume.Strategy,
			}
		},
	)
}

var _ = Describe("navigation", Ordered, func() {
	var (
		root         string
		jroot        string
		fromJSONPath string
		jobsChOut    nav.TraverseItemJobStream
		outputCh     async.OutputStream[nav.TraverseOutput]
	)

	BeforeAll(func() {
		root = musico()
		jroot = helpers.JoinCwd("Test", "json")
		fromJSONPath = helpers.Path(jroot, "resume-state.json")
		outputCh = nil

		if err := Use(func(o *UseOptions) {
			o.Tag = DefaultLanguage.Get()
		}); err != nil {
			Fail(err.Error())
		}
	})

	BeforeEach(func() {
		if err := Use(func(o *UseOptions) {
			o.Tag = DefaultLanguage.Get()
		}); err != nil {
			Fail(err.Error())
		}
		jobsChOut = make(nav.TraverseItemJobStream, DefaultJobsChSize)
	})

	DescribeTable("async",
		func(ctx SpecContext, entry *asyncOkTE) {
			defer leaktest.Check(GinkgoT())()

			var (
				wgex async.WaitGroupEx
			)

			path := helpers.Path(root, "RETRO-WAVE")
			session := getSession(&entry.asyncTE, root, path, fromJSONPath)
			runner := session.Init()

			if entry.operator != nil {
				entry.operator(runner)
			}

			wgex = async.NewAnnotatedWaitGroup("🍂 traversal")
			wgex.Add(1, navigatorRoutineName)
			_, err := runner.Run(&nav.AsyncInfo{
				Ctx:                  ctx,
				NavigatorRoutineName: navigatorRoutineName,
				Adder:                wgex,
				Quitter:              wgex,
				WaitAQ:               wgex,
				JobsChanOut:          jobsChOut,
			})

			var consumer *Consumer[nav.TraverseOutput]
			if outputCh != nil {
				consumer = StartConsumer[nav.TraverseOutput](
					ctx,
					wgex,
					outputCh,
				)
			}
			wgex.Wait("👾 test-main")

			if consumer != nil {
				fmt.Printf("---> 📌📌📌 consumer.count: '%v'\n", consumer.Count)
			}

			Expect(err).To(BeNil())
		},
		func(entry *asyncOkTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'", entry.given, entry.should)
		},

		Entry(nil, &asyncOkTE{
			asyncTE: asyncTE{
				given:  "PrimarySession WithCPUPool",
				should: "run with context",
				operator: func(r nav.NavigationRunner) nav.NavigationRunner {
					return r.WithCPUPool()
				},
			},
		}, SpecTimeout(time.Second*2)),

		Entry(nil, &asyncOkTE{
			asyncTE: asyncTE{
				given:  "PrimarySession WithPool",
				should: "run with context",
				operator: func(r nav.NavigationRunner) nav.NavigationRunner {
					return r.WithPool(3)
				},
			},
		}, SpecTimeout(time.Second*2)),

		Entry(nil, &asyncOkTE{
			asyncTE: asyncTE{
				given:  "Fastward Resume WithCPUPool(universal: listen pending(logged)",
				should: "run with context",
				operator: func(r nav.NavigationRunner) nav.NavigationRunner {
					return r.WithCPUPool()
				},
				// 🔥 panic: send on closed channel; this is intermittent
				// probably a race condition
				//
				resume: &asyncResumeTE{
					Strategy: nav.ResumeStrategyFastwardEn,
					Listen:   nav.ListenPending,
				},
			},
		}, SpecTimeout(time.Second*2)),

		Entry(nil, &asyncOkTE{
			asyncTE: asyncTE{
				given:  "Spawn Resume WithPool(universal: listen not active/deaf)",
				should: "run with context",
				operator: func(r nav.NavigationRunner) nav.NavigationRunner {
					return r.WithPool(3)
				},
				resume: &asyncResumeTE{
					Strategy: nav.ResumeStrategySpawnEn,
					Listen:   nav.ListenDeaf,
				},
			},
		}, SpecTimeout(time.Second*2)),

		Entry(nil, &asyncOkTE{
			asyncTE: asyncTE{
				given:  "PrimarySession Consume",
				should: "output should be externally consumed",
				operator: func(r nav.NavigationRunner) nav.NavigationRunner {
					outputCh = nav.CreateTraverseOutputCh(3)
					return r.WithPool(4).Consume(outputCh)
				},
			},
		}, SpecTimeout(time.Second*2)),

		Entry(nil, &asyncOkTE{
			asyncTE: asyncTE{
				given:  "Fastward Resume Consume(universal: listen pending(logged)",
				should: "output should be externally consumed",
				operator: func(r nav.NavigationRunner) nav.NavigationRunner {
					outputCh = nav.CreateTraverseOutputCh(3)
					return r.WithPool(4).Consume(outputCh)
				},
				// 🔥 panic: send on closed channel;
				//
				resume: &asyncResumeTE{
					Strategy: nav.ResumeStrategyFastwardEn,
					Listen:   nav.ListenPending,
				},
			},
		}, SpecTimeout(time.Second*2)),
	)

	DescribeTable(
		"errors",
		func(ctx SpecContext, entry *asyncErrorTE) {
			defer leaktest.Check(GinkgoT())()

			defer func() {
				pe := recover()
				if err, ok := pe.(error); !ok {
					Fail(fmt.Sprintf("panic is not an error (%v)", err))
				} else if !strings.Contains(err.Error(),
					entry.fragment) {
					Fail(fmt.Sprintf("🔥 error (%v), does not contain expected fragment (%v)",
						err.Error(), entry.fragment))
				}
			}()

			if entry.fragment == "" {
				Fail("🔥 invalid test; error fragment not specified")
			}

			var wgex async.WaitGroupEx

			path := helpers.Path(root, "RETRO-WAVE")
			session := getSession(&entry.asyncTE, root, path, fromJSONPath)
			runner := session.Init()

			if entry.operator != nil {
				entry.operator(runner)
			}

			wgex = async.NewAnnotatedWaitGroup("🍂 traversal")
			wgex.Add(1, navigatorRoutineName)
			_, _ = runner.Run(&nav.AsyncInfo{
				Ctx:                  ctx,
				NavigatorRoutineName: navigatorRoutineName,
				Adder:                wgex,
				Quitter:              wgex,
				WaitAQ:               wgex,
				JobsChanOut:          jobsChOut,
			})

			Fail("❌ expected panic due to invalid async traversal setup")
		},
		func(entry *asyncErrorTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'", entry.given, entry.should)
		},

		Entry(nil, &asyncErrorTE{
			asyncTE: asyncTE{
				given:  "PrimarySession Consume, missing no of workers",
				should: "panic",
				operator: func(r nav.NavigationRunner) nav.NavigationRunner {
					outputCh = nav.CreateTraverseOutputCh(3)
					return r.Consume(outputCh)
				},
			},
			fragment: "worker pool acceleration not active",
		}, SpecTimeout(time.Second*2)),

		Entry(nil, &asyncErrorTE{
			asyncTE: asyncTE{
				given:  "Fastward Resume Consume(universal: listen pending(logged), missing no of workers",
				should: "panic",
				operator: func(r nav.NavigationRunner) nav.NavigationRunner {
					outputCh = nav.CreateTraverseOutputCh(3)
					return r.Consume(outputCh)
				},
				resume: &asyncResumeTE{
					Strategy: nav.ResumeStrategyFastwardEn,
					Listen:   nav.ListenPending,
				},
			},
			fragment: "worker pool acceleration not active",
		}, SpecTimeout(time.Second*2)),

		Entry(nil, &asyncErrorTE{
			asyncTE: asyncTE{
				given:  "Spawn Resume Consume(universal: listen not active/deaf), WithPool after Consume",
				should: "output should be externally consumed",
				operator: func(r nav.NavigationRunner) nav.NavigationRunner {
					outputCh = nav.CreateTraverseOutputCh(3)
					return r.Consume(outputCh).WithPool(4)
				},
				resume: &asyncResumeTE{
					Strategy: nav.ResumeStrategySpawnEn,
					Listen:   nav.ListenDeaf,
				},
			},
			fragment: "worker pool acceleration not active",
		}, SpecTimeout(time.Second*2)),
	)
})
