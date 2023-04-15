package utils_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/extendio/xfs/utils"
)

type slugInfo struct {
	count utils.VarProp[int]
}

type sizeable interface {
	Measure() int
}

type nugget struct {
	size int
}

func (n *nugget) Measure() int {
	return n.size
}

type fnProperty func() float64
type slice []int
type dictionary map[int]string
type intChan chan int
type widget struct {
	// variable properties
	//
	colour   utils.VarProp[string]
	slug     utils.VarProp[*slugInfo]
	quantity utils.VarProp[sizeable]
	gold     utils.VarProp[nugget]
	fn       utils.VarProp[fnProperty]
	points   utils.VarProp[slice]
	numbers  utils.VarProp[dictionary]
	tunnel   utils.VarProp[intChan]
	fraction utils.RwProp[float64]

	// putter variable properties
	//
	rank utils.PutProp[int]

	// const properties
	//
	colourRo utils.RoProp[string]
}

const pi5dp = 3.14159

func pi() float64 {
	return pi5dp
}

var evens = dictionary{
	2: "Two",
	4: "Four",
	6: "Six",
}

var _ = Describe("Property", func() {

	Context("given: string property", func() {
		Context("VarProp", func() {
			Context("Get", func() {
				It("should: retrieve property value", func() {
					w := &widget{
						colour: utils.VarProp[string]{Field: "red"},
					}
					Expect(w.colour.Get()).To(Equal("red"))
				})
			})

			Context("Set", func() {
				It("should: set property value", func() {
					w := &widget{
						colour: utils.VarProp[string]{Field: "red"},
					}
					w.colour.Set("blue")
					Expect(w.colour.Get()).To(Equal("blue"))
				})
			})

			Context("IsNone", func() {
				When("string value is previously defined", func() {
					It("should: return false without panic", func() {
						w := &widget{
							colour: utils.VarProp[string]{Field: "red"},
						}
						Expect(w.colour.IsNone()).To(BeFalse())
					})
				})

				When("string value is unassigned", func() {
					It("should: return false without panic", func() {
						w := &widget{
							colour: utils.VarProp[string]{},
						}
						Expect(w.colour.IsNone()).To(BeTrue())
					})
				})
			})

			Context("RoRef", func() {
				It("should: get read only interface", func() {
					w := &widget{
						colour: utils.VarProp[string]{Field: "red"},
					}
					RoRef := w.colour.RoRef()
					Expect(RoRef.Get()).To(Equal("red"))

					w.colour.Set("blue")
					Expect(RoRef.Get()).To(Equal("blue"))
				})
			})
		})

		Context("ConstProp", func() {
			Context("Get", func() {
				It("should: retrieve property value", func() {
					w := &widget{
						colourRo: utils.NewRoProp("red"),
					}
					Expect(w.colourRo.Get()).To(Equal("red"))
				})
			})

			Context("IsNone", func() {
				When("string value is previously defined", func() {
					It("should: return false without panic", func() {
						w := &widget{
							colourRo: utils.RoPropFactory[string]{}.New("red"),
						}
						Expect(w.colourRo.IsNone()).To(BeFalse())
					})
				})
			})
		})
	})

	Context("given: pointer to struct property (with scalar [count])", func() {
		Context("VarProp", func() {
			Context("Get", func() {
				It("should: retrieve property value", func() {
					w := &widget{
						slug: utils.VarProp[*slugInfo]{
							Field: &slugInfo{
								count: utils.VarProp[int]{Field: 42},
							},
						},
					}
					Expect(w.slug.Get().count.Get()).To(Equal(42))
				})
			})

			Context("Set", func() {
				It("should: set property value", func() {
					w := &widget{}
					w.slug.Set(&slugInfo{
						count: utils.VarProp[int]{Field: 42},
					})
					Expect(w.slug.Get().count.Get()).To(Equal(42))
				})
			})

			Context("IsNone", func() {
				When("pointer value is previously defined", func() {
					It("should: return false without panic", func() {
						w := &widget{
							slug: utils.VarProp[*slugInfo]{
								Field: &slugInfo{
									count: utils.VarProp[int]{Field: 42},
								},
							},
						}
						Expect(w.slug.IsNone()).To(BeFalse())
						Expect(w.slug.Get().count.IsNone()).To(BeFalse())
					})
				})

				When("pointer value is unassigned", func() {
					It("should: return true without panic", func() {
						w := &widget{}
						Expect(w.slug.IsNone()).To(BeTrue())
					})
				})
			})
		})
	})

	Context("given: int property", func() {
		Context("PutProp/(putVarProp)", func() {
			Context("Get", func() {
				It("should: retrieve property value", func() {
					factory := utils.PutPropFactory[int]{}
					w := &widget{
						rank: factory.New(77, func(value int) {}),
					}
					Expect(w.rank.Get()).To(Equal(77))
				})
			})

			Context("Set", func() {
				It("should: set property value", func() {
					factory := utils.PutPropFactory[int]{}
					w := &widget{
						rank: factory.New(77, func(value int) {}),
					}
					w.rank.Set(88)
					Expect(w.rank.Get()).To(Equal(88))
				})
			})

			Context("Put", func() {
				It("should: set property value via putter", func() {
					factory := utils.PutPropFactory[int]{}
					var another int
					w := &widget{
						rank: factory.New(77, func(value int) {
							another = value
						}),
					}
					w.rank.Put(88)
					Expect(another).To(Equal(88))
				})
			})

			Context("IsNone", func() {
				When("int value is previously defined", func() {
					It("should: return false without panic", func() {
						factory := utils.PutPropFactory[int]{}
						w := &widget{
							rank: factory.New(77, func(value int) {}),
						}
						Expect(w.rank.IsNone()).To(BeFalse())
					})
				})

				When("value is explicitly set to it's zero value", func() {
					Context("Is Zeroable", func() {
						It("should: return false without panic", func() {
							factory := utils.PutPropFactory[int]{Zeroable: true}
							w := &widget{
								rank: factory.New(0, func(value int) {}),
							}
							Expect(w.rank.IsNone()).To(BeFalse())
						})
					})

					Context("Is NOT Zeroable", func() {
						It("should: return false without panic", func() {
							factory := utils.PutPropFactory[int]{}
							w := &widget{
								rank: factory.New(0, func(value int) {}),
							}
							Expect(w.rank.IsNone()).To(BeTrue())
						})
					})
				})
			})
		})
	})

	Context("given: float64 property", func() {
		Context("VarProp", func() {
			Context("IsNone", func() {
				When("flat64 value is previously defined", func() {
					It("should: return false without panic", func() {
						w := &widget{
							fraction: utils.NewRwPropZ(float64(0.12345)),
						}
						Expect(w.fraction.IsNone()).To(BeFalse())
					})
				})

				When("value is explicitly set to it's zero value", func() {
					Context("Is Zeroable", func() {
						It("should: return false without panic", func() {
							factory := utils.RwPropFactory[float64]{Zeroable: true}
							w := &widget{
								fraction: factory.New(float64(0)),
							}
							Expect(w.fraction.IsNone()).To(BeFalse())
						})
					})

					Context("Is NOT Zeroable", func() {
						It("should: return false without panic", func() {
							factory := utils.RwPropFactory[float64]{}
							w := &widget{
								fraction: factory.New(float64(0)),
							}
							Expect(w.fraction.IsNone()).To(BeTrue())
						})
					})
				})
			})
		})
	})

	Context("given: interface property", func() {
		Context("VarProp", func() {
			Context("Get", func() {
				It("should: retrieve property value", func() {
					w := &widget{
						quantity: utils.VarProp[sizeable]{Field: &nugget{66}},
					}
					Expect(w.quantity.Get().Measure()).To(Equal(66))
				})
			})

			Context("Set", func() {
				It("should: set property value", func() {
					w := &widget{
						quantity: utils.VarProp[sizeable]{Field: &nugget{66}},
					}
					w.quantity.Set(&nugget{99})
					Expect(w.quantity.Get().Measure()).To(Equal(99))
				})
			})

			Context("IsNone", func() {
				When("interface value is previously defined", func() {
					It("should: return false without panic", func() {
						w := &widget{
							quantity: utils.VarProp[sizeable]{Field: &nugget{66}},
						}
						Expect(w.quantity.IsNone()).To(BeFalse())
					})
				})

				When("pointer value is unassigned", func() {
					It("should: return true without panic", func() {
						w := &widget{}
						Expect(w.quantity.IsNone()).To(BeTrue())
					})
				})
			})
		})
	})

	Context("given: struct property", func() {
		Context("VarProp", func() {
			Context("Get", func() {
				It("should: retrieve copy of property value", func() {
					w := &widget{
						gold: utils.VarProp[nugget]{Field: nugget{66}},
					}
					clone := w.gold.Get()
					Expect(clone.size).To(Equal(66))

					clone.size = 42
					Expect(w.gold.Get().size).To(Equal(66))
				})
			})

			Context("Set", func() {
				It("should: set struct property value", func() {
					w := &widget{
						gold: utils.VarProp[nugget]{},
					}
					w.gold.Set(nugget{66})
					Expect(w.gold.Get().size).To(Equal(66))
				})
			})

			Context("IsNone", func() {
				When("struct value is previously defined", func() {
					It("should: return false without panic", func() {
						w := &widget{
							gold: utils.VarProp[nugget]{Field: nugget{66}},
						}
						Expect(w.gold.IsNone()).To(BeFalse())
					})
				})

				When("struct value is unassigned", func() {
					It("should: return true without panic", func() {
						w := &widget{}
						Expect(w.slug.IsNone()).To(BeTrue())
					})
				})
			})
		})
	})

	Context("given: function property", func() {
		Context("VarProp", func() {
			Context("Get", func() {
				It("should: retrieve function property value", func() {
					w := &widget{
						fn: utils.VarProp[fnProperty]{Field: pi},
					}
					Expect(w.fn.Get()()).To(Equal(pi5dp))
				})
			})

			Context("Set", func() {
				It("should: set function property value", func() {
					w := &widget{
						fn: utils.VarProp[fnProperty]{},
					}
					w.fn.Set(pi)
					Expect(w.fn.Get()()).To(Equal(pi5dp))
				})
			})

			Context("IsNone", func() {
				When("struct value is previously defined", func() {
					It("should: return false without panic", func() {
						w := &widget{
							fn: utils.VarProp[fnProperty]{Field: pi},
						}
						Expect(w.fn.IsNone()).To(BeFalse())
					})
				})

				When("struct value is unassigned", func() {
					It("should: return true without panic", func() {
						w := &widget{}
						Expect(w.fn.IsNone()).To(BeTrue())
					})
				})
			})
		})
	})

	Context("given: slice property", func() {
		Context("VarProp", func() {
			Context("Get", func() {
				It("should: retrieve slice property value", func() {
					w := &widget{
						points: utils.VarProp[slice]{Field: slice{2, 4, 6, 8}},
					}
					Expect(w.points.Get()).To(Equal(slice{2, 4, 6, 8}))
				})
			})

			Context("Set", func() {
				It("should: set slice property value", func() {
					w := &widget{}
					w.points.Set(slice{2, 4, 6, 8})
					Expect(w.points.Get()).To(Equal(slice{2, 4, 6, 8}))
				})
			})

			Context("IsNone", func() {
				When("slice value is previously defined", func() {
					It("should: return false without panic", func() {
						w := &widget{
							points: utils.VarProp[slice]{Field: slice{2, 4, 6, 8}},
						}
						Expect(w.points.IsNone()).To(BeFalse())
					})
				})

				When("slice value is unassigned", func() {
					It("should: return true without panic", func() {
						w := &widget{}
						Expect(w.points.IsNone()).To(BeTrue())
					})
				})
			})
		})
	})

	Context("given: map property", func() {
		Context("VarProp", func() {
			Context("Get", func() {
				It("should: retrieve map property value", func() {
					w := &widget{
						numbers: utils.VarProp[dictionary]{Field: evens},
					}
					Expect(w.numbers.Get()[2]).To(Equal("Two"))
				})
			})

			Context("Set", func() {
				It("should: set map property value", func() {
					w := &widget{}
					w.numbers.Set(evens)
					// check equivalence, not identity
					//
					Expect(w.numbers.Get()).To(Equal(dictionary{
						2: "Two",
						4: "Four",
						6: "Six",
					}))
				})
			})

			Context("IsNone", func() {
				When("slice value is previously defined", func() {
					It("should: return false without panic", func() {
						w := &widget{
							numbers: utils.VarProp[dictionary]{Field: evens},
						}
						Expect(w.numbers.IsNone()).To(BeFalse())
					})
				})

				When("slice value is unassigned", func() {
					It("should: return true without panic", func() {
						w := &widget{}
						Expect(w.numbers.IsNone()).To(BeTrue())
					})
				})
			})
		})
	})

	Context("given: channel property", func() {
		Context("VarProp", func() {
			Context("Get", func() {
				It("should: retrieve channel property value", func() {
					w := &widget{
						tunnel: utils.VarProp[intChan]{Field: make(chan int)},
					}
					Expect(w.tunnel.Get()).NotTo(BeNil())
				})
			})

			Context("Set", func() {
				It("should: set map property value", func() {
					w := &widget{}
					w.tunnel.Set(make(chan int))
					Expect(w.tunnel.Get()).NotTo(BeNil())
				})
			})

			Context("IsNone", func() {
				When("channel value is previously defined", func() {
					It("should: return false without panic", func() {
						w := &widget{
							tunnel: utils.VarProp[intChan]{Field: make(chan int)},
						}
						Expect(w.tunnel.IsNone()).To(BeFalse())
					})
				})

				When("channel value is unassigned", func() {
					It("should: return true without panic", func() {
						w := &widget{}
						Expect(w.tunnel.IsNone()).To(BeTrue())
					})
				})
			})
		})
	})
})
