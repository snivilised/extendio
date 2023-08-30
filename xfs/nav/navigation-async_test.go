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
	"github.com/snivilised/lorax/boost"

	. "github.com/snivilised/extendio/i18n"
)

var (
	navigatorRoutineName = boost.GoRoutineName("✨ observable-navigator")
)

type (
	boostResumeTE struct {
		Strategy nav.ResumeStrategyEnum
		Listen   nav.ListeningState
	}

	operatorFunc func(r nav.NavigationRunner) nav.NavigationRunner

	boostTE struct {
		given    string
		should   string
		operator operatorFunc
		resume   *boostResumeTE
	}

	boostOkTE struct {
		boostTE
	}

	boostErrorTE struct {
		boostTE
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
	waitAQ      boost.AnnotatedWgAQ
	RoutineName boost.GoRoutineName
	OutputChIn  boost.OutputStream[O]
	Count       int
}

func StartConsumer[O any](
	ctx context.Context,
	waitAQ boost.AnnotatedWgAQ,
	outputChIn boost.OutputStream[O],
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

func getSession(entry *boostTE, root, path, resumeJSONPath string) nav.TraverseSession {
	getOptions := func(o *nav.TraverseOptions) {
		o.Store.Subscription = nav.SubscribeFolders
		o.Store.DoExtend = true
		o.Callback = boostCallback("boost primary session")
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
				o.Callback = boostCallback(fmt.Sprintf("%v/%v", entry.given, entry.should))
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
		outputCh     boost.OutputStream[nav.TraverseOutput]
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

	DescribeTable("boost",
		func(ctx SpecContext, entry *boostOkTE) {
			defer leaktest.Check(GinkgoT())()

			var (
				wgan boost.WaitGroupAn
			)

			path := helpers.Path(root, "RETRO-WAVE")
			session := getSession(&entry.boostTE, root, path, fromJSONPath)
			runner := session.Init()

			if entry.operator != nil {
				entry.operator(runner)
			}

			wgan = boost.NewAnnotatedWaitGroup("🍂 traversal")
			wgan.Add(1, navigatorRoutineName)
			_, err := runner.Run(&nav.AsyncInfo{
				Ctx:                  ctx,
				NavigatorRoutineName: navigatorRoutineName,
				WaitAQ:               wgan,
				JobsChanOut:          jobsChOut,
			})

			var consumer *Consumer[nav.TraverseOutput]
			if outputCh != nil {
				consumer = StartConsumer[nav.TraverseOutput](
					ctx,
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
		func(entry *boostOkTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'", entry.given, entry.should)
		},

		Entry(nil, &boostOkTE{
			boostTE: boostTE{
				given:  "PrimarySession WithCPUPool",
				should: "run with context",
				operator: func(r nav.NavigationRunner) nav.NavigationRunner {
					return r.WithCPUPool()
				},
			},
		}, SpecTimeout(time.Second*2)),

		Entry(nil, &boostOkTE{
			boostTE: boostTE{
				given:  "PrimarySession WithPool",
				should: "run with context",
				operator: func(r nav.NavigationRunner) nav.NavigationRunner {
					return r.WithPool(3)
				},
			},
		}, SpecTimeout(time.Second*2)),

		Entry(nil, &boostOkTE{
			boostTE: boostTE{
				given:  "Fastward Resume WithCPUPool(universal: listen pending(logged)",
				should: "run with context",
				operator: func(r nav.NavigationRunner) nav.NavigationRunner {
					return r.WithCPUPool()
				},
				// 🔥 panic: send on closed channel; this is intermittent
				// probably a race condition
				//
				resume: &boostResumeTE{
					Strategy: nav.ResumeStrategyFastwardEn,
					Listen:   nav.ListenPending,
				},
			},
		}, SpecTimeout(time.Second*2)),

		Entry(nil, &boostOkTE{
			boostTE: boostTE{
				given:  "Spawn Resume WithPool(universal: listen not active/deaf)",
				should: "run with context",
				operator: func(r nav.NavigationRunner) nav.NavigationRunner {
					return r.WithPool(3)
				},
				resume: &boostResumeTE{
					Strategy: nav.ResumeStrategySpawnEn,
					Listen:   nav.ListenDeaf,
				},
			},
		}, SpecTimeout(time.Second*2)),

		Entry(nil, &boostOkTE{
			boostTE: boostTE{
				given:  "PrimarySession Consume",
				should: "output should be externally consumed",
				operator: func(r nav.NavigationRunner) nav.NavigationRunner {
					outputCh = nav.CreateTraverseOutputCh(3)
					return r.WithPool(4).Consume(outputCh)
				},
			},
		}, SpecTimeout(time.Second*2)),

		Entry(nil, &boostOkTE{
			boostTE: boostTE{
				given:  "Fastward Resume Consume(universal: listen pending(logged)",
				should: "output should be externally consumed",
				operator: func(r nav.NavigationRunner) nav.NavigationRunner {
					outputCh = nav.CreateTraverseOutputCh(3)
					return r.WithPool(4).Consume(outputCh)
				},
				// 🔥 panic: send on closed channel;
				//
				resume: &boostResumeTE{
					Strategy: nav.ResumeStrategyFastwardEn,
					Listen:   nav.ListenPending,
				},
			},
		}, SpecTimeout(time.Second*2)),
	)

	DescribeTable(
		"errors",
		func(ctx SpecContext, entry *boostErrorTE) {
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

			var wgan boost.WaitGroupAn

			path := helpers.Path(root, "RETRO-WAVE")
			session := getSession(&entry.boostTE, root, path, fromJSONPath)
			runner := session.Init()

			if entry.operator != nil {
				entry.operator(runner)
			}

			wgan = boost.NewAnnotatedWaitGroup("🍂 traversal")
			wgan.Add(1, navigatorRoutineName)
			_, _ = runner.Run(&nav.AsyncInfo{
				Ctx:                  ctx,
				NavigatorRoutineName: navigatorRoutineName,
				WaitAQ:               wgan,
				JobsChanOut:          jobsChOut,
			})

			Fail("❌ expected panic due to invalid boost traversal setup")
		},
		func(entry *boostErrorTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'", entry.given, entry.should)
		},

		Entry(nil, &boostErrorTE{
			boostTE: boostTE{
				given:  "PrimarySession Consume, missing no of workers",
				should: "panic",
				operator: func(r nav.NavigationRunner) nav.NavigationRunner {
					outputCh = nav.CreateTraverseOutputCh(3)
					return r.Consume(outputCh)
				},
			},
			fragment: "worker pool acceleration not active",
		}, SpecTimeout(time.Second*2)),

		Entry(nil, &boostErrorTE{
			boostTE: boostTE{
				given:  "Fastward Resume Consume(universal: listen pending(logged), missing no of workers",
				should: "panic",
				operator: func(r nav.NavigationRunner) nav.NavigationRunner {
					outputCh = nav.CreateTraverseOutputCh(3)
					return r.Consume(outputCh)
				},
				resume: &boostResumeTE{
					Strategy: nav.ResumeStrategyFastwardEn,
					Listen:   nav.ListenPending,
				},
			},
			fragment: "worker pool acceleration not active",
		}, SpecTimeout(time.Second*2)),

		Entry(nil, &boostErrorTE{
			boostTE: boostTE{
				given:  "Spawn Resume Consume(universal: listen not active/deaf), WithPool after Consume",
				should: "output should be externally consumed",
				operator: func(r nav.NavigationRunner) nav.NavigationRunner {
					outputCh = nav.CreateTraverseOutputCh(3)
					return r.Consume(outputCh).WithPool(4)
				},
				resume: &boostResumeTE{
					Strategy: nav.ResumeStrategySpawnEn,
					Listen:   nav.ListenDeaf,
				},
			},
			fragment: "worker pool acceleration not active",
		}, SpecTimeout(time.Second*2)),
	)
})
