package nav

// Tail extracts the end part of a string, starting from the offset
func Tail(input string, offset int) string {
	asRunes := []rune(input)

	if offset >= len(asRunes) {
		return ""
	}

	return string(asRunes[offset:])
}

// Difference returns the difference between a child path and a parent path
// Designed to be used with paths created from the file system rather than
// custom created or user provided input. For this reason, if there is no
// relationship between the parent and child paths provided then a panic
// may occur.
func Difference(parent string, child string) string {
	return Tail(child, len(parent))
}

func RootItemSubPath(info *SubPathInfo) string {
	return Difference(info.Root, info.Item.Path)
}

func RootParentSubPath(info *SubPathInfo) string {

	if info.Item.Extension.NodeScope == TopScopeEn {
		return ""
	}
	return Difference(info.Root, info.Item.Extension.Parent)
}
