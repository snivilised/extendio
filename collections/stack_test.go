package collections_test

import (
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ginkgo ok
	. "github.com/onsi/gomega"    //nolint:revive // gomega ok

	"github.com/snivilised/extendio/collections"
)

var _ = Describe("Stack", func() {
	Context("Push", func() {
		It("🧪 should: add element to stack", func() {
			stack := collections.NewStackWith([]string{"north", "east", "south"})
			stack.Push("west")

			const (
				size    = uint(4)
				current = "west"
			)
			Expect(stack).To(BeInCorrectState(size, current))
		})
	})

	Context("Pop", func() {
		When("stack is empty", func() {
			It("🧪 should: return empty error", func() {
				stack := collections.NewStack[string]()
				_, err := stack.Pop()
				Expect(err).To(MatchError("internal: stack is empty"))

			})
		})
		When("stack is NOT empty", func() {
			It("🧪 should: remove top item", func() {
				stack := collections.NewStackWith([]string{
					"red", "orange", "yellow", "green", "blue", "indigo", "violet",
				})
				actualPop, _ := stack.Pop()

				const (
					size        = uint(6)
					expectedPop = "violet"
				)
				withExpectedPop := &WithExpectedPop{
					stack:  stack,
					popped: expectedPop,
				}

				Expect(withExpectedPop).To(HavePopped(size, actualPop))
			})
		})
	})

	Context("MustPop", func() {
		When("stack is NOT empty", func() {
			It("🧪 should: remove top item", func() {
				stack := collections.NewStackWith([]string{
					"red", "orange", "yellow", "green", "blue", "indigo", "violet",
				})
				actualPop := stack.MustPop()

				const (
					size        = uint(6)
					expectedPop = "violet"
				)
				withExpectedPop := &WithExpectedPop{
					stack:  stack,
					popped: expectedPop,
				}

				Expect(withExpectedPop).To(HavePopped(size, actualPop))
			})
		})

		When("stack is empty", func() {
			It("🧪 should: panic", func() {
				stack := collections.NewStack[string]()

				Expect(func() {
					stack.MustPop()
				}).To(PanicWith(collections.NewStackIsEmptyNativeError()))
			})
		})
	})

	Context("Current", func() {
		When("stack is empty", func() {
			It("🧪 should: return empty error", func() {
				stack := collections.NewStack[string]()
				_, err := stack.Current()
				Expect(err).To(MatchError("internal: stack is empty"))
			})
		})

		When("stack not empty", func() {
			It("🧪 should: return correct current value", func() {
				stack := collections.NewStackWith([]string{"north", "east", "south", "west"})
				const (
					size    = uint(4)
					current = "west"
				)
				Expect(stack).To(BeInCorrectState(size, current))
			})
		})
	})

	Context("Content", func() {
		It("🧪 should: return inner slice", func() {
			with := []string{"red", "orange", "yellow", "green", "blue", "indigo", "violet"}
			stack := collections.NewStackWith(with)
			content := stack.Content()
			Expect(content).To(Equal(with))
		})
	})
})
