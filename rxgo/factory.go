package rxgo

// FromChannel creates a cold observable from a channel.
func FromChannel[T any](next <-chan Item[T], opts ...Option[T]) Observable[T] {
	option := parseOptions(opts...)
	ctx := option.buildContext(emptyContext)

	return &ObservableImpl[T]{
		parent:   ctx,
		iterable: newChannelIterable(next, opts...),
	}
}

// Just creates an Observable with the provided items.
func Just[T any](items ...T) func(opts ...Option[T]) Observable[T] {
	return func(opts ...Option[T]) Observable[T] {
		return &ObservableImpl[T]{
			iterable: newJustIterable(items...)(opts...),
		}
	}
}
