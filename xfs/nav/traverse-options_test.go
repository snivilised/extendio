package nav_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("TraverseOptions", Ordered, func() {

	var (
		o *nav.TraverseOptions
	)

	BeforeEach(func() {
		o = nav.GetDefaultOptions()
	})

	Context("clone", func() {
		When("given: options", func() {
			It("should: return a deep copy", func() {

				cloneCount, sourceCount := 0, 0

				o.Notify.OnBegin = func(root string) {
					sourceCount++
				}
				clone := o.Clone()
				Expect(clone).NotTo(BeNil())

				clone.Subscription = nav.SubscribeFiles
				Expect(o.Subscription).To(Equal(nav.SubscribeAny))

				clone.Behaviours.SubPath.KeepTrailingSep = false
				Expect(o.Behaviours.SubPath.KeepTrailingSep).To(BeTrue())

				clone.Filters.Current = &nav.RegexFilter{
					Filter: nav.Filter{
						Name:    "test filter",
						Pattern: "foo bar",
					},
				}
				Expect(o.Filters.Current).To(BeNil())
				o.Notify.OnBegin("/foo-bar")

				clone.Notify.OnBegin = func(root string) {
					cloneCount++
				}
				clone.Notify.OnBegin("/foo-bar")

				Expect(sourceCount).To(Equal(1), "")
				Expect(cloneCount).To(Equal(1))
			})
		})
	})
})
