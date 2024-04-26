package collections_test

import (
	"fmt"

	. "github.com/onsi/gomega" //nolint:revive // gomega ok

	"github.com/onsi/gomega/types"
	"github.com/snivilised/extendio/collections"
)

// === HaveSize

type HaveSizeMatcher struct {
	size uint
}

func HaveSize(size uint) types.GomegaMatcher {
	return &HaveSizeMatcher{
		size: size,
	}
}

func (m *HaveSizeMatcher) Match(actual interface{}) (bool, error) {
	stack, ok := actual.(*collections.Stack[string])

	if !ok {
		return false, fmt.Errorf("matcher expected a *collections.Stack[T] (%T)", stack)
	}

	return stack.Size() == m.size, nil
}

func (m *HaveSizeMatcher) FailureMessage(_ interface{}) string {
	return fmt.Sprintf("🔥 Expected stack to have size: %v\n", m.size)
}

func (m *HaveSizeMatcher) NegatedFailureMessage(_ interface{}) string {
	return fmt.Sprintf("🔥 Expected stack NOT to have size: %v\n", m.size)
}

// === HaveCurrent

type HaveCurrentMatcher struct {
	current string
}

func HaveCurrent(current string) types.GomegaMatcher {
	return &HaveCurrentMatcher{
		current: current,
	}
}

func (m *HaveCurrentMatcher) Match(actual interface{}) (bool, error) {
	stack, ok := actual.(*collections.Stack[string])

	if !ok {
		return false, fmt.Errorf("matcher expected a *collections.Stack[T] (%T)", stack)
	}

	current, _ := stack.Current()

	return current == m.current, nil
}

func (m *HaveCurrentMatcher) FailureMessage(_ interface{}) string {
	return fmt.Sprintf("🔥 Expected stack to have current value of: %v\n", m.current)
}

func (m *HaveCurrentMatcher) NegatedFailureMessage(_ interface{}) string {
	return fmt.Sprintf("🔥 Expected stack NOT to have current value of: %v\n", m.current)
}

// === BeInCorrectState

func BeInCorrectState(size uint, current string) types.GomegaMatcher {
	return And(
		HaveSize(size),
		HaveCurrent(current),
	)
}

type HavePoppedMatcher struct {
	size       uint
	actualItem string
}

type WithExpectedPop struct {
	stack  *collections.Stack[string]
	popped string
}

func HavePopped(size uint, actual string) types.GomegaMatcher {
	return &HavePoppedMatcher{
		size:       size,
		actualItem: actual,
	}
}

func (m *HavePoppedMatcher) Match(actual interface{}) (bool, error) {
	expectation, ok := actual.(*WithExpectedPop)

	if !ok {
		return false, fmt.Errorf("matcher expected a *ExpectedPop (%T)", expectation)
	}

	result := expectation.stack.Size() == m.size && m.actualItem == expectation.popped

	return result, nil
}

func (m *HavePoppedMatcher) FailureMessage(_ interface{}) string {
	return fmt.Sprintf("🔥 Expected stack to\n\thave size: %v\n\tand popped item: %v\n",
		m.size, m.actualItem,
	)
}

func (m *HavePoppedMatcher) NegatedFailureMessage(_ interface{}) string {
	return fmt.Sprintf("🔥 Expected stack NOT to\n\thave size: %v\n\tand popped item: %v\n",
		m.size, m.actualItem,
	)
}
