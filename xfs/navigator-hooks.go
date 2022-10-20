package xfs

import (
	"io/fs"
	"sort"
	"strings"

	"github.com/samber/lo"
)

// TODO: change error to LocalisableError
type ReadDirectoryHookFn func(dirname string) ([]fs.DirEntry, error)
type SortEntriesHookFn func(entries []fs.DirEntry, custom ...any) ([]fs.DirEntry, error)
type FilterEntriesHookFn func(entries []fs.DirEntry, info *FilterInfo, custom ...any) ([]fs.DirEntry, error)

type TraverseHooks struct {
	ReadDirectory ReadDirectoryHookFn
	Sort          SortEntriesHookFn
	Filter        FilterEntriesHookFn
}

func FilterHookFn(entries []fs.DirEntry, info *FilterInfo, custom ...any) ([]fs.DirEntry, error) {

	filtered := lo.Filter(entries, func(entry fs.DirEntry, index int) bool {
		info.Filter.IsMatch(entry.Name(), info.ActualScope)
		return false
	})
	return filtered, nil
}

func CaseSensitiveSortHookFn(entries []fs.DirEntry, custom ...any) ([]fs.DirEntry, error) {
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	return entries, nil
}

func CaseInSensitiveSortHookFn(entries []fs.DirEntry, custom ...any) ([]fs.DirEntry, error) {
	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})

	return entries, nil
}
