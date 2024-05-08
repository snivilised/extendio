// Copyright 2021 Mustafa Turan. All rights reserved.
// Use of this source code is governed by a Apache License 2.0 license that can
// be found in the bus.LICENSE file.

package bus

import (
	"context"
	"fmt"
	"regexp"
	"time"
)

type (
	// Broker controls the message bus
	Broker struct {
		gen      Next
		topics   map[string][]Handler
		handlers map[string]Handler
	}

	// Next is a sequential unique id generator func type
	Next func() string

	// IDGenerator is a sequential unique id generator interface
	IDGenerator interface {
		Generate() string
	}

	// Message is data structure for any logs
	Message struct {
		ID         string      // identifier
		TxID       string      // transaction identifier
		Topic      string      // topic name
		Source     string      // source of the message
		OccurredAt time.Time   // creation time in nanoseconds
		Data       interface{} // actual message data
	}

	// Handler is a receiver for message reference with the given regex pattern
	Handler struct {
		key string

		// handler func to process messages
		Handle func(ctx context.Context, m Message)

		// topic matcher as regex pattern
		Matcher string
	}

	// MessageOption is a function type to mutate message fields
	MessageOption = func(Message) Message

	ctxKey int8

	BrokerError struct {
		message string
	}
)

const (
	// CtxKeyTxID tx id context key
	CtxKeyTxID = ctxKey(116)

	// CtxKeySource source context key
	CtxKeySource = ctxKey(117)

	empty = ""
)

// New inits a new bus
func New(g IDGenerator) (*Broker, error) {
	if g == nil {
		return nil, fmt.Errorf("bus: Next() id generator func can't be nil") // TODO: i18n
	}

	return &Broker{
		gen:      g.Generate,
		topics:   make(map[string][]Handler),
		handlers: make(map[string]Handler),
	}, nil
}

// WithID returns an option to set message's id field
func WithID(id string) MessageOption {
	return func(m Message) Message {
		m.ID = id

		return m
	}
}

// WithTxID returns an option to set message's txID field
func WithTxID(txID string) MessageOption {
	return func(m Message) Message {
		m.TxID = txID

		return m
	}
}

// WithSource returns an option to set message's source field
func WithSource(source string) MessageOption {
	return func(e Message) Message {
		e.Source = source

		return e
	}
}

// WithOccurredAt returns an option to set message's occurredAt field
func WithOccurredAt(t time.Time) MessageOption {
	return func(e Message) Message {
		e.OccurredAt = t

		return e
	}
}

// Emit inits a new message and delivers to the interested in handlers with
// sync safety
func (b *Broker) Emit(ctx context.Context,
	topic string, data interface{},
) error {
	handlers, ok := b.topics[topic]

	if !ok {
		return &BrokerError{
			message: fmt.Sprintf("topics(%s) not found", topic), // TODO: i18n
		}
	}

	source, _ := ctx.Value(CtxKeySource).(string)
	txID, _ := ctx.Value(CtxKeyTxID).(string)
	if txID == empty {
		txID = b.gen()
		ctx = context.WithValue(ctx, CtxKeyTxID, txID)
	}

	e := Message{
		ID:         b.gen(),
		Topic:      topic,
		Data:       data,
		OccurredAt: time.Now(),
		TxID:       txID,
		Source:     source,
	}

	for _, h := range handlers {
		h.Handle(ctx, e)
	}

	return nil
}

// EmitWithOpts inits a new message and delivers to the interested in handlers
// with sync safety and options
func (b *Broker) EmitWithOpts(ctx context.Context,
	topic string, data interface{}, opts ...MessageOption,
) error {
	handlers, ok := b.topics[topic]

	if !ok {
		return &BrokerError{
			message: fmt.Sprintf("topics(%s) not found", topic), // TODO: i18n
		}
	}

	e := Message{Topic: topic, Data: data}
	for _, o := range opts {
		e = o(e)
	}

	if e.TxID == empty {
		e.TxID = b.gen()
	}
	if e.ID == empty {
		e.ID = b.gen()
	}
	if e.OccurredAt.IsZero() {
		e.OccurredAt = time.Now()
	}

	for _, h := range handlers {
		h.Handle(ctx, e)
	}

	return nil
}

