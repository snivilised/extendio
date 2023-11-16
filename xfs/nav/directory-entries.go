package nav

import (
	"io/fs"

	"github.com/samber/lo"
	"github.com/snivilised/extendio/i18n"
)

// DirectoryEntryOrderEnum determines what order a directories
// entries are invoked for.
type DirectoryEntryOrderEnum uint

const (
	// DirectoryEntryOrderFoldersFirstEn invoke folders first
	DirectoryEntryOrderFoldersFirstEn DirectoryEntryOrderEnum = iota

	// DirectoryEntryOrderFilesFirstEn invoke files first
	DirectoryEntryOrderFilesFirstEn
)

type newDirectoryEntriesParams struct {
	o       *TraverseOptions
	entries []fs.DirEntry
}

func newDirectoryEntries(params *newDirectoryEntriesParams) *DirectoryEntries {
	instance := DirectoryEntries{
		Options: params.o,
	}

	instance.arrange(params.entries)

	return &instance
}

// DirectoryEntries represents the contents of a directory's contents and
// handles sorting order which by default is different between various
// operating systems. This abstraction removes the differences in sorting
// behaviour on different platforms.
type DirectoryEntries struct {
	Options *TraverseOptions
	Folders []fs.DirEntry
	Files   []fs.DirEntry
}

// All returns the contents of a directory respecting the directory sorting
// order defined in the traversal options.
func (e *DirectoryEntries) All() []fs.DirEntry {
	result := []fs.DirEntry{}

	switch e.Options.Store.Behaviours.Sort.DirectoryEntryOrder {
	case DirectoryEntryOrderFoldersFirstEn:
		result = append(e.Folders, e.Files...) //nolint:gocritic // no alternative known
	case DirectoryEntryOrderFilesFirstEn:
		result = append(e.Files, e.Folders...) //nolint:gocritic // no alternative known
	}

	return result
}

func (e *DirectoryEntries) arrange(entries []fs.DirEntry) {
	grouped := lo.GroupBy(entries, func(item fs.DirEntry) bool {
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

func (e *DirectoryEntries) sort(entries []fs.DirEntry) {
	if err := e.Options.Hooks.Sort(entries); err != nil {
		panic(i18n.NewSortFnFailedError())
	}
}

func newEmptyDirectoryEntries(o *TraverseOptions) *DirectoryEntries {
	return &DirectoryEntries{
		Options: o,
		Files:   []fs.DirEntry{},
		Folders: []fs.DirEntry{},
	}
}
