package nav_test

import (
	"errors"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	"github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/internal/helpers"
	"github.com/snivilised/extendio/internal/log"
	"github.com/snivilised/extendio/xfs/nav"
	"github.com/snivilised/extendio/xfs/utils"
)

type recordingMap map[string]int
type recordingScopeMap map[string]nav.FilterScopeBiEnum
type recordingOrderMap map[string]int

type directoryQuantities struct {
	files    uint
	folders  uint
	children map[string]int
}

type naviTE struct {
	message       string
	should        string
	relative      string
	extended      bool
	once          bool
	visit         bool
	caseSensitive bool
	subscription  nav.TraverseSubscription
	callback      *nav.LabelledTraverseCallback
	mandatory     []string
	prohibited    []string
	expectedNoOf  directoryQuantities
}

type skipTE struct {
	naviTE
	skipAt       string
	prohibit     string
	all          bool
	expectedNoOf directoryQuantities
}

type sampleTE struct {
	naviTE
	sampleType   nav.SampleTypeEnum
	reverse      bool
	filter       *filterTE
	noOf         nav.EntryQuantities
	each         nav.EachDirectoryItemPredicate
	while        nav.WhileDirectoryPredicate
	expectedNoOf directoryQuantities
}

type listenTE struct {
	naviTE
	listenDefs *nav.ListenDefinitions
	incStart   bool
	incStop    bool
	mute       bool
}

type filterTE struct {
	naviTE
	name            string
	pattern         string
	scope           nav.FilterScopeBiEnum
	negate          bool
	expectedErr     error
	errorContains   string
	ifNotApplicable nav.TriStateBoolEnum
}

type marshalTE struct {
	naviTE
	errorContains string
	format        nav.PersistenceFormatEnum
}

type scopeTE struct {
	naviTE
	expectedScopes recordingScopeMap
}

type sortTE struct {
	filterTE
	expectedOrder []string
	order         nav.DirectoryContentsOrderEnum
}

type activeTE struct {
	resumeAt    string
	listenState nav.ListeningState
}

type resumeTE struct {
	naviTE
	active         activeTE
	clientListenAt string
	profile        string
	log            bool
}

type resumeTestProfile struct {
	filtered   bool
	prohibited map[string]string
	mandatory  []string
}

func musico() string {
	if current, err := os.Getwd(); err == nil {
		parent, _ := filepath.Split(current)
		grand := filepath.Dir(parent)
		great := filepath.Dir(grand)
		result := filepath.Join(great, "Test", "data", "MUSICO")

		utils.Must(helpers.Ensure(result))

		return result
	}

	panic("could not get root path")
}

func logo() nav.LoggingOptions {
	return nav.LoggingOptions{
		Enabled:         true,
		Path:            helpers.Log(),
		TimeStampFormat: "2006-01-02 15:04:05",
		Rotation: nav.LogRotationOptions{
			MaxSizeInMb:    5,
			MaxNoOfBackups: 1,
			MaxAgeInDays:   7,
		},
	}
}

const IsExtended = true
const NotExtended = false

func begin(em string) nav.BeginHandler {
	return func(state *nav.NavigationState) {
		state.Logger.Get().Info("💧 Beginning Traversal (client side)",
			log.String("Root", state.Root.Get()),
			log.Uint("Foo", 42),
			log.Int("Bar", 13),
			log.Float64("Pi", float64(math.Pi)),
		)

		GinkgoWriter.Printf(
			"---> %v [traverse-navigator-test:BEGIN], root: '%v'\n", em, state.Root,
		)
	}
}

func universalCallback(name string, extended bool) *nav.LabelledTraverseCallback {
	ex := lo.Ternary(extended, "-EX", "")

	return &nav.LabelledTraverseCallback{
		Label: "test universal callback",
		Fn: func(item *nav.TraverseItem) error {
			depth := lo.TernaryF(extended,
				func() int { return item.Extension.Depth },
				func() int { return 9999 },
			)
			GinkgoWriter.Printf(
				"---> 🌊 UNIVERSAL//%v-CALLBACK%v: (depth:%v) '%v'\n", name, ex, depth, item.Path,
			)

			if extended {
				Expect(item.Extension).NotTo(BeNil(), helpers.Reason(item.Path))
			} else {
				Expect(item.Extension).To(BeNil(), helpers.Reason(item.Path))
			}
			return nil
		},
	}
}

func universalCallbackNoAssert(name string, extended bool) *nav.LabelledTraverseCallback { //nolint:unparam // for future use
	ex := lo.Ternary(extended, "-EX", "")

	return &nav.LabelledTraverseCallback{
		Label: "test universal callback",
		Fn: func(item *nav.TraverseItem) error {
			depth := lo.TernaryF(extended,
				func() int { return item.Extension.Depth },
				func() int { return 9999 },
			)
			GinkgoWriter.Printf(
				"---> 🌊 UNIVERSAL//%v-CALLBACK%v: (depth:%v) '%v'\n", name, ex, depth, item.Path,
			)

			return nil
		},
	}
}

