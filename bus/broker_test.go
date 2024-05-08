package bus_test

// Copyright 2021 Mustafa Turan. All rights reserved.
// Use of this source code is governed by a Apache License 2.0 license that can
// be found in the bus.LICENSE file.

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/extendio/bus"
)

const (
	topicCommentCreated = "comment.created"
	topicCommentDeleted = "comment.deleted"
	topicCommentUpdated = "comment.updated"
	topicUserCreated    = "user.created"
	topicUserDeleted    = "user.deleted"
	topicUserUpdated    = "user.updated"
	fake                = "a-fake-id"
)

type handlerTE struct {
	given            string
	should           string
	handler          bus.Handler
	handlerKey       string
	handlerLookupKey string
	expected         []string
}

var _ = Describe("Bus", func() {
	Context("foo", func() {
		It("🧪 should: ", func() {
			_, err := bus.New(&bus.Sequential{
				Format: "message: '%05d'",
			})

			Expect(err).To(Succeed())
		})
	})

	Context("CtxKeyTxID", func() {
		// TestCtxKeyTxID
		It("🧪 should: return predefined id", func() {
			Expect(bus.CtxKeyTxID).To(BeEquivalentTo(116))
		})
	})

	Context("New", func() {
		// TestNew...
		It("🧪 should: run with valid generator", func() {
			var fn bus.Next = func() string { return fake }
			_, err := bus.New(fn)
			Expect(err).To(Succeed())
		})
	})

	Context("Emit", Ordered, func() {
		var b *bus.Broker

		BeforeEach(func() {
			b = setup(topicCommentCreated, topicCommentDeleted)
		})

		AfterEach(func() {
			tearDown(b, topicCommentCreated, topicCommentDeleted)
		})

		Context("no opts", func() {
			When("tbd", func() {
				It("🧪 should: correctly assigns fields", func() {
					ctx := context.Background()
					ctx = context.WithValue(ctx, bus.CtxKeyTxID, "tx")
					ctx = context.WithValue(ctx, bus.CtxKeySource, "source")
					err := b.Emit(ctx, topicCommentDeleted, "my comment")

					Expect(err).To(Succeed())
				})
			})

			When("txID is empty", func() {
				It("🧪 should: update txID", func() {
					ctx := context.Background()
					err := b.Emit(ctx, topicCommentDeleted, "my comment")

					Expect(err).To(Succeed())
				})
			})

			When("with handler", func() {
				It("🧪 should: tbd", func() {
					ctx := context.Background()
					registerFakeHandler(b, "test")

					err := b.Emit(ctx, topicCommentCreated, "my comment with handler")
					if err != nil {
						Fail(fmt.Sprintf("emit failed: %v", err))
					}
					b.DeregisterHandler("test")
				})
			})

			When("with unknown topic", func() {
				It("🧪 should: tbd", func() {
					ctx := context.Background()
					err := b.Emit(ctx, topicCommentUpdated, "my comment")

					Expect(err).NotTo(Succeed())
					Expect(err.Error()).To(Equal("topics(comment.updated) not found"))
				})
			})
		})

		Context("with opts", func() {
			When("tbd", func() {
				It("🧪 should: correctly assigns fields", func() {
					ctx := context.Background()
					ctx = context.WithValue(ctx, bus.CtxKeyTxID, "tx")
					ctx = context.WithValue(ctx, bus.CtxKeySource, "source")
					err := b.EmitWithOpts(ctx, topicCommentDeleted, "my comment",
						bus.WithTxID("tx"),
						bus.WithID("id"),
						bus.WithSource("source"),
						bus.WithOccurredAt(time.Now()),
					)

					Expect(err).To(Succeed())
				})
			})

			When("txID is empty", func() {
				It("🧪 should: update txID", func() {
					ctx := context.Background()
					err := b.EmitWithOpts(ctx, topicCommentDeleted, "my comment")

					Expect(err).To(Succeed())
				})
			})

			When("with handler", func() {
				It("🧪 should: tbd", func() {
					ctx := context.Background()
					registerFakeHandler(b, "test")

					err := b.EmitWithOpts(ctx, topicCommentCreated, "my comment with handler")
					if err != nil {
						Fail(fmt.Sprintf("emit failed: %v", err))
					}
					b.DeregisterHandler("test")
				})
			})

			When("with unknown topic", func() {
				It("🧪 should: tbd", func() {
					ctx := context.Background()
					err := b.EmitWithOpts(ctx, topicCommentUpdated, "my comment")

					Expect(err).NotTo(Succeed())
					Expect(err.Error()).To(Equal("topics(comment.updated) not found"))
				})
			})
		})
	})

	Context("Topics", Ordered, func() {
		When("register", func() {
			It("🧪 should: return registered topics", func() {
				topics := []string{topicUserCreated, topicUserDeleted}
				b := setup(topics...)
				defer tearDown(b, topics...)

				Expect(b.Topics()).To(ContainElements(topics))
			})
		})
	})

	Context("register", Ordered, func() {
		var (
			b      *bus.Broker
			topics []string
		)

		BeforeAll(func() {
			topics = []string{topicUserCreated, topicUserDeleted}
		})

		BeforeEach(func() {
			b = setup(topics...)
		})

		AfterEach(func() {
			tearDown(b, topics...)
		})

		When("register topics", func() {
			It("🧪 should: ", func() {
				b.RegisterTopics(topics...)
				Expect(b.Topics()).To(ContainElements(topics))
			})
		})

		When("topic already registered", func() {
			It("🧪 should: not register a topic twice", func() {
				Expect(b.Topics()).To(HaveLen(2))
				b.RegisterTopics(topicUserCreated)
				Expect(b.Topics()).To(HaveLen(2))
				Expect(b.Topics()).To(ContainElements(topics))
			})
		})
	})

	When("deregister", func() {
		It("🧪 should: ", func() {
			topics := []string{topicUserCreated, topicUserDeleted, topicUserUpdated}
			b := setup(topics...)
			defer tearDown(b)

			b.RegisterTopics(topics...)
			b.DeregisterTopics(topicUserCreated, topicUserUpdated)
			Expect(b.Topics()).To(ContainElements([]string{topicUserDeleted}))
		})
	})

	Context("handler", func() {
		When("keys", func() {
			It("🧪 should: ", func() {
				b := setup(topicCommentCreated, topicCommentDeleted)
				defer tearDown(b, topicCommentCreated, topicCommentDeleted)
				defer b.DeregisterHandler("test.key.1")
				defer b.DeregisterHandler("test.key.2")

				h := fakeHandler(".*")
				b.RegisterHandler("test.key.1", h)
				b.RegisterHandler("test.key.2", h)

				expected := []string{"test.key.1", "test.key.2"}
				Expect(b.HandlerKeys()).To(ContainElements(expected))
			})
		})

		Context("HandlerTopicSubscriptions", Ordered, func() {
			DescribeTable("look up key",
				func(entry *handlerTE) {
					b := setup(topicCommentCreated, topicCommentDeleted)
					defer tearDown(b, topicCommentCreated, topicCommentDeleted)

					b.RegisterHandler(entry.handlerKey, entry.handler)

					Expect(b.HandlerTopicSubscriptions(entry.handlerLookupKey)).To(
						ContainElements(entry.expected),
					)
				},

				func(entry *handlerTE) string {
					return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'", entry.given, entry.should)
				},

				Entry(nil, &handlerTE{
					given:            "handler defined for registered topic",
					should:           "return subscriptions for handler",
					handler:          fakeHandler(".*"),
					handlerKey:       "test.handler.1",
					handlerLookupKey: "test.handler.1",
					expected:         []string{topicCommentCreated, topicCommentDeleted},
				}),

				Entry(nil, &handlerTE{
					given:            "handler defined for un-registered topic",
					should:           "return no subscriptions",
					handler:          fakeHandler(topicUserUpdated),
					handlerKey:       "test.handler.2",
					handlerLookupKey: "test.handler.2",
					expected:         []string{},
				}),

				Entry(nil, &handlerTE{
					given:            "lookup key requested does not exist",
					should:           "return no subscriptions",
					handler:          fakeHandler(".*"),
					handlerKey:       "test.handler.3",
					handlerLookupKey: "test.handler.NA",
					expected:         []string{},
				}),
			)
		})

		Context("Register Handler", Ordered, func() {
			var (
				b      *bus.Broker
				topics []string
			)

			BeforeAll(func() {
				topics = []string{topicCommentCreated, topicCommentDeleted}
				b = setup(topics...)
				h := fakeHandler(".*created$")
				b.RegisterHandler("test.handler", h)
			})

			AfterAll(func() {
				tearDown(b, topics...)
				b.DeregisterHandler("test.handler")
			})

			When("tbd", func() {
				It("🧪 should: register handler key", func() {
					Expect(isHandlerKeyExists(b, "test.handler")).To(BeTrue())
				})
			})

			When("topic is matched", func() {
				It("🧪 should: add handler references to the matched topics", func() {
					Expect(isTopicHandler(b, topicCommentCreated, "test.handler")).To(BeTrue())
				})
			})

			When("when topic is not matched", func() {
				It("🧪 should: add handler references to the matched topics", func() {
					Expect(isTopicHandler(b, topicCommentDeleted, "test.handler")).To(BeFalse())
				})
			})
		})
	})
})

