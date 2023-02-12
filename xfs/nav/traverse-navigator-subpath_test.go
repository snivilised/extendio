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
		root = origin()
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
				navigator := nav.NavigatorFactory{}.Construct(func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🛡️")
					o.Store.Subscription = nav.SubscribeAny
					o.Store.Behaviours.SubPath.KeepTrailingSep = true
					o.Store.DoExtend = true
					o.Callback = nav.LabelledTraverseCallback{
						Label: "test sub-path callback",
						Fn: func(item *nav.TraverseItem) *LocalisableError {
							if expected, ok := expectations[item.Extension.Name]; ok {
								GinkgoWriter.Printf("---> 🧩 SUB-PATH-CALLBACK(with): '%v', name: '%v', scope: '%v'\n",
									item.Extension.SubPath, item.Extension.Name, item.Extension.NodeScope,
								)
								Expect(item.Extension.SubPath).To(Equal(expected), reason(item.Extension.Name))
							}

							return nil
						},
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
					navigator := nav.NavigatorFactory{}.Construct(func(o *nav.TraverseOptions) {
						o.Notify.OnBegin = begin("🛡️")
						o.Store.Subscription = nav.SubscribeAny
						o.Hooks.FolderSubPath = nav.RootItemSubPath
						o.Hooks.FileSubPath = nav.RootItemSubPath
						o.Store.Behaviours.SubPath.KeepTrailingSep = true
						o.Store.DoExtend = true
						o.Callback = nav.LabelledTraverseCallback{
							Label: "test sub-path callback",
							Fn: func(item *nav.TraverseItem) *LocalisableError {
								if expected, ok := expectations[item.Extension.Name]; ok {
									GinkgoWriter.Printf("---> 🧩🧩 SUB-PATH-CALLBACK(with): '%v', name: '%v', scope: '%v'\n",
										item.Extension.SubPath, item.Extension.Name, item.Extension.NodeScope,
									)
									Expect(item.Extension.SubPath).To(Equal(expected), reason(item.Extension.Name))
								}

								return nil
							},
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
				navigator := nav.NavigatorFactory{}.Construct(func(o *nav.TraverseOptions) {
					o.Notify.OnBegin = begin("🛡️")
					o.Store.Behaviours.SubPath.KeepTrailingSep = false
					o.Store.Subscription = nav.SubscribeAny
					o.Store.DoExtend = true
					o.Callback = nav.LabelledTraverseCallback{
						Label: "test sub-path callback",
						Fn: func(item *nav.TraverseItem) *LocalisableError {
							if expected, ok := expectations[item.Extension.Name]; ok {
								GinkgoWriter.Printf("---> 🧩 SUB-PATH-CALLBACK(without): '%v', name: '%v', scope: '%v'\n",
									item.Extension.SubPath, item.Extension.Name, item.Extension.NodeScope,
								)
								Expect(item.Extension.SubPath).To(Equal(expected), reason(item.Extension.Name))
							}

							return nil
						},
					}
				})
				path := path(root, "RETRO-WAVE")
				navigator.Walk(path)
			})
		})
	})
})
