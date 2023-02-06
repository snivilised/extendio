package nav

import (
	"path/filepath"
	"strings"
)

// navigationPeriscope: depth and scope manager
type navigationPeriscope struct {
	_offset uint
	_depth  uint
}

func (p *navigationPeriscope) scope(isLeaf bool) FilterScopeEnum {

	result := ScopeIntermediateEn

	// Root=0
	// Top=1
	//
	depth := p.depth()

	switch {
	case isLeaf && depth == 0:
		result = ScopeRootEn | ScopeLeafEn
	case depth == 0:
		result = ScopeRootEn
	case isLeaf && depth == 1:
		result = ScopeTopEn | ScopeLeafEn
	case depth == 1:
		result = ScopeTopEn
	case isLeaf:
		result = ScopeLeafEn
	}

	return result
}

func (p *navigationPeriscope) depth() uint {
	return p._offset + p._depth - uint(1)
}

func (p *navigationPeriscope) difference(root, current string) {
	rootSize := len(strings.Split(root, string(filepath.Separator)))
	currentSize := len(strings.Split(current, string(filepath.Separator)))

	if rootSize > currentSize {
		panic("navigationPeriscope: internal error, root path can't be longer than current path")
	}

	p._offset = uint(currentSize) - uint(rootSize)
}

func (p *navigationPeriscope) descend() {
	p._depth++
}

func (p *navigationPeriscope) ascend() {
	p._depth--
}
