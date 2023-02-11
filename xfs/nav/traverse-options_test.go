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

				o.Notify.OnBegin = func(state *nav.NavigationState) {
					sourceCount++
				}
				clone := o.Clone()
				Expect(clone).NotTo(BeNil())

				clone.Store.Subscription = nav.SubscribeFiles
				Expect(o.Store.Subscription).To(Equal(nav.SubscribeAny))

				clone.Store.Behaviours.SubPath.KeepTrailingSep = false
				Expect(o.Store.Behaviours.SubPath.KeepTrailingSep).To(BeTrue())

				clone.Store.FilterDefs = &nav.FilterDefinitions{
					Node: nav.FilterDef{
						Type:        nav.FilterTypeRegexEn,
						Description: "test filter",
						Source:      "foo bar",
					},
				}
				state := &nav.NavigationState{}
				o.Notify.OnBegin(state)

				clone.Notify.OnBegin = func(state *nav.NavigationState) {
					cloneCount++
				}
				clone.Notify.OnBegin(state)

				Expect(sourceCount).To(Equal(1), "")
				Expect(cloneCount).To(Equal(1))
			})
		})
	})
})
