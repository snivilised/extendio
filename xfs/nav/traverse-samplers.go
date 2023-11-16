package nav

import (
	"io/fs"

	"github.com/samber/lo"
)

type subSetFn func(entries []fs.DirEntry, n int) []fs.DirEntry

func firstSampler(entries []fs.DirEntry, n int) []fs.DirEntry {
	result := entries[:(min(n, len(entries)))]

	return result
}

func lastSampler(entries []fs.DirEntry, n int) []fs.DirEntry {
	return entries[len(entries)-(min(n, len(entries))):]
}

func getSubSetSampler(noOf *SampleNoOf, fn subSetFn) SampleCallback {
	return func(entries *DirectoryEntries) {
		o := entries.Options

		switch o.Store.Subscription {
		case SubscribeAny:
			entries.Folders = lo.TernaryF(noOf.Folders > 0,
				func() []fs.DirEntry {
					return fn(entries.Folders, int(noOf.Folders))
				},
				func() []fs.DirEntry {
					return entries.Folders
				},
			)

			entries.Files = lo.TernaryF(noOf.Files > 0,
				func() []fs.DirEntry {
					return fn(entries.Files, int(noOf.Files))
				},
				func() []fs.DirEntry {
					return entries.Files
				},
			)

		case SubscribeFolders:
			entries.Folders = fn(entries.Folders, int(noOf.Folders))

		case SubscribeFoldersWithFiles:
			entries.Folders = fn(entries.Folders, int(noOf.Folders))

		case SubscribeFiles:
			entries.Files = fn(entries.Files, int(noOf.Files))

		default:
		}
	}
}

// GetFirstSampler obtains a sampler function which gets the first
// n entries of a directory, where n is either the number of files
// or folders which is determined by the subscription type.
// To use the sampler feature, the client will find they will need to
// use the options push model so that the SampleNoOf instance required
// can be provided. The push model requires the use of ProvidedOptions
// instead of using the pull model via OptionsFn callback.
func GetFirstSampler(noOf *SampleNoOf) SampleCallback {
	return getSubSetSampler(noOf, firstSampler)
}

// GetLastSampler obtains a sampler function which gets the last
// n entries of a directory, where n is either the number of files
// or folders which is determined by the subscription type.
// To use the sampler feature, the client will find they will need to
// use the options push model so that the SampleNoOf instance required
// can be provided. The push model requires the use of ProvidedOptions
// instead of using the pull model via OptionsFn callback.
func GetLastSampler(noOf *SampleNoOf) SampleCallback {
	return getSubSetSampler(noOf, lastSampler)
}
