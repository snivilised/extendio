package nav

import (
	"github.com/samber/lo"
	"github.com/snivilised/extendio/xfs/utils"
)

type navigationFrame struct {
	root        utils.VarProp[string]
	currentPath utils.VarProp[string]
	listener    *navigationListener
	raw         *LabelledTraverseCallback // un-decorated (except for filter) client callback
	client      *LabelledTraverseCallback // decorate-able client callback
	filters     *NavigationFilters
	notifiers   notificationsSink
	periscope   *navigationPeriscope
	metrics     *NavigationMetrics
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
// needs a custom ListenTriggers instance, therefore it requires a push.
//

func (f *navigationFrame) decorate(_ string, decorator *LabelledTraverseCallback) {
	// this method doesn't do much, but it needs to be made explicit because it
	// is easy to setup the callback decoration chain incorrectly resulting in
	// stack overflow due to infinite recursion. Its easy to search when decoration is
	// occurring in the code base, just search for decorate or go to references.
	//
	f.client = decorator
}

func (f *navigationFrame) save(active *ActiveState) {
	active.Root = f.root.Get()
	active.NodePath = f.currentPath.Get()
	active.Depth = f.periscope.depth()
	f.metrics.save(active)
}

func (f *navigationFrame) collate() *TraverseResult {
	return &TraverseResult{
		Metrics: f.metrics,
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

func (f *navigationFrame) proxy(item *TraverseItem, compoundCounts *compoundCounters) error {
	// proxy is the correct way to invoke the client callback, because it takes into
	// account any active decorations such as listening and filtering. It should be noted
	// that the Callback on the options represents the client defined function which
	// can be decorated. Only the callback on the frame should ever be invoked.
	//
	err := f.invoke(item, compoundCounts)

	return lo.Ternary(item.Error != nil, item.Error, err)
}

func (f *navigationFrame) invoke(item *TraverseItem, compoundCounts *compoundCounters) error {
	f.currentPath.Set(item.Path)
	err := f.client.Fn(item)
	f.track(item, compoundCounts)

	return err
}

func (f *navigationFrame) track(item *TraverseItem, compoundCounts *compoundCounters) {
	isDirectory := item.IsDir()

	if item.Error == nil {
		if item.filteredOut {
			metricEn := lo.Ternary(isDirectory,
				MetricNoFoldersFilteredOutEn,
				MetricNoFilesFilteredOutEn,
			)
			f.metrics.tick(metricEn)
		} else {
			metricEn := lo.Ternary(isDirectory,
				MetricNoFoldersInvokedEn,
				MetricNoFilesInvokedEn,
			)
			f.metrics.tick(metricEn)
		}

		if compoundCounts != nil {
			f.metrics.post(
				MetricNoChildFilesFoundEn,
				compoundCounts.filteredIn,
			)
			f.metrics.post(
				MetricNoChildFilesFilteredOutEn,
				compoundCounts.filteredOut,
			)
		}
	}
}
