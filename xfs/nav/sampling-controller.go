package nav

import "github.com/snivilised/extendio/internal/lo"

type samplingController struct {
	o        *TraverseOptions
	fn       samplerControllerFunc
	adapters samplingAdaptersCollection
	iterator *directoryEntryWhileIt
}

func (c *samplingController) init(ns *NavigationState) {
	filter := lo.TernaryF(ns.Filters == nil,
		func() TraverseFilter {
			return nil
		},
		func() TraverseFilter {
			return ns.Filters.Node
		},
	)

	forward := !c.o.Store.Sampling.SampleInReverse
	samplingType := c.o.Store.Sampling.SampleType

	if (samplingType == SampleTypeFilterEn) && filter == nil {
		panic("sampling type is filter, but no filter defined")
	}

	if samplingType == SampleTypeCustomEn {
		if c.o.Sampler.Custom.Each == nil {
			panic("sampling type is custom, but no each predicate defined")
		}

		if c.o.Sampler.Custom.While == nil {
			panic("sampling type is custom, but no where predicate defined")
		}
	}

	if (samplingType == SampleTypeFilterEn) || (samplingType == SampleTypeCustomEn) {
		iterator := newDirectoryEntryWhileIt(&newDirectoryEntryWhileItParams{
			o:                c.o,
			forward:          forward,
			useCustomSampler: samplingType == SampleTypeCustomEn,
			filter:           filter,
			adapters:         c.adapters,
			each:             c.o.Sampler.Custom.Each,
			while:            c.o.Sampler.Custom.While,
		})

		c.iterator = iterator
	}
}

func (c *samplingController) initInspector(navigator inspector) {
	if c.iterator != nil {
		c.iterator.initInspector(navigator)
	}
}

func (c *samplingController) sample(contents *DirectoryContents,
	navi *NavigationInfo,
	tp *traverseParams,
) {
	c.fn(&samplerControllerFuncParams{
		contents:     contents,
		subscription: c.o.Store.Subscription,
		noOf:         &c.o.Store.Sampling.NoOf,
		forward:      !c.o.Store.Sampling.SampleInReverse,
		navi:         navi,
		tp:           tp,
		adapters:     c.adapters,
		iterator:     c.iterator,
	})
}
