package nav

import (
	"fmt"
	"io/fs"

	"github.com/google/uuid"
	"github.com/snivilised/extendio/internal/log"
	"github.com/snivilised/extendio/xfs/utils"
	"github.com/snivilised/lorax/async"
)

type navigator struct {
	o     *TraverseOptions
	agent *navigationAgent
	log   utils.RoProp[log.Logger]
}

func (n *navigator) options() *TraverseOptions {
	return n.o
}

func (n *navigator) ensync(frame *navigationFrame, ai *AsyncInfo) {
	decorated := frame.client
	decorator := &LabelledTraverseCallback{
		Label: "async decorator",
		Fn: func(item *TraverseItem) error {
			fmt.Printf("---> 🐬 ASYNC-CALLBACK: '%v' \n", item.Path)

			var err error
			select {
			case <-ai.Ctx.Done():
				err = fs.SkipDir
			default:
				j := TraverseItemJob{
					ID: fmt.Sprintf("JOB-ID:%v", uuid.NewString()),
					Input: TraverseItemInput{
						Item: item,
						Fn:   decorated.Fn,
					},
					SequenceNo: -999,
				}

				select {
				case <-ai.Ctx.Done():
					err = fs.SkipDir

				case ai.JobsChanOut <- async.Job[TraverseItemInput](j):
				}
			}

			return err
		},
	}

	frame.decorate("async decorator", decorator)
}

func (n *navigator) logger() log.Logger {
	return n.log.Get()
}

func (n *navigator) descend(navi *NavigationInfo) {
	navi.Frame.periscope.descend()
	navi.Frame.notifiers.descend.invoke(navi.Item)
}

func (n *navigator) ascend(navi *NavigationInfo) {
	navi.Frame.periscope.ascend()
	navi.Frame.notifiers.ascend.invoke(navi.Item)
}

func (n *navigator) finish() error {
	return n.log.Get().Sync()
}
