package nav

import (
	_ "github.com/snivilised/extendio/translate"
)

type navigatorController struct {
	impl navigatorImpl
}

func (c *navigatorController) Walk(root string) *TraverseResult {
	o := c.impl.options()
	frame := navigationFrame{
		Root:   root,
		client: o.Callback,
	}

	bootstrapFilter(o, &frame)
	bootstrapListener(o, &frame)

	state := &NavigationState{Root: root, Filters: frame.filters}
	o.Notify.OnBegin(state)

	result := &TraverseResult{
		Error: c.impl.top(&frame),
	}
	o.Notify.OnEnd(result)

	return result
}
