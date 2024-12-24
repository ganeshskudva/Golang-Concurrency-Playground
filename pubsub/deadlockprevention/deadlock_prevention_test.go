package pubsub

import (
	"fmt"
	"sync"
	"testing"
)

func TestPubSubDeadlockPrevention(t *testing.T) {
	ps := NewPubSub[string]()

	const numPublishers = 5
	const numSubscribers = 10
	const numMessages = 100

	// Subscribe multiple subscribers
	subscribers := make([]chan string, numSubscribers)
	for i := 0; i < numSubscribers; i++ {
		subscribers[i] = ps.Subscribe("test-topic")
	}

	// Goroutines to process messages for each subscriber
	var subscriberWg sync.WaitGroup
	for _, ch := range subscribers {
		subscriberWg.Add(1)
		go func(c chan string) {
			defer subscriberWg.Done()
			for msg := range c {
				_ = msg // Simulate message processing
			}
		}(ch)
	}

	// Goroutines to publish messages
	var publisherWg sync.WaitGroup
	for i := 0; i < numPublishers; i++ {
		publisherWg.Add(1)
		go func(id int) {
			defer publisherWg.Done()
			for j := 0; j < numMessages; j++ {
				ps.Publish("test-topic", fmt.Sprintf("Publisher %d: Message %d", id, j))
			}
		}(i)
	}

	// Wait for publishers to finish
	publisherWg.Wait()

	// Unsubscribe all subscribers
	for _, ch := range subscribers {
		ps.Unsubscribe("test-topic", ch)
	}

	// Shutdown the PubSub system
	ps.Shutdown()

	// Wait for subscribers to finish
	subscriberWg.Wait()

	fmt.Println("Test completed without deadlocks.")
}
