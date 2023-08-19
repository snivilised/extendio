package nav

import (
	"context"
	"sync"

	"github.com/snivilised/lorax/async"
)

type TraverseItemInput struct {
	Item *TraverseItem
	Fn   TraverseCallback
}
type TraverseItemJob async.Job[TraverseItemInput]
type TraverseItemJobStream async.JobStream[TraverseItemInput]
type TraverseItemJobStreamR async.JobStreamR[TraverseItemInput]
type TraverseItemJobStreamW async.JobStreamW[TraverseItemInput]

type TraverseOutput struct {
	Item  *TraverseItem
	Error error
}

type TraverseItemJobOutput async.JobOutput[TraverseOutput]
type TraverseItemStream async.JobStream[TraverseItemJobOutput]
type TraverseItemStreamR async.JobStreamR[TraverseItemJobOutput]
type TraverseItemStreamW async.JobStreamW[TraverseItemJobOutput]

type AsyncInfo struct {
	Ctx          context.Context
	Wg           *sync.WaitGroup
	JobsChanOut  TraverseItemJobStream
	OutputsChOut async.OutputStreamW[TraverseOutput] // consume this???
	NoWorkers    int
}
