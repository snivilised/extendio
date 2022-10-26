package xfs

import (
	"errors"
	"io/fs"

	"github.com/samber/lo"
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
		dirs := lo.Filter(entries, func(de fs.DirEntry, i int) bool {
			return de.Type().IsDir()
		})

		var err error
		if err = n.o.Hooks.Sort(dirs); err != nil {
			panic(LocalisableError{
				Inner: errors.New("folder navigator sort function failed"),
			})
		}

		return n.agent.traverse(&agentTraverseParams{
			impl:    n,
			entries: dirs,
			parent:  currentItem,
			frame:   frame,
		})
	}
}
