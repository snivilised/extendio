package nav

import (
	"fmt"
)

// PolyFilter is a dual filter that allows files and folders to be filtered
// independently. The Folder filter only applies when the current item
// is a file. This is because, filtering doesn't affect navigation, it only
// controls wether the client callback is invoked or not. That is to say, if
// a particular folder fails to pass a filter, the callback will not be
// invoked for that folder, but we still descend into it and navigate its
// children. This is the reason why the poly filter is only active when the
// the current item is a filter as the client callback will only be invoked
// for the file if its parent folder passes the poly folder filter and
// the file passes the poly file filter.
type PolyFilter struct {
	// File is the filter that applies to a file. Note that the client does
	// not have to set the File scope as this is enforced automatically as
	// well as ensuring that the Folder scope has not been set. The client is
	// still free to set other scopes.
	File TraverseFilter

	// Folder is the filter that applies to a folder. Note that the client does
	// not have to set the Folder scope as this is enforced automatically as
	// well as ensuring that the File scope has not been set. The client is
	// still free to set other scopes.
	Folder TraverseFilter
}

// Description
func (f *PolyFilter) Description() string {
	return fmt.Sprintf("Poly - FILE: '%v', FOLDER: '%v'",
		f.File.Description(), f.Folder.Description(),
	)
}

// Validate ensures that both filters definition are valid, panics when invalid
func (f *PolyFilter) Validate() {
	f.File.Validate()
	f.Folder.Validate()
}

// Source returns the Sources of both the File and Folder filters separated
// by a '##'
func (f *PolyFilter) Source() string {
	return fmt.Sprintf("%v##%v",
		f.File.Source(), f.Folder.Source(),
	)
}

// IsMatch returns true if the current item is a file and both the current
// file matches the poly file filter and the file's parent folder matches
// the poly folder filter. Returns true of the current item is a folder.
func (f *PolyFilter) IsMatch(item *TraverseItem) bool {
	if !item.IsDirectory() {
		return f.Folder.IsMatch(item.Parent) && f.File.IsMatch(item)
	}

	return true
}

// IsApplicable returns the result of applying IsApplicable to
// the poly Filter filter if the current item is a file, returns false
// for folders.
func (f *PolyFilter) IsApplicable(item *TraverseItem) bool {
	if !item.IsDirectory() {
		return f.File.IsApplicable(item)
	}

	return false
}

// Scope is a bitwise OR combination of both filters
func (f *PolyFilter) Scope() FilterScopeBiEnum {
	return f.File.Scope() | f.Folder.Scope()
}
