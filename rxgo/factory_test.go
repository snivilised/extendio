package rxgo_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/matchers"
	"github.com/onsi/gomega/types"

	"github.com/snivilised/extendio/rxgo"
)

var _ = Describe("Factory", func() {
	Context("Just", func() {
		When("item is scalar", func() {
			It("🧪 should: create observable", func() {
				observable := rxgo.Just(1, 2, 3)

				actual := RxObservable[int]{
					context:    context.Background(),
					observable: observable(),
				}

				Expect(actual).Should(&AndMatcher{
					Matchers: []types.GomegaMatcher{
						MatchHasItems(1, 2, 3),
						MatchHasNoError[int](),
					},
				})
			})
		})

		When("item is custom struct", func() {
			It("🧪 should: create observable", func() {
				type customer struct {
					id int
				}

				sequence := []customer{{id: 1}, {id: 2}, {id: 3}}
				observable := rxgo.Just(sequence...)
				actual := RxObservable[customer]{
					context:    context.Background(),
					observable: observable(),
				}
				expected := []customer{{id: 1}, {id: 2}, {id: 3}}

				Expect(actual).Should(&AndMatcher{
					Matchers: []types.GomegaMatcher{
						MatchHasItems(expected...),
						MatchHasNoError[customer](),
					},
				})
			})
		})
	})
})
