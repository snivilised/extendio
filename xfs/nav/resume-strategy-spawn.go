package nav

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/samber/lo"
	"github.com/snivilised/extendio/xfs/utils"
	"go.uber.org/zap"
)

const (
	_FOLLOWING_SIBLINGS = true
)

type spawnStrategy struct {
	baseStrategy
}

func (s *spawnStrategy) init(params *strategyInitParams) {

}

func (s *spawnStrategy) resume(info *strategyResumeInfo) *TraverseResult {
	s.nc.frame.root.Set(info.ps.Active.Root)
	resumeAt := s.ps.Active.NodePath

	s.nc.logger().Info("spawn resume",
		zap.String("root-path", info.ps.Active.Root),
		zap.String("resume-at-path", resumeAt),
	)

	return s.conclude(&concludeInfo{
		active:    info.ps.Active,
		root:      info.ps.Active.Root,
		current:   resumeAt,
		inclusive: true,
	})
}

type concludeInfo struct {
	active    *ActiveState
	root      string
	current   string
	inclusive bool
}

func (s *spawnStrategy) conclude(conclusion *concludeInfo) *TraverseResult {
	if conclusion.current == conclusion.active.Root {
		// reach the top, so we're done
		//
		return &TraverseResult{}
	}
	parent, child := utils.SplitParent(conclusion.current)

	following := s.following(&followingParams{
		parent:    parent,
		anchor:    child,
		order:     s.o.Store.Behaviours.Sort.DirectoryEntryOrder,
		inclusive: conclusion.inclusive,
	})
	following.siblings.sort(&following.siblings.Files)
	following.siblings.sort(&following.siblings.Folders)

	compoundResult := s.seed(&seedParams{
		frame:      s.nc.frame,
		parent:     parent,
		entries:    following.siblings.all(),
		conclusion: conclusion,
	})

	if !utils.IsNil(compoundResult.Error) {
		return compoundResult
	}
	conclusion.current = parent
	conclusion.inclusive = false

	return compoundResult.merge(s.conclude(conclusion))
}

type seedParams struct {
	frame      *navigationFrame
	parent     string
	entries    *[]fs.DirEntry
	conclusion *concludeInfo
}

func (s *spawnStrategy) seed(params *seedParams) *TraverseResult {
	params.frame.link(&linkParams{
		root:    params.conclusion.root,
		current: params.conclusion.current,
	})

	compoundResult := &TraverseResult{}
	for _, entry := range *params.entries {
		topPath := filepath.Join(params.parent, entry.Name())
		result := s.nc.impl.top(params.frame, topPath)
		compoundResult.merge(result)

		if result.Error != nil {
			return compoundResult
		}
	}
	return compoundResult
}

type shard struct {
	siblings *DirectoryEntries
}

type followingParams struct {
	parent    string
	anchor    string
	order     DirectoryEntryOrderEnum
	inclusive bool
}

func (s *spawnStrategy) following(params *followingParams) *shard {

	entries, err := s.o.Hooks.ReadDirectory(params.parent)

	if err != nil {
		panic(fmt.Sprintf("following failed to read contents of directory: '%v'",
			params.parent),
		)
	}

	groups := lo.GroupBy(entries, func(item fs.DirEntry) bool {

		if params.inclusive {
			return item.Name() >= params.anchor
		}
		return item.Name() > params.anchor
	})
	siblings := groups[_FOLLOWING_SIBLINGS]

	de := s.deFactory.new(
		&directoryEntriesFactoryParams{
			o:       s.o,
			order:   params.order,
			entries: &siblings,
		},
	)

	return &shard{siblings: de}
}
