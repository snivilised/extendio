package xfs

import "github.com/samber/lo"

// TraverseOptions customise the way a directory tree in traversed
type TraverseOptions struct {
	Subscription    TraverseSubscription // defines which node types are visited
	IsCaseSensitive bool                 // case sensitive traversal order
	Extend          bool                 // request an extended result
	WithMetrics     bool                 // request metrics in TraversalResult
	Callback        TraverseCallback     // traversal callback (universal, folders, files)
	OnDescend       TraverseCallback     // callback to invoke as a folder is descended (before children)
	OnAscend        TraverseCallback     // callback to invoke as a folder is ascended (after children)

	Hooks traverseHooks
}
type TraverseOptionFn func(o *TraverseOptions) // functional traverse options

func composeTraverseOptions(fn ...TraverseOptionFn) *TraverseOptions {
	options := TraverseOptions{
		Subscription:    SubscribeAny,
		IsCaseSensitive: false,
		Extend:          false,
		Hooks: traverseHooks{
			ReadDirectory: readDir,
		},
	}

	for _, functionalOption := range fn {
		functionalOption(&options)
	}

	if options.Hooks.Sort == nil {
		options.Hooks.Sort = lo.Ternary(options.IsCaseSensitive,
			CaseSensitiveSortHookFn, CaseInSensitiveSortHookFn,
		)
	}

	return &options
}
