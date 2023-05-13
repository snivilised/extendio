This implementation is based upon the original version of rxgo at https://github.com/ReactiveX/RxGo, based upon generics rather than relying on reflection and interface{}. It is intended to replace this local version with the official version when and if it becomes available and for this to happen with minimal issues, the design here has to mirror the original as closely as possible.


- Issues:

* justIterable.Observe: whilst implementing this, it was discovered there was a potential
bug. Since I converted this to a generic, it exposed the incorrect invocation of SendItems.
The original code was this:

func (i *justIterable) Observe(opts ...Option) <-chan Item {
	option := parseOptions(append(i.opts, opts...)...)
	next := option.buildChannel()

	go SendItems(option.buildContext(emptyContext), next, CloseChannel, i.items)
	return next
}

but the signature of SendItems is:

func SendItems(ctx context.Context, ch chan<- Item, strategy CloseChannelStrategy, items ...interface{})

This tells me that items is supposed to be variadic, not a single entity. Invoking SendItems with
i.items does not match the intention of SendItems, but because the parameter type of items is
interface{}, anything can be passed in; its not type-safe, hence it compiles, but I suspect this
is a subtle bug.
