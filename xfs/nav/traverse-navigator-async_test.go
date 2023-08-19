package nav_test

import (
	"fmt"
	"sync"
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

type operatorFunc func(r nav.NavigationRunner) nav.NavigationRunner

type asyncResumeTE struct {
	Strategy nav.ResumeStrategyEnum
	Listen   nav.ListeningState
}

type asyncTE struct {
	given    string
	should   string
	operator operatorFunc
	resume   *asyncResumeTE
}

const (
	// we use a large job queue size to prevent blocking as these unit
	// tests don't have a consumer
	JobsChSize    = 50
	OutputsChSize = 20
)

// TODO: rename this file to navigation-async_test.go
var _ = Describe("navigation", Ordered, func() {
	var (
		root         string
		jroot        string
		fromJSONPath string
		jobsChOut    nav.TraverseItemJobStream
		outputsChIn  async.OutputStreamW[nav.TraverseOutput]
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
		jobsChOut = make(nav.TraverseItemJobStream, JobsChSize)
		outputsChIn = make(async.OutputStreamW[nav.TraverseOutput], OutputsChSize)
	})

	DescribeTable("async",
		func(ctx SpecContext, entry *asyncTE) {
			defer leaktest.Check(GinkgoT())()

			var wg sync.WaitGroup

			path := helpers.Path(root, "RETRO-WAVE")
			optionFn := func(o *nav.TraverseOptions) {
				o.Store.Subscription = nav.SubscribeFolders
				o.Store.DoExtend = true
				o.Callback = asyncCallback("WithCPUPool/primary session")
				o.Notify.OnBegin = begin("🛡️")
			}

			session := lo.TernaryF(entry.resume == nil,
				func() nav.TraverseSession {
					return &nav.PrimarySession{
						Path:     path,
						OptionFn: optionFn,
					}
				},
				func() nav.TraverseSession {
					restorer := func(o *nav.TraverseOptions, active *nav.ActiveState) {
						// synthetic assignments: The client should not perform these
						// types of assignments. Only being done here for testing purposes
						// to avoid the need to create many restore files
						// (eg resume-state.json) for different test cases.
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
						Path:     fromJSONPath,
						Restorer: restorer,
						Strategy: entry.resume.Strategy,
					}
				},
			)

			runner := session.Init()
			if entry.operator != nil {
				entry.operator(runner)
			}

			_, err := runner.Run(&nav.AsyncInfo{
				Ctx:          ctx,
				Wg:           &wg,
				JobsChanOut:  jobsChOut,
				OutputsChOut: outputsChIn,
			})

			wg.Wait()
			Expect(err).To(BeNil())
		},
		func(entry *asyncTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'", entry.given, entry.should)
		},

		Entry(nil, &asyncTE{
			given:  "PrimarySession WithCPUPool",
			should: "run with context",
			operator: func(r nav.NavigationRunner) nav.NavigationRunner {
				return r.WithCPUPool()
			},
		}, SpecTimeout(time.Second*2)),

		Entry(nil, &asyncTE{
			given:  "PrimarySession WithPool",
			should: "run with context",
			operator: func(r nav.NavigationRunner) nav.NavigationRunner {
				return r.WithPool(3)
			},
		}, SpecTimeout(time.Second*2)),

		XEntry(nil, &asyncTE{
			// 🔥 panic: send on closed channel
			//
			resume: &asyncResumeTE{
				Strategy: nav.ResumeStrategyFastwardEn,
				Listen:   nav.ListenPending,
			},
			given:  "Fastward Resume WithCPUPool(universal: listen pending(logged)",
			should: "run with context",
			operator: func(r nav.NavigationRunner) nav.NavigationRunner {
				return r.WithCPUPool()
			},
		}, SpecTimeout(time.Second*2)),

		Entry(nil, &asyncTE{
			resume: &asyncResumeTE{
				Strategy: nav.ResumeStrategySpawnEn,
				Listen:   nav.ListenDeaf,
			},
			given:  "Spawn Resume WithPool(universal: listen not active/deaf)",
			should: "run with context",
			operator: func(r nav.NavigationRunner) nav.NavigationRunner {
				return r.WithPool(3)
			},
		}, SpecTimeout(time.Second*2)),
	)
})
