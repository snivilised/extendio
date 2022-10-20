package xfs

import (
	"io/fs"
	"path/filepath"

	"github.com/samber/lo"
)

type childAgent struct {
	DO_INVOKE bool // this should be considered const
	options   *TraverseOptions
}

type agentTraverseInfo struct {
	core    navigatorCore
	entries []fs.DirEntry
	parent  *TraverseItem
	frame   *navigationFrame
}

func (a *childAgent) read(item *TraverseItem) ([]fs.DirEntry, error) {
	// this method was spun out from notify, as there needs to be a separation
	// between these pieces of functionality to support 'extension'; ie we
	// need to read the contents of an items contents to determine the properties
	// created for the extension.
	//
	return a.options.Hooks.ReadDirectory(item.Path)
}

type notifyInfo struct {
	item    *TraverseItem
	entries []fs.DirEntry
	readErr error
}

func (a *childAgent) notify(ni *notifyInfo) (bool, *LocalisableError) {

	exit := false
	if ni.readErr != nil {

		if a.DO_INVOKE {
			item2 := ni.item.Clone()
			item2.Error = &LocalisableError{Inner: ni.readErr}

			// Second call, to report ReadDir error
			//
			if le := a.options.Callback(item2); le != nil {
				if ni.readErr == fs.SkipDir && (item2.Entry != nil && item2.Entry.IsDir()) {
					ni.readErr = nil
				}
				return true, &LocalisableError{Inner: ni.readErr}
			}
		} else {
			return true, &LocalisableError{Inner: ni.readErr}
		}
	}

	return exit, nil
}

func (a *childAgent) traverse(ti *agentTraverseInfo) *LocalisableError {
	for _, entry := range ti.entries {
		path := filepath.Join(ti.parent.Path, entry.Name())
		info, err := entry.Info()
		le := lo.Ternary(err == nil, nil, &LocalisableError{Inner: err})
		child := TraverseItem{Path: path, Info: info, Entry: entry, Error: le}

		if le = ti.core.traverse(&child, ti.frame); le != nil {
			if le.Inner == fs.SkipDir {
				break
			}
			return le
		}
	}
	return nil
}
