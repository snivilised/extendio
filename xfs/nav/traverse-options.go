package nav

import (
	"github.com/mohae/deepcopy"
	"github.com/samber/lo"
)

// A note about TraverseOptions/MARSHAL_EXCLUSIONS:
// TraverseOptionsAsJSON needs to be kept in sync with TraverseOptions,
// so when the former changes, then so should the latter. The only exception
// to this rule is if the new option is not serialisable, eg it is a function
// or an interface. A unit test (traverse-options-marshall) has been defined
// to enforce this so if an addition is made to TraverseOptions without the
// corresponding change to TraverseOptionsAsJSON, then the test should fail.
// If a non serialisable option is being added, then its name should be added
// to GetMarshalOptionsExclusions.
//

// SubPathBehaviour
type SubPathBehaviour struct {
	KeepTrailingSep bool
}

// SortBehaviour
type SortBehaviour struct {
	// case sensitive traversal order
	//
	IsCaseSensitive bool

	// DirectoryEntryOrder defines whether a folder's files or directories
	// should be navigated first.
	//
	DirectoryEntryOrder DirectoryEntryOrderEnum
}

// NavigationBehaviours
type NavigationBehaviours struct {
	// SubPath, behaviours relating to handling of sub-path calculation
	//
	SubPath SubPathBehaviour

	// Sort, behaviours relating to sorting of a folder's directory entries.
	//
	Sort SortBehaviour

	// Listen, behaviours relating to listen functionality.
	//
	Listen ListenBehaviour
}

// Notifications
type Notifications struct {
	// OnBegin invoked at beginning of traversal
	///
	OnBegin BeginHandler

	// OnEnd invoked at end of traversal
	//
	OnEnd EndHandler

	// OnDescend handler to invoke as a folder is descended (before children)
	//
	OnDescend AscendancyHandler

	// OnAscend handler to invoke as a folder is ascended (after children)
	OnAscend AscendancyHandler

	// OnStart handler invoked when start listening condition met if enabled
	//
	OnStart ListenHandler

	// OnStop handler invoked when finish listening condition met if enabled
	//
	OnStop ListenHandler
}

type FilterDefinitions struct {
	// Current denotes the filter object that represents the current file system item
	// being visited.
	//
	Current FilterDef

	// Children denotes the compound filter that is applied to the direct descendants
	// of the current file system item being visited.
	//
	Children CompoundFilterDef
}

type NavigationFilters struct {
	// Current denotes the filter object that represents the Current file system item
	// being visited.
	//
	Current TraverseFilter

	// Compound denotes the Compound filter that is applied to the direct descendants
	// of the current file system item being visited.
	//
	Compound CompoundTraverseFilter
}

// NavigationState carries information about navigation that client may be
// interested in and permitted to access, as opposed to the navigationFrame
// which is meant for internal purposes only.
type NavigationState struct {
	Root    string
	Filters NavigationFilters
}

// PersistOptions contains options for persisting traverse options
type PersistOptions struct {
	Format   PersistenceFormatEnum
	Restorer OptionsRestorer `json:"-"`
}

// TraverseOptions customise the way a directory tree is traversed
type TraverseOptions struct {
	// Subscription defines which node types are visited
	//
	Subscription TraverseSubscription

	// DoExtend request an extended result.
	//
	DoExtend bool

	// WithMetrics request metrics in TraversalResult.
	//
	WithMetrics bool

	// Callback function to invoke for every item visited in the file system.
	//
	Callback TraverseCallback `json:"-"`

	// Notify collection of notification function.
	//
	Notify Notifications `json:"-"`

	// TraverseHooks collection of hook functions, that can be overridden.
	//
	Hooks TraverseHooks `json:"-"`

	// Behaviours collection of behaviours that adjust the way navigation occurs,
	// that can be tweaked by the client.
	//
	Behaviours NavigationBehaviours

	// Listen options that control when listening state starts and finishes.
	//
	Listen ListenOptions `json:"-"`

	// FilterDefs definitions of filters that restricts for which file system nodes the
	// Callback is invoked for.
	//
	FilterDefs FilterDefinitions

	// Persist contains options for persisting traverse options
	//
	Persist PersistOptions
}

func GetMarshalOptionsExclusions() []string {
	return []string{"Callback", "Notify", "Hooks", "Listen", "Persist"}
}

// TraverseOptionFn functional traverse options
type TraverseOptionFn func(o *TraverseOptions)

func composeTraverseOptions(fn ...TraverseOptionFn) *TraverseOptions {
	o := GetDefaultOptions()

	for _, functionalOption := range fn {
		functionalOption(o)
	}

	if o.Hooks.Sort == nil {
		o.Hooks.Sort = lo.Ternary(o.Behaviours.Sort.IsCaseSensitive,
			CaseSensitiveSortHookFn, CaseInSensitiveSortHookFn,
		)
	}

	if o.Hooks.Extend == nil {
		o.Hooks.Extend = lo.Ternary(o.DoExtend, DefaultExtendHookFn, nullExtendHookFn)
	}

	return o
}

func (o *TraverseOptions) Clone() *TraverseOptions {
	clone := deepcopy.Copy(o)
	return clone.(*TraverseOptions)
}

func GetDefaultOptions() *TraverseOptions {
	return &TraverseOptions{
		Subscription: SubscribeAny,
		DoExtend:     false,
		Notify: Notifications{
			OnBegin:   func(state *NavigationState) {},
			OnEnd:     func(result *TraverseResult) {},
			OnDescend: func(item *TraverseItem) {},
			OnAscend:  func(item *TraverseItem) {},
		},
		Hooks: TraverseHooks{
			QueryStatus:   LstatHookFn,
			ReadDirectory: ReadEntries,
			FolderSubPath: RootParentSubPath,
			FileSubPath:   RootParentSubPath,
			InitFilters:   InitFiltersHookFn,
		},
		Behaviours: NavigationBehaviours{
			SubPath: SubPathBehaviour{
				KeepTrailingSep: true,
			},
			Sort: SortBehaviour{
				IsCaseSensitive:     false,
				DirectoryEntryOrder: DirectoryEntryOrderFoldersFirstEn,
			},
			Listen: ListenBehaviour{
				InclusiveStart: true,
				InclusiveStop:  false,
			},
		},
		Listen: ListenOptions{
			Start: nil,
			Stop:  nil,
		},
		FilterDefs: FilterDefinitions{},
		Persist: PersistOptions{
			Format: PersistInJSONEn,
		},
	}
}