// Topics lists the all registered topics
func (b *Broker) Topics() []string {
	topics, index := make([]string, len(b.topics)), 0

	for topic := range b.topics {
		topics[index] = topic
		index++
	}

	return topics
}

// RegisterTopics registers topics and fulfills handlers
func (b *Broker) RegisterTopics(topics ...string) {
	for _, n := range topics {
		b.registerTopic(n)
	}
}

// DeregisterTopics deletes topic
func (b *Broker) DeregisterTopics(topics ...string) {
	for _, n := range topics {
		b.deregisterTopic(n)
	}
}

// TopicHandlerKeys returns all handlers for the topic
func (b *Broker) TopicHandlerKeys(topic string) []string {
	handlers := b.topics[topic]

	keys := make([]string, len(handlers))

	for i, h := range handlers {
		keys[i] = h.key
	}

	return keys
}

// HandlerKeys returns list of registered handler keys
func (b *Broker) HandlerKeys() []string {
	keys, index := make([]string, len(b.handlers)), 0

	for k := range b.handlers {
		keys[index] = k
		index++
	}

	return keys
}

// HandlerTopicSubscriptions returns all topic subscriptions of the handler
func (b *Broker) HandlerTopicSubscriptions(handlerKey string) []string {
	return b.handlerTopicSubscriptions(handlerKey)
}

// RegisterHandler re/register the handler to the registry
func (b *Broker) RegisterHandler(key string, h Handler) {
	h.key = key
	b.registerHandler(h)
}

// DeregisterHandler deletes handler from the registry
func (b *Broker) DeregisterHandler(key string) {
	b.deregisterHandler(key)
}

// Generate is an implementation of IDGenerator for bus.Next fn type
func (n Next) Generate() string {
	return n()
}

func (b *Broker) registerHandler(h Handler) {
	b.deregisterHandler(h.key)
	b.handlers[h.key] = h
	for _, t := range b.handlerTopicSubscriptions(h.key) {
		b.registerTopicHandler(t, h)
	}
}

func (b *Broker) deregisterHandler(handlerKey string) {
	if _, ok := b.handlers[handlerKey]; ok {
		for _, t := range b.handlerTopicSubscriptions(handlerKey) {
			b.deregisterTopicHandler(t, handlerKey)
		}
		delete(b.handlers, handlerKey)
	}
}

func (b *Broker) registerTopicHandler(topic string, h Handler) {
	b.topics[topic] = append(b.topics[topic], h)
}

func (b *Broker) deregisterTopicHandler(topic, handlerKey string) {
	l := len(b.topics[topic])
	for i, h := range b.topics[topic] {
		if h.key == handlerKey {
			b.topics[topic][i] = b.topics[topic][l-1]
			b.topics[topic] = b.topics[topic][:l-1]
			break
		}
	}
}

func (b *Broker) registerTopic(topic string) {
	if _, ok := b.topics[topic]; ok {
		return
	}

	b.topics[topic] = b.buildHandlers(topic)
}

func (b *Broker) deregisterTopic(topic string) {
	delete(b.topics, topic)
}

func (b *Broker) buildHandlers(topic string) []Handler {
	handlers := make([]Handler, 0)
	for _, h := range b.handlers {
		if matched, _ := regexp.MatchString(h.Matcher, topic); matched {
			handlers = append(handlers, h)
		}
	}

	return handlers
}

func (b *Broker) handlerTopicSubscriptions(handlerKey string) []string {
	var subscriptions []string
	h, ok := b.handlers[handlerKey]
	if !ok {
		return subscriptions
	}

	for topic := range b.topics {
		if matched, _ := regexp.MatchString(h.Matcher, topic); matched {
			subscriptions = append(subscriptions, topic)
		}
	}
	return subscriptions
}

func (e *BrokerError) Error() string {
	return e.message
}

type Sequential struct {
	Format string
	id     int
}

func (g *Sequential) Generate() string {
	g.id++

	return fmt.Sprintf(g.Format, g.id)
}
