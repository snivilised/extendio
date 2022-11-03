package nav_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/snivilised/extendio/translate"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("TraverseNavigatorSubpath", Ordered, func() {
	var root string

	BeforeAll(func() {
		root = cwd()
	})

	Context("sub-path", func() {
		When("KeepTrailingSep set to true", func() {
			It("should: calculate subpath WITH trailing separator", func() {

				expectations := map[string]string{
					"RETRO-WAVE":                   "",
					"Chromatics":                   normalise("/"),
					"Night Drive":                  normalise("/Chromatics/"),
					"A1 - The Telephone Call.flac": normalise("/Chromatics/Night Drive/"),
				}
				navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🛡️")
					o.Subscription = nav.SubscribeAny
					o.Behaviours.SubPath.KeepTrailingSep = true
					o.DoExtend = true
					o.Callback = func(item *nav.TraverseItem) *LocalisableError {
						if expected, ok := expectations[item.Extension.Name]; ok {
							Expect(item.Extension.SubPath).To(Equal(expected), reason(item.Extension.Name))
							GinkgoWriter.Printf("---> 🧩 SUB-PATH-CALLBACK(with): '%v', name: '%v', scope: '%v'\n",
								item.Extension.SubPath, item.Extension.Name, item.Extension.NodeScope,
							)
						}

						return nil
					}
				})
				path := path(root, "RETRO-WAVE")
				navigator.Walk(path)
			})

			When("using RootItemSubPath", func() {
				It("should: calculate subpath WITH trailing separator", func() {

					expectations := map[string]string{
						"edm":                         "",
						"_segments.def.infex.txt":     normalise("/_segments.def.infex.txt"),
						"Orbital 2 (The Brown Album)": normalise("/ELECTRONICA/Orbital/Orbital 2 (The Brown Album)"),
						"03 - Lush 3-1.flac":          normalise("/ELECTRONICA/Orbital/Orbital 2 (The Brown Album)/03 - Lush 3-1.flac"),
					}
					navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
						o.Notify.OnBegin = begin("🛡️")
						o.Subscription = nav.SubscribeAny
						o.Hooks.FolderSubPath = nav.RootItemSubPath
						o.Hooks.FileSubPath = nav.RootItemSubPath
						o.Behaviours.SubPath.KeepTrailingSep = true
						o.DoExtend = true
						o.Callback = func(item *nav.TraverseItem) *LocalisableError {
							if expected, ok := expectations[item.Extension.Name]; ok {
								Expect(item.Extension.SubPath).To(Equal(expected), reason(item.Extension.Name))
								GinkgoWriter.Printf("---> 🧩🧩 SUB-PATH-CALLBACK(with): '%v', name: '%v', scope: '%v'\n",
									item.Extension.SubPath, item.Extension.Name, item.Extension.NodeScope,
								)
							}

							return nil
						}
					})
					path := path(root, "edm")
					navigator.Walk(path)
				})
			})
		})

		When("KeepTrailingSep set to false", func() {
			It("should: calculate subpath WITHOUT trailing separator", func() {
				expectations := map[string]string{
					"RETRO-WAVE":            "",
					"Electric Youth":        "",
					"Innerworld":            normalise("/Electric Youth"),
					"A1 - Before Life.flac": normalise("/Electric Youth/Innerworld"),
				}
				navigator := nav.NewNavigator(func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🛡️")
					o.Behaviours.SubPath.KeepTrailingSep = false
					o.Subscription = nav.SubscribeAny
					o.DoExtend = true
					o.Callback = func(item *nav.TraverseItem) *LocalisableError {
						if expected, ok := expectations[item.Extension.Name]; ok {
							Expect(item.Extension.SubPath).To(Equal(expected), reason(item.Extension.Name))
							GinkgoWriter.Printf("---> 🧩 SUB-PATH-CALLBACK(without): '%v', name: '%v', scope: '%v'\n",
								item.Extension.SubPath, item.Extension.Name, item.Extension.NodeScope,
							)
						}

						return nil
					}
				})
				path := path(root, "RETRO-WAVE")
				navigator.Walk(path)
			})
		})
	})
})
