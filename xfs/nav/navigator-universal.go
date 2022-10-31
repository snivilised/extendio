package nav

import (
	"io/fs"

	. "github.com/snivilised/extendio/translate"
)

type universalNavigator struct {
	navigator
}

func (n *universalNavigator) top(frame *navigationFrame) *LocalisableError {

	return n.agent.top(&agentTopParams{
		impl:  n,
		frame: frame,
	})
}

func (n *universalNavigator) traverse(currentItem *TraverseItem, frame *navigationFrame) *LocalisableError {
	defer func() {
		n.ascend(&NavigationParams{Options: n.o, Item: currentItem, Frame: frame})
	}()
	navi := &NavigationParams{Options: n.o, Item: currentItem, Frame: frame}
	n.descend(navi)
	entries, readErr := n.agent.read(currentItem)
	n.o.Hooks.Extend(navi, entries)

	if le := n.o.Callback(currentItem); le != nil || (currentItem.Entry != nil && !currentItem.Entry.IsDir()) {
		if le != nil && le.Inner == fs.SkipDir && currentItem.Entry.IsDir() {
			// Successfully skipped directory
			//
			le = nil
		}
		return le
	}

	if exit, err := n.agent.notify(&agentNotifyParams{
		item: currentItem, entries: entries, readErr: readErr,
	}); exit {
		return err
	} else {
		var err error
		if err = n.o.Hooks.Sort(entries); err != nil {
			panic(UNIVERSAL_NAV_SORT_L_ERR)
		}

		return n.agent.traverse(&agentTraverseParams{
			impl:    n,
			entries: entries,
			parent:  currentItem,
			frame:   frame,
		})
	}
}
