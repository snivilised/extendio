package nav

import "github.com/samber/lo"

// SubPathBehaviour
type SubPathBehaviour struct {
	KeepTrailingSep bool
}

// SortBehaviours
type SortBehaviours struct {
	IsCaseSensitive bool // case sensitive traversal order
}

// NavigationBehaviours
type NavigationBehaviours struct {
	SubPath SubPathBehaviour
	Sort    SortBehaviours
}

type Notifications struct {
	OnBegin   BeginHandler      // invoked at beginning of traversal
	OnEnd     EndHandler        // invoked at end of traversal
	OnDescend AscendancyHandler // handler to invoke as a folder is descended (before children)
	OnAscend  AscendancyHandler // handler to invoke as a folder is ascended (after children)
}

// TraverseOptions customise the way a directory tree is traversed
type TraverseOptions struct {
	Subscription TraverseSubscription // defines which node types are visited
	DoExtend     bool                 // request an extended result
	WithMetrics  bool                 // request metrics in TraversalResult
	Callback     TraverseCallback     // traversal callback (universal, folders, files)
	Notify       Notifications
	Hooks        TraverseHooks
	Behaviours   NavigationBehaviours
}
type TraverseOptionFn func(o *TraverseOptions) // functional traverse options

func composeTraverseOptions(fn ...TraverseOptionFn) *TraverseOptions {
	o := TraverseOptions{
		Subscription: SubscribeAny,
		DoExtend:     false,
		Notify: Notifications{
			OnBegin:   func(root string) {},
			OnEnd:     func(result *TraverseResult) {},
			OnDescend: func(item *TraverseItem) {},
			OnAscend:  func(item *TraverseItem) {},
		},
		Hooks: TraverseHooks{
			QueryStatus:   LstatHookFn,
			ReadDirectory: ReadEntries,
			FolderSubPath: RootParentSubPath,
			FileSubPath:   RootParentSubPath,
		},
		Behaviours: NavigationBehaviours{
			SubPath: SubPathBehaviour{
				KeepTrailingSep: true,
			},
			Sort: SortBehaviours{
				IsCaseSensitive: false,
			},
		},
	}

	for _, functionalOption := range fn {
		functionalOption(&o)
	}

	if o.Hooks.Sort == nil {
		o.Hooks.Sort = lo.Ternary(o.Behaviours.Sort.IsCaseSensitive,
			CaseSensitiveSortHookFn, CaseInSensitiveSortHookFn,
		)
	}

	if o.Hooks.Extend == nil {
		o.Hooks.Extend = lo.Ternary(o.DoExtend, DefaultExtendHookFn, nullExtendHookFn)
	}

	return &o
}
