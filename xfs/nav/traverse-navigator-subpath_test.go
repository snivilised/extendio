package nav_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/snivilised/extendio/internal/helpers"
	. "github.com/snivilised/extendio/translate"
	"github.com/snivilised/extendio/xfs/nav"
)

var _ = Describe("TraverseNavigatorSubpath", Ordered, func() {
	var root string

	BeforeAll(func() {
		root = musico()
	})

	Context("sub-path", func() {
		When("KeepTrailingSep set to true", func() {
			It("should: calculate subpath WITH trailing separator", func() {

				expectations := map[string]string{
					"RETRO-WAVE":                   "",
					"Chromatics":                   helpers.Normalise("/"),
					"Night Drive":                  helpers.Normalise("/Chromatics/"),
					"A1 - The Telephone Call.flac": helpers.Normalise("/Chromatics/Night Drive/"),
				}
				path := helpers.Path(root, "RETRO-WAVE")
				session := &nav.PrimarySession{
					Path: path,
				}
				_ = session.Configure(func(o *nav.TraverseOptions) {
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
								Expect(item.Extension.SubPath).To(Equal(expected),
									helpers.Reason(item.Extension.Name),
								)
							}

							return nil
						},
					}
				}).Run()
			})

			When("using RootItemSubPath", func() {
				It("should: calculate subpath WITH trailing separator", func() {

					expectations := map[string]string{
						"edm":                         "",
						"_segments.def.infex.txt":     helpers.Normalise("/_segments.def.infex.txt"),
						"Orbital 2 (The Brown Album)": helpers.Normalise("/ELECTRONICA/Orbital/Orbital 2 (The Brown Album)"),
						"03 - Lush 3-1.flac":          helpers.Normalise("/ELECTRONICA/Orbital/Orbital 2 (The Brown Album)/03 - Lush 3-1.flac"),
					}
					path := helpers.Path(root, "edm")
					session := &nav.PrimarySession{
						Path: path,
					}
					_ = session.Configure(func(o *nav.TraverseOptions) {
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
									Expect(item.Extension.SubPath).To(Equal(expected), helpers.Reason(item.Extension.Name))
								}

								return nil
							},
						}
					}).Run()
				})
			})
		})

		When("KeepTrailingSep set to false", func() {
			It("should: calculate subpath WITHOUT trailing separator", func() {
				expectations := map[string]string{
					"RETRO-WAVE":            "",
					"Electric Youth":        "",
					"Innerworld":            helpers.Normalise("/Electric Youth"),
					"A1 - Before Life.flac": helpers.Normalise("/Electric Youth/Innerworld"),
				}
				path := helpers.Path(root, "RETRO-WAVE")
				session := &nav.PrimarySession{
					Path: path,
				}
				_ = session.Configure(func(o *nav.TraverseOptions) {
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
								Expect(item.Extension.SubPath).To(Equal(expected), helpers.Reason(item.Extension.Name))
							}

							return nil
						},
					}
				}).Run()
			})
		})
	})
})
