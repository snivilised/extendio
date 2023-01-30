package nav

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"reflect"

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

func (s *spawnStrategy) init(params *strategyInitParams) {
	// TODO: set the depth and other appropriate properties on the frame
	//
	params.frame.depth = params.ps.Active.Depth
}

func (s *spawnStrategy) resume(info *strategyResumeInfo) *TraverseResult {
	info.nc.root(func() string {
		return info.ps.Active.Root
	})
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

	// TODO: ensure that child is related to parent and also, not the / path
	//
	return s.conclude(&concludeInfo{
		active:      info.ps.Active,
		currentPath: resumeAt,
	})
}

type seedParams struct {
	frame   *navigationFrame
	parent  string
	entries *[]fs.DirEntry
}

type concludeInfo struct {
	active      *ActiveState
	currentPath string
}

func (s *spawnStrategy) conclude(conclusion *concludeInfo) *TraverseResult {
	fmt.Printf("   👾 conclude: '%v' \n", conclusion.currentPath)

	result := &TraverseResult{}
	if conclusion.currentPath == conclusion.active.Root {
		// TODO: need to make sure that the term active in 'resume' scenarios
		// is a legacy item from the previous session, not the current
		// session. Perhaps, whenever active is legacy, it should be wrapped
		// in a container struct called 'legacy'.
		//
		// reach the top, so we're done
		//
		fmt.Printf("   👽 conclude - completed at: 🧿'%v' \n", conclusion.currentPath)

		return result
	}

	parent, child := utils.SplitParent(conclusion.currentPath)
	fmt.Printf("   💥 - parent: 🧿'%v' \n", parent)
	fmt.Printf("   💥 - child: 🧿'%v' \n", child)

	following := s.following(&followingParams{
		parent: parent,
		anchor: child,
		order:  s.o.Store.Behaviours.Sort.DirectoryEntryOrder,
	})
	following.siblings.sort(&following.siblings.Files)
	following.siblings.sort(&following.siblings.Folders)

	seedsFn := func() *LocalisableError {

		return s.seed(&seedParams{
			frame:   s.nc.frame,
			parent:  parent,
			entries: following.siblings.all(),
		})
	}

	result.Error = s.run(&coreSequence{
		seedsFn,
	})
	fmt.Println("   =====================================")

	if !reflect.ValueOf(result.Error).IsNil() {
		return result
	}
	conclusion.currentPath = parent

	return s.conclude(conclusion)
}

func (s *spawnStrategy) seed(params *seedParams) *LocalisableError {
	fmt.Print("   🎈seeds: ")
	for _, entry := range *params.entries {
		fmt.Printf("'%v', ", entry.Name())
	}
	fmt.Println("")

	for _, entry := range *params.entries {
		topPath := filepath.Join(params.parent, entry.Name())
		le := s.nc.impl.top(params.frame, topPath)

		if le != nil {
			return le
		}
	}
	return nil
}

type shard struct {
	siblings *directoryEntries
}

type followingParams struct {
	parent string
	anchor string
	order  DirectoryEntryOrderEnum
}

func (s *spawnStrategy) following(params *followingParams) *shard {

	entries, err := s.o.Hooks.ReadDirectory(params.parent)

	if err != nil {
		panic(fmt.Sprintf("siblingsFollowing failed to read contents of directory: '%v'",
			params.parent),
		)
	}

	groups := lo.GroupBy(entries, func(item fs.DirEntry) bool {
		return item.Name() >= params.anchor
	})

	siblings := groups[_FOLLOWING_SIBLINGS]
	de := directoryEntries{
		Options: s.o,
		Order:   params.order,
	}
	de.arrange(&siblings)

	return &shard{siblings: &de}
}

type coreSequence []func() *LocalisableError

func (s *spawnStrategy) run(sequence *coreSequence) *LocalisableError {
	var le *LocalisableError
	for _, fn := range *sequence {
		if le = fn(); le != nil {
			break
		}
	}

	return le
}
