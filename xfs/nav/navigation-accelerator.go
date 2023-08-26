package nav

import (
	"github.com/snivilised/lorax/async"
)

type navigationAccelerator struct {
	noWorkers int
	pool      *async.WorkerPool[TraverseItemInput, TraverseOutput]
}

func (a *navigationAccelerator) start(ai *AsyncInfo) {
	a.pool = async.NewWorkerPool[TraverseItemInput, TraverseOutput](
		&async.NewWorkerPoolParams[TraverseItemInput, TraverseOutput]{
			NoWorkers: a.noWorkers,
			Exec:      traverseExecutive,
			JobsCh:    ai.JobsChanOut,
			Quitter:   ai.Wgex,
		})

	// We are handing over ownership of this channel (ai.OutputsChIn) to the pool as
	// its go routine will write to it and knows when no more data is available
	// and thus knows when to close it.
	//
	go a.pool.Start(ai.Ctx, ai.OutputsChOut)

	ai.Adder.Add(1, a.pool.RoutineName)
}

func traverseExecutive(job async.Job[TraverseItemInput]) (async.JobOutput[TraverseOutput], error) {
	err := job.Input.Fn(job.Input.Item)

	return async.JobOutput[TraverseOutput]{
		Payload: TraverseOutput{
			Item:  job.Input.Item,
			Error: err,
		},
	}, err
}
