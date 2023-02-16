package nav

import "github.com/snivilised/extendio/xfs/utils"

type navigationFrame struct {
	root        utils.VarProp[string]
	currentPath utils.VarProp[string]
	listener    *navigationListener
	raw         LabelledTraverseCallback // un-decorated (except for filter) client callback
	client      LabelledTraverseCallback // decorate-able client callback
	filters     *NavigationFilters
	notifiers   notificationsSink
	periscope   *navigationPeriscope
	metrics     *navigationMetrics
}

// attach/decorate
// ===============
// decorate: Really means decorate, ie there is a new decorator wrapper around the
// existing client. Requires new stack push. When the listener is initialised, it
// should be using a decorate method. However, when a filter is initialised, there
// should not be a stack push as that has nothing to do with filtering. Therefore,
// there are 2 types of decoration, with/without stack push. So you could say there
// is a filter decorate and a listen decorate. !!! But, frame.decorate is simple
// and does not do the stack push.
//
// attach: Depends on existing listener. The listener stack accounts for behaviour
// changes as pushes and pops occur. WE DON'T need a new decorator. attach/detach
// should issue push/pop to the stack.
//
// the 3 scenarios:
// - No listener: The client is considered raw, but is still subject to filtering
// - Listen active: The client is subject to the listen state: listener.current and
// as the state changes, then so does the decorator behaviour.
// - Resume: Uses the listen feature, but because there is a resume specific state,
// behaviour changes occur because of this particular state. However, the stack still
// needs a custom ListenOptions instance, therefore it requires a push.
//

func (f *navigationFrame) decorate(label string, decorator *LabelledTraverseCallback) *LabelledTraverseCallback {
	// this method doesn't do much, but it needs to be made explicit because it
	// is easy to setup the callback decoration chain incorrectly resulting in
	// stack overflow due to infinite recursion. Its easy to search when decoration is
	// occurring in the code base, just search for decorate or go to references.
	//
	previous := f.client
	f.client = *decorator

	return &previous
}

func (f *navigationFrame) save(active *ActiveState) {

	active.Root = f.root.Get()
	active.NodePath = f.currentPath.Get()
	active.Depth = f.periscope.depth()
	f.metrics.save(active)
}

func (f *navigationFrame) collate() *TraverseResult {

	return &TraverseResult{
		Metrics: &f.metrics._metrics,
	}
}

type linkParams struct {
	root    string
	current string
}

func (f *navigationFrame) link(params *linkParams) {
	// Combines information gleaned from the previous traversal that was
	// interrupted, into the resume traversal.
	//
	f.periscope.difference(params.root, params.current)
}

func (f *navigationFrame) reset() {
	f.metrics = navigationMetricsFactory{}.new()
}
