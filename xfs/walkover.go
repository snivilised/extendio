package xfs

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/samber/lo"
)

func WalkOver(root string) {
	err := WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		const p, l = 0, 1

		parent, leaf := filepath.Split(path)
		emojis := lo.Ternary(entry.IsDir(), []string{"🍉", "🍂"}, []string{"🧊", "🍃"})

		fmt.Printf("---> '%s %s' (%s %s)\n", emojis[p], parent, emojis[l], leaf)
		return nil
	})
	if err != nil {
		fmt.Printf("---> 💥ERROR: '%v' ...\n", err)
	}
}

func WalkOverOnlyDirectories(root string) {
	err := WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if !entry.IsDir() {
			return nil
		}
		const p, l = 0, 1

		parent, leaf := filepath.Split(path)
		emojis := lo.Ternary(entry.IsDir(), []string{"🍉", "🍂"}, []string{"🧊", "🍃"})

		fmt.Printf("---> '%s %s' (%s %s)\n", emojis[p], parent, emojis[l], leaf)
		return nil
	})
	if err != nil {
		fmt.Printf("---> 💥ERROR: '%v' ...\n", err)
	}
}
