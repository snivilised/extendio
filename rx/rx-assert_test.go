//nolint:gocritic // commented out code
package rx_test

import (
	"context"

	"github.com/snivilised/extendio/rx"
)

// AssertPredicate is a custom predicate based on the items.
type AssertPredicate[T any] func(items []T) error

// RxAssert[T] lists the Observable assertions.
type RxAssert[T any] interface {
	apply(*rxAssert[T])
	itemsToBeChecked() (bool, []T)
	itemsNoOrderedToBeChecked() (bool, []T)
	noItemsToBeChecked() bool
	someItemsToBeChecked() bool
	raisedErrorToBeChecked() (bool, error)
	raisedErrorsToBeChecked() (bool, []error)
	raisedAnErrorToBeChecked() (bool, error)
	notRaisedErrorToBeChecked() bool
	itemToBeChecked() (bool, T)
	noItemToBeChecked() (bool, T)
	customPredicatesToBeChecked() (bool, []AssertPredicate[T])
}

type rxAssert[T any] struct {
	f                       func(*rxAssert[T])
	checkHasItems           bool
	checkHasNoItems         bool
	checkHasSomeItems       bool
	items                   []T
	checkHasItemsNoOrder    bool
	itemsNoOrder            []T
	checkHasRaisedError     bool
	err                     error
	checkHasRaisedErrors    bool
	errs                    []error
	checkHasRaisedAnError   bool
	checkHasNotRaisedError  bool
	checkHasItem            bool
	item                    T
	checkHasNoItem          bool
	checkHasCustomPredicate bool
	customPredicates        []AssertPredicate[T]
}

func (ass *rxAssert[T]) apply(do *rxAssert[T]) {
	ass.f(do)
}

func (ass *rxAssert[T]) itemsToBeChecked() (res bool, items []T) {
	return ass.checkHasItems, ass.items
}

func (ass *rxAssert[T]) itemsNoOrderedToBeChecked() (res bool, items []T) {
	return ass.checkHasItemsNoOrder, ass.itemsNoOrder
}

func (ass *rxAssert[T]) noItemsToBeChecked() bool {
	return ass.checkHasNoItems
}

func (ass *rxAssert[T]) someItemsToBeChecked() bool {
	return ass.checkHasSomeItems
}

func (ass *rxAssert[T]) raisedErrorToBeChecked() (bool, error) {
	return ass.checkHasRaisedError, ass.err
}

func (ass *rxAssert[T]) raisedErrorsToBeChecked() (bool, []error) {
	return ass.checkHasRaisedErrors, ass.errs
}

func (ass *rxAssert[T]) raisedAnErrorToBeChecked() (bool, error) {
	return ass.checkHasRaisedAnError, ass.err
}

func (ass *rxAssert[T]) notRaisedErrorToBeChecked() bool {
	return ass.checkHasNotRaisedError
}

func (ass *rxAssert[T]) itemToBeChecked() (res bool, item T) {
	return ass.checkHasItem, ass.item
}

func (ass *rxAssert[T]) noItemToBeChecked() (res bool, item T) {
	return ass.checkHasNoItem, ass.item
}

func (ass *rxAssert[T]) customPredicatesToBeChecked() (bool, []AssertPredicate[T]) {
	return ass.checkHasCustomPredicate, ass.customPredicates
}

func newAssertion[T any](f func(*rxAssert[T])) *rxAssert[T] {
	return &rxAssert[T]{
		f: f,
	}
}

// HasItems checks that the observable produces the corresponding items.
func HasItems[T any](items ...T) RxAssert[T] {
	return newAssertion(func(a *rxAssert[T]) {
		a.checkHasItems = true
		a.items = items
	})
}

// HasItem checks if a single or optional single has a specific item.
func HasItem[T any](i T) RxAssert[T] {
	return newAssertion(func(a *rxAssert[T]) {
		a.checkHasItem = true
		a.item = i
	})
}

// HasItemsNoOrder checks that an observable produces the corresponding
// items regardless of the order.
func HasItemsNoOrder[T any](items ...T) RxAssert[T] {
	return newAssertion(func(a *rxAssert[T]) {
		a.checkHasItemsNoOrder = true
		a.itemsNoOrder = items
	})
}

// IsNotEmpty checks that the observable produces some items.
func IsNotEmpty[T any]() RxAssert[T] {
	return newAssertion(func(a *rxAssert[T]) {
		a.checkHasSomeItems = true
	})
}

// IsEmpty checks that the observable has not produce any item.
func IsEmpty[T any]() RxAssert[T] {
	return newAssertion(func(a *rxAssert[T]) {
		a.checkHasNoItems = true
	})
}

