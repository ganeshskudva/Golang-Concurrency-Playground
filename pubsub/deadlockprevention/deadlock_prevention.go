package pubsub

import (
	"fmt"
	"sync"
)

// PubSub manages publishers and subscribers for any message type
type PubSub[T any] struct {
	subscribers map[string]map[chan T]struct{} // Map of topics to subscribers
	mu          sync.RWMutex                   // Read-Write lock for synchronizing access
}

// NewPubSub initializes a new PubSub instance
func NewPubSub[T any]() *PubSub[T] {
	return &PubSub[T]{
		subscribers: make(map[string]map[chan T]struct{}),
	}
}

// Subscribe adds a subscriber to a specific topic
// Returns a channel through which the subscriber will receive messages
func (ps *PubSub[T]) Subscribe(topic string) chan T {
	ch := make(chan T, 100) // Buffered channel to prevent blocking
	ps.mu.Lock()
	defer ps.mu.Unlock()

	// Initialize the topic if it doesn't exist
	if ps.subscribers[topic] == nil {
		ps.subscribers[topic] = make(map[chan T]struct{})
	}

	// Add the subscriber channel to the topic
	ps.subscribers[topic][ch] = struct{}{}
	return ch
}

// Publish sends a message to all subscribers of a topic
// Ensures that message delivery does not hold locks for an extended duration
func (ps *PubSub[T]) Publish(topic string, message T) {
	ps.mu.RLock()
	subscriberChannels := ps.getSubscriberChannels(topic) // Snapshot of subscribers
	ps.mu.RUnlock()                                       // Release lock early

	var wg sync.WaitGroup
	for _, ch := range subscriberChannels {
		wg.Add(1)
		go func(c chan T) {
			defer wg.Done()
			select {
			case c <- message: // Deliver message
			default: // Drop message if subscriber channel is full
				fmt.Println("Subscriber is too slow. Dropping message.")
			}
		}(ch)
	}
	wg.Wait() // Wait for all deliveries to complete
}

// getSubscriberChannels safely retrieves a snapshot of subscriber channels for a topic
// This avoids holding locks during message delivery
func (ps *PubSub[T]) getSubscriberChannels(topic string) []chan T {
	subscribers, exists := ps.subscribers[topic]
	if !exists {
		return nil
	}

	channels := make([]chan T, 0, len(subscribers))
	for ch := range subscribers {
		channels = append(channels, ch)
	}
	return channels
}

// Unsubscribe removes a subscriber from a specific topic
// Closes the channel to signal the subscriber that no more messages will be sent
func (ps *PubSub[T]) Unsubscribe(topic string, ch chan T) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if subscribers, ok := ps.subscribers[topic]; ok {
		if _, exists := subscribers[ch]; exists {
			delete(subscribers, ch)
			close(ch) // Close the channel to clean up resources
		}

		// Remove the topic if no subscribers remain
		if len(subscribers) == 0 {
			delete(ps.subscribers, topic)
		}
	}
}

// Shutdown gracefully shuts down the PubSub system
// Closes all channels and cleans up the internal data structure
func (ps *PubSub[T]) Shutdown() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	// Iterate over all topics and their subscribers
	for topic, subscribers := range ps.subscribers {
		for ch := range subscribers {
			close(ch) // Close each subscriber channel
		}
		delete(ps.subscribers, topic) // Remove the topic
	}
}
