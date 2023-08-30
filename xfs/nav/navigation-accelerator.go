package nav

import (
	"github.com/snivilised/lorax/async"
)

type navigationAccelerator struct {
	active      bool
	noWorkers   int
	outputChOut async.OutputStream[TraverseOutput]
	pool        *async.WorkerPool[TraverseItemInput, TraverseOutput]
}

func (a *navigationAccelerator) start(ai *AsyncInfo) {
	a.pool = async.NewWorkerPool[TraverseItemInput, TraverseOutput](
		&async.NewWorkerPoolParams[TraverseItemInput, TraverseOutput]{
			NoWorkers: a.noWorkers,
			Exec:      workerExecutive,
			JobsCh:    ai.JobsChanOut,
			Quitter:   ai.Quitter,
		})

	// We are handing over ownership of this channel (ai.OutputsChIn) to the pool as
	// its go routine will write to it, knows when no more data is available
	// and thus knows when to close it.
	//
	ai.Adder.Add(1, a.pool.RoutineName)

	go a.pool.Start(ai.Ctx, a.outputChOut)
}

func workerExecutive(job async.Job[TraverseItemInput]) (async.JobOutput[TraverseOutput], error) {
	err := job.Input.Fn(job.Input.Item)

	return async.JobOutput[TraverseOutput]{
		Payload: TraverseOutput{
			Item:  job.Input.Item,
			Error: err,
		},
	}, err
}
