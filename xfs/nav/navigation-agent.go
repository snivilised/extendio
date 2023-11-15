package nav

import (
	"errors"
	"io/fs"
	"path/filepath"

	xi18n "github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/xfs/utils"
)

type agentFactory struct{}
type agentFactoryParams struct {
	doInvoke  bool
	o         *TraverseOptions
	deFactory directoryEntriesFactory
	handler   fileSystemErrorHandler
}

func (agentFactory) new(params *agentFactoryParams) *navigationAgent {
	instance := navigationAgent{
		doInvoke: utils.NewRoProp(params.doInvoke),
		o:        params.o,
		handler:  params.handler,
	}

	return &instance
}

type navigationAgent struct {
	doInvoke utils.RoProp[bool]
	o        *TraverseOptions
	handler  fileSystemErrorHandler
}

type agentTopParams struct {
	impl  navigatorImpl
	frame *navigationFrame
	top   string
}

func (a *navigationAgent) top(params *agentTopParams) (*TraverseResult, error) {
	params.frame.reset()

	info, err := a.o.Hooks.QueryStatus(params.top)

	var (
		le error
	)

	if err != nil {
		le = a.handler.accept(&fileSystemErrorParams{
			err:   err,
			path:  params.top,
			info:  info,
			agent: a,
			frame: params.frame,
		})
	} else {
		item := &TraverseItem{
			Path: params.top, Info: info,
			Children: []fs.DirEntry{},
		}

		_, le = params.impl.traverse(&traverseParams{
			item:  item,
			frame: params.frame,
		})
	}

	result := params.frame.collate()
	result.err = le

	return result, result.err
}

func (a *navigationAgent) read(
	path string,
	order DirectoryEntryOrderEnum,
) (*DirectoryEntries, error) {
	// this method was spun out from notify, as there needs to be a separation
	// between these pieces of functionality to support 'extension'; ie we
	// need to read the contents of an items contents to determine the properties
	// created for the extension.
	//
	entries, err := a.o.Hooks.ReadDirectory(path)

	deFactory := directoryEntriesFactory{}
	de := deFactory.new(&directoryEntriesFactoryParams{
		o:       a.o,
		order:   order,
		entries: entries,
	})

	return de, err
}

type agentNotifyParams struct {
	frame   *navigationFrame
	item    *TraverseItem
	entries []fs.DirEntry
	readErr error
}

func (a *navigationAgent) notify(params *agentNotifyParams) (SkipTraversal, error) {
	skip := SkipTraversalNoneEn

	if params.readErr != nil {
		if a.doInvoke.Get() {
			item2 := params.item.clone()
			item2.Error = xi18n.NewThirdPartyErr(params.readErr)

			// Second call, to report ReadDir error
			//
			if le := params.frame.proxy(item2, nil); le != nil {
				if errors.Is(params.readErr, fs.SkipAll) && (item2.Entry != nil && item2.Entry.IsDir()) {
					params.readErr = nil
				}

				return SkipTraversalAllEn, xi18n.NewThirdPartyErr(params.readErr)
			}
		} else {
			return SkipTraversalAllEn, xi18n.NewThirdPartyErr(params.readErr)
		}
	}

	return skip, nil
}

type agentTraverseParams struct {
	impl     navigatorImpl
	contents []fs.DirEntry
	parent   *TraverseItem
	frame    *navigationFrame
}

var dontSkipTraverseItem *TraverseItem

func (a *navigationAgent) traverse(params *agentTraverseParams) (*TraverseItem, error) {
	for _, entry := range params.contents {
		path := filepath.Join(params.parent.Path, entry.Name())
		info, e := entry.Info()

		if skipItem, err := params.impl.traverse(&traverseParams{
			item: &TraverseItem{
				Path:     path,
				Info:     info,
				Entry:    entry,
				Error:    e,
				Children: []fs.DirEntry{},
				Parent:   params.parent,
			},
			frame: params.frame,
		}); skipItem == dontSkipTraverseItem {
			if err != nil {
				if errors.Is(err, fs.SkipDir) {
					// The returning of the parent traverse item by the child, denotes
					// a skip; params.parent is the skipItem. So when a child item
					// returns a SkipDir error and return's it parent item, what we're
					// saying is that we want to skip processing all successive siblings
					// but continue traversal. The skipItem indicates we're skipping
					// the remaining processing of all of the parent item's remaining children.
					// (see the ✨ below ...)
					//
					return params.parent, err
				}

				return dontSkipTraverseItem, err
			}
		} else if err != nil {
			// ✨ ... we skip processing all the remaining children for
			// this item, but still continue the overall traversal.
			//
			switch {
			case errors.Is(err, fs.SkipDir):
				continue
			case errors.Is(err, fs.SkipAll):
				break
			default:
				return dontSkipTraverseItem, err
			}
		}
	}

	return dontSkipTraverseItem, nil
}
