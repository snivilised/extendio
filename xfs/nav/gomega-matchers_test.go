package nav_test

import (
	"fmt"

	. "github.com/onsi/gomega/types"
	"github.com/snivilised/extendio/xfs/nav"
)

// === MatchCurrentRegexFilter ===
//

type IsCurrentRegexMatchMatcher struct {
	filter interface{}
}

func MatchCurrentRegexFilter(expected interface{}) GomegaMatcher {
	return &IsCurrentRegexMatchMatcher{
		filter: expected,
	}
}

func (m *IsCurrentRegexMatchMatcher) Match(actual interface{}) (bool, error) {
	item, itemOk := actual.(*nav.TraverseItem)
	if !itemOk {
		return false, fmt.Errorf("matcher expected a *TraverseItem (%T)", item)
	}

	filter, filterOk := m.filter.(*nav.RegexFilter)
	if !filterOk {
		return false, fmt.Errorf("matcher expected a *RegexFilter (%T)", filter)
	}

	return filter.IsMatch(item), nil
}

func (m *IsCurrentRegexMatchMatcher) FailureMessage(actual interface{}) string {
	item, _ := actual.(*nav.TraverseItem)
	filter, _ := m.filter.(*nav.RegexFilter)

	return fmt.Sprintf("🔥 Expected\n\t%v\nto match regex\n\t%v\n", item.Extension.Name, filter.Source())
}

func (m *IsCurrentRegexMatchMatcher) NegatedFailureMessage(actual interface{}) string {
	item, _ := actual.(*nav.TraverseItem)
	filter, _ := m.filter.(*nav.RegexFilter)

	return fmt.Sprintf("🔥 Expected\n\t%v\nNOT to match regex\n\t%v\n", item.Extension.Name, filter.Source())
}

// === MatchCurrentGlobFilter ===
//

type IsCurrentGlobMatchMatcher struct {
	filter interface{}
}

func MatchCurrentGlobFilter(expected interface{}) GomegaMatcher {
	return &IsCurrentGlobMatchMatcher{
		filter: expected,
	}
}

func (m *IsCurrentGlobMatchMatcher) Match(actual interface{}) (bool, error) {
	item, itemOk := actual.(*nav.TraverseItem)
	if !itemOk {
		return false, fmt.Errorf("matcher expected a *TraverseItem (%T)", item)
	}

	filter, filterOk := m.filter.(*nav.GlobFilter)
	if !filterOk {
		return false, fmt.Errorf("matcher expected a *GlobFilter (%T)", filter)
	}

	return filter.IsMatch(item), nil
}

func (m *IsCurrentGlobMatchMatcher) FailureMessage(actual interface{}) string {
	item, _ := actual.(*nav.TraverseItem)
	filter, _ := m.filter.(*nav.GlobFilter)

	return fmt.Sprintf("🔥 Expected\n\t%v\nto match glob\n\t%v\n", item.Extension.Name, filter.Source())
}

func (m *IsCurrentGlobMatchMatcher) NegatedFailureMessage(actual interface{}) string {
	item, _ := actual.(*nav.TraverseItem)
	filter, _ := m.filter.(*nav.GlobFilter)

	return fmt.Sprintf("🔥 Expected\n\t%v\nNOT to match glob\n\t%v\n", item.Extension.Name, filter.Source())
}
