package nav

import (
	"context"

	"github.com/snivilised/lorax/boost"
)

type TraverseItemInput struct {
	Item *TraverseItem
	Fn   TraverseCallback
}
type TraverseItemJob = boost.Job[TraverseItemInput]
type TraverseItemJobStream = boost.JobStream[TraverseItemInput]
type TraverseItemJobStreamR = boost.JobStreamR[TraverseItemInput]
type TraverseItemJobStreamW = boost.JobStreamW[TraverseItemInput]

type TraverseOutput struct {
	Item  *TraverseItem
	Error error
}
type TraverseItemOutput = boost.JobOutput[TraverseOutput]
type TraverseItemOutputStream = boost.OutputStream[TraverseOutput]
type TraverseItemOutputStreamR = boost.OutputStreamR[TraverseOutput]
type TraverseItemOutputStreamW = boost.OutputStreamW[TraverseOutput]

type AsyncInfo struct {
	Context context.Context

	// this doesn't seem right, the client shouldn't have to specify
	// the routine name for the navigator; should be a readonly prop
	// of the navigator. Perhaps, it can be overridden by the user
	// here, but the navigator should have an internally defined default.
	//
	NavigatorRoutineName boost.GoRoutineName
	WaitAQ               boost.AnnotatedWgAQ

	// why are we passing in the jobs output channel here, rather than...
	// (perhaps a new operator, like how Consume() op takes the outputCh?)
	// The difference here though is that JobsChanOut is mandatory but
	// Consume & outputCh are optional.
	//
	JobsChanOut TraverseItemJobStream
}
