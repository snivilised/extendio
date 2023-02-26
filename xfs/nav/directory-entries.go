package nav

import (
	"io/fs"

	"github.com/samber/lo"
	. "github.com/snivilised/extendio/i18n"
)

type DirectoryEntryOrderEnum uint

const (
	DirectoryEntryOrderFoldersFirstEn DirectoryEntryOrderEnum = iota
	DirectoryEntryOrderFilesFirstEn
)

type directoryEntriesFactory struct{}

type directoryEntriesFactoryParams struct {
	o       *TraverseOptions
	order   DirectoryEntryOrderEnum
	entries *[]fs.DirEntry
}

func (directoryEntriesFactory) new(params *directoryEntriesFactoryParams) *DirectoryEntries {

	instance := DirectoryEntries{
		Options: params.o,
		Order:   params.order,
	}
	instance.arrange(params.entries)
	return &instance
}

// DirectoryEntries
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

	if e.Folders == nil {
		e.Folders = []fs.DirEntry{}
	}
	if e.Files == nil {
		e.Files = []fs.DirEntry{}
	}
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
		panic(NewSortFnFailedError())
	}
}
