package nav_test

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	. "github.com/snivilised/extendio/translate"
	"github.com/snivilised/extendio/xfs/nav"
)

type recordingMap map[string]int
type recordingScopeMap map[string]nav.FilterScopeEnum
type recordingOrderMap map[string]int

type naviTE struct {
	message            string
	relative           string
	extended           bool
	once               bool
	visit              bool
	caseSensitive      bool
	subscription       nav.TraverseSubscription
	callback           nav.TraverseCallback
	mandatory          []string
	prohibited         []string
	expectedNoChildren map[string]int
}

type skipTE struct {
	naviTE
	skip    string
	exclude string
}

type listenTE struct {
	naviTE
	start    nav.Listener
	stop     nav.Listener
	incStart bool
	incStop  bool
	mute     bool
}

type filterTE struct {
	naviTE
	name            string
	pattern         string
	scope           nav.FilterScopeEnum
	negate          bool
	expectedErr     error
	errorContains   string
	ifNotApplicable bool
}

type scopeTE struct {
	naviTE
	expectedScopes recordingScopeMap
}

type sortTE struct {
	filterTE
	expectedOrder []string
	order         nav.DirectoryEntryOrderEnum
}

func cwd() string { // oops, this is a bad name for the function, cwd is incorrect
	if current, err := os.Getwd(); err == nil {
		parent, _ := filepath.Split(current)
		grand := filepath.Dir(parent)
		great := filepath.Dir(grand)
		return filepath.Join(great, "Test", "data", "MUSICO")
	}
	panic("could not get root path")
}

func joinCwd(segments ...string) string {
	if current, err := os.Getwd(); err == nil {
		parent, _ := filepath.Split(current)
		grand := filepath.Dir(parent)
		great := filepath.Dir(grand)
		all := append([]string{great}, segments...)
		return filepath.Join(all...)
	}
	panic("could not get root path")

}

const IsExtended = true
const NotExtended = false

func normalise(p string) string {
	return strings.ReplaceAll(p, "/", string(filepath.Separator))
}

func reason(name string) string {
	return fmt.Sprintf("❌ for item named: '%v'", name)
}

func begin(em string) nav.BeginHandler {
	return func(state *nav.NavigationState) {
		GinkgoWriter.Printf(
			"---> %v [traverse-navigator-test:BEGIN], root: '%v'\n", em, state.Root,
		)
	}
}

func path(parent, relative string) string {
	segments := strings.Split(relative, "/")
	return filepath.Join(append([]string{parent}, segments...)...)
}

func universalCallback(name string, extended bool) nav.TraverseCallback {

	ex := lo.Ternary(extended, "-EX", "")
	return func(item *nav.TraverseItem) *LocalisableError {
		depth := lo.TernaryF(extended,
			func() uint { return item.Extension.Depth },
			func() uint { return 9999 },
		)
		GinkgoWriter.Printf(
			"---> 🌊 UNIVERSAL//%v-CALLBACK%v: (depth:%v) '%v'\n", name, ex, depth, item.Path,
		)

		if extended {
			Expect(item.Extension).NotTo(BeNil(), reason(item.Path))
		} else {
			Expect(item.Extension).To(BeNil(), reason(item.Path))
		}
		return nil
	}
}

func foldersCallback(name string, extended bool) nav.TraverseCallback {

	ex := lo.Ternary(extended, "-EX", "")
	return func(item *nav.TraverseItem) *LocalisableError {
		depth := lo.TernaryF(extended,
			func() uint { return item.Extension.Depth },
			func() uint { return 9999 },
		)
		actualNoChildren := len(item.Children)
		GinkgoWriter.Printf(
			"---> ☀️ FOLDERS//%v-CALLBACK%v: (depth:%v, children:%v) '%v'\n",
			name, ex, depth, actualNoChildren, item.Path,
		)
		Expect(item.Info.IsDir()).To(BeTrue())
		// Expect(actualNoChildren).To(Equal(expectedNoChildren))

		if extended {
			Expect(item.Extension).NotTo(BeNil(), reason(item.Path))
		} else {
			Expect(item.Extension).To(BeNil(), reason(item.Path))
		}

		return nil
	}
}

func filesCallback(name string, extended bool) nav.TraverseCallback {

	ex := lo.Ternary(extended, "-EX", "")
	return func(item *nav.TraverseItem) *LocalisableError {
		GinkgoWriter.Printf("---> 🌙 FILES//%v-CALLBACK%v: '%v'\n", name, ex, item.Path)
		Expect(item.Info.IsDir()).To(BeFalse())

		if extended {
			Expect(item.Extension).NotTo(BeNil(), reason(item.Path))
		}
		return nil
	}
}

// === scope

func universalScopeCallback(name string) nav.TraverseCallback {

	return func(item *nav.TraverseItem) *LocalisableError {
		GinkgoWriter.Printf("---> 🌠 UNIVERSAL//%v-CALLBACK-EX item-scope: (%v) '%v'\n",
			name, item.Extension.NodeScope, item.Extension.Name,
		)
		Expect(item.Extension).NotTo(BeNil(), reason(item.Extension.Name))
		return nil
	}
}

func foldersScopeCallback(name string) nav.TraverseCallback {

	return func(item *nav.TraverseItem) *LocalisableError {
		GinkgoWriter.Printf("---> 🌟 FOLDERS//%v-CALLBACK-EX item-scope: (%v) '%v'\n",
			name, item.Extension.NodeScope, item.Extension.Name,
		)
		Expect(item.Extension).NotTo(BeNil(), reason(item.Extension.Name))
		Expect(item.Info.IsDir()).To(BeTrue())
		return nil
	}
}

