package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/snivilised/extendio/i18n"
	"github.com/snivilised/extendio/xfs/nav"
)

func Path(parent, relative string) string {
	segments := strings.Split(relative, "/")
	return filepath.Join(append([]string{parent}, segments...)...)
}

func Normalise(p string) string {
	return strings.ReplaceAll(p, "/", string(filepath.Separator))
}

func Reason(name string) string {
	return fmt.Sprintf("❌ for item named: '%v'", name)
}

func BecauseQuantity(name string, expected, actual int) string {
	return fmt.Sprintf("❌ incorrect no of items for: '%v', expected: '%v', actual: '%v'",
		name, expected, actual,
	)
}

func JoinCwd(segments ...string) string {
	if current, err := os.Getwd(); err == nil {
		parent, _ := filepath.Split(current)
		grand := filepath.Dir(parent)
		great := filepath.Dir(grand)
		all := append([]string{great}, segments...)

		return filepath.Join(all...)
	}

	panic("could not get root path")
}

func Root() string {
	if current, err := os.Getwd(); err == nil {
		return current
	}

	panic("could not get root path")
}

func Repo(relative string) string {
	_, filename, _, _ := runtime.Caller(0) //nolint:dogsled // use of 3 _ is out of our control

	return Path(filepath.Dir(filename), relative)
}

func Log() string {
	if current, err := os.Getwd(); err == nil {
		parent, _ := filepath.Split(current)
		grand := filepath.Dir(parent)
		great := filepath.Dir(grand)

		return filepath.Join(great, "Test", "test.log")
	}

	panic("could not get root path")
}

type CustomFilter struct {
	Name  string
	Value string
}

func (f *CustomFilter) Description() string {
	return f.Name
}

func (f *CustomFilter) Validate() {}

func (f *CustomFilter) Source() string {
	return f.Value
}

func (f *CustomFilter) IsMatch(item *nav.TraverseItem) bool {
	return f.Value == item.Extension.Name
}

func (f *CustomFilter) IsApplicable(_ *nav.TraverseItem) bool {
	return true
}

func (f *CustomFilter) Scope() nav.FilterScopeBiEnum {
	return nav.ScopeAllEn
}

type DummyCreator struct {
	Invoked bool
}

func (dc *DummyCreator) Create(_ *i18n.LanguageInfo, _ string) (*i18n.Localizer, error) {
	dc.Invoked = true

	return &i18n.Localizer{}, nil
}
