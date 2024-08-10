package nav

import (
	"path/filepath"

	"github.com/snivilised/extendio/internal/lo"
)

// Tail extracts the end part of a string, starting from the offset
func Tail(input string, offset int) string {
	asRunes := []rune(input)

	if offset >= len(asRunes) {
		return ""
	}

	return string(asRunes[offset:])
}

// difference returns the difference between a child path and a parent path
// Designed to be used with paths created from the file system rather than
// custom created or user provided input. For this reason, if there is no
// relationship between the parent and child paths provided then a panic
// may occur.
func difference(parent, child string) string {
	return Tail(child, len(parent))
}

// RootItemSubPathHookFn
func RootItemSubPathHookFn(info *SubPathInfo) string {
	return difference(info.Root, info.Item.Path)
}

// RootParentSubPathHookFn
func RootParentSubPathHookFn(info *SubPathInfo) string {
	if info.Item.Extension.NodeScope == ScopeTopEn {
		return lo.Ternary(info.Behaviour.KeepTrailingSep, string(filepath.Separator), "")
	}

	return difference(info.Root, info.Item.Extension.Parent)
}
