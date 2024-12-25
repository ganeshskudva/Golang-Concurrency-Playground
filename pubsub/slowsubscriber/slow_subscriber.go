package pubsub

import (
	"fmt"
	"sync"
)

// PubSub manages publishers and subscribers for any message type
type PubSub[T any] struct {
	subscribers map[string]map[chan T]struct{} // Map of topics to a set of subscriber channels
	mu          sync.RWMutex                   // Read-Write lock to manage concurrent access
}

// NewPubSub initializes a new PubSub instance for a specific type
func NewPubSub[T any]() *PubSub[T] {
	return &PubSub[T]{
		subscribers: make(map[string]map[chan T]struct{}),
	}
}

// Subscribe adds a new subscriber to a specific topic
func (ps *PubSub[T]) Subscribe(topic string) chan T {
	ch := make(chan T, 100) // Buffered channel for slow subscriber handling
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.subscribers[topic] == nil {
		ps.subscribers[topic] = make(map[chan T]struct{})
	}
	ps.subscribers[topic][ch] = struct{}{}
	return ch
}

// Publish sends a message to all subscribers of a topic
// Handles slow subscribers by dropping messages if the buffer is full
func (ps *PubSub[T]) Publish(topic string, message T) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for ch := range ps.subscribers[topic] {
		select {
		case ch <- message: // Deliver message
		default: // Drop message if channel is full
			fmt.Printf("Dropping message for a slow subscriber: %v\n", message)
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

// Shutdown gracefully shuts down the PubSub system by closing all channels
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
