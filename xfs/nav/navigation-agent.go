package nav

import (
	"io/fs"
	"path/filepath"

	"github.com/samber/lo"
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
	doInvoke  utils.RoProp[bool]
	o         *TraverseOptions
	deFactory directoryEntriesFactory
	handler   fileSystemErrorHandler
}

type agentTopParams struct {
	impl  navigatorImpl
	frame *navigationFrame
	top   string
}

func (a *navigationAgent) top(params *agentTopParams) (*TraverseResult, error) {
	params.frame.reset()

	info, err := a.o.Hooks.QueryStatus(params.top)

	var le error

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

		le = params.impl.traverse(&traverseParams{
			item:  item,
			frame: params.frame,
		})
	}

	result := params.frame.collate()
	result.err = le

	return result, result.err
}

func (a *navigationAgent) read(path string, order DirectoryEntryOrderEnum) (*DirectoryEntries, error) {
	// this method was spun out from notify, as there needs to be a separation
	// between these pieces of functionality to support 'extension'; ie we
	// need to read the contents of an items contents to determine the properties
	// created for the extension.
	//
	entries, err := a.o.Hooks.ReadDirectory(path)

	de := a.deFactory.new(&directoryEntriesFactoryParams{
		o:       a.o,
		order:   order,
		entries: &entries,
	})

	return de, err
}

type agentNotifyParams struct {
	frame   *navigationFrame
	item    *TraverseItem
	entries []fs.DirEntry
	readErr error
}

func (a *navigationAgent) notify(params *agentNotifyParams) (bool, error) {
	exit := false

	if params.readErr != nil {
		if a.doInvoke.Get() {
			item2 := params.item.clone()
			item2.Error = xi18n.NewThirdPartyErr(params.readErr)

			// Second call, to report ReadDir error
			//
			if le := a.proxy(item2, params.frame); le != nil {
				if QuerySkipDirError(params.readErr) && (item2.Entry != nil && item2.Entry.IsDir()) {
					params.readErr = nil
				}

				return true, xi18n.NewThirdPartyErr(params.readErr)
			}
		} else {
			return true, xi18n.NewThirdPartyErr(params.readErr)
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

func (a *navigationAgent) traverse(params *agentTraverseParams) error {
	for _, entry := range *params.contents {
		path := filepath.Join(params.parent.Path, entry.Name())
		info, err := entry.Info()

		var le error

		if le != nil {
			le = xi18n.NewThirdPartyErr(err)
		}

		child := TraverseItem{
			Path: path, Info: info, Entry: entry, Error: le,
			Children: []fs.DirEntry{},
		}

		if le = params.impl.traverse(&traverseParams{
			item:  &child,
			frame: params.frame,
		}); le != nil {
			if QuerySkipDirError(le) {
				break
			}

			return le
		}
	}

	return nil
}

func (a *navigationAgent) proxy(item *TraverseItem, frame *navigationFrame) error {
	// proxy is the correct way to invoke the client callback, because it takes into
	// account any active decorations such as listening and filtering. It should be noted
	// that the Callback on the options represents the client defined function which
	// can be decorated. Only the callback on the frame should ever be invoked.
	//
	frame.currentPath.Set(item.Path)
	clientErr := frame.client.Fn(item)

	if !item.skip && item.Error == nil {
		metricEn := lo.Ternary(item.IsDir(), MetricNoFoldersEn, MetricNoFilesEn)
		frame.metrics.tick(metricEn)
	}

	if clientErr == nil && item.Error == nil {
		return nil
	}

	var resultErr error

	switch {
	case item.Error != nil:
		resultErr = item.Error
	default:
		resultErr = clientErr
	}

	return resultErr
}
