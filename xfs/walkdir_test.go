package xfs_test

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	"github.com/snivilised/extendio/xfs"
)

// new function to create will be Traverse
// TRAP CTRL-C: https://nathanleclaire.com/blog/2014/08/24/handling-ctrl-c-interrupt-signal-in-golang-programs/

var _ = Describe("WalkDir", Ordered, func() {

	var root string

	BeforeEach(func() {
		Skip("comprehension only")
	})

	BeforeAll(func() {
		if current, err := os.Getwd(); err == nil {
			parent, _ := filepath.Split(current)
			root = filepath.Join(parent, "Test", "data", "MUSICO")
		}
	})

	Context("MUSICO", func() {
		It("🧪 should: walk", func() {
			GinkgoWriter.Printf("---> 🔰 ROOT-PATH: '%v' ...\n", root)
			xfs.WalkOver(root)
			Expect(true)
		})
	})

	Context("Comprehension", func() {
		Context("Walk", func() {
			Context("Walk", func() {
				It("🧪 should: Walk (standard)", func() {
					// Walk navigates directories on a sorted basis, except that case
					// is significant, so the order is probably not what is expected
					// or desired. Eg "PROGRESSIVE-ROCK" is processed before "metal".
					// In addition, it invokes the callback for every item and calling
					// Lstat on every file. So "Walk" is not acceptable.
					//
					rock := filepath.Join(root, "rock")
					if err := filepath.Walk(rock, func(path string, info os.FileInfo, err error) error {
						GinkgoWriter.Printf("---> 💎 WALK-PATH: '%v' ...\n", path)

						return nil
					}); err != nil {
						Expect(false)
					}
				})
			})

			Context("Walk", func() {
				It("🧪 should: WalkDir (standard)", func() {
					// WalkDir invokes the callback for files, so not suitable
					//
					dream := filepath.Join(root, "DREAM-POP")

					if err := filepath.WalkDir(dream, func(path string, d fs.DirEntry, err error) error {
						GinkgoWriter.Printf("---> 💎 WALK-PATH: '%v' ...\n", path)

						return nil
					}); err != nil {
						Expect(false)
					}
				})
			})

			Context("WalkDir", func() {
				It("🧪 should: WalkDir (standard)", func() {
					// WalkDir invokes the callback for files, so not suitable
					//
					dream := filepath.Join(root, "DREAM-POP")

					if err := filepath.WalkDir(dream, func(path string, d fs.DirEntry, err error) error {
						GinkgoWriter.Printf("---> 📚 WALKDIR-PATH: '%v' ...\n", path)

						return nil
					}); err != nil {
						Expect(false)
					}
				})
			})
		})

		Context("ReadDir", func() {
			It("🧪 should: read contents associated with a directory", func() {
				GinkgoWriter.Printf("---> 🔰 ROOT-PATH (all-entries): '%v' ...\n", root)
				if d, err := os.Open(root); err == nil {
					defer d.Close()

					if entries, err := d.ReadDir(-1); err == nil {
						for _, e := range entries {
							GinkgoWriter.Printf("---> 💠 ENTRY: '%v' ...\n", e.Name())
						}
					}
				}
			})

			It("should: exclude file entries (unsorted)", func() {
				GinkgoWriter.Printf("---> 🔰 ROOT-PATH (directories-only): '%v' ...\n", root)

				if d, err := os.Open(root); err == nil {
					defer d.Close()

					if entries, err := d.ReadDir(-1); err == nil {
						dirs := lo.Filter(entries, func(de fs.DirEntry, i int) bool {
							return de.Type().IsDir()
						})

						for _, d := range dirs {
							GinkgoWriter.Printf("---> 🧊 DIRECTORY-ENTRY: '%v' ...\n", d.Name())
						}
					}
				}
			})

			It("should: exclude file entries (sorted, case sensitive)", func() {
				GinkgoWriter.Printf("---> 🔰 ROOT-PATH (directories-only): '%v' ...\n", root)

				if d, err := os.Open(root); err == nil {
					defer d.Close()

					if entries, err := d.ReadDir(-1); err == nil {
						dirs := lo.Filter(entries, func(de fs.DirEntry, i int) bool {
							return de.Type().IsDir()
						})
						sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })

						for _, d := range dirs {
							GinkgoWriter.Printf("---> 🧊 DIRECTORY-ENTRY: '%v' ...\n", d.Name())
						}
					}
				}
			})

			It("should: exclude file entries (sorted, case insensitive)", func() {
				GinkgoWriter.Printf("---> 🔰 ROOT-PATH (directories-only): '%v' ...\n", root)

				if d, err := os.Open(root); err == nil {
					defer d.Close()

					if entries, err := d.ReadDir(-1); err == nil {
						dirs := lo.Filter(entries, func(de fs.DirEntry, i int) bool {
							return de.Type().IsDir()
						})

						sort.Slice(dirs, func(i, j int) bool {
							return strings.ToLower(dirs[i].Name()) < strings.ToLower(dirs[j].Name())
						})

						for _, d := range dirs {
							GinkgoWriter.Printf("---> 🧊 DIRECTORY-ENTRY: '%v' ...\n", d.Name())
						}
					}
				}
			})
		})

		Context("Readdirnames", func() {
			It("🧪 should: read names of items in a directory", func() {
				GinkgoWriter.Printf("---> 🔰 ROOT-PATH: '%v' ...\n", root)

				if d, err := os.Open(root); err == nil {
					defer d.Close()

					if names, err := d.Readdirnames(-1); err == nil {
						for _, n := range names {
							GinkgoWriter.Printf("---> ➕ ENTRY-NAME: '%v' ...\n", n)
						}
					}
				}
			})
		})
	})

	Context("File Mode Bit Patterns", func() {
		It("🧪 should: show bit patterns of various FileMode values", func() {
			values := []fs.FileMode{fs.ModeDir, fs.ModeAppend, fs.ModeType, fs.ModePerm}

			for _, val := range values {
				raw := uint32(val)
				GinkgoWriter.Printf("---> 🤖 ModeDir: '%v' (%v)\n", val, raw)
			}
		})
	})

	Context("FilterScopeEnum", func() {
		It("🧪 should: show bit patterns of various FilterScopeEnum values", func() {
			values := []xfs.FilterScopeEnum{
				xfs.LeafNodes, xfs.TopNodes, xfs.IntermediateNodes, xfs.FolderNodes, xfs.FileNodes, xfs.AllNodes,
			}

			for _, val := range values {
				raw := uint32(val)
				GinkgoWriter.Printf("---> 🐉 FilterScopeEnum: '%v' (%v)\n", val, raw)
			}
		})

		It("🧪 should: show bit pattern when multiple node types defined", func() {
			values := []xfs.FilterScopeEnum{
				xfs.LeafNodes | xfs.TopNodes,
			}
			for _, val := range values {
				raw := uint32(val)
				GinkgoWriter.Printf("---> 🐦 Multi - FilterScopeEnum: '%v' (%v)\n", val, raw)
			}
		})
	})
})
