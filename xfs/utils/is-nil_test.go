package utils_test

import (
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ginkgo ok
	. "github.com/onsi/gomega"    //nolint:revive // gomega ok

	"github.com/snivilised/extendio/xfs/utils"
)

type blob struct{}

var _ = Describe("IsNil", func() {

	When("received item is not nil", func() {
		Context("pointer to struct", func() {
			It("🧪 should: return false", func() {
				item := &blob{}
				utils.IsNil(item)
				Expect(utils.IsNil(item)).To(BeFalse())
			})
		})

		Context("interface", func() {
			It("🧪 should: return false", func() {
				var item interface{} = &blob{}
				utils.IsNil(item)
				Expect(utils.IsNil(item)).To(BeFalse())
			})
		})

		Context("struct", func() {
			It("🧪 should: return false", func() {
				item := blob{}
				utils.IsNil(item)
				Expect(utils.IsNil(item)).To(BeFalse())
			})
		})
	})

	When("received item is a nil", func() {
		Context("pointer to struct", func() {
			It("🧪 should: return true", func() {
				var item *blob
				utils.IsNil(item)
				Expect(utils.IsNil(item)).To(BeTrue())
			})
		})

		Context("interface", func() {
			It("🧪 should: return true", func() {
				var item interface{}
				utils.IsNil(item)
				Expect(utils.IsNil(item)).To(BeTrue())
			})
		})
	})
})
