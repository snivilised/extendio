package nav

import (
	"context"
	"fmt"
	"runtime"

	"github.com/samber/lo"
	"github.com/snivilised/lorax/boost"
)

const (
	MinNoWorkers = 1
	MaxNoWorkers = 100
)

type RunnerOperators interface {
	WithCPUPool() NavigationRunner
	WithPool(now int) NavigationRunner
	Consume(outputCh boost.OutputStream[TraverseOutput]) NavigationRunner
}

type NavigationRunner interface {
	RunnerOperators
	Run(ai ...*AsyncInfo) (*TraverseResult, error)
}

type sessionRunner struct {
	session     TraverseSession
	accelerator navigationAccelerator
}

// Run executes the traversal session
func (r *sessionRunner) Run(ai ...*AsyncInfo) (*TraverseResult, error) {
	if r.accelerator.active {
		var fakeCancel context.CancelFunc // This is not real, don't invoke it
		defer r.accelerator.finish(fakeCancel, ai[0])

		r.accelerator.start(ai[0])

		return r.session.Run(ai[0])
	}

	return r.session.Run()
}

func (r *sessionRunner) WithPool(now int) NavigationRunner {
	r.accelerator.active = true
	if now >= MinNoWorkers && now <= MaxNoWorkers {
		r.accelerator.noWorkers = now
	} else {
		// TODO: turn this into an i18n error
		panic(fmt.Errorf("no of workers requested (%v) is out of range ('%v' - '%v')",
			now, MinNoWorkers, MaxNoWorkers),
		)
	}

	return r
}

func (r *sessionRunner) WithCPUPool() NavigationRunner {
	r.accelerator.active = true
	r.accelerator.noWorkers = runtime.NumCPU()

	return r
}

func CreateTraverseOutputCh(outputChSize int) boost.OutputStream[TraverseOutput] {
	return lo.TernaryF(outputChSize > 0,
		func() boost.OutputStream[TraverseOutput] {
			return make(boost.OutputStream[TraverseOutput], outputChSize)
		},
		func() boost.OutputStream[TraverseOutput] {
			return nil
		},
	)
}

func (r *sessionRunner) Consume(outputCh boost.OutputStream[TraverseOutput]) NavigationRunner {
	if !r.accelerator.active {
		// TODO: turn this into an i18n error
		panic(fmt.Errorf(
			"worker pool acceleration not active; ensure With(CPU)Pool specified before Consume",
		))
	}

	r.accelerator.outputChOut = outputCh

	return r
}

type primaryRunner struct {
	sessionRunner
}

type resumeRunner struct {
	sessionRunner
}