func filesScopeCallback(name string) nav.TraverseCallback {

	return func(item *nav.TraverseItem) *LocalisableError {
		GinkgoWriter.Printf("---> 🌬️ FILES//%v-CALLBACK-EX item-scope: (%v) '%v'\n",
			name, item.Extension.NodeScope, item.Extension.Name,
		)
		Expect(item.Extension).NotTo(BeNil(), reason(item.Extension.Name))
		Expect(item.Info.IsDir()).To(BeFalse())
		return nil
	}
}

// === sort

func universalSortCallback(name string) nav.TraverseCallback {

	return func(item *nav.TraverseItem) *LocalisableError {
		GinkgoWriter.Printf("---> 💚 UNIVERSAL//%v-SORT-CALLBACK-EX(scope:%v, depth:%v) '%v'\n",
			name, item.Extension.NodeScope, item.Extension.Depth, item.Extension.Name,
		)
		Expect(item.Extension).NotTo(BeNil(), reason(item.Extension.Name))
		return nil
	}
}

func foldersSortCallback(name string) nav.TraverseCallback {

	return func(item *nav.TraverseItem) *LocalisableError {
		GinkgoWriter.Printf("---> 💜 FOLDERS//%v-SORT-CALLBACK-EX '%v'\n",
			name, item.Extension.Name,
		)
		Expect(item.Extension).NotTo(BeNil(), reason(item.Extension.Name))
		Expect(item.Info.IsDir()).To(BeTrue())
		return nil
	}
}

func filesSortCallback(name string) nav.TraverseCallback {

	return func(item *nav.TraverseItem) *LocalisableError {
		GinkgoWriter.Printf("---> 💙 FILES//%v-SORT-CALLBACK-EX '%v'\n",
			name, item.Extension.Name,
		)
		Expect(item.Extension).NotTo(BeNil(), reason(item.Extension.Name))
		Expect(item.Info.IsDir()).To(BeFalse())
		return nil
	}
}

func universalDepthCallback(name string, maxDepth uint) nav.TraverseCallback {

	return func(item *nav.TraverseItem) *LocalisableError {
		if item.Extension.Depth <= maxDepth {
			GinkgoWriter.Printf("---> 💚 UNIVERSAL//%v-SORT-CALLBACK-EX(scope:%v, depth:%v) '%v'\n",
				name, item.Extension.NodeScope, item.Extension.Depth, item.Extension.Name,
			)
		}
		Expect(item.Extension).NotTo(BeNil(), reason(item.Extension.Name))
		return nil
	}
}

func foldersCaseSensitiveCallback(first, second string) nav.TraverseCallback {
	recording := recordingMap{}

	return func(item *nav.TraverseItem) *LocalisableError {
		recording[item.Path] = len(item.Children)

		GinkgoWriter.Printf("---> ☀️ CASE-SENSITIVE-CALLBACK: '%v'\n", item.Path)
		Expect(item.Info.IsDir()).To(BeTrue())

		if strings.HasSuffix(item.Path, second) {
			GinkgoWriter.Printf("---> 💧 FIRST: '%v', 💧 SECOND: '%v'\n", first, second)

			paths := lo.Keys(recording)
			_, found := lo.Find(paths, func(s string) bool {
				return strings.HasSuffix(s, first)
			})

			Expect(found).To(BeTrue())
		}

		return nil
	}
}

// === skip

func skipFolderCallback(skip, exclude string) nav.TraverseCallback {

	return func(item *nav.TraverseItem) *LocalisableError {
		GinkgoWriter.Printf(
			"---> ♻️ ON-NAVIGATOR-SKIP-CALLBACK(skip:%v): '%v'\n", skip, item.Path,
		)

		Expect(strings.HasSuffix(item.Path, exclude)).To(BeFalse())

		return lo.Ternary(strings.HasSuffix(item.Path, skip),
			&LocalisableError{Inner: fs.SkipDir}, nil,
		)
	}
}

func subscribes(subscription nav.TraverseSubscription, de fs.DirEntry) bool {

	any := (subscription == nav.SubscribeAny)
	files := (subscription == nav.SubscribeFiles) && (!de.IsDir())
	folders := ((subscription == nav.SubscribeFolders) || subscription == nav.SubscribeFoldersWithFiles) && (de.IsDir())

	return any || files || folders
}

// === errors

type errorTE struct {
	naviTE
}

func readDirFakeError(dirname string) ([]fs.DirEntry, error) {

	entries := []fs.DirEntry{}
	err := errors.New("fake read error")
	return entries, err
}

func readDirFakeErrorAt(name string) func(dirname string) ([]fs.DirEntry, error) {

	return func(dirname string) ([]fs.DirEntry, error) {
		if strings.HasSuffix(dirname, name) {
			return readDirFakeError(dirname)
		}

		return nav.ReadEntries(dirname)
	}
}

func errorCallback(name string, extended bool, hasError bool) nav.TraverseCallback {

	ex := lo.Ternary(extended, "-EX", "")
	return func(item *nav.TraverseItem) *LocalisableError {
		GinkgoWriter.Printf("---> 🔥 %v-CALLBACK%v: '%v'\n", name, ex, item.Path)

		if extended {
			Expect(item.Extension).NotTo(BeNil(), reason(item.Path))
		} else {
			Expect(item.Extension).To(BeNil(), reason(item.Path))
		}
		if hasError {
			Expect(item.Error).ToNot(BeNil())
		}
		return item.Error
	}
}
