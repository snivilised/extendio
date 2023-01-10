package nav

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/samber/lo"
	. "github.com/snivilised/extendio/translate"
	"github.com/snivilised/extendio/xfs/utils"
)

type executionSequence []func() *LocalisableError

// behaves like the universal navigator with a few exceptions
type spawnerImpl struct {
	navigator
	ps *persistState
}

func (s *spawnerImpl) top(frame *navigationFrame) *LocalisableError {
	fmt.Printf(">>> spawnerImpl.top 🎪 \n")

	info, err := s.o.Hooks.QueryStatus(s.ps.Active.NodePath)

	if err != nil {
		return &LocalisableError{
			Inner: err,
		}
	}
	indicator := lo.Ternary(info.IsDir(), "📁 DIRECTORY", "📜 FILE")
	fmt.Printf("   🧿 root: '%v' \n", frame.root)
	fmt.Printf("   🧿 resume-at(%v): '%v' \n", indicator, s.ps.Active.NodePath)

	parent, child := utils.SplitParent(s.ps.Active.NodePath)
	fmt.Printf("   - parent: 🧿'%v', child: 🧿'%v' \n", parent, child)

	fracture := s.siblingsFollowing(&followingInfo{
		parent: parent,
		order:  s.o.Store.Behaviours.Sort.DirectoryEntryOrder,
		anchor: child,
	})
	fracture.siblings.sort(&fracture.siblings.Files)
	fracture.siblings.sort(&fracture.siblings.Folders)
	sorted := fracture.siblings.all()

	sequence := executionSequence{
		func() *LocalisableError {
			return s.files(&topSpawnParams{
				frame:   frame,
				parent:  parent,
				entries: sorted,
			})
		},
	}

	var le *LocalisableError
	for _, fn := range sequence {
		if le = fn(); le != nil {
			break
		}
	}

	return le
}

type topSpawnParams struct {
	frame   *navigationFrame
	parent  string
	entries *[]fs.DirEntry
}

func (s *spawnerImpl) files(params *topSpawnParams) *LocalisableError {
	for _, entry := range *params.entries {
		if !entry.IsDir() {
			topPath := filepath.Join(params.parent, entry.Name())
			le := s.agent.top(&agentTopParams{
				impl:  s,
				frame: params.frame,
				top:   topPath,
			})

			if le != nil {
				return le
			}
		}
	}
	return nil
}
func (s *spawnerImpl) folders(params *topSpawnParams) *LocalisableError {
	for _, entry := range *params.entries {
		if entry.IsDir() {
			topPath := filepath.Join(params.parent, entry.Name())
			le := s.agent.top(&agentTopParams{
				impl:  s,
				frame: params.frame,
				top:   topPath,
			})

			if le != nil {
				return le
			}
		}
	}
	return nil
}

func (s *spawnerImpl) traverse(params *traverseParams) *LocalisableError {
	fmt.Printf(">>> spawnerImpl:traverse ✈️ \n")
	// get the children for this item, then split those children according to breakpoint
	//
	return nil
}

type followingInfo struct {
	parent string
	order  DirectoryEntryOrderEnum
	anchor string
}

type fractureInfo struct {
	siblings *directoryEntries
}

func (s *spawnerImpl) siblingsFollowing(following *followingInfo) *fractureInfo {

	// s.agent.read()
	entries, err := s.o.Hooks.ReadDirectory(following.parent)

	if err != nil {
		panic(fmt.Sprintf("siblingsFollowing failed to read contents of directory: '%v'",
			following.parent),
		)
	}

	groups := lo.GroupBy(entries, func(item fs.DirEntry) bool {
		return item.Name() >= following.anchor
	})

	siblings := groups[true]
	de := directoryEntries{
		Options: s.o,
		Order:   following.order,
	}
	de.arrange(&siblings)

	return &fractureInfo{siblings: &de}
}

func (s *spawnerImpl) spawn(params *spawnParams) *LocalisableError {

	return nil
}

func (s *spawnerImpl) seed(params *seedParams) *LocalisableError {
	return nil
}
