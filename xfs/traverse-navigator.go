package xfs

type TraverseSubscription uint

const (
	_ TraverseSubscription = iota
	SubscribeAny
	SubscribeFolders
	SubscribeFiles
)

type TraverseCallback func(item *TraverseItem) *LocalisableError

type TraverseOptions struct {
	Subscription    TraverseSubscription
	IsCaseSensitive bool // case sensitive traversal order
	Extend          bool // request an extended result
	WithMetrics     bool
	Callback        TraverseCallback
	OnDescend       TraverseCallback
	OnAscend        TraverseCallback
}
type TraverseOptionFn func(o *TraverseOptions)

func composeTraverseOptions(fn ...TraverseOptionFn) TraverseOptions {
	options := TraverseOptions{
		Subscription:    SubscribeAny,
		IsCaseSensitive: false,
		Extend:          false,
	}

	for _, functionalOption := range fn {
		functionalOption(&options)
	}
	return options
}

type TraverseResult struct {
}

type TraverseNavigator interface {
	Walk(root string) *TraverseResult
}

type navigatorSubject interface {
	top(root string) *TraverseResult
	traverse(currentItem *TraverseItem) *TraverseResult
}
