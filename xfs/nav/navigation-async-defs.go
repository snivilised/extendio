package nav

import (
	"context"

	"github.com/snivilised/lorax/boost"
)

type TraverseItemInput struct {
	Item *TraverseItem
	Fn   TraverseCallback
}
type TraverseItemJob boost.Job[TraverseItemInput]
type TraverseItemJobStream boost.JobStream[TraverseItemInput]
type TraverseItemJobStreamR boost.JobStreamR[TraverseItemInput]
type TraverseItemJobStreamW boost.JobStreamW[TraverseItemInput]

type TraverseOutput struct {
	Item  *TraverseItem
	Error error
}
type TraverseItemOutput boost.JobOutput[TraverseOutput]
type TraverseItemOutputStream boost.OutputStream[TraverseOutput]
type TraverseItemOutputStreamR boost.OutputStreamR[TraverseOutput]
type TraverseItemOutputStreamW boost.OutputStreamW[TraverseOutput]

type AsyncInfo struct {
	Ctx                  context.Context
	NavigatorRoutineName boost.GoRoutineName
	Adder                boost.AnnotatedWgAdder
	Quitter              boost.AnnotatedWgQuitter
	WaitAQ               boost.AnnotatedWgAQ
	JobsChanOut          TraverseItemJobStream
}
