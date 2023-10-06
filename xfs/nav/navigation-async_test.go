package nav_test

import (
	"context"
	"fmt"
	"time"

	"github.com/fortytw2/leaktest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"

	"github.com/snivilised/extendio/internal/helpers"
	"github.com/snivilised/extendio/xfs/nav"
	"github.com/snivilised/lorax/boost"

	. "github.com/snivilised/extendio/i18n"
)

var (
	navigatorRoutineName = boost.GoRoutineName("✨ observable-navigator")
)

type (
	asyncResumeTE struct {
		Strategy nav.ResumeStrategyEnum
		Listen   nav.ListeningState
	}

	acceleratorFunc func(a nav.AccelerationOperators) nav.AccelerationOperators

	asyncTE struct {
		given  string
		should string
		acc    acceleratorFunc
		resume *asyncResumeTE
	}

	asyncOkTE struct {
		asyncTE
	}

	asyncErrorTE struct {
		asyncTE
		fragment string
		timeout  time.Duration
	}
)

const (
	// we use a large job queue size to prevent blocking as these unit
	// tests don't have a consumer
	DefaultJobsChSize    = 50
	DefaultOutputsChSize = 50
)

type Consumer[O any] struct {
	waitAQ      boost.AnnotatedWgAQ
	RoutineName boost.GoRoutineName
	OutputChIn  boost.JobOutputStream[O]
	Count       int
}

func StartConsumer[O any](
	ctx context.Context,
	waitAQ boost.AnnotatedWgAQ,
	outputChIn boost.JobOutputStream[O],
) *Consumer[O] {
	consumer := &Consumer[O]{
		waitAQ:      waitAQ,
		RoutineName: boost.GoRoutineName("💠 consumer"),
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

func getRunner(entry *asyncTE, root, path, resumeJSONPath string) nav.NavigationRunner {
	return lo.TernaryF(entry.resume == nil,
		func() nav.NavigationRunner {
			return nav.New().Primary(&nav.Prime{
				Path: path,
				OptionsFn: func(o *nav.TraverseOptions) {
					o.Store.Subscription = nav.SubscribeFolders
					o.Store.DoExtend = true
					o.Callback = boostCallback("boost primary session")
					o.Notify.OnBegin = begin("🛡️")
				},
			})
		},
		func() nav.NavigationRunner {
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
				o.Callback = boostCallback(fmt.Sprintf("%v/%v", entry.given, entry.should))
			}

			return nav.New().Resume(&nav.Resumption{
				RestorePath: resumeJSONPath,
				Restorer:    restorer,
				Strategy:    entry.resume.Strategy,
			})
		},
	)
}

var _ = Describe("navigation", Ordered, func() {
	var (
		root         string
		jroot        string
		fromJSONPath string
		jobsChOut    nav.TraverseItemJobStream
		outputCh     boost.JobOutputStream[nav.TraverseOutput]
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
		func(ctxSpec SpecContext, entry *asyncOkTE) {
			defer leaktest.Check(GinkgoT())()

			var (
				consumer *Consumer[nav.TraverseOutput]
			)

			ctx, cancel := context.WithCancel(ctxSpec)
			path := helpers.Path(root, "RETRO-WAVE")
			runner := getRunner(&entry.asyncTE, root, path, fromJSONPath)
			wgan := boost.NewAnnotatedWaitGroup("🍂 traversal")
			wgan.Add(1, navigatorRoutineName)

			runner.WithPool(
				&nav.AsyncInfo{
					Context:              ctx,
					NavigatorRoutineName: navigatorRoutineName,
					WaitAQ:               wgan,
					JobsChanOut:          jobsChOut,
				},
			)

			if entry.acc != nil {
				entry.acc(runner)
			}

			_, err := runner.Run(ctx, cancel)

			if outputCh != nil {
				consumer = StartConsumer[nav.TraverseOutput](
					ctxSpec,
					wgan,
					outputCh,
				)
			}
			wgan.Wait("👾 test-main")

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
				acc: func(a nav.AccelerationOperators) nav.AccelerationOperators {
					return a // the default is like CPUPool
				},
			},
		}, SpecTimeout(time.Second*2)),

		Entry(nil, &asyncOkTE{
			asyncTE: asyncTE{
				given:  "PrimarySession WithPool",
				should: "run with context",
				acc: func(a nav.AccelerationOperators) nav.AccelerationOperators {
					return a.NoW(3)
				},
			},
		}, SpecTimeout(time.Second*2)),

		Entry(nil, &asyncOkTE{
			asyncTE: asyncTE{
				given:  "Fastward Resume WithCPUPool(universal: listen pending(logged)",
				should: "run with context",
				acc: func(a nav.AccelerationOperators) nav.AccelerationOperators {
					return a
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
				acc: func(a nav.AccelerationOperators) nav.AccelerationOperators {
					return a.NoW(3)
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
				should: "enable output to be consumed externally",
				acc: func(a nav.AccelerationOperators) nav.AccelerationOperators {
					return a.NoW(4).Consume(outputCh)
				},
			},
		}, SpecTimeout(time.Second*2)),

		Entry(nil, &asyncOkTE{
			asyncTE: asyncTE{
				given:  "Fastward Resume Consume(universal: listen pending(logged)",
				should: "enable output to be consumed externally",
				acc: func(a nav.AccelerationOperators) nav.AccelerationOperators {
					outputCh = nav.CreateTraverseOutputCh(3)
					return a.NoW(4).Consume(outputCh)
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
})