func foldersCallback(name string, extended bool) *nav.LabelledTraverseCallback {
	ex := lo.Ternary(extended, "-EX", "")

	return &nav.LabelledTraverseCallback{
		Label: "folders callback",
		Fn: func(item *nav.TraverseItem) error {
			depth := lo.TernaryF(extended,
				func() int { return item.Extension.Depth },
				func() int { return 9999 },
			)
			actualNoChildren := len(item.Children)
			GinkgoWriter.Printf(
				"---> ☀️ FOLDERS//%v-CALLBACK%v: (depth:%v, children:%v) '%v'\n",
				name, ex, depth, actualNoChildren, item.Path,
			)
			Expect(item.IsDir()).To(BeTrue())

			if extended {
				Expect(item.Extension).NotTo(BeNil(), helpers.Reason(item.Path))
			} else {
				Expect(item.Extension).To(BeNil(), helpers.Reason(item.Path))
			}

			return nil
		},
	}
}

func filesCallback(name string, extended bool) *nav.LabelledTraverseCallback {
	ex := lo.Ternary(extended, "-EX", "")

	return &nav.LabelledTraverseCallback{
		Label: "files callback",
		Fn: func(item *nav.TraverseItem) error {
			GinkgoWriter.Printf("---> 🌙 FILES//%v-CALLBACK%v: '%v'\n", name, ex, item.Path)
			Expect(item.IsDir()).To(BeFalse())

			if extended {
				Expect(item.Extension).NotTo(BeNil(), helpers.Reason(item.Path))
			}
			return nil
		},
	}
}

// === scope

func universalScopeCallback(name string) *nav.LabelledTraverseCallback {
	return &nav.LabelledTraverseCallback{
		Label: "test universal callback",
		Fn: func(item *nav.TraverseItem) error {
			GinkgoWriter.Printf("---> 🌠 UNIVERSAL//%v-CALLBACK-EX item-scope: (%v) '%v'\n",
				name, item.Extension.NodeScope, item.Extension.Name,
			)
			Expect(item.Extension).NotTo(BeNil(), helpers.Reason(item.Extension.Name))
			return nil
		},
	}
}

func foldersScopeCallback(name string) *nav.LabelledTraverseCallback {
	return &nav.LabelledTraverseCallback{
		Label: "test folders callback",
		Fn: func(item *nav.TraverseItem) error {
			GinkgoWriter.Printf("---> 🌟 FOLDERS//%v-CALLBACK-EX item-scope: (%v) '%v'\n",
				name, item.Extension.NodeScope, item.Extension.Name,
			)
			Expect(item.Extension).NotTo(BeNil(), helpers.Reason(item.Extension.Name))
			Expect(item.IsDir()).To(BeTrue())
			return nil
		},
	}
}

func filesScopeCallback(name string) *nav.LabelledTraverseCallback {
	return &nav.LabelledTraverseCallback{
		Label: "test files callback",
		Fn: func(item *nav.TraverseItem) error {
			GinkgoWriter.Printf("---> 🌬️ FILES//%v-CALLBACK-EX item-scope: (%v) '%v'\n",
				name, item.Extension.NodeScope, item.Extension.Name,
			)
			Expect(item.Extension).NotTo(BeNil(), helpers.Reason(item.Extension.Name))
			Expect(item.IsDir()).To(BeFalse())
			return nil
		},
	}
}

// === sort

func universalSortCallback(name string) *nav.LabelledTraverseCallback {
	return &nav.LabelledTraverseCallback{
		Label: "test universal callback",
		Fn: func(item *nav.TraverseItem) error {
			GinkgoWriter.Printf("---> 💚 UNIVERSAL//%v-SORT-CALLBACK-EX(scope:%v, depth:%v) '%v'\n",
				name, item.Extension.NodeScope, item.Extension.Depth, item.Extension.Name,
			)
			Expect(item.Extension).NotTo(BeNil(), helpers.Reason(item.Extension.Name))
			return nil
		},
	}
}

func foldersSortCallback(name string) *nav.LabelledTraverseCallback {
	return &nav.LabelledTraverseCallback{
		Label: "test folders sort callback",
		Fn: func(item *nav.TraverseItem) error {
			GinkgoWriter.Printf("---> 💜 FOLDERS//%v-SORT-CALLBACK-EX '%v'\n",
				name, item.Extension.Name,
			)
			Expect(item.Extension).NotTo(BeNil(), helpers.Reason(item.Extension.Name))
			Expect(item.IsDir()).To(BeTrue())
			return nil
		},
	}
}

func filesSortCallback(name string) *nav.LabelledTraverseCallback {
	return &nav.LabelledTraverseCallback{
		Label: "test files sort callback",
		Fn: func(item *nav.TraverseItem) error {
			GinkgoWriter.Printf("---> 💙 FILES//%v-SORT-CALLBACK-EX '%v'\n",
				name, item.Extension.Name,
			)
			Expect(item.Extension).NotTo(BeNil(), helpers.Reason(item.Extension.Name))
			Expect(item.IsDir()).To(BeFalse())
			return nil
		},
	}
}

