package pubsub

import (
	"fmt"
	"sync"

	"golang.org/x/time/rate"
)

// PubSub manages publishers and subscribers for any message type
type PubSub[T any] struct {
	subscribers map[string]map[chan T]struct{} // Map of topics to a set of subscriber channels
	mu          sync.RWMutex                   // Read-Write lock to manage concurrent access
	limiter     *rate.Limiter                  // Rate limiter for publishers
}

// NewPubSub initializes a new PubSub instance for a specific type with a rate limit.
// limit: maximum number of messages allowed per second
// burst: maximum burst size
func NewPubSub[T any](limit rate.Limit, burst int) *PubSub[T] {
	return &PubSub[T]{
		subscribers: make(map[string]map[chan T]struct{}),
		limiter:     rate.NewLimiter(limit, burst),
	}
}

// Subscribe adds a new subscriber to a specific topic.
// Returns a channel through which the subscriber will receive messages.
func (ps *PubSub[T]) Subscribe(topic string) chan T {
	// Create a buffered channel to prevent blocking during message delivery
	ch := make(chan T, 100) // Buffered channel size is set to 100 for high throughput
	ps.mu.Lock()
	defer ps.mu.Unlock()

	// Initialize the topic in the map if it doesn't exist
	if ps.subscribers[topic] == nil {
		ps.subscribers[topic] = make(map[chan T]struct{})
	}

	// Add the subscriber channel to the topic
	ps.subscribers[topic][ch] = struct{}{}
	return ch
}

// Publish sends a message to all subscribers of a given topic.
// Ensures rate limiting for publishers.
func (ps *PubSub[T]) Publish(topic string, message T) {
	// Enforce rate limiting
	if !ps.limiter.Allow() {
		fmt.Println("Rate limit exceeded. Dropping message:", message)
		return
	}

	ps.mu.RLock() // Acquire read lock to allow concurrent publishing
	defer ps.mu.RUnlock()

	var wg sync.WaitGroup // WaitGroup to ensure all goroutines finish

	// Iterate over all subscriber channels for the topic
	for ch := range ps.subscribers[topic] {
		wg.Add(1)
		// Deliver message to each subscriber in a separate goroutine
		go func(c chan T) {
			defer wg.Done()
			// Send the message or drop it if the channel is full
			select {
			case c <- message:
				// Message successfully delivered
			default:
				// Channel is full; drop the message to avoid blocking
				fmt.Println("Subscriber is too slow. Dropping message.")
			}
		}(ch)
	}

	// Wait for all goroutines to complete
	wg.Wait()
}

// Unsubscribe removes a subscriber from a specific topic.
// The channel is closed to signal the subscriber that no more messages will be sent.
func (ps *PubSub[T]) Unsubscribe(topic string, ch chan T) {
	ps.mu.Lock() // Acquire write lock to modify the subscriber map
	defer ps.mu.Unlock()

	// Check if the topic exists
	if subscribers, ok := ps.subscribers[topic]; ok {
		// Remove the subscriber channel if it exists
		if _, exists := subscribers[ch]; exists {
			delete(subscribers, ch)
			close(ch) // Close the channel to clean up resources
		}
		// If no subscribers remain for the topic, remove the topic
		if len(subscribers) == 0 {
			delete(ps.subscribers, topic)
		}
	}
}

// Shutdown gracefully shuts down the PubSub system by closing all channels.
// This signals all subscribers that no more messages will be sent.
func (ps *PubSub[T]) Shutdown() {
	ps.mu.Lock() // Acquire write lock to prevent new subscriptions/publishing
	defer ps.mu.Unlock()

	// Iterate over all topics
	for topic, subscribers := range ps.subscribers {
		// Close all channels for each topic
		for ch := range subscribers {
			close(ch)
		}
		// Remove the topic from the map
		delete(ps.subscribers, topic)
	}
}