func setup(topicNames ...string) *bus.Broker {
	var fn bus.Next = func() string { return fake }

	b, _ := bus.New(fn)
	b.RegisterTopics(topicNames...)

	return b
}

func tearDown(b *bus.Broker, topicNames ...string) {
	b.DeregisterTopics(topicNames...)
}

func fakeHandler(matcher string) bus.Handler {
	return bus.Handler{
		Handle: func(context.Context, bus.Message) {}, Matcher: matcher,
	}
}

func registerFakeHandler(b *bus.Broker, key string) {
	fn := func(_ context.Context, e bus.Message) {
		Expect(e.ID).To(Equal(fake))
		Expect(e.Topic).To(Equal(topicCommentCreated))
		Expect(e.Data).To(Equal("my comment with handler"))
		Expect(e.OccurredAt.Before(time.Now())).To(BeTrue())
	}
	h := bus.Handler{Handle: fn, Matcher: ".*created$"}
	b.RegisterHandler(key, h)
}

func isTopicHandler(b *bus.Broker, topicName, handlerKey string) bool {
	for _, th := range b.TopicHandlerKeys(topicName) {
		if handlerKey == th {
			return true
		}
	}

	return false
}

func isHandlerKeyExists(b *bus.Broker, key string) bool {
	for _, k := range b.HandlerKeys() {
		if k == key {
			return true
		}
	}

	return false
}
