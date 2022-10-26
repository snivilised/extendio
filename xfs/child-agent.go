package xfs

import (
	"errors"
	"io/fs"
	"path/filepath"

	"github.com/samber/lo"
)

type childAgent struct {
	DO_INVOKE bool // this should be considered const
	o         *TraverseOptions
}

type agentTopParams struct {
	impl  navigatorImpl
	frame *navigationFrame
}

func (a *childAgent) top(params *agentTopParams) *LocalisableError {
	info, err := a.o.Hooks.QueryStatus(params.frame.Root)
	var le *LocalisableError = nil
	if err != nil {
		// top level stat error
		//
		item := &TraverseItem{Path: params.frame.Root, Info: info, Error: &LocalisableError{Inner: err}}
		le = params.impl.options().Callback(item)
	} else {
		// traverse from top
		//
		if info.IsDir() {
			item := &TraverseItem{Path: params.frame.Root, Info: info}
			le = params.impl.traverse(item, params.frame)
		} else {
			NOT_A_DIRECTORY_L_ERROR := &LocalisableError{Inner: errors.New("Not a directory")}

			if a.DO_INVOKE {
				item := &TraverseItem{
					Path: params.frame.Root, Info: info, Error: NOT_A_DIRECTORY_L_ERROR,
				}
				params.impl.options().Hooks.Extend(&NavigationParams{
					Options: params.impl.options(), Item: item, Frame: params.frame,
				}, []fs.DirEntry{})
				le = params.impl.options().Callback(item)
			} else {
				le = NOT_A_DIRECTORY_L_ERROR
			}
		}
	}
	if (le != nil) && (le.Inner == fs.SkipDir) {
		return nil
	}
	return le
}

func (a *childAgent) read(item *TraverseItem) ([]fs.DirEntry, error) {
	// this method was spun out from notify, as there needs to be a separation
	// between these pieces of functionality to support 'extension'; ie we
	// need to read the contents of an items contents to determine the properties
	// created for the extension.
	//
	return a.o.Hooks.ReadDirectory(item.Path)
}

type agentNotifyParams struct {
	item    *TraverseItem
	entries []fs.DirEntry
	readErr error
}

func (a *childAgent) notify(params *agentNotifyParams) (bool, *LocalisableError) {

	exit := false
	if params.readErr != nil {

		if a.DO_INVOKE {
			item2 := params.item.Clone()
			item2.Error = &LocalisableError{Inner: params.readErr}

			// Second call, to report ReadDir error
			//
			if le := a.o.Callback(item2); le != nil {
				if params.readErr == fs.SkipDir && (item2.Entry != nil && item2.Entry.IsDir()) {
					params.readErr = nil
				}
				return true, &LocalisableError{Inner: params.readErr}
			}
		} else {
			return true, &LocalisableError{Inner: params.readErr}
		}
	}

	return exit, nil
}

type agentTraverseParams struct {
	impl    navigatorImpl
	entries []fs.DirEntry
	parent  *TraverseItem
	frame   *navigationFrame
}

func (a *childAgent) traverse(params *agentTraverseParams) *LocalisableError {
	for _, entry := range params.entries {
		path := filepath.Join(params.parent.Path, entry.Name())
		info, err := entry.Info()
		le := lo.Ternary(err == nil, nil, &LocalisableError{Inner: err})
		child := TraverseItem{Path: path, Info: info, Entry: entry, Error: le}

		if le = params.impl.traverse(&child, params.frame); le != nil {
			if le.Inner == fs.SkipDir {
				break
			}
			return le
		}
	}
	return nil
}
