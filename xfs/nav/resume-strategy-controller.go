package nav

import (
	"github.com/samber/lo"
)

type resumeStrategyController struct {
	nc       *navigationController
	ps       *persistState
	strategy resumeStrategy
}

func (c *resumeStrategyController) run(ai ...*AsyncInfo) (*TraverseResult, error) {
	return c.strategy.resume(&strategyResumeInfo{
		nc: c.nc,
		ps: c.ps,
		ai: lo.TernaryF(len(ai) > 0, func() *AsyncInfo {
			return ai[0]
		}, func() *AsyncInfo {
			return nil
		}),
	})
}

func (c *resumeStrategyController) Save(path string) error {
	return c.nc.save(path)
}

func (c *resumeStrategyController) detach(frame *navigationFrame) {
	c.strategy.detach(frame)
}

func (c *resumeStrategyController) finish() error {
	return c.strategy.finish()
}
