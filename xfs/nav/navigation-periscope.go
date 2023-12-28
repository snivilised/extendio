package nav

import (
	"path/filepath"
	"strings"
)

// navigationPeriscope: depth and scope manager
type navigationPeriscope struct {
	_offset int
	_depth  int
}

func (p *navigationPeriscope) scope(isLeaf bool) FilterScopeBiEnum {
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

func (p *navigationPeriscope) depth() int {
	return p._offset + p._depth - 1
}

func (p *navigationPeriscope) difference(root, current string) {
	rootSize := len(strings.Split(root, string(filepath.Separator)))
	currentSize := len(strings.Split(current, string(filepath.Separator)))

	if rootSize > currentSize {
		panic(NewInvalidPeriscopeRootPathNativeError(root, current))
	}

	p._offset = currentSize - rootSize
}

func (p *navigationPeriscope) descend(max uint) bool {
	if max > 0 && p._depth > int(max) {
		return false
	}

	p._depth++

	return true
}

func (p *navigationPeriscope) ascend() {
	p._depth--
}
