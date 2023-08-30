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
type TraverseItemOutput async.JobOutput[TraverseOutput]
type TraverseItemOutputStream async.OutputStream[TraverseOutput]
type TraverseItemOutputStreamR async.OutputStreamR[TraverseOutput]
type TraverseItemOutputStreamW async.OutputStreamW[TraverseOutput]

type AsyncInfo struct {
	Ctx                  context.Context
	NavigatorRoutineName async.GoRoutineName
	Adder                async.AnnotatedWgAdder
	Quitter              async.AnnotatedWgQuitter
	WaitAQ               async.AnnotatedWgAQ
	JobsChanOut          TraverseItemJobStream
}
