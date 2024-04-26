package utils_test

import (
	"fmt"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ginkgo ok
	. "github.com/onsi/gomega"    //nolint:revive // gomega ok
	"github.com/snivilised/extendio/internal/helpers"

	"github.com/snivilised/extendio/xfs/utils"
)

func path(parent, relative string) string {
	segments := strings.Split(relative, "/")
	return filepath.Join(append([]string{parent}, segments...)...)
}

var _ = Describe("Exists Utils", Ordered, func() {
	var repo string

	BeforeAll(func() {
		repo = helpers.Repo("../..")
		Expect(utils.FolderExists(repo)).To(BeTrue())
	})

	DescribeTable("Exists",
		func(_, relative string, expected bool) {
			path := path(repo, relative)

			GinkgoWriter.Printf("---> 🔰 FULL-PATH: '%v'\n", path)
			Expect(utils.Exists(path)).To(Equal(expected))
		},

		func(message, _ string, _ bool) string {
			return fmt.Sprintf("🥣 message: '%v'", message)
		},
		Entry(nil, "folder exists", "/", true),
		Entry(nil, "file exists", "README.md", true),
		Entry(nil, "does not exist", "foo-bar", false),
	)

	DescribeTable("FolderExists",
		func(_, relative string, expected bool) {
			path := path(repo, relative)
			GinkgoWriter.Printf("---> 🔰 FULL-PATH: '%v'\n", path)

			Expect(utils.FolderExists(path)).To(Equal(expected))
		},
		func(message, _ string, _ bool) string {
			return fmt.Sprintf("🍤 message: '%v'", message)
		},
		Entry(nil, "folder exists", "/", true),
		Entry(nil, "folder does not exist", "foo-bar", false),
		Entry(nil, "exists as file", "README.md", false),
	)

	DescribeTable("FileExists",
		func(_, relative string, expected bool) {
			path := path(repo, relative)
			GinkgoWriter.Printf("---> 🔰 FULL-PATH: '%v'\n", path)

			Expect(utils.FileExists(path)).To(Equal(expected))
		},
		func(message, _ string, _ bool) string {
			return fmt.Sprintf("🍤 message: '%v'", message)
		},
		Entry(nil, "file exists", "README.md", true),
		Entry(nil, "file does not exist", "foo-bar", false),
		Entry(nil, "does not exist as file", "Test", false),
	)
})
