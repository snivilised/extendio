package nav

import (
	"context"

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
	Ctx                  context.Context
	NavigatorRoutineName async.GoRoutineName
	Wgex                 async.WaitGroupEx
	Adder                async.AssistedAdder
	Quitter              async.AssistedQuitter
	JobsChanOut          TraverseItemJobStream
	OutputsChOut         async.OutputStreamW[TraverseOutput] // consume this???
}
