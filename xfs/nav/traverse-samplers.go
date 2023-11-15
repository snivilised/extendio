package nav

import (
	"io/fs"

	"github.com/samber/lo"
)

type subSetFn func(entries []fs.DirEntry, n int) []fs.DirEntry

func firstSampler(entries []fs.DirEntry, n int) []fs.DirEntry {
	return entries[:len(entries)-(min(n, len(entries)))]
}

func lastSampler(entries []fs.DirEntry, n int) []fs.DirEntry {
	return entries[len(entries)-(min(n, len(entries))):]
}

func getSubSetSampler(noOf *SampleNoOf, fn subSetFn) SampleCallback {
	return func(entries *DirectoryEntries) *DirectoryEntries {
		o := entries.Options
		clone := entries.Sample()

		switch o.Store.Subscription {
		case SubscribeAny:
			clone.Folders = lo.Ternary(noOf.Folders > 0,
				fn(entries.Folders, int(noOf.Folders)),
				entries.Folders,
			)

			clone.Files = lo.Ternary(noOf.Files > 0,
				fn(entries.Files, int(noOf.Files)),
				entries.Files,
			)

		case SubscribeFolders:
		case SubscribeFoldersWithFiles:
			clone.Folders = fn(entries.Folders, int(noOf.Folders))
		case SubscribeFiles:
			clone.Files = fn(entries.Files, int(noOf.Files))

		default:
		}

		return clone
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
