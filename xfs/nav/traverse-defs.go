package nav

import (
	"context"
	"io/fs"

	"github.com/snivilised/extendio/internal/log"
	"github.com/snivilised/extendio/xfs/utils"
)

// ExtendedItem provides extended information if the client requests
// it by setting the DoExtend boolean in the traverse options.
type ExtendedItem struct {
	Depth     int               // traversal depth relative to the root
	IsLeaf    bool              // defines whether this node a leaf node
	Name      string            // derived as the leaf segment from filepath.Split
	Parent    string            // derived as the directory from filepath.Split
	SubPath   string            // represents the path between the root and the current item
	NodeScope FilterScopeBiEnum // type of folder corresponding to the Filter Scope
	Custom    any               // to be set and used by the client
}

// TraverseItem info provided for each file system entity encountered
// during traversal. The root item does not have a DirEntry because it is
// not created as a result of a readDir invoke. Therefore, the client has
// to know that when its function is called back, they will be no DirEntry
// for the root entity.
type TraverseItem struct {
	Path        string
	Entry       fs.DirEntry   // contains a FileInfo via Info() function
	Info        fs.FileInfo   // optional file info instance
	Extension   *ExtendedItem // extended information about the file system node, if requested
	Error       error
	Children    []fs.DirEntry
	filteredOut bool
}

// clone makes shallow copy of TraverseItem (except the error).
func (ti *TraverseItem) clone() *TraverseItem {
	return &TraverseItem{
		Path:      ti.Path,
		Entry:     ti.Entry,
		Info:      ti.Info,
		Extension: ti.Extension,
		Children:  ti.Children,
	}
}

func (ti *TraverseItem) IsDir() bool {
	if !utils.IsNil(ti.Entry) {
		return ti.Entry.IsDir()
	} else if !utils.IsNil(ti.Info) {
		return ti.Info.IsDir()
	}
	// only get here in error scenario, because neither Entry or Info is set
	//
	return false
}

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
type AscendancyHandler func(_ *TraverseItem)

// BeginHandler life cycle event handler, invoked before start of traversal
type BeginHandler func(_ *NavigationState)

// EndHandler life cycle event handler, invoked at end of traversal
type EndHandler func(_ *TraverseResult)

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
	item  *TraverseItem
	frame *navigationFrame
}

type navigatorImpl interface {
	options() *TraverseOptions
	logger() log.Logger
	ensync(_ context.Context, _ context.CancelFunc, _ *navigationFrame, _ *AsyncInfo)
	top(_ *navigationFrame, _ string) (*TraverseResult, error)
	traverse(_ *traverseParams) error
	finish() error
}

// NavigationInfo
type NavigationInfo struct {
	Options *TraverseOptions
	Item    *TraverseItem
	Frame   *navigationFrame
}

// SubPathInfo
type SubPathInfo struct {
	Root      string
	Item      *TraverseItem
	Behaviour *SubPathBehaviour
}

// ClientLogger
type ClientLogger interface {
	Debug(_ string, _ ...log.Field)
	Info(_ string, _ ...log.Field)
	Warn(_ string, _ ...log.Field)
	Error(_ string, _ ...log.Field)
}

type TriStateBoolEnum uint

const (
	TriStateBoolUnsetEn TriStateBoolEnum = iota
	TriStateBoolTrueEn
	TriStateBoolFalseEn
)
