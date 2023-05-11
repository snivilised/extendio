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
