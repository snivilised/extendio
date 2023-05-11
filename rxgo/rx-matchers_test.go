//nolint:revive // foo bar baz
package rxgo_test

import (
	"context"
	"fmt"

	. "github.com/onsi/gomega/types"
	"github.com/snivilised/extendio/rxgo"
)

// RxObservable
type RxObservable[T any] struct {
	// TODO: think of a better name for RxObservable
	//
	context    context.Context
	observable rxgo.Observable[T]
}

type RxMatcher struct {
	Name string
}

func (m *RxMatcher) FailureMessage(actual interface{}) string {
	// TODO: improve this message, using actual
	//
	return fmt.Sprintf("🔥 Expected\n\t\n%v to match\n\t\n", m.Name)
}

func (m *RxMatcher) NegatedFailureMessage(actual interface{}) string {
	// TODO: improve this message, using actual
	//
	return fmt.Sprintf("🔥 Expected\n\t\n%v NOT to match\n\t\n", m.Name)
}

type HasItemsMatcher[T any] struct {
	RxMatcher
	ass *rxAssert[T]
}

func MatchHasItems[T any](expected ...T) GomegaMatcher {
	if ass, ok := HasItems(expected...).(*rxAssert[T]); ok {
		return &HasItemsMatcher[T]{
			RxMatcher: RxMatcher{
				Name: "HasItems",
			},
			ass: ass,
		}
	}

	// TODO: should this be a panic instead?
	//
	return nil
}

func (m *HasItemsMatcher[T]) Match(actual interface{}) (bool, error) {
	matcher, ok := actual.(RxObservable[T])
	if !ok {
		return false, fmt.Errorf("❌ expected an RxObservable (%T)", matcher)
	}

	return Assert[T](matcher.context, matcher.observable, m.ass)
}
