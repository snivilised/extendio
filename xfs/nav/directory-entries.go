package nav

import (
	"io/fs"

	"github.com/samber/lo"
)

type DirectoryEntryOrderEnum uint

const (
	DirectoryEntryOrderFoldersFirstEn DirectoryEntryOrderEnum = iota
	DirectoryEntryOrderFilesFirstEn
)

type DirectoryEntries struct {
	Options *TraverseOptions
	Order   DirectoryEntryOrderEnum
	Folders []fs.DirEntry
	Files   []fs.DirEntry
}

func (e *DirectoryEntries) arrange(entries *[]fs.DirEntry) {
	grouped := lo.GroupBy(*entries, func(item fs.DirEntry) bool {
		return item.IsDir()
	})

	e.Folders = grouped[true]
	e.Files = grouped[false]
}

func (e *DirectoryEntries) all() *[]fs.DirEntry {

	result := []fs.DirEntry{}
	switch e.Order {
	case DirectoryEntryOrderFoldersFirstEn:
		result = append(e.Folders, e.Files...)
	case DirectoryEntryOrderFilesFirstEn:
		result = append(e.Files, e.Folders...)
	}

	return &result
}

func (e *DirectoryEntries) sort(entries *[]fs.DirEntry) {
	if err := e.Options.Hooks.Sort(*entries); err != nil {
		panic(SORT_L_ERR)
	}
}
