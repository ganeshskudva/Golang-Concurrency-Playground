package pubsub

import (
	"fmt"
	"sync"
	"testing"
)

func TestHighThroughputPubSub(t *testing.T) {
	ps := NewPubSub[string]()

	const numPublishers = 10
	const numSubscribers = 10000
	const numMessages = 1000

	// Subscribe multiple subscribers
	subscribers := make([]chan string, numSubscribers)
	for i := 0; i < numSubscribers; i++ {
		subscribers[i] = ps.Subscribe("high-throughput")
	}

	// Goroutines to process subscriber messages
	var wg sync.WaitGroup
	for i, ch := range subscribers {
		wg.Add(1)
		go func(id int, ch chan string) {
			defer wg.Done()
			for msg := range ch {
				_ = msg // Process message (e.g., print/log in real scenarios)
			}
		}(i, ch)
	}

	// Start multiple publishers
	var publisherWg sync.WaitGroup
	for i := 0; i < numPublishers; i++ {
		publisherWg.Add(1)
		go func(id int) {
			defer publisherWg.Done()
			for j := 0; j < numMessages; j++ {
				ps.Publish("high-throughput", fmt.Sprintf("Publisher %d: Message %d", id, j))
			}
		}(i)
	}

	// Wait for all publishers to finish
	publisherWg.Wait()

	// Unsubscribe and shutdown
	for _, ch := range subscribers {
		ps.Unsubscribe("high-throughput", ch)
	}
	ps.Shutdown()

	// Wait for all subscriber goroutines to finish
	wg.Wait()
}
