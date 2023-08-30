package nav

import "github.com/snivilised/lorax/boost"

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

	go a.pool.Start(ai.Ctx, a.outputChOut)
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
