package nav

type navigationFrame struct {
	Root     string
	NodePath string
	Depth    uint
	listener *navigationListener
	raw      LabelledTraverseCallback // un-decorated (except for filter) client callback
	client   LabelledTraverseCallback // decorate-able client callback
	filters  *NavigationFilters
}

func (f *navigationFrame) decorate(label string, decorator LabelledTraverseCallback) *LabelledTraverseCallback {
	// this method doesn't do much, but it needs to be made explicit because it
	// is easy to setup the callback decoration chain incorrectly resulting in
	// stack overflow due to infinite recursion. Its easy to search when decoration is
	// occurring in the code base, just search for decorate or go to references.
	//
	previous := f.client
	f.client = decorator

	return &previous
}
