package nav

import (
	"io/fs"

	"github.com/samber/lo"
	"github.com/snivilised/extendio/i18n"
)

// DirectoryContentsOrderEnum determines what order a directories
// entries are invoked for.
type DirectoryContentsOrderEnum uint

const (
	// DirectoryContentsOrderFoldersFirstEn invoke folders first
	DirectoryContentsOrderFoldersFirstEn DirectoryContentsOrderEnum = iota

	// DirectoryContentsOrderFilesFirstEn invoke files first
	DirectoryContentsOrderFilesFirstEn
)

type newDirectoryContentsParams struct {
	o       *TraverseOptions
	entries []fs.DirEntry
}

func newDirectoryContents(params *newDirectoryContentsParams) *DirectoryContents {
	instance := DirectoryContents{
		Options: params.o,
	}

	instance.arrange(params.entries)

	return &instance
}

// DirectoryContents represents the contents of a directory's contents and
// handles sorting order which by default is different between various
// operating systems. This abstraction removes the differences in sorting
// behaviour on different platforms.
type DirectoryContents struct {
	Options *TraverseOptions
	Folders []fs.DirEntry
	Files   []fs.DirEntry
}

// All returns the contents of a directory respecting the directory sorting
// order defined in the traversal options.
func (e *DirectoryContents) All() []fs.DirEntry {
	result := []fs.DirEntry{}

	switch e.Options.Store.Behaviours.Sort.DirectoryEntryOrder {
	case DirectoryContentsOrderFoldersFirstEn:
		result = append(e.Folders, e.Files...) //nolint:gocritic // no alternative known
	case DirectoryContentsOrderFilesFirstEn:
		result = append(e.Files, e.Folders...) //nolint:gocritic // no alternative known
	}

	return result
}

func (e *DirectoryContents) arrange(entries []fs.DirEntry) {
	grouped := lo.GroupBy(entries, func(entry fs.DirEntry) bool {
		return entry.IsDir()
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

func (e *DirectoryContents) sort(entries []fs.DirEntry) {
	if err := e.Options.Hooks.Sort(entries); err != nil {
		panic(i18n.NewSortFnFailedError())
	}
}

func (e *DirectoryContents) clear() {
	e.Files = []fs.DirEntry{}
	e.Folders = []fs.DirEntry{}
}

func newEmptyDirectoryEntries(o *TraverseOptions, prealloc ...*EntryQuantities) *DirectoryContents {
	return lo.TernaryF(len(prealloc) == 0,
		func() *DirectoryContents {
			return &DirectoryContents{
				Options: o,
				Files:   []fs.DirEntry{},
				Folders: []fs.DirEntry{},
			}
		},
		func() *DirectoryContents {
			return &DirectoryContents{
				Options: o,
				Files:   make([]fs.DirEntry, 0, prealloc[0].Files),
				Folders: make([]fs.DirEntry, 0, prealloc[0].Folders),
			}
		},
	)
}
