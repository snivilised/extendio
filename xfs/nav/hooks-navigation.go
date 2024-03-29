package nav

import (
	"io/fs"
	"os"
	"sort"
	"strings"
)

// Lstat function signature that enables the default t be overridden
type QueryStatusHookFn func(path string) (fs.FileInfo, error)

// ReadDirectoryHookFn hook function to define implementation of how a directory's
// entries are read. A default implementation is preset, so does not have to be set
// by the client.
type ReadDirectoryHookFn func(dirname string) ([]fs.DirEntry, error)

// SortEntriesHookFn hook function to define how directory entries are sorted. Does not
// have to be set explicitly. This will be set according to the IsCaseSensitive on
// the TraverseOptions, but can be overridden if needed.
type SortEntriesHookFn func(entries []fs.DirEntry, custom ...any) error

// FilterInitHookFn
type FilterInitHookFn func(o *TraverseOptions, frame *navigationFrame)

// ExtendHookFn
type ExtendHookFn func(navi *NavigationInfo, entries *DirectoryContents)

// SubPathHookFn
type SubPathHookFn func(info *SubPathInfo) string

// TraverseHooks defines the suite of items that can be customised by the client
type TraverseHooks struct {
	QueryStatus   QueryStatusHookFn
	ReadDirectory ReadDirectoryHookFn
	FolderSubPath SubPathHookFn
	FileSubPath   SubPathHookFn
	InitFilters   FilterInitHookFn
	Sort          SortEntriesHookFn
	Extend        ExtendHookFn
}

// Lstat is the default Query Status hook function
func LstatHookFn(path string) (fs.FileInfo, error) {
	return os.Lstat(path)
}

// CaseSensitiveSortHookFn hook function for case sensitive directory traversal. A
// directory of "a" will be visited after a sibling directory "B".
func CaseSensitiveSortHookFn(entries []fs.DirEntry, _ ...any) error {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	return nil
}

// CaseInSensitiveSortHookFn hook function for case insensitive directory traversal. A
// directory of "a" will be visited before a sibling directory "B".
func CaseInSensitiveSortHookFn(entries []fs.DirEntry, _ ...any) error {
	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})

	return nil
}
