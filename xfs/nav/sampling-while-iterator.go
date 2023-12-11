package nav

import (
	"io/fs"
	"path/filepath"

	"github.com/samber/lo"
	"github.com/snivilised/extendio/collections"
)

// WhileDirectoryPredicate determines when to terminate the loop
type WhileDirectoryPredicate func(fi *FilteredInfo) bool

// EachDirectoryItemPredicate callback to invoke for each child item
type EachDirectoryItemPredicate func(childItem *TraverseItem) bool

// EnoughAlready when using the universal navigator within a sampling operation, set
// these accordingly from inside the while predicate to indicate wether the iteration
// loop should continue to consider more entries to be included in the sample. So
// set Files/Folders flags to true, when enough of those items have been included.
type EnoughAlready struct {
	Files   bool
	Folders bool
}

// FilteredInfo used within the sampling process during a traversal; more specifically,
// they should be set inside the while predicate. Note, the Enough field is only
// appropriate when using the universal navigator.
type FilteredInfo struct {
	Counts EntryQuantities
	Enough EnoughAlready
}

// directoryEntryWhileIt
type directoryEntryWhileIt struct {
	o         *TraverseOptions
	forward   bool
	each      EachDirectoryItemPredicate
	while     WhileDirectoryPredicate
	zero      fs.DirEntry
	adapters  samplingAdaptersCollection
	navigator inspector
	iterator  collections.Iterator[fs.DirEntry]
	tp        *traverseParams
	universal bool
}

type newDirectoryEntryWhileItParams struct {
	o                *TraverseOptions
	useCustomSampler bool
	filter           TraverseFilter
	adapters         samplingAdaptersCollection
	forward          bool
	each             EachDirectoryItemPredicate
	while            WhileDirectoryPredicate
}

// newDirectoryEntryWhileIt represents the predefined while iterator, which
// uses the defined filter to create a sample of the directory entries. A
// client can if they wish define their own custom while iterator that
// does not use a filter. The while iterator that is returned implements
// the looping mechanism used to sample directory entries. It uses a while
// and each predicate to control the iteration process and select items
// into the sample set.
//
// To define the behaviour of the directory while iterator, the client needs to
// provide their own each and while predicates and set them in options at
// Store.Samplers.Custom;
// - each predicate; this is the function that is invoked for each
// child TraverseItem to be considered for the sample.
// - while predicate; a function which controls how long the loop will
// continue to run. When the predicate returns false, then sampling stops
// for the current directory.
func newDirectoryEntryWhileIt(
	params *newDirectoryEntryWhileItParams,
) *directoryEntryWhileIt {
	var (
		subscription = params.o.Store.Subscription
		whit         = &directoryEntryWhileIt{
			o:         params.o,
			forward:   params.forward,
			adapters:  params.adapters,
			universal: subscription == SubscribeAny,
		}
	)

	if params.useCustomSampler {
		whit.each = params.each
		whit.while = params.while
	} else {
		filter := params.filter
		whit.each = func(childItem *TraverseItem) bool {
			return filter.IsMatch(childItem)
		}

		noOf := params.o.Store.Sampling.NoOf

		whit.while = func(fi *FilteredInfo) bool {
			return whit.iterator.Valid() && (!whit.adapters[subscription].isFull(fi, &noOf))
		}
	}

	return whit
}

func (i *directoryEntryWhileIt) isNil() bool {
	// 📚 It is ok to use == operator with nil interface here because
	// the iterator interface has not been set to anything so the value of the
	// interface itself is nil as well as the value it points to.
	//
	return i.iterator == nil
}

func (i *directoryEntryWhileIt) initInspector(navigator inspector) {
	i.navigator = navigator
}

func (i *directoryEntryWhileIt) withParams(tp *traverseParams) {
	i.tp = tp
}

func (i *directoryEntryWhileIt) start(entries []fs.DirEntry) {
	if i.isNil() {
		i.iterator = lo.TernaryF(i.forward,
			func() collections.Iterator[fs.DirEntry] {
				return collections.ForwardIt[fs.DirEntry](entries, i.zero)
			},
			func() collections.Iterator[fs.DirEntry] {
				return collections.ReverseIt[fs.DirEntry](entries, i.zero)
			},
		)
	} else {
		i.iterator.Reset(entries)
	}
}

// sample creates a sample with a new collection of directories entries
func (i *directoryEntryWhileIt) sample(entries []fs.DirEntry, processingFiles bool) []fs.DirEntry {
	i.start(entries)

	result := i.loop()

	return lo.Ternary[[]fs.DirEntry](processingFiles, result.Files, result.Folders)
}

func (i *directoryEntryWhileIt) samples(
	sourceEntries *DirectoryContents,
) (files, folders []fs.DirEntry) {
	i.start(sourceEntries.All())

	result := i.loop()

	files = result.Files
	folders = result.Folders

	return
}

func (i *directoryEntryWhileIt) loop() *DirectoryContents {
	result := newEmptyDirectoryEntries(i.o, &i.o.Store.Sampling.NoOf)
	parent := i.tp.current

	var fi FilteredInfo

	for entry := i.iterator.Start(); i.while(&fi); entry = i.iterator.Next() {
		if entry == nil {
			break
		}

		info, err := entry.Info()

		if i.universal {
			if fi.Enough.Files && !info.IsDir() {
				break
			}

			if fi.Enough.Folders && info.IsDir() {
				break
			}
		}

		path := filepath.Join(parent.Path, entry.Name())
		child := &TraverseItem{
			Path:   path,
			Info:   info,
			Entry:  entry,
			Error:  err,
			Parent: parent,
		}

		stash := i.navigator.inspect(&traverseParams{ // preview
			current: child,
			frame:   i.tp.frame,
			navi: &NavigationInfo{
				Options: i.tp.navi.Options,
				Item:    child,
				frame:   i.tp.frame,
			},
		})

		if i.each(child) {
			i.navigator.keep(stash)

			if entry.IsDir() {
				result.Folders = append(result.Folders, entry)
				fi.Counts.Folders++
			} else {
				result.Files = append(result.Files, entry)
				fi.Counts.Files++
			}
		}
	}

	return result
}
