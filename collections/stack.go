package collections

// NewStack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		content: []T{},
	}
}

func NewStackWith[T any](with []T) *Stack[T] {
	return &Stack[T]{
		content: with,
	}
}

// Stack
type Stack[T any] struct {
	content []T
}

// Push
func (s *Stack[T]) Push(item T) {
	s.content = append(s.content, item)
}

// Pop
func (s *Stack[T]) Pop() (T, error) {
	if s.IsEmpty() {
		var zero T
		return zero, NewStackIsEmptyNativeError()
	}

	item := s.pop()

	return item, nil
}

// MustPop
func (s *Stack[T]) MustPop() T {
	if s.IsEmpty() {
		panic(NewStackIsEmptyNativeError())
	}

	return s.pop()
}

// Current
func (s *Stack[T]) Current() (T, error) {
	if s.IsEmpty() {
		var zero T
		return zero, NewStackIsEmptyNativeError()
	}

	return s.content[s.top()], nil
}

// Size
func (s *Stack[T]) Size() uint {
	return uint(len(s.content))
}

// IsEmpty
func (s *Stack[T]) IsEmpty() bool {
	return len(s.content) == 0
}

// Content
func (s *Stack[T]) Content() []T {
	return s.content
}

func (s *Stack[T]) top() int {
	return len(s.content) - 1
}

func (s *Stack[T]) pop() T {
	t := s.top()
	item := s.content[t]
	s.content = s.content[:t]

	return item
}
