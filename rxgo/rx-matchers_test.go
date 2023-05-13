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

type RxMatcher[T any] struct {
	Name string
	ass  *rxAssert[T]
}

func (m *RxMatcher[T]) Match(actual interface{}) (bool, error) {
	matcher, ok := actual.(RxObservable[T])
	if !ok {
		return false, fmt.Errorf("❌ expected an %T (%T)", matcher, actual)
	}

	return Assert[T](matcher.context, matcher.observable, m.ass)
}

func (m *RxMatcher[T]) FailureMessage(actual interface{}) string {
	// TODO: improve this message, using actual
	//
	return fmt.Sprintf("🔥 Expected\n\t\n%v to match\n\t\n", m.Name)
}

func (m *RxMatcher[T]) NegatedFailureMessage(actual interface{}) string {
	// TODO: improve this message, using actual
	//
	return fmt.Sprintf("🔥 Expected\n\t\n%v NOT to match\n\t\n", m.Name)
}

type HasItemsMatcher[T any] struct {
	RxMatcher[T]
}

func MatchHasItems[T any](expected ...T) GomegaMatcher {
	if ass, ok := HasItems(expected...).(*rxAssert[T]); ok {
		return &HasItemsMatcher[T]{
			RxMatcher: RxMatcher[T]{
				Name: "HasItems",
				ass:  ass,
			},
		}
	}

	panic("invalid expected in MatchHasItems test")
}

type HasNoErrorMatcher[T any] struct {
	RxMatcher[T]
}

func MatchHasNoError[T any](expected ...T) GomegaMatcher {
	if ass, ok := HasNoError[T]().(*rxAssert[T]); ok {
		return &HasNoErrorMatcher[T]{
			RxMatcher: RxMatcher[T]{
				Name: "HasNoError",
				ass:  ass,
			},
		}
	}

	panic("invalid expected in MatchHasNoError test")
}
