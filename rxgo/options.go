package rxgo

import (
	"context"

	"github.com/teivah/onecontext"
)

var emptyContext context.Context

type Option[T any] interface {
	apply(*funcOption[T])
	buildChannel() chan Item[T]
	buildContext(parent context.Context) context.Context
	isConnectable() bool
	isConnectOperation() bool
}

type funcOption[T any] struct {
	f                    func(*funcOption[T])
	isBuffer             bool
	buffer               int
	ctx                  context.Context
	observation          ObservationStrategy
	pool                 int
	backPressureStrategy BackPressureStrategy
	onErrorStrategy      OnErrorStrategy
	propagate            bool
	connectable          bool
	connectOperation     bool
	serialized           func(T) int
}

func parseOptions[T any](opts ...Option[T]) Option[T] {
	o := new(funcOption[T])
	for _, opt := range opts {
		opt.apply(o)
	}

	return o
}

func (fdo *funcOption[T]) apply(do *funcOption[T]) {
	fdo.f(do)
}

func (fdo *funcOption[T]) buildChannel() chan Item[T] {
	if fdo.isBuffer {
		return make(chan Item[T], fdo.buffer)
	}

	return make(chan Item[T])
}

func (fdo *funcOption[T]) buildContext(parent context.Context) context.Context {
	if fdo.ctx != nil && parent != nil {
		ctx, _ := onecontext.Merge(fdo.ctx, parent)
		return ctx
	}

	if fdo.ctx != nil {
		return fdo.ctx
	}

	if parent != nil {
		return parent
	}

	return context.Background()
}

func (fdo *funcOption[T]) isConnectable() bool {
	return fdo.connectable
}

func (fdo *funcOption[T]) isConnectOperation() bool {
	return fdo.connectOperation
}
