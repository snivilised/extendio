package xfs

import (
	"io/fs"
	"sort"
	"strings"

	"github.com/samber/lo"
)

// ReadDirectoryHookFn hook function to define implementation of how a directory's
// entries are read. A default implementation is preset, so does not have to be set
// by the client.
type ReadDirectoryHookFn func(dirname string) ([]fs.DirEntry, error)

// SortEntriesHookFn hook function to define how directory entries are sorted. Does not
// have to be set explicitly. This will be set according to the IsCaseSensitive on
// the TraverseOptions, but can be overridden if needed.
type SortEntriesHookFn func(entries []fs.DirEntry, custom ...any) error

// FilterEntriesHookFn hook function.
type FilterEntriesHookFn func(entries []fs.DirEntry, info *FilterInfo, custom ...any) ([]fs.DirEntry, error)

// ExtendHookFn
type ExtendHookFn func(ei *navigationInfo, descendants []fs.DirEntry) error

type traverseHooks struct {
	ReadDirectory ReadDirectoryHookFn
	Sort          SortEntriesHookFn
	Filter        FilterEntriesHookFn
	Extend        ExtendHookFn
}

// FilterHookFn is the default Filter hook function.
func FilterHookFn(entries []fs.DirEntry, info *FilterInfo, custom ...any) ([]fs.DirEntry, error) {

	filtered := lo.Filter(entries, func(entry fs.DirEntry, index int) bool {
		info.Filter.IsMatch(entry.Name(), info.ActualScope)
		return false
	})
	return filtered, nil
}

// CaseSensitiveSortHookFn hook function for case sensitive directory traversal. A
// directory of "a" will be visited after a sibling directory "B".
func CaseSensitiveSortHookFn(entries []fs.DirEntry, custom ...any) error {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	return nil
}

// CaseInSensitiveSortHookFn hook function for case insensitive directory traversal. A
// directory of "a" will be visited before a sibling directory "B".
func CaseInSensitiveSortHookFn(entries []fs.DirEntry, custom ...any) error {
	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})

	return nil
}
