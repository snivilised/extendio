package nav

import (
	_ "github.com/snivilised/extendio/translate"
)

type navigatorController struct {
	impl navigatorImpl
}

func (c *navigatorController) Walk(root string) *TraverseResult {
	o := c.impl.options()
	o.Notify.OnBegin(root)
	frame := navigationFrame{ // inject the frame in from new-navigator
		Root:   root,
		client: o.Callback,
	}
	bootstrapListener(o, &frame)

	result := &TraverseResult{
		Error: c.impl.top(&frame),
	}
	o.Notify.OnEnd(result)

	return result
}