// HasError checks that the observable has produce a specific error.
func HasError[T any](err error) RxAssert[T] {
	return newAssertion(func(a *rxAssert[T]) {
		a.checkHasRaisedError = true
		a.err = err
	})
}

// HasAnError checks that the observable has produce an error.
func HasAnError[T any]() RxAssert[T] {
	return newAssertion(func(a *rxAssert[T]) {
		a.checkHasRaisedAnError = true
	})
}

// HasErrors checks that the observable has produce a set of errors.
func HasErrors[T any](errs ...error) RxAssert[T] {
	return newAssertion(func(a *rxAssert[T]) {
		a.checkHasRaisedErrors = true
		a.errs = errs
	})
}

// HasNoError checks that the observable has not raised any error.
func HasNoError[T any]() RxAssert[T] {
	return newAssertion(func(a *rxAssert[T]) {
		a.checkHasRaisedError = true
	})
}

// CustomPredicate checks a custom predicate.
func CustomPredicate[T any](predicate AssertPredicate[T]) RxAssert[T] {
	return newAssertion(func(a *rxAssert[T]) {
		if !a.checkHasCustomPredicate {
			a.checkHasCustomPredicate = true
			a.customPredicates = make([]AssertPredicate[T], 0)
		}
		a.customPredicates = append(a.customPredicates, predicate)
	})
}

func parseAssertions[T any](assertions ...RxAssert[T]) RxAssert[T] {
	ass := new(rxAssert[T])
	for _, assertion := range assertions {
		assertion.apply(ass)
	}

	return ass
}

// Assert asserts the result of an iterable against a list of assertions.
func Assert[T any](_ context.Context, _ rx.Iterable[T], _ ...RxAssert[T]) {
	// 	ass := parseAssertions(assertions...)

	// 	got := make([]interface{}, 0)
	// 	errs := make([]error, 0)

	// 	observe := iterable.Observe()
	// loop:
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			break loop
	// 		case item, ok := <-observe:
	// 			if !ok {
	// 				break loop
	// 			}
	// 			if item.Error() {
	// 				errs = append(errs, item.E)
	// 			} else {
	// 				got = append(got, item.V)
	// 			}
	// 		}
	// 	}

	// 	if checked, predicates := ass.customPredicatesToBeChecked(); checked {
	// 		for _, predicate := range predicates {
	// 			err := predicate(got)
	// 			if err != nil {
	// 				Fail(err.Error())
	// 			}
	// 		}
	// 	}
	// 	if checkHasItems, expectedItems := ass.itemsToBeChecked(); checkHasItems {
	// 		Expect(1).To(Equal(1)) // REMOVE ME
	// 		assert.Equal(t, expectedItems, got)
	// 	}
	// 	if checkHasItemsNoOrder, itemsNoOrder := ass.itemsNoOrderedToBeChecked(); checkHasItemsNoOrder {
	// 		m := make(map[interface{}]interface{})
	// 		for _, v := range itemsNoOrder {
	// 			m[v] = nil
	// 		}

	//		for _, v := range got {
	//			delete(m, v)
	//		}
	//		if len(m) != 0 {
	//			assert.Fail(t, "missing elements", "%v", got)
	//		}
	//	}
	//
	//	if checkHasItem, value := ass.itemToBeChecked(); checkHasItem {
	//		length := len(got)
	//		if length != 1 {
	//			assert.FailNow(t, "wrong number of items", "expected 1, got %d", length)
	//		}
	//		assert.Equal(t, value, got[0])
	//	}
	//
	//	if ass.noItemsToBeChecked() {
	//		assert.Equal(t, 0, len(got))
	//	}
	//
	//	if ass.someItemsToBeChecked() {
	//		assert.NotEqual(t, 0, len(got))
	//	}
	//
	//	if checkHasRaisedError, expectedError := ass.raisedErrorToBeChecked(); checkHasRaisedError {
	//		if expectedError == nil {
	//			assert.Equal(t, 0, len(errs))
	//		} else {
	//			if len(errs) == 0 {
	//				assert.FailNow(t, "no error raised", "expected %v", expectedError)
	//			}
	//			assert.Equal(t, expectedError, errs[0])
	//		}
	//	}
	//
	//	if checkHasRaisedErrors, expectedErrors := ass.raisedErrorsToBeChecked(); checkHasRaisedErrors {
	//		assert.Equal(t, expectedErrors, errs)
	//	}
	//
	//	if checkHasRaisedAnError, expectedError := ass.raisedAnErrorToBeChecked(); checkHasRaisedAnError {
	//		assert.Nil(t, expectedError)
	//	}
	//
	//	if ass.notRaisedErrorToBeChecked() {
	//		assert.Equal(t, 0, len(errs))
	//	}
}
