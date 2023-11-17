package nav

type (
	samplerControllerFuncParams struct {
		contents     *DirectoryContents
		subscription TraverseSubscription
		noOf         *EntryQuantities
		forward      bool
		navi         *NavigationInfo
		tp           *traverseParams
		adapters     samplingAdaptersCollection
		iterator     *directoryEntryWhileIt
	}

	// samplerControllerFunc is the root sampling handler and distinguishes between
	// sampling with a slice or sampling by using an iterator (filter/custom)
	samplerControllerFunc func(
		params *samplerControllerFuncParams,
	)

	// sliceSamplerFunc invokes the selected slice function fn, to create
	// sample. The client will specify the first or last slice function
	// depending on the sampling option.
	sliceSamplerFunc func(contents *DirectoryContents, noOf *EntryQuantities, fn sliceEntriesFunc)

	// filterSamplerIteratorFunc
	filterSamplerIteratorFunc func(contents *DirectoryContents, iterator *directoryEntryWhileIt)

	// sampleIsFullFunc
	sampleIsFullFunc func(fi *FilteredInfo, noOf *EntryQuantities) bool

	// samplingAdaptersCollection
	samplingAdaptersCollection map[TraverseSubscription]samplingAdapters

	// samplingAdapters allows sampling behaviour to be customised according to
	// the type of subscription
	samplingAdapters struct {
		slice        sliceSamplerFunc
		filterWithIt filterSamplerIteratorFunc
		isFull       sampleIsFullFunc
	}
)
