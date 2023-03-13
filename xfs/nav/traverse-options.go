package nav

import (
	"path/filepath"

	"github.com/mohae/deepcopy"
	"github.com/samber/lo"
	"github.com/snivilised/extendio/internal/log"
	"github.com/snivilised/extendio/xfs/utils"
)

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
	// Node denotes the filter object that represents the current file system item
	// being visited.
	//
	Node FilterDef

	// Children denotes the compound filter that is applied to the direct descendants
	// of the current file system item being visited.
	//
	Children CompoundFilterDef
}

type ListenDefinitions struct {
	// Node denotes the filter object that represents the current file system item
	// being visited.
	//
	StartAt *FilterDef
	StopAt  *FilterDef
}

type NavigationFilters struct {
	// Node denotes the filter object that represents the Node file system item
	// being visited.
	//
	Node TraverseFilter

	// Children denotes the Compound filter that is applied to the direct descendants
	// of the current file system item being visited.
	//
	Children CompoundTraverseFilter
}

// NavigationState carries information about navigation that client may be
// interested in and permitted to access, as opposed to the navigationFrame
// which is meant for internal purposes only.
type NavigationState struct {
	Root    *utils.VarProp[string]
	Filters *NavigationFilters
}

// PersistOptions contains options for persisting traverse options
type PersistOptions struct {
	Format PersistenceFormatEnum
}

type LogRotationOptions struct {
	// Max size of a log file, before it is re-cycled
	MaxSizeInMb int

	// Max number of legacy log files that can exist before being deleted
	MaxNoOfBackups int
	MaxAgeInDays   int
}

type LoggingOptions struct {
	Enabled         bool
	Path            string
	TimeStampFormat string
	Level           log.Level
	Rotation        LogRotationOptions
}

// OptionsStore represents that part of options that is directly
// persist-able.
type OptionsStore struct {
	// Subscription defines which node types are visited
	//
	Subscription TraverseSubscription

	// DoExtend request an extended result.
	//
	DoExtend bool

	// Behaviours collection of behaviours that adjust the way navigation occurs,
	// that can be tweaked by the client.
	//
	Behaviours NavigationBehaviours

	// FilterDefs definitions of filters that restricts for which file system nodes the
	// Callback is invoked for.
	//
	FilterDefs *FilterDefinitions

	// ListenDefs definitions to define when listening starts and stops. Filters
	// can be used to define the Listeners.
	//
	ListenDefs ListenDefinitions

	// Logging options
	//
	Logging LoggingOptions
}

// TraverseOptions customise the way a directory tree is traversed
type TraverseOptions struct {
	Store OptionsStore

	// Callback function to invoke for every item visited in the file system.
	//
	Callback LabelledTraverseCallback `json:"-"`

	// Notify collection of notification function.
	//
	Notify Notifications `json:"-"`

	// TraverseHooks collection of hook functions, that can be overridden.
	//
	Hooks TraverseHooks `json:"-"`

	// Listen options that control when listening state starts and finishes.
	//
	Listen ListenTriggers `json:"-"`

	// Persist contains options for persisting traverse options
	//
	Persist PersistOptions `json:"-"`
}

// TraverseOptionFn functional traverse options
type TraverseOptionFn func(o *TraverseOptions)

func composeTraverseOptions(fn ...TraverseOptionFn) *TraverseOptions {
	o := GetDefaultOptions()

	for _, functionalOption := range fn {
		functionalOption(o)
	}
	o.afterUserOptions()

	return o
}

func (o *TraverseOptions) afterUserOptions() {
	if o.Hooks.Sort == nil {
		o.Hooks.Sort = lo.Ternary(o.Store.Behaviours.Sort.IsCaseSensitive,
			CaseSensitiveSortHookFn, CaseInSensitiveSortHookFn,
		)
	}

	if o.Hooks.Extend == nil {
		o.Hooks.Extend = lo.Ternary(o.Store.DoExtend, DefaultExtendHookFn, nullExtendHookFn)
	}
}

func (o *TraverseOptions) Clone() *TraverseOptions {
	clone := deepcopy.Copy(o)
	return clone.(*TraverseOptions)
}

func GetDefaultOptions() *TraverseOptions {
	return &TraverseOptions{
		Store: OptionsStore{
			Subscription: SubscribeAny,
			DoExtend:     false,
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
			Logging: LoggingOptions{
				Path:            filepath.Join("~", "snivilised.extendio.nav.log"),
				TimeStampFormat: "2006-01-02 15:04:05",
				Level:           log.InfoLevel,
				Rotation: LogRotationOptions{
					MaxSizeInMb:    50,
					MaxNoOfBackups: 3,
					MaxAgeInDays:   28,
				},
			},
		},
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
		Listen: ListenTriggers{
			Start: nil,
			Stop:  nil,
		},
		Persist: PersistOptions{
			Format: PersistInJSONEn,
		},
	}
}

func (o *TraverseOptions) useExtendHook() {
	o.Store.DoExtend = true
	o.Hooks.Extend = lo.Ternary(o.Store.DoExtend, DefaultExtendHookFn, nullExtendHookFn)
}
