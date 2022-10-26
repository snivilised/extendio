package xfs

import "github.com/samber/lo"

// TraverseOptions customise the way a directory tree in traversed
type TraverseOptions struct {
	Subscription    TraverseSubscription // defines which node types are visited
	IsCaseSensitive bool                 // case sensitive traversal order
	DoExtend        bool                 // request an extended result
	WithMetrics     bool                 // request metrics in TraversalResult
	Callback        TraverseCallback     // traversal callback (universal, folders, files)
	OnDescend       AscendancyHandler    // handler to invoke as a folder is descended (before children)
	OnAscend        AscendancyHandler    // handler to invoke as a folder is ascended (after children)

	Hooks TraverseHooks
}
type TraverseOptionFn func(o *TraverseOptions) // functional traverse options

func composeTraverseOptions(fn ...TraverseOptionFn) *TraverseOptions {
	o := TraverseOptions{
		Subscription:    SubscribeAny,
		IsCaseSensitive: false,
		DoExtend:        false,
		OnDescend:       func(item *TraverseItem) {},
		OnAscend:        func(item *TraverseItem) {},
		Hooks: TraverseHooks{
			QueryStatus:   LstatHookFn,
			ReadDirectory: ReadEntries,
		},
	}

	for _, functionalOption := range fn {
		functionalOption(&o)
	}

	if o.Hooks.Sort == nil {
		o.Hooks.Sort = lo.Ternary(o.IsCaseSensitive,
			CaseSensitiveSortHookFn, CaseInSensitiveSortHookFn,
		)
	}

	if o.Hooks.Extend == nil {
		o.Hooks.Extend = lo.Ternary(o.DoExtend, DefaultExtendHookFn, nullExtendHookFn)
	}

	return &o
}
