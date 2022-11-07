package nav

import (
	"io/fs"
	"path/filepath"

	"github.com/samber/lo"
	. "github.com/snivilised/extendio/translate"
)

type agent struct {
	DO_INVOKE bool // this should be considered const
	o         *TraverseOptions
}

type agentTopParams struct {
	impl  navigatorImpl
	frame *navigationFrame
}

func (a *agent) top(params *agentTopParams) *LocalisableError {
	info, err := a.o.Hooks.QueryStatus(params.frame.Root)
	var le *LocalisableError = nil
	if err != nil {
		item := &TraverseItem{
			Path: params.frame.Root, Info: info, Error: &LocalisableError{Inner: err},
			Children: []fs.DirEntry{},
		}
		le = a.proxy(item, params.frame)
	} else {
		if info.IsDir() {
			item := &TraverseItem{
				Path: params.frame.Root, Info: info,
				Children: []fs.DirEntry{},
			}
			le = params.impl.traverse(item, params.frame)
		} else {

			if a.DO_INVOKE {
				item := &TraverseItem{
					Path: params.frame.Root, Info: info, Error: &NOT_DIRECTORY_L_ERR,
					Children: []fs.DirEntry{},
				}
				params.impl.options().Hooks.Extend(&NavigationInfo{
					Options: params.impl.options(), Item: item, Frame: params.frame,
				}, []fs.DirEntry{})
				le = a.proxy(item, params.frame)
			} else {
				le = &NOT_DIRECTORY_L_ERR
			}
		}
	}
	if (le != nil) && (le.Inner == fs.SkipDir) {
		return nil
	}
	return le
}

func (a *agent) read(item *TraverseItem, order DirectoryEntryOrderEnum) (*DirectoryEntries, error) {
	// this method was spun out from notify, as there needs to be a separation
	// between these pieces of functionality to support 'extension'; ie we
	// need to read the contents of an items contents to determine the properties
	// created for the extension.
	//
	entries, err := a.o.Hooks.ReadDirectory(item.Path)

	de := DirectoryEntries{
		Options: a.o,
		Order:   order,
	}
	de.arrange(&entries)

	return &de, err
}

type agentNotifyParams struct {
	frame   *navigationFrame
	item    *TraverseItem
	entries []fs.DirEntry
	readErr error
}

func (a *agent) notify(params *agentNotifyParams) (bool, *LocalisableError) {

	exit := false
	if params.readErr != nil {

		if a.DO_INVOKE {
			item2 := params.item.Clone()
			item2.Error = &LocalisableError{Inner: params.readErr}

			// Second call, to report ReadDir error
			//
			if le := a.proxy(item2, params.frame); le != nil {
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
	impl     navigatorImpl
	contents *[]fs.DirEntry
	parent   *TraverseItem
	frame    *navigationFrame
}

func (a *agent) traverse(params *agentTraverseParams) *LocalisableError {
	for _, entry := range *params.contents {
		path := filepath.Join(params.parent.Path, entry.Name())
		info, err := entry.Info()
		le := lo.Ternary(err == nil, nil, &LocalisableError{Inner: err})
		child := TraverseItem{
			Path: path, Info: info, Entry: entry, Error: le,
			Children: []fs.DirEntry{},
		}

		if le = params.impl.traverse(&child, params.frame); le != nil {
			if le.Inner == fs.SkipDir {
				break
			}
			return le
		}
	}
	return nil
}

func (a *agent) proxy(currentItem *TraverseItem, frame *navigationFrame) *LocalisableError {
	// proxy is the correct way to invoke the client callback, because it takes into
	// account any active decorations such as listening and filtering. It should be noted
	// that the Callback on the options represents the client defined function which
	// can be decorated. Only the callback on the frame should ever be invoked.
	//
	return frame.client(currentItem)
}
