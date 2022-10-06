package xfs

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/samber/lo"
)

func Walkover(root string) {
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		parent, leaf := filepath.Split(path)

		emoji := lo.Ternary(entry.IsDir(), "🍉", "🧊")
		fmt.Printf("---> '%s %s' (🌿 %s)\n", emoji, parent, leaf)
		return nil
	})
	if err != nil {
		fmt.Printf("---> 💥ERROR: '%v' ...\n", err)
	}
}
