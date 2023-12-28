package nav

import (
	"context"
	"fmt"
	"io/fs"
	"strings"

	"github.com/google/uuid"
	"github.com/snivilised/extendio/internal/log"
	"github.com/snivilised/extendio/xfs/utils"
)

type navigator struct {
	o                    *TraverseOptions
	agent                *navigationAgent
	log                  utils.RoProp[log.Logger]
	samplingActive       bool
	filteringActive      bool
	samplingFilterActive bool
	samplingCtrl         *samplingController
}

func (n *navigator) options() *TraverseOptions {
	return n.o
}

func (n *navigator) init(ns *NavigationState) {
	if n.samplingActive {
		adapters := createSamplingAdapters()
		n.samplingCtrl = &samplingController{
			o:        n.o,
			fn:       getSamplerControllerFunc(n.o),
			adapters: adapters,
		}

		samplingType := n.o.Store.Sampling.SampleType

		if (samplingType == SampleTypeFilterEn) || (samplingType == SampleTypeCustomEn) {
			n.samplingCtrl.init(ns)
		}
	}
}

func (n *navigator) ensync(
	ctx context.Context,
	_ context.CancelFunc, // we don't need this here; only the worker pool needs it!
	frame *navigationFrame,
	ai *AsyncInfo,
) {
	decorated := frame.client
	decorator := &LabelledTraverseCallback{
		Label: "boost decorator",
		Fn: func(item *TraverseItem) error {
			defer func() {
				if pe := recover(); pe != nil {
					if err, ok := pe.(error); ok || strings.Contains(err.Error(),
						"send on closed channel") {
						n.logger().Error("☠️☠️☠️ send on closed channel",
							log.String("item-path", item.Path),
						)
					} else {
						panic(pe)
					}
				}
			}()

			var err error
			select {
			case <-ctx.Done():
				err = fs.SkipDir
			default:
				job := TraverseItemJob{
					ID: fmt.Sprintf("JOB-ID:%v", uuid.NewString()),
					Input: TraverseItemInput{
						Item:  item,
						Fn:    decorated.Fn,
						Label: decorated.Label,
					},
					SequenceNo: -999,
				}

				select {
				case <-ctx.Done():
					err = fs.SkipDir

				case ai.JobsChanOut <- job:
					//
					// intermittent panic: send on closed channel, in fastward resume scenarios
					// 'gr:observable-navigator'
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

func (n *navigator) descend(navi *NavigationInfo) bool {
	if !navi.frame.periscope.descend(n.o.Store.Behaviours.Cascade.Depth) {
		return false
	}

	navi.frame.notifiers.descend.invoke(navi.Item)

	return true
}

func (n *navigator) ascend(navi *NavigationInfo, permit bool) {
	if permit {
		navi.frame.periscope.ascend()
		navi.frame.notifiers.ascend.invoke(navi.Item)
	}
}

func (n *navigator) finish() error {
	return n.log.Get().Sync()
}

func (n *navigator) keep(stash *inspection) {
	n.agent.keep(stash)
}
