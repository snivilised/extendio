package nav

import (
	"log/slog"

	"github.com/mohae/deepcopy"
	"github.com/snivilised/extendio/internal/lo"
	"github.com/snivilised/extendio/xfs/utils"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
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
	DirectoryEntryOrder DirectoryContentsOrderEnum
}

type CascadeBehaviour struct {
	// Depth sets a maximum traversal depth
	//
	Depth uint

	// NoRecurse is an alternative to using Depth, but limits the traversal
	// to just the path specified by the user. Since the raison d'etre
	// of the navigator is to recursively process a directory tree, using
	// NoRecurse would appear to be contrary to its natural behaviour. However
	// there are clear usage scenarios where a client needs to process
	// only the files in a specified directory.
	//
	NoRecurse bool
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

	// Cascade controls how deep to navigate
	//
	Cascade CascadeBehaviour
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
	Logger  *slog.Logger
}

// PersistOptions contains options for persisting traverse options
type PersistOptions struct {
	Format PersistenceFormatEnum
}

type LogRotationOptions struct {
	// MaxSizeInMb, max size of a log file, before it is re-cycled
	MaxSizeInMb int

	// MaxNoOfBackups, max number of legacy log files that can exist
	// before being deleted
	MaxNoOfBackups int

	// MaxAgeInDays, max no of days before old log file is deleted
	MaxAgeInDays int
}

type MonitorOptions struct {
	Log *slog.Logger
}

// EntryQuantities contains specification of no of files and folders
// used in various contexts, but primarily sampling.
type EntryQuantities struct {
	Files   uint
	Folders uint
}

// SampleTypeEnum determines the type of sampling to use
type SampleTypeEnum uint

const (
	SampleTypeUnsetEn SampleTypeEnum = iota
	SampleTypeSliceEn
	SampleTypeFilterEn
	SampleTypeCustomEn
)

type SamplingOptions struct {
	// SampleInReverse determines the direction of iteration for the sampling
	// operation
	SampleInReverse bool

	// SampleType the type of sampling to use
	SampleType SampleTypeEnum

	// NoOf specifies number of items required in each sample (only applies
	// when not using Custom iterator options)
	NoOf EntryQuantities
}

type SamplingIteratorOptions struct {
	// Each enables customisation of the sampling functionality, instead of using
	// the defined filter. A directory's contents is sampled according to this
	// function. The function receives the TraverseItem being considered and should
	// return true to include in the sample, false otherwise.
	Each EachDirectoryItemPredicate

	// While enables customisation of the sampling functionality, instead of using
	// the defined filter. The sampling loop will continue to run while this
	// condition is true. The predicate function should return false once condition
	// has been met to complete the sample. Of course, the loop will only run while
	// there are still remaining items to consider (ie there are no more entries
	// to consider for the current traverse item).
	While WhileDirectoryPredicate
}

// SamplerOptions
type SamplerOptions struct {
	// Custom allows the client to customise how a directory's contents are sampled.
	// The default way to sample is either by slicing the directory's contents or
	// by using the filter to select either the first/last n entries (using the
	// SamplingOptions). If the client requires an alternative way of creating a
	// sample, eg to take all files greater than a certain size, then this
	// can be achieved by specifying Each and While inside Custom.
	Custom SamplingIteratorOptions
}

// OptionsStore represents that part of options that is directly
// persist-able.
type OptionsStore struct {
	// Subscription defines which node types are visited
	//
	Subscription TraverseSubscription

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

	// Sampling options
	//
	Sampling SamplingOptions
}

// TraverseOptions customise the way a directory tree is traversed
type TraverseOptions struct {
	Store OptionsStore

	// Callback function to invoke for every item visited in the file system.
	//
	Callback *LabelledTraverseCallback `json:"-"`

	// Notify collection of notification function.
	//
	Notify Notifications `json:"-"`

	// TraverseHooks collection of hook functions, that can be overridden.
	//
	Hooks TraverseHooks `json:"-"`

	// Persist contains options for persisting traverse options
	//
	Persist PersistOptions `json:"-"`

	// Sampler defines options for sampling directory entries. There are
	// multiple ways of performing sampling. The client can either:
	// A) Use one of the four predefined functions see (SamplerOptions.Fn)
	// B) Use a Custom iterator. When setting the Custom iterator properties
	//
	Sampler SamplerOptions `json:"-"`

	// Monitor contains externally provided logger
	//
	Monitor MonitorOptions `json:"-"`
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

func (o *TraverseOptions) isFilteringActive() bool {
	if o.Store.FilterDefs != nil {
		patternDefined := o.Store.FilterDefs.Node.Pattern != ""
		customDefined := o.Store.FilterDefs.Node.Custom != nil
		polyDefined := o.Store.FilterDefs.Node.Poly != nil

		return patternDefined || customDefined || polyDefined
	}

	return false
}

func (o *TraverseOptions) afterUserOptions() {
	if o.Hooks.Sort == nil {
		o.Hooks.Sort = lo.Ternary(o.Store.Behaviours.Sort.IsCaseSensitive,
			CaseSensitiveSortHookFn, CaseInSensitiveSortHookFn,
		)
	}

	if o.Hooks.Extend == nil {
		o.Hooks.Extend = DefaultExtendHookFn
	}

	noEach := o.Sampler.Custom.Each == nil && o.Sampler.Custom.While != nil
	noWhile := o.Sampler.Custom.Each != nil && o.Sampler.Custom.While == nil

	if o.Store.Behaviours.Cascade.NoRecurse {
		o.Store.Behaviours.Cascade.Depth = 1
	}

	if noEach || noWhile {
		panic("invalid SamplingIteratorOptions (set both or neither: Each, While)")
	}
}

func (o *TraverseOptions) Clone() *TraverseOptions {
	clone := deepcopy.Copy(o)
	return clone.(*TraverseOptions)
}

// GetDefaultOptions
func GetDefaultOptions() *TraverseOptions {
	return &TraverseOptions{
		Store: OptionsStore{
			Subscription: SubscribeAny,
			Behaviours: NavigationBehaviours{
				SubPath: SubPathBehaviour{
					KeepTrailingSep: true,
				},
				Sort: SortBehaviour{
					IsCaseSensitive:     false,
					DirectoryEntryOrder: DirectoryContentsOrderFoldersFirstEn,
				},
				Listen: ListenBehaviour{
					InclusiveStart: true,
					InclusiveStop:  false,
				},
			},
		},
		Notify: Notifications{
			OnBegin:   func(_ *NavigationState) {},
			OnEnd:     func(_ *TraverseResult) {},
			OnDescend: func(_ *TraverseItem) {},
			OnAscend:  func(_ *TraverseItem) {},
		},
		Hooks: TraverseHooks{
			QueryStatus:   LstatHookFn,
			ReadDirectory: ReadEntriesHookFn,
			FolderSubPath: RootParentSubPathHookFn,
			FileSubPath:   RootParentSubPathHookFn,
			InitFilters:   InitFiltersHookFn,
		},
		Persist: PersistOptions{
			Format: PersistInJSONEn,
		},
		Monitor: MonitorOptions{
			Log: slog.New(zapslog.NewHandler(
				zapcore.NewNopCore(), nil),
			),
		},
	}
}

func (o *TraverseOptions) useExtendHook() {
	o.Hooks.Extend = DefaultExtendHookFn
}
