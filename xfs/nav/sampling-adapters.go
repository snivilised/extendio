package nav

import (
	"io/fs"

	"github.com/snivilised/extendio/internal/lo"
)

const samplingFiles = true

func createSamplingAdapters() samplingAdaptersCollection {
	universal := samplingAdapters{
		slice: func(contents *DirectoryContents, noOf *EntryQuantities, fn sliceEntriesFunc) {
			contents.Folders = lo.TernaryF(noOf.Folders > 0,
				func() []fs.DirEntry {
					return fn(contents.Folders, int(noOf.Folders))
				},
				func() []fs.DirEntry {
					return contents.Folders
				},
			)

			contents.Files = lo.TernaryF(noOf.Files > 0,
				func() []fs.DirEntry {
					return fn(contents.Files, int(noOf.Files))
				},
				func() []fs.DirEntry {
					return contents.Files
				},
			)
		},
		filterWithIt: func(contents *DirectoryContents, iterator *directoryEntryWhileIt) {
			contents.Files, contents.Folders = iterator.samples(contents)
		},
		isFull: func(fi *FilteredInfo, noOf *EntryQuantities) bool {
			if fi.Counts.Files == noOf.Files {
				fi.Enough.Files = true
			}
			if fi.Counts.Folders == noOf.Folders {
				fi.Enough.Folders = true
			}

			return fi.Enough.Files && fi.Enough.Folders
		},
	}

	folders := samplingAdapters{
		slice: func(contents *DirectoryContents, noOf *EntryQuantities, fn sliceEntriesFunc) {
			contents.Folders = fn(contents.Folders, int(noOf.Folders))
		},
		filterWithIt: func(contents *DirectoryContents, iterator *directoryEntryWhileIt) {
			contents.Folders = iterator.sample(contents.Folders, !samplingFiles)
		},
		isFull: func(fi *FilteredInfo, noOf *EntryQuantities) bool {
			return fi.Counts.Folders == noOf.Folders
		},
	}

	files := samplingAdapters{
		slice: func(contents *DirectoryContents, noOf *EntryQuantities, fn sliceEntriesFunc) {
			contents.Files = fn(contents.Files, int(noOf.Files))
		},
		filterWithIt: func(contents *DirectoryContents, iterator *directoryEntryWhileIt) {
			contents.Files = iterator.sample(contents.Files, samplingFiles)
		},
		isFull: func(fi *FilteredInfo, noOf *EntryQuantities) bool {
			return fi.Counts.Files == noOf.Files
		},
	}

	return samplingAdaptersCollection{
		SubscribeAny:              universal,
		SubscribeFolders:          folders,
		SubscribeFoldersWithFiles: folders,
		SubscribeFiles:            files,
	}
}
