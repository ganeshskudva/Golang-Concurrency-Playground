package pubsub

import (
	"fmt"
	"sync"
)

// PubSub manages publishers and subscribers for any message type
type PubSub[T any] struct {
	subscribers map[string]map[chan T]struct{}
	mu          sync.RWMutex
}

// NewPubSub initializes a new generic PubSub
func NewPubSub[T any]() *PubSub[T] {
	return &PubSub[T]{
		subscribers: make(map[string]map[chan T]struct{}),
	}
}

// Subscribe adds a subscriber to a specific topic
func (ps *PubSub[T]) Subscribe(topic string) chan T {
	ch := make(chan T, 10) // Buffered channel
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.subscribers[topic] == nil {
		ps.subscribers[topic] = make(map[chan T]struct{})
	}
	ps.subscribers[topic][ch] = struct{}{}
	return ch
}

// Publish sends a message to all subscribers of a topic
func (ps *PubSub[T]) Publish(topic string, message T) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for ch := range ps.subscribers[topic] {
		select {
		case ch <- message:
		default: // Avoid blocking if subscriber is slow
			fmt.Println("Subscriber is too slow. Dropping message:", message)
		}
	}
}

// Unsubscribe removes a subscriber from a specific topic
func (ps *PubSub[T]) Unsubscribe(topic string, ch chan T) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if subscribers, ok := ps.subscribers[topic]; ok {
		if _, exists := subscribers[ch]; exists {
			delete(subscribers, ch)
			close(ch) // Close channel to signal subscriber
		}
		if len(subscribers) == 0 {
			delete(ps.subscribers, topic)
		}
	}
}

// Shutdown gracefully shuts down PubSub by closing all channels
func (ps *PubSub[T]) Shutdown() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for topic, subscribers := range ps.subscribers {
		for ch := range subscribers {
			close(ch)
		}
		delete(ps.subscribers, topic)
	}
}
