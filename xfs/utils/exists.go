package utils

import (
	"os"
)

// Exists provides a simple way to determine whether the item identified by a
// path actually exists either as a file or a folder
func Exists(path string) bool {
	result := false
	if _, err := os.Stat(path); err == nil {
		result = true
	}
	return result
}

// FileExists provides a simple way to determine whether the item identified by a
// path actually exists as a file
func FileExists(path string) bool {
	result := false
	if info, err := os.Lstat(path); err == nil {
		result = !info.IsDir()
	}

	return result
}

// FileExists provides a simple way to determine whether the item identified by a
// path actually exists as a folder
func FolderExists(path string) bool {
	result := false
	if info, err := os.Lstat(path); err == nil {
		result = info.IsDir()
	}

	return result
}
