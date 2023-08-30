package nav

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/google/uuid"
	"github.com/snivilised/extendio/internal/log"
	"github.com/snivilised/extendio/xfs/utils"
	"github.com/snivilised/lorax/boost"
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
		Label: "boost decorator",
		Fn: func(item *TraverseItem) error {
			defer func() {
				pe := recover()
				if err, ok := pe.(error); !ok || !strings.Contains(err.Error(),
					"send on closed channel") {
					fmt.Printf("---> ☠️☠️☠️ ENSYNC-NAV-CALLBACK(panic on close): '%v' (err:'%v')\n",
						item.Path, pe,
					)
				}
			}()
			fmt.Printf("---> 🐬 ENSYNC-NAV-CALLBACK: '%v' \n", item.Path)

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

				case ai.JobsChanOut <- boost.Job[TraverseItemInput](j):
					//
					// intermittent panic: send on closed channel, in fastward resume scenarios
					// 'gr:observable-navigator'

					fmt.Printf("-->> 🍆🍆 sending job(%v)\n", j.ID)
				}
			}

			return err
		},
	}

	frame.decorate("boost decorator", decorator)
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
