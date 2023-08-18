package utils_test

import (
	"fmt"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/extendio/xfs/utils"
)

type RPEntry struct {
	given  string
	should string
	path   string
	expect string
}

var fakeHome = filepath.Join(string(filepath.Separator), "home", "rabbitweed")
var fakeAbs = filepath.Join(string(filepath.Separator), "home", "rabbitweed", "music", "xpander")
var fakeUpAbs = filepath.Join(string(filepath.Separator), "home", "rabbitweed", "music")

func fakeHomeResolver() string {
	return fakeHome
}

func fakeAbsResolver(path string) (string, error) {
	if strings.HasPrefix(path, "..") {
		return filepath.Join(fakeUpAbs, path[2:]), nil
	}

	if strings.HasPrefix(path, ".") {
		return filepath.Join(fakeAbs, path[1:]), nil
	}

	return path, nil
}

var _ = Describe("ResolvePath", func() {
	DescribeTable("Overrides provided",
		func(entry *RPEntry) {
			overrides := utils.ResolveOverrides{
				HomeFunc: fakeHomeResolver,
				AbsFunc:  fakeAbsResolver,
			}

			if filepath.Separator == '/' {
				actual := utils.ResolvePath(entry.path, overrides)
				Expect(actual).To(Equal(entry.expect))
			} else {
				normalisedPath := strings.ReplaceAll(entry.path, "/", string(filepath.Separator))
				normalisedExpect := strings.ReplaceAll(entry.expect, "/", string(filepath.Separator))

				actual := utils.ResolvePath(normalisedPath, overrides)
				Expect(actual).To(Equal(normalisedExpect))
			}
		},
		func(entry *RPEntry) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'", entry.given, entry.should)
		},

		Entry(nil, &RPEntry{
			given:  "path is a valid absolute path",
			should: "return path unmodified",
			path:   "/home/rabbitweed/foo",
			expect: "/home/rabbitweed/foo",
		}),
		Entry(nil, &RPEntry{
			given:  "path contains leading ~",
			should: "replace ~ with home path",
			path:   "~/foo",
			expect: "/home/rabbitweed/foo",
		}),
		Entry(nil, &RPEntry{
			given:  "path is relative to cwd",
			should: "replace ~ with home path",
			path:   "./foo",
			expect: "/home/rabbitweed/music/xpander/foo",
		}),
		Entry(nil, &RPEntry{
			given:  "path is relative to parent",
			should: "replace ~ with home path",
			path:   "../foo",
			expect: "/home/rabbitweed/music/foo",
		}),
	)

	When("No overrides provided", func() {
		Context("and: home", func() {
			It("🧪 should: not fail", func() {
				utils.ResolvePath("~/")
			})
		})

		Context("and: abs", func() {
			It("🧪 should: not fail", func() {
				utils.ResolvePath("./")
			})
		})
	})
})
