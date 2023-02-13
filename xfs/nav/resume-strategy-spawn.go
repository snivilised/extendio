package nav

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/samber/lo"
	. "github.com/snivilised/extendio/translate"
	"github.com/snivilised/extendio/xfs/utils"
)

const (
	_FOLLOWING_SIBLINGS = true
)

type spawnStrategy struct {
	baseStrategy
}

func (s *spawnStrategy) resume(info *strategyResumeInfo) *TraverseResult {
	s.nc.frame.root.Set(info.ps.Active.Root)
	resumeAt := s.ps.Active.NodePath

	statusInfo, err := s.o.Hooks.QueryStatus(resumeAt)

	if err != nil {
		return &TraverseResult{
			Error: LocalisableError{
				Inner: err,
			},
		}
	}

	indicator := lo.Ternary(statusInfo.IsDir(), "📁 DIRECTORY", "📜 FILE")
	fmt.Printf("   🧿 root: '%v' \n", info.ps.Active.Root)
	fmt.Printf("   🧿 resume-at(%v): '%v' \n", indicator, resumeAt)
	fmt.Println("   =====================================")

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
	fmt.Printf("   👾 conclude: '%v' \n", conclusion.current)

	if conclusion.current == conclusion.active.Root {
		// reach the top, so we're done
		//
		fmt.Printf("   👽 conclude - completed at: 🧿'%v' \n", conclusion.current)

		return &TraverseResult{}
	}

	parent, child := utils.SplitParent(conclusion.current)
	fmt.Printf("   💥 - parent: 🧿'%v' \n", parent)
	fmt.Printf("   💥 - child: 🧿'%v' \n", child)

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
	fmt.Println("   =====================================")

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
	fmt.Print("   🎈seeds: ")
	for _, entry := range *params.entries {
		fmt.Printf("'%v', ", entry.Name())
	}
	fmt.Println("")

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
	siblings *directoryEntries
}

type followingParams struct {
	parent    string
	anchor    string
	order     DirectoryEntryOrderEnum
	inclusive bool
}

func (s *spawnStrategy) following(params *followingParams) *shard {

	entries, err := s.o.Hooks.ReadDirectory(params.parent)

	// TODO: This should not be a panic
	//
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

	de := s.deFactory.construct(
		&directoryEntriesFactoryParams{
			o:       s.o,
			order:   params.order,
			entries: &siblings,
		},
	)

	return &shard{siblings: de}
}
