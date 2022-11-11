package nav_test

import (
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("TraverseOptionsIo", Ordered, func() {
	var (
		o            *nav.TraverseOptions
		root         string
		fromJsonPath string
		toJsonPath   string
	)

	BeforeAll(func() {
		root = joinCwd("Test", "json")
		fromJsonPath = strings.Join([]string{root, "options.json"}, string(filepath.Separator))
		toJsonPath = strings.Join([]string{root, "test-options-marshal.json"}, string(filepath.Separator))
	})

	BeforeEach(func() {
		o = nav.GetDefaultOptions()
	})

	Context("Marshal", func() {
		Context("given: correct config", func() {
			It("🧪 should: write options in JSON", func() {
				o.Persist.Restorer = func(o *nav.TraverseOptions) {
					GinkgoWriter.Printf("---> marshaller ...\n")
				}

				err := o.Marshal(toJsonPath)
				Expect(err).To(BeNil())
			})
		})

		When("restorer function undefined", func() {
			It("🧪 should: panic", func() {
				defer func() {
					pe := recover()
					if err, ok := pe.(error); ok {
						Expect(strings.Contains(err.Error(), "missing restorer function")).To(BeTrue())
					}
				}()

				if err := o.Marshal(toJsonPath); err != nil {
					GinkgoWriter.Printf("---> 🔥🔥🔥 marshall error: '%v'\n", err)
				}
				Fail("❌ expected panic due to missing restorer function")
			})
		})

		When("persist format undefined", func() {
			It("🧪 should: panic", func() {
				defer func() {
					pe := recover()
					if err, ok := pe.(error); ok {
						Expect(strings.Contains(err.Error(), "unknown marshal format")).To(BeTrue())
					}
				}()

				o.Persist.Format = nav.PersistInUndefinedEn
				o.Persist.Restorer = func(o *nav.TraverseOptions) {}
				_ = o.Marshal(toJsonPath)
				Fail("❌ expected panic due to unknown marshal format")
			})
		})
	})

	Context("Unmarshal", func() {
		Context("given: correct config", func() {
			It("🧪 should: write options in JSON", func() {
				o.Persist.Restorer = func(o *nav.TraverseOptions) {
					GinkgoWriter.Printf("---> marshaller ...\n")
				}

				err := o.Unmarshal(fromJsonPath)
				Expect(err).To(BeNil())
			})
		})
	})
})
