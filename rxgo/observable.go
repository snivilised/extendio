package rxgo

import "context"

type Observable[T any] interface {
	Iterable[T]
}

type ObservableImpl[T any] struct {
	parent   context.Context
	iterable Iterable[T]
}

func (o *ObservableImpl[T]) Observe(opts ...Option[T]) <-chan Item[T] {
	return o.iterable.Observe(opts...)
}
