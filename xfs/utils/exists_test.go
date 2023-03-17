package utils_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/snivilised/extendio/internal/helpers"

	"github.com/snivilised/extendio/xfs/utils"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func musico() string {
	if current, err := os.Getwd(); err == nil {
		parent, _ := filepath.Split(current)
		grand := filepath.Dir(parent)
		great := filepath.Dir(grand)
		result := filepath.Join(great, "Test", "data", "MUSICO")
		must(helpers.Ensure(result))

		return result
	}
	panic("could not get root path")
}

func path(parent, relative string) string {
	segments := strings.Split(relative, "/")
	return filepath.Join(append([]string{parent}, segments...)...)
}

var _ = Describe("Exists Utils", Ordered, func() {
	var root, heavy string

	BeforeAll(func() {
		root = musico()
		heavy = filepath.Join(root, "rock", "metal", "dark", "HEAVY-METAL")
	})

	DescribeTable("Exists",
		func(message, relative string, expected bool) {
			path := path(heavy, relative)

			GinkgoWriter.Printf("---> 🔰 FULL-PATH: '%v'\n", path)
			Expect(utils.Exists(path)).To(Equal(expected))
		},

		func(message, relative string, expected bool) string {
			return fmt.Sprintf("🥣 message: '%v'", message)
		},
		Entry(nil, "existing folder", "Mötley Crüe/Theatre of Pain", true),
		Entry(nil, "existing file", "Mötley Crüe/Theatre of Pain/01 - City Boy Blues.flac", true),
		Entry(nil, "missing", "Mötley Crüe/Insomnia", false),
	)

	DescribeTable("FolderExists",
		func(message, relative string, expected bool) {
			path := path(heavy, relative)
			GinkgoWriter.Printf("---> 🔰 FULL-PATH: '%v'\n", path)

			Expect(utils.FolderExists(path)).To(Equal(expected))
		},
		func(message, relative string, expected bool) string {
			return fmt.Sprintf("🍤 message: '%v'", message)
		},
		Entry(nil, "folder exists", "Mötley Crüe/Theatre of Pain", true),
		Entry(nil, "folder does not exist", "Mötley Crüe/Theatre of Pain/Insomnia", false),

		Entry(nil, "item exists as file", "Mötley Crüe/Theatre of Pain/01 - City Boy Blues.flac", false),
	)

	DescribeTable("FileExists",
		func(message, relative string, expected bool) {
			path := path(heavy, relative)
			GinkgoWriter.Printf("---> 🔰 FULL-PATH: '%v'\n", path)

			Expect(utils.FileExists(path)).To(Equal(expected))
		},
		func(message, relative string, expected bool) string {
			return fmt.Sprintf("🍤 message: '%v'", message)
		},
		Entry(nil, "file exists", "Mötley Crüe/Theatre of Pain/01 - City Boy Blues.flac", true),
		Entry(nil, "file does not exist", "Mötley Crüe/Theatre of Pain/Insomnia", false),

		Entry(nil, "item exists as folder", "Mötley Crüe/Theatre of Pain", false),
	)
})
