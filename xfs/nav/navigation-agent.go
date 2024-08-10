package nav

import (
	"errors"
	"io/fs"
	"path/filepath"

	"github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/internal/lo"
	"github.com/snivilised/extendio/xfs/utils"
)

type newAgentParams struct {
	doInvoke             bool
	o                    *TraverseOptions
	handler              fileSystemErrorHandler
	samplingFilterActive bool
}

func newAgent(params *newAgentParams) *navigationAgent {
	instance := navigationAgent{
		doInvoke:             utils.NewRoProp(params.doInvoke),
		o:                    params.o,
		handler:              params.handler,
		cache:                make(inspectCache),
		samplingFilterActive: params.samplingFilterActive,
	}

	return &instance
}

type navigationAgent struct {
	doInvoke             utils.RoProp[bool]
	o                    *TraverseOptions
	handler              fileSystemErrorHandler
	cache                inspectCache
	samplingFilterActive bool
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
		item := newTraverseItem(
			params.top,
			nil,
			info,
			nil,
			nil,
		)

		_, le = params.impl.traverse(&traverseParams{
			current: item,
			frame:   params.frame,
		})
	}

	result := params.frame.collate()
	result.err = le

	return result, result.err
}

func (a *navigationAgent) read(
	path string,
) (*DirectoryContents, error) {
	// this method was spun out from notify, as there needs to be a separation
	// between these pieces of functionality to support 'extension'; ie we
	// need to read the contents of an items contents to determine the properties
	// created for the extension.
	//
	entries, err := a.o.Hooks.ReadDirectory(path)
	de := newDirectoryContents(&newDirectoryContentsParams{
		o:       a.o,
		entries: entries,
	})

	return de, err
}

type agentNotifyParams struct {
	frame   *navigationFrame
	current *TraverseItem
	entries []fs.DirEntry
	readErr error
}

func (a *navigationAgent) notify(params *agentNotifyParams) (SkipTraversal, error) {
	skip := SkipNoneTraversalEn

	if params.readErr != nil {
		if a.doInvoke.Get() {
			clone := params.current.clone()
			clone.Error = i18n.NewThirdPartyErr(params.readErr)

			// Second call, to report ReadDir error
			//
			if le := params.frame.proxy(clone, nil); le != nil {
				if errors.Is(params.readErr, fs.SkipAll) && (clone.IsDirectory()) {
					params.readErr = nil
				}

				return SkipAllTraversalEn, i18n.NewThirdPartyErr(params.readErr)
			}
		} else {
			return SkipAllTraversalEn, i18n.NewThirdPartyErr(params.readErr)
		}
	}

	return skip, nil
}

type agentTraverseParams struct {
	impl    navigatorImpl
	entries []fs.DirEntry
	parent  *TraverseItem
	frame   *navigationFrame
}

var dontSkipTraverseItem *TraverseItem

func (a *navigationAgent) traverse(params *agentTraverseParams) (*TraverseItem, error) {
	for _, entry := range params.entries {
		path := filepath.Join(params.parent.Path, entry.Name())
		info, e := entry.Info()

		var current *TraverseItem

		if a.samplingFilterActive {
			inspection, found := a.cache[path]
			current = lo.TernaryF(found,
				func() *TraverseItem {
					return inspection.current
				},
				func() *TraverseItem {
					return nil
				},
			)
		}

		if current == nil {
			current = newTraverseItem(
				path,
				entry,
				info,
				params.parent,
				e,
			)
		}

		if skipItem, err := params.impl.traverse(&traverseParams{
			current: current,
			frame:   params.frame,
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

func (a *navigationAgent) keep(stash *inspection) {
	a.cache[stash.current.key()] = stash
	stash.current.filtered()
}
