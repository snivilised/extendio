package rxgo_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
	"github.com/onsi/gomega/types"

	"github.com/snivilised/extendio/rxgo"
)

// TODO: check to see if we need to use async ginkgo types

var _ = Describe("Item", func() {
	Context("Variadic", func() {
		Context("AndMatcher", func() {
			It("🧪 should: Send Items", func() {
				ch := make(chan rxgo.Item[int], 3)
				actual := RxObservable[int]{
					context:    context.Background(),
					observable: rxgo.FromChannel(ch),
				}

				go rxgo.SendItemsV(context.Background(), ch, rxgo.CloseChannel, 1, 2, 3)
				Expect(actual).Should(&AndMatcher{
					Matchers: []types.GomegaMatcher{
						MatchHasItems(1, 2, 3),
					},
				})
			})
		})

		Context("Singular MatchHasItems", func() {
			It("🧪 should: Send Items", func() {
				ch := make(chan rxgo.Item[int], 3)
				actual := RxObservable[int]{
					context:    context.Background(),
					observable: rxgo.FromChannel(ch),
				}

				go rxgo.SendItems(context.Background(), ch, rxgo.CloseChannel, 1, 2, 3)
				Expect(actual).Should(MatchHasItems(1, 2, 3))
			})
		})
	})

	Context("Variadic with Error", func() {
		It("🧪 should: (can't send an error into a channel)", Pending)
	})
})
