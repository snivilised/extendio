package nav

import (
	"io/fs"

	"github.com/samber/lo"
)

func sampleWithSliceController(params *samplerControllerFuncParams,
) {
	params.adapters[params.subscription].slice(
		params.contents,
		params.noOf,
		lo.Ternary(params.forward, firstBySliceSampler, lastBySliceSampler),
	)
}

func sampleWithIteratorController(params *samplerControllerFuncParams) {
	params.iterator.withParams(params.tp)

	params.adapters[params.subscription].filterWithIt(
		params.contents,
		params.iterator,
	)
}

func getSamplerControllerFunc(o *TraverseOptions) samplerControllerFunc {
	switch o.Store.Sampling.SampleType {
	case SampleTypeSliceEn:
		return sampleWithSliceController

	case SampleTypeFilterEn, SampleTypeCustomEn:
		return sampleWithIteratorController

	case SampleTypeUnsetEn:
	}

	panic("sampling type not set")
}

type sliceEntriesFunc func(entries []fs.DirEntry, n int) []fs.DirEntry

// firstBySliceSampler sampler function that creates a subset of
// a directory's entries using a slice expression. The subset extracted
// is the first n items of the slice.
func firstBySliceSampler(entries []fs.DirEntry, n int) []fs.DirEntry {
	return entries[:(min(n, len(entries)))]
}

// lastBySliceSampler sampler function that creates a subset of
// a directory's entries using a slice expression. The subset extracted
// is the last n items of the slice.
func lastBySliceSampler(entries []fs.DirEntry, n int) []fs.DirEntry {
	return entries[len(entries)-(min(n, len(entries))):]
}
