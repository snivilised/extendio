package nav

import (
	"context"
	"io/fs"
	"log/slog"

	"github.com/snivilised/extendio/xfs/utils"
)

// TraverseSubscription type to define traversal subscription (for which file system
// items the client defined callback are invoked for).
type TraverseSubscription uint

const (
	_                         TraverseSubscription = iota
	SubscribeAny                                   // invoke callback for files and folders
	SubscribeFolders                               // invoke callback for folders only
	SubscribeFoldersWithFiles                      // invoke callback for folders only but include files
	SubscribeFiles                                 // invoke callback for files only
)

// TraverseCallback defines traversal callback function signature.
type TraverseCallback func(item *TraverseItem) error

type LabelledTraverseCallback struct {
	Label string
	Fn    TraverseCallback
}

// AscendancyHandler defines the signatures of ascend/descend handlers
type AscendancyHandler func(item *TraverseItem)

// BeginHandler life cycle event handler, invoked before start of traversal
type BeginHandler func(ns *NavigationState)

// EndHandler life cycle event handler, invoked at end of traversal
type EndHandler func(result *TraverseResult)

// TraverseResult the result of the traversal process.
type TraverseResult struct {
	Session Session
	Metrics *NavigationMetrics
	err     error
}

func (r *TraverseResult) merge(other *TraverseResult) (*TraverseResult, error) {
	if !utils.IsNil(other.err) {
		r.err = other.err
	}

	if other.Metrics != nil {
		if r.Metrics == nil {
			r.Metrics = other.Metrics
		} else {
			for k, v := range other.Metrics.collection {
				(r.Metrics.collection)[k].Count += v.Count
			}
		}
	}

	return r, r.err
}

type Prime struct {
	Path            string
	OptionsFn       TraverseOptionFn
	ProvidedOptions *TraverseOptions
}

// Resumption
type Resumption struct {
	RestorePath string
	Restorer    PersistenceRestorer
	Strategy    ResumeStrategyEnum
}

type syncable interface {
	ensync(ctx context.Context, cancel context.CancelFunc, ai *AsyncInfo)
}

// TraverseNavigator interface to the main traverse instance.
type TraverseNavigator interface {
	syncable
	walk(_ string) (*TraverseResult, error)
	save(_ string) error
	finish() error
}

type traverseParams struct {
	current *TraverseItem
	frame   *navigationFrame
	navi    *NavigationInfo
}

type inspector interface {
	inspect(params *traverseParams) *inspection
	keep(stash *inspection)
}

type navigatorImpl interface {
	inspector
	options() *TraverseOptions
	logger() *slog.Logger
	init(ns *NavigationState)
	ensync(ctx context.Context, cancel context.CancelFunc, frame *navigationFrame, ai *AsyncInfo)
	top(frame *navigationFrame, root string) (*TraverseResult, error)
	traverse(params *traverseParams) (*TraverseItem, error)
	finish() error
}

// NavigationInfo
type NavigationInfo struct {
	Options *TraverseOptions
	Item    *TraverseItem
	frame   *navigationFrame
}

// SubPathInfo
type SubPathInfo struct {
	Root      string
	Item      *TraverseItem
	Behaviour *SubPathBehaviour
}

type TriStateBoolEnum uint

const (
	TriStateBoolUnsetEn TriStateBoolEnum = iota
	TriStateBoolTrueEn
	TriStateBoolFalseEn
)

type SkipTraversal uint

const (
	SkipNoneTraversalEn SkipTraversal = iota
	SkipDirTraversalEn
	SkipAllTraversalEn
)

type inspection struct {
	current        *TraverseItem
	contents       *DirectoryContents
	entries        []fs.DirEntry
	readErr        error
	isDir          bool
	compoundCounts *compoundCounters
}

func (i *inspection) clearContents(o *TraverseOptions) {
	if i.contents != nil {
		i.contents.clear()
	} else {
		i.contents = newEmptyDirectoryEntries(o)
	}
}

type itemSubPath = string
type inspectCache map[itemSubPath]*inspection
