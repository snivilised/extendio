package rx_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/extendio/rx"
)

var _ = Describe("Item", func() {
	Context("foo", func() {
		XIt("should: ...", func() {
			var i rx.Item[int]
			_ = i
			// Expect(item).Should(MatchCurrentGlobFilter(filter))
			Expect(1).To(Equal(1))
		})
	})
})
