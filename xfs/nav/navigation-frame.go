package nav

type navigationFrame struct {
	Root     string
	Depth    uint
	listener *navigationListener
	client   TraverseCallback
}

func (f *navigationFrame) decorate(label string, decorator TraverseCallback) {
	// this method doesn't do much, but it needs to be made explicit because it
	// is easy to setup the callback decoration chain incorrectly resulting in
	// stack overflow due to infinite recursion. Its easy to search when decoration is
	// occurring in the code base, just search for decorate or go to references.
	//
	// fmt.Printf(">>>> decorate: '%v'\n", label)
	f.client = decorator
}
