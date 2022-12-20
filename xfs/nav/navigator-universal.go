package nav

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/samber/lo"
	. "github.com/snivilised/extendio/translate"
	"github.com/snivilised/extendio/xfs/utils"
)

type universalNavigator struct {
	navigator
}

func (n *universalNavigator) top(frame *navigationFrame) *LocalisableError {

	return n.agent.top(&agentTopParams{
		impl:  n,
		frame: frame,
		top:   frame.root,
	})
}

func (n *universalNavigator) traverse(params *traverseParams) *LocalisableError {
	defer func() {
		n.ascend(&NavigationInfo{Options: n.o, Item: params.currentItem, Frame: params.frame})
	}()
	navi := &NavigationInfo{Options: n.o, Item: params.currentItem, Frame: params.frame}
	n.descend(navi)
	entries, readErr := n.agent.read(params.currentItem.Path, n.o.Store.Behaviours.Sort.DirectoryEntryOrder)
	// Files and Folders need to be sorted independently to preserve the navigation order
	// stipulated by .Behaviours.Sort.DirectoryEntryOrder
	//
	entries.sort(&entries.Files)
	entries.sort(&entries.Folders)
	sorted := entries.all()
	n.o.Hooks.Extend(navi, *sorted)

	if le := n.agent.proxy(params.currentItem, params.frame); le != nil ||
		(params.currentItem.Entry != nil && !params.currentItem.Entry.IsDir()) {
		if le != nil && le.Inner == fs.SkipDir && params.currentItem.Entry.IsDir() {
			// Successfully skipped directory
			//
			le = nil
		}
		return le
	}

	if exit, err := n.agent.notify(&agentNotifyParams{
		frame: params.frame, item: params.currentItem, entries: *sorted, readErr: readErr,
	}); exit {
		return err
	} else {

		return n.agent.traverse(&agentTraverseParams{
			impl:     n,
			contents: sorted,
			parent:   params.currentItem,
			frame:    params.frame,
		})
	}
}

func (n *universalNavigator) spawn(params *spawnParams) *LocalisableError {
	fmt.Printf(">>> universalNavigator.spawn(top) 🎪 \n")

	info, err := n.o.Hooks.QueryStatus(params.active.NodePath)

	if err != nil {
		return &LocalisableError{
			Inner: err,
		}
	}

	indicator := lo.Ternary(info.IsDir(), "📁 DIRECTORY", "📜 FILE")
	fmt.Printf("   🧿 root: '%v' \n", params.frame.root)
	fmt.Printf("   🧿 resume-at(%v): '%v' \n", indicator, params.anchor)

	parent, child := utils.SplitParent(params.anchor)
	fmt.Printf("   - parent: 🧿'%v', child: 🧿'%v' \n", parent, child)

	fracture, _ := n.agent.siblingsFollowing(&followingInfo{
		parent: parent,
		order:  n.o.Store.Behaviours.Sort.DirectoryEntryOrder,
		anchor: child,
	})
	fracture.siblings.sort(&fracture.siblings.Files)
	fracture.siblings.sort(&fracture.siblings.Folders)
	sorted := fracture.siblings.all()

	sequence := executionSequence{
		func() *LocalisableError {
			return n.files(&topSpawnParams{
				frame:   params.frame,
				parent:  parent,
				entries: sorted,
			})
		},
	}

	var le *LocalisableError
	for _, fn := range sequence {
		if le = fn(); le != nil {
			break
		}
	}

	return nil
}

func (n *universalNavigator) files(params *topSpawnParams) *LocalisableError {
	for _, entry := range *params.entries {
		if !entry.IsDir() {
			topPath := filepath.Join(params.parent, entry.Name())
			le := n.agent.top(&agentTopParams{
				impl:  n,
				frame: params.frame,
				top:   topPath,
			})

			if le != nil {
				return le
			}
		}
	}
	return nil
}

func (n *universalNavigator) seed(params *seedParams) *LocalisableError {
	return nil
}
