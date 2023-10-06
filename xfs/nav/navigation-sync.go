package nav

import (
	"context"
	"fmt"
	"reflect"

	"github.com/snivilised/lorax/boost"
)

type NavigationSync interface {
	Run(callback sessionCallback, nc syncable, args ...any) (*TraverseResult, error)
}

type baseSync struct {
	session TraverseSession
}

func (s *baseSync) extract(args ...any) (bool, context.Context, context.CancelFunc) {
	// NB: extracting a context does not necessarily mean we want to run accelerated. Down the line,
	// it may be required that even for inline scenario's, we allow the client to pass in
	// a context for cancellation reasons.
	//
	var (
		ctx       context.Context
		cancel    context.CancelFunc
		extracted = len(args) > 0
	)

	for _, a := range args {
		switch argument := a.(type) {
		case context.Context:
			ctx = argument
		case context.CancelFunc:
			cancel = argument

		default:
			// TODO: convert to i18n error
			//
			panic(
				fmt.Errorf("extract found invalid type found in 'Run' arguments (val: '%v', type: '%v')",
					a,
					reflect.TypeOf(a).String(),
				),
			)
		}
	}

	return extracted, ctx, cancel
}

type inlineSync struct {
	baseSync
}

func (s *inlineSync) Run(callback sessionCallback, _ syncable, _ ...any) (*TraverseResult, error) {
	return callback()
}

type acceleratedSync struct {
	baseSync
	ai          *AsyncInfo
	noWorkers   int
	outputChOut boost.JobOutputStream[TraverseOutput]
	pool        *boost.WorkerPool[TraverseItemInput, TraverseOutput]
}

func (s *acceleratedSync) Run(callback sessionCallback, nc syncable, args ...any) (*TraverseResult, error) {
	defer s.finish(s.ai)

	extracted, ctx, cancel := s.extract(args...)

	if !extracted {
		panic("failed to obtain context")
	}

	nc.ensync(ctx, cancel, s.ai)
	s.start(ctx, cancel)

	return callback()
}

func (s *acceleratedSync) start(ctx context.Context, cancel context.CancelFunc) {
	s.pool = boost.NewWorkerPool[TraverseItemInput, TraverseOutput](
		&boost.NewWorkerPoolParams[TraverseItemInput, TraverseOutput]{
			NoWorkers: s.noWorkers,
			Exec:      workerExecutive,
			JobsCh:    s.ai.JobsChanOut,
			WaitAQ:    s.ai.WaitAQ,
		})

	// We are handing over ownership of this channel (ai.OutputsChIn) to the pool as
	// its go routine will write to it, knows when no more data is available
	// and thus knows when to close it.
	//
	s.ai.WaitAQ.Add(1, s.pool.RoutineName)

	go s.pool.Start(ctx, cancel, s.outputChOut)
}

func (s *acceleratedSync) finish(
	ai *AsyncInfo,
) {
	fmt.Printf("---> observable navigator 😈😈😈 defer session.finish (CLOSE(JobsChanOut)/QUIT)\n")
	close(ai.JobsChanOut) // ⚠️ fastward: intermittent panic on close
	ai.WaitAQ.Done(ai.NavigatorRoutineName)
}

func workerExecutive(job boost.Job[TraverseItemInput]) (boost.JobOutput[TraverseOutput], error) {
	err := job.Input.Fn(job.Input.Item)

	return boost.JobOutput[TraverseOutput]{
		Payload: TraverseOutput{
			Item:  job.Input.Item,
			Error: err,
		},
	}, err
}