func universalDepthCallback(name string, maxDepth int) *nav.LabelledTraverseCallback {
	return &nav.LabelledTraverseCallback{
		Label: "test universal depth callback",
		Fn: func(item *nav.TraverseItem) error {
			if item.Extension.Depth <= maxDepth {
				GinkgoWriter.Printf("---> 💚 UNIVERSAL//%v-SORT-CALLBACK-EX(scope:%v, depth:%v) '%v'\n",
					name, item.Extension.NodeScope, item.Extension.Depth, item.Extension.Name,
				)
			}
			Expect(item.Extension).NotTo(BeNil(), helpers.Reason(item.Extension.Name))
			return nil
		},
	}
}

func foldersCaseSensitiveCallback(first, second string) *nav.LabelledTraverseCallback {
	recording := make(recordingMap)

	return &nav.LabelledTraverseCallback{
		Label: "test folders case sensitive callback",
		Fn: func(item *nav.TraverseItem) error {
			recording[item.Path] = len(item.Children)

			GinkgoWriter.Printf("---> ☀️ CASE-SENSITIVE-CALLBACK: '%v'\n", item.Path)
			Expect(item.IsDir()).To(BeTrue())

			if strings.HasSuffix(item.Path, second) {
				GinkgoWriter.Printf("---> 💧 FIRST: '%v', 💧 SECOND: '%v'\n", first, second)

				paths := lo.Keys(recording)
				_, found := lo.Find(paths, func(s string) bool {
					return strings.HasSuffix(s, first)
				})

				Expect(found).To(BeTrue())
			}

			return nil
		},
	}
}

// === skip

func skipDirFolderCallback(skip, exclude string) *nav.LabelledTraverseCallback {
	return &nav.LabelledTraverseCallback{
		Label: "test skip folder callback",
		Fn: func(item *nav.TraverseItem) error {
			GinkgoWriter.Printf(
				"---> ♻️ ON-NAVIGATOR-SKIP-DIR-CALLBACK(skip:%v): '%v'\n", skip, item.Path,
			)

			Expect(strings.HasSuffix(item.Path, exclude)).To(BeFalse())

			return lo.Ternary(strings.HasSuffix(item.Path, skip),
				fs.SkipDir, nil,
			)
		},
	}
}

func skipAllFolderCallback(skip, exclude string) *nav.LabelledTraverseCallback {
	return &nav.LabelledTraverseCallback{
		Label: "test skipAll folder callback",
		Fn: func(item *nav.TraverseItem) error {
			GinkgoWriter.Printf(
				"---> ♻️ ON-NAVIGATOR-SKIP-ALL-CALLBACK(skip:%v): '%v'\n", skip, item.Path,
			)

			Expect(strings.HasSuffix(item.Path, exclude)).To(BeFalse())

			return lo.Ternary(strings.HasSuffix(item.Path, skip),
				fs.SkipAll, nil,
			)
		},
	}
}

func boostCallback(name string) *nav.LabelledTraverseCallback {
	return &nav.LabelledTraverseCallback{
		Label: "test boost callback",
		Fn: func(item *nav.TraverseItem) error {
			fmt.Printf("---> ⏩ ON-boost-CALLBACK(%v) '%v'\n", name, item.Path)

			return nil
		},
	}
}
func subscribes(subscription nav.TraverseSubscription, de fs.DirEntry) bool {
	isAnySubscription := (subscription == nav.SubscribeAny)

	files := (subscription == nav.SubscribeFiles) && (!de.IsDir())
	folders := ((subscription == nav.SubscribeFolders) || subscription == nav.SubscribeFoldersWithFiles) && (de.IsDir())

	return isAnySubscription || files || folders
}

// === errors

type errorTE struct {
	naviTE
}

func readDirFakeError(_ string) ([]fs.DirEntry, error) {
	entries := []fs.DirEntry{}
	path := "/foo/bar"
	reason := errors.New("access denied")
	err := i18n.NewFailedToReadDirectoryContentsError(path, reason)

	return entries, err
}

func readDirFakeErrorAt(name string) func(dirname string) ([]fs.DirEntry, error) {
	return func(dirname string) ([]fs.DirEntry, error) {
		if strings.HasSuffix(dirname, name) {
			return readDirFakeError(dirname)
		}

		return nav.ReadEntriesHookFn(dirname)
	}
}

func errorCallback(name string, extended, hasError bool) *nav.LabelledTraverseCallback {
	ex := lo.Ternary(extended, "-EX", "")

	return &nav.LabelledTraverseCallback{
		Label: "test error callback",
		Fn: func(item *nav.TraverseItem) error {
			GinkgoWriter.Printf("---> 🔥 %v-CALLBACK%v: '%v'\n", name, ex, item.Path)

			Expect(item.Extension).NotTo(BeNil(), helpers.Reason(item.Path))
			if hasError {
				Expect(item.Error).ToNot(BeNil())
			}
			return item.Error
		},
	}
}
