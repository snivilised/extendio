package nav

import (
	"github.com/samber/lo"
)

type resumeController struct {
	navigator *navigatorController
	ps        *persistState
	strategy  resumeStrategy
}

func (c *resumeController) Continue(ai ...*AsyncInfo) (*TraverseResult, error) {
	return c.strategy.resume(&strategyResumeInfo{
		nc: c.navigator,
		ps: c.ps,
		ai: lo.TernaryF(len(ai) > 0, func() *AsyncInfo {
			return ai[0]
		}, func() *AsyncInfo {
			return nil
		}),
	})
}

func (c *resumeController) Save(path string) error {
	return c.navigator.save(path)
}

func (c *resumeController) detach(frame *navigationFrame) {
	c.strategy.detach(frame)
}

func (c *resumeController) finish() error {
	return c.strategy.finish()
}
