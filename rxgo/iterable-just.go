package rxgo

type justIterable[T any] struct {
	items []T
	opts  []Option[T]
}

func newJustIterable[T any](items ...T) func(opts ...Option[T]) Iterable[T] {
	return func(opts ...Option[T]) Iterable[T] {
		return &justIterable[T]{
			items: items,
			opts:  opts,
		}
	}
}

func (i *justIterable[T]) Observe(opts ...Option[T]) <-chan Item[T] {
	option := parseOptions(append(i.opts, opts...)...)
	next := option.buildChannel()

	// this is weird, in the original, i.items is passed into SendItems without being
	// spread! Because the original was using ...interface{} anything satisfies this
	// signature, including a slice. I suspect this was a mistake, because it doesn't
	// appear to match the intention of SendItems and is only exposed here because of
	// the type safety of generics over the use of interface{}. I wonder if the original
	// unit test(s) of justIterable were deficient.
	//
	go SendItems(option.buildContext(emptyContext), next, CloseChannel, i.items...)

	return next
}
