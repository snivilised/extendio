package nav

import (
	"io/fs"
	"path/filepath"

	"github.com/samber/lo"
	. "github.com/snivilised/extendio/translate"
)

type agentFactory struct{}
type agentFactoryParams struct {
	doInvoke  bool
	o         *TraverseOptions
	deFactory *directoryEntriesFactory
}

func (*agentFactory) construct(params *agentFactoryParams) *navigationAgent {
	instance := navigationAgent{
		DO_INVOKE: params.doInvoke,
		o:         params.o,
	}
	instance.deFactory = &directoryEntriesFactory{}

	return &instance
}

type navigationAgent struct {
	DO_INVOKE bool // this should be considered const
	o         *TraverseOptions
	deFactory *directoryEntriesFactory
}

type agentTopParams struct {
	impl  navigatorImpl
	frame *navigationFrame
	top   string
}

func (a *navigationAgent) top(params *agentTopParams) *LocalisableError {
	info, err := a.o.Hooks.QueryStatus(params.top)
	var le *LocalisableError = nil
	if err != nil {
		item := &TraverseItem{
			Path: params.top, Info: info, Error: &LocalisableError{Inner: err},
			Children: []fs.DirEntry{},
		}
		le = a.proxy(item, params.frame)
	} else {
		if info.IsDir() {
			item := &TraverseItem{
				Path: params.top, Info: info,
				Children: []fs.DirEntry{},
			}

			le = params.impl.traverse(&traverseParams{
				currentItem: item,
				frame:       params.frame,
			})
		} else {

			if a.DO_INVOKE {
				// TODO: this might be a problem. We must not treat the top entity being a file
				// as an error for spawn-resume scenarios
				//
				item := &TraverseItem{
					Path: params.top, Info: info, Error: &NOT_DIRECTORY_L_ERR,
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

func (a *navigationAgent) read(path string, order DirectoryEntryOrderEnum) (*directoryEntries, error) {
	// this method was spun out from notify, as there needs to be a separation
	// between these pieces of functionality to support 'extension'; ie we
	// need to read the contents of an items contents to determine the properties
	// created for the extension.
	//
	entries, err := a.o.Hooks.ReadDirectory(path)

	de := directoryEntries{
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

func (a *navigationAgent) notify(params *agentNotifyParams) (bool, *LocalisableError) {

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

func (a *navigationAgent) traverse(params *agentTraverseParams) *LocalisableError {
	for _, entry := range *params.contents {
		path := filepath.Join(params.parent.Path, entry.Name())
		info, err := entry.Info()
		le := lo.Ternary(err == nil, nil, &LocalisableError{Inner: err})
		child := TraverseItem{
			Path: path, Info: info, Entry: entry, Error: le,
			Children: []fs.DirEntry{},
		}

		if le = params.impl.traverse(&traverseParams{
			currentItem: &child,
			frame:       params.frame,
		}); le != nil {
			if le.Inner == fs.SkipDir {
				break
			}
			return le
		}
	}
	return nil
}

func (a *navigationAgent) proxy(currentItem *TraverseItem, frame *navigationFrame) *LocalisableError {
	// proxy is the correct way to invoke the client callback, because it takes into
	// account any active decorations such as listening and filtering. It should be noted
	// that the Callback on the options represents the client defined function which
	// can be decorated. Only the callback on the frame should ever be invoked.
	//
	frame.nodePath = currentItem.Path
	return frame.client.Fn(currentItem)
}

func (a *navigationAgent) siblingsFollowing(info *followingInfo) (*fractureInfo, error) {

	de, le := a.read(info.parent, info.order)

	if le != nil {
		return nil, le
	}
	// is this subscription dependent? if so then, this functionality needs to be moved
	// to the navigator-[subscription]
	//
	entries := de.all()

	groups := lo.GroupBy(*entries, func(item fs.DirEntry) bool {
		return item.Name() >= info.anchor
	})
	following := groups[true]

	siblingsDe := a.deFactory.construct(&directoryEntriesFactoryParams{
		o:       a.o,
		order:   info.order,
		entries: &following,
	})

	return &fractureInfo{siblings: siblingsDe}, nil
}
