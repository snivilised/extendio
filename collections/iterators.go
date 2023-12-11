package collections

// IteratorCtrl represents a narrow view of the iterator that exposes
// just the method, Valid which indicates when all items in the
// sequence have been consumed. The purpose of this is to allow a
// client abstraction, which implements the looping logic to provide
// a condition which halts its iteration. There are 2 scenarios:
// 1) the client wants to iterate the entire sequence; in this case
// the client just needs to continue the iteration until Valid
// returns false.
// 2) the client only wants to iterate the sequence until a certain
// other condition arises; in this case, the client combines the
// result of Valid() with another predicate within the for statement.
type IteratorCtrl[T any] interface {
	Valid() bool
}

// Iterator represents an iterator over a slice. The underlying slice may
// be empty. When created, the iterator does not point to a valid
// slice entry. To begin iteration, the client invokes Start. At any stage
// after the iteration has started, the iterator, points to the current item
// in the sequence. The client can query the validity of the current item
// using the Valid method. To obtain successive elements, the client invokes
// the Next method and this can be continued until Valid returns false. Once
// Valid returns false, the client should no longer call Next. Doing so in
// this scenario, returns the zero value for the value T.
// It should also be noted that the iterator does not in itself implement
// the loop operation, it merely provides the logic to enable the client
// to implement the looping operation.
type Iterator[T any] interface {
	IteratorCtrl[T]

	// Start returns the first element of the sequence and moves the current
	// position to the next item.
	Start() T

	// Next moves the iterator to the next item in the sequence then returns
	// that item. If the iterator is already at the end, then the zero
	// value of T is returned. However, when used properly, this situation
	// should never occur, as Valid would indicate that iteration with
	// Next should no longer occur.
	Next() T

	// Reset is designed to be used in high frequency applications. The client
	// can reuse this iterator for a new collection rather that having to throw this
	// instance away and create a new one. This helps to reduce the number of
	// allocations in a high frequency application.
	Reset(entries []T)
}

// ForwardIt creates a forward iterator over a non empty slice. If the provided
// slice is empty, then a nil iterator is returned.
//
// The zero value represents the value that is returned if the Next method on the
// iterator is incorrectly invoked after Valid has returned false.
// If the collection contains interfaces, or pointers just pass in nil as the
// zero value.
//
// If the collection contains scalars, pass in 0 cast to the appropriate type,
// eg int32(0). It doesn't matter if 0 is a valid value in the collection,
// because this value is only ever return in an invalid scenario, ie, calling
// Next after Valid has returned false. This is preferable than generating a
// panic. If the collection contains structs, then pass in an empty struct
// as the nil value.
func ForwardIt[T any](elements []T, zero T) Iterator[T] {
	// 📚 NB: it is not possible to obtain the type of a generic parameter at runtime
	// using reflection. Generics in Go are a compile-time feature, and type information
	// is generally not available at runtime due to the language's design principles.
	// This is why we need the client to pass in the zero value manually.
	//
	return &forwardIterator[T]{
		baseIterator: baseIterator[T]{
			zero:      zero,
			container: elements,
			current:   -1,
		},
	}
}

// ReverseIt creates a reverse iterator over a non empty slice. If the provided
// slice is empty, then a nil iterator is returned. (NB: please remember to check
// for a nil interface correctly; see the helper function IsNil in utils).
func ReverseIt[T any](elements []T, zero T) Iterator[T] {
	return &reverseIterator[T]{
		baseIterator: baseIterator[T]{
			zero:      zero,
			container: elements,
			current:   len(elements),
		},
	}
}

type baseIterator[T any] struct {
	zero      T
	container []T
	current   int
}

// Valid returns true if the current position of the iterator points
// to a valid entry in the sequence. When the iterator reaches the
// end of the sequence, Valid returns false.
func (i *baseIterator[T]) Valid() bool {
	return i.current >= 0 && i.current < len(i.container)
}

// forwardIterator navigates the sequence from the start (index 0) to the
// end (index len-1)
type forwardIterator[T any] struct {
	baseIterator[T]
}

// Start returns the first element of the sequence and moves the current
// position to the next item.
func (i *forwardIterator[T]) Start() T {
	if len(i.container) == 0 {
		return i.zero
	}

	const initial = 0
	i.current = initial

	return i.container[i.current]
}

// Next moves the iterator to the next item in the sequence then returns
// that item. If the iterator is already at the end, then the zero
// value of T is returned. However, when used properly, this situation
// should never occur, as Valid would indicate that iteration with
// Next should no longer occur.
func (i *forwardIterator[T]) Next() T {
	i.current++
	if !i.Valid() {
		return i.zero
	}

	return i.container[i.current]
}

// Reset is designed to be used in high frequency applications. The client
// can reuse this iterator for a new collection rather that having throw this
// instance away and create a new one. This helps to reduce the number of
// allocations in a high frequency application.
func (i *forwardIterator[T]) Reset(entries []T) {
	i.container = entries
	i.current = -1
}

type reverseIterator[T any] struct {
	baseIterator[T]
}

// Start returns the last element of the sequence and moves the current
// position to the next item.
func (i *reverseIterator[T]) Start() T {
	if len(i.container) == 0 {
		return i.zero
	}

	const offset = 1
	i.current = len(i.container) - offset

	return i.container[i.current]
}

// Next moves the iterator to the next item in the sequence then returns
// that item. If the iterator is already at the end, then the zero
// value of T is returned. However, when used properly, this situation
// should never occur, as Valid would indicate that iteration with
// Next should no longer occur.
func (i *reverseIterator[T]) Next() T {
	i.current--
	if !i.Valid() {
		return i.zero
	}

	return i.container[i.current]
}

// Reset is designed to be used in high frequency applications. The client
// can reuse this iterator for a new collection rather that having throw this
// instance away and create a new one. This helps to reduce the number of
// allocations in a high frequency application.
func (i *reverseIterator[T]) Reset(entries []T) {
	i.container = entries
	i.current = len(i.container)
}
