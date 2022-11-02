package nav

import (
	"io/fs"

	"github.com/samber/lo"
	. "github.com/snivilised/extendio/translate"
)

type foldersNavigator struct {
	navigator
}

func (n *foldersNavigator) top(frame *navigationFrame) *LocalisableError {

	return n.agent.top(&agentTopParams{
		impl:  n,
		frame: frame,
	})
}

func (n *foldersNavigator) traverse(currentItem *TraverseItem, frame *navigationFrame) *LocalisableError {
	defer func() {
		n.ascend(&NavigationParams{Options: n.o, Item: currentItem, Frame: frame})
	}()
	navi := &NavigationParams{Options: n.o, Item: currentItem, Frame: frame}
	n.descend(navi)
	entries, readErr := n.agent.read(currentItem)
	n.o.Hooks.Extend(navi, entries)

	if le := n.agent.proxy(currentItem, frame); le != nil || (currentItem.Entry != nil && !currentItem.Entry.IsDir()) {
		if le != nil && le.Inner == fs.SkipDir && currentItem.Entry.IsDir() {
			// Successfully skipped directory
			//
			le = nil
		}
		return le
	}

	if exit, err := n.agent.notify(&agentNotifyParams{
		frame: frame, item: currentItem, entries: entries, readErr: readErr,
	}); exit {
		return err
	} else {
		dirs := lo.Filter(entries, func(de fs.DirEntry, i int) bool {
			return de.Type().IsDir()
		})

		var err error
		if err = n.o.Hooks.Sort(dirs); err != nil {
			panic(FOLDERS_NAV_SORT_L_ERR)
		}

		return n.agent.traverse(&agentTraverseParams{
			impl:    n,
			entries: dirs,
			parent:  currentItem,
			frame:   frame,
		})
	}
}