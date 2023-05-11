package rxgo

import (
	"context"
	"fmt"

	"github.com/samber/lo"
)

// CloseChannelStrategy indicates a strategy on whether to close a channel.
type CloseChannelStrategy uint32

const (
	// LeaveChannelOpen indicates to leave the channel open after completion.
	LeaveChannelOpen CloseChannelStrategy = iota
	// CloseChannel indicates to close the channel open after completion.
	CloseChannel
)

type Item[T any] struct {
	V T
	E error
}

// Of creates an item from a value.
func Of[T any](i T) Item[T] {
	return Item[T]{V: i}
}

// Error checks if an item is an error.
func (i Item[T]) Error() bool {
	return i.E != nil
}

// SendItemsV is an utility function that send a list of items and indicate
// a strategy on whether to close the channel once the function completes.
func SendItems[T any](ctx context.Context, ch chan<- Item[T], strategy CloseChannelStrategy, items ...T) {
	if strategy == CloseChannel {
		defer close(ch)
	}

	send(ctx, ch, items...)
}

func send[T any](ctx context.Context, ch chan<- Item[T], items ...T) {
	// This is only the basic implementation. It does not yet support sending
	// a slice or a channel. Support for these don't need to be added unless
	// explicitly required.
	//
	for _, currentItem := range items {
		_ = Of(currentItem).SendContext(ctx, ch)
	}
}

// SendItemsV (verbose) is an utility function that send a list of items and indicate
// a strategy on whether to close the channel once the function completes.
// This will eventually be removed as its only required for debugging purposes
func SendItemsV[T any](ctx context.Context, ch chan<- Item[T], strategy CloseChannelStrategy, items ...T) {
	if strategy == CloseChannel {
		defer close(ch)
	}

	sendV(ctx, ch, items...)
}

// This will eventually be removed as its only required for debugging purposes
func sendV[T any](ctx context.Context, ch chan<- Item[T], items ...T) {
	fmt.Println("")

	// This is only the basic implementation. It does not yet support sending
	// a slice or a channel. Support for these don't need to be added unless
	// explicitly required.
	//
	for _, currentItem := range items {
		result := Of(currentItem).SendContext(ctx, ch)
		indicator := lo.Ternary(result, "✔️", "❌")

		fmt.Printf("===> 💘 sending item: '%v' (%v)\n", currentItem, indicator)
	}
}

// SendContext sends an item and blocks until it is sent or a context canceled.
// It returns a boolean to indicate whether the item was sent.
func (i Item[T]) SendContext(ctx context.Context, ch chan<- Item[T]) bool {
	select {
	case <-ctx.Done(): // Context's done channel has the highest priority
		return false
	default:
		select {
		case <-ctx.Done():
			return false
		case ch <- i:
			return true
		}
	}
}
