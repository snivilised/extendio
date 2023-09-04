package nav

import (
	"context"
	"fmt"

	"github.com/snivilised/lorax/boost"
)

type navigationAccelerator struct {
	active      bool
	noWorkers   int
	outputChOut boost.OutputStream[TraverseOutput]
	pool        *boost.WorkerPool[TraverseItemInput, TraverseOutput]
}

func (a *navigationAccelerator) start(ai *AsyncInfo) {
	a.pool = boost.NewWorkerPool[TraverseItemInput, TraverseOutput](
		&boost.NewWorkerPoolParams[TraverseItemInput, TraverseOutput]{
			NoWorkers: a.noWorkers,
			Exec:      workerExecutive,
			JobsCh:    ai.JobsChanOut,
			WaitAQ:    ai.WaitAQ,
		})

	// We are handing over ownership of this channel (ai.OutputsChIn) to the pool as
	// its go routine will write to it, knows when no more data is available
	// and thus knows when to close it.
	//
	ai.WaitAQ.Add(1, a.pool.RoutineName)

	go a.pool.Start(ai.Context, a.outputChOut)
}

func (a *navigationAccelerator) finish(_ context.CancelFunc,
	ai *AsyncInfo,
) {
	// TODO: parentCancel, for the time being don't invoke the cancel func
	//
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
