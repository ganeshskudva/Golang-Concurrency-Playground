package pubsub

import (
	"fmt"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestRateLimiter(t *testing.T) {
	// Set up the PubSub with a rate limit of 2 messages per second and a burst of 2
	ps := NewPubSub[string](rate.Limit(2), 2)

	// Subscribe to a topic
	sub := ps.Subscribe("test-topic")

	// Channel to collect received messages
	received := make(chan string, 10)

	// Goroutine to process messages
	go func() {
		for msg := range sub {
			received <- msg
		}
	}()

	// Publish messages faster than the rate limit
	for i := 0; i < 10; i++ {
		ps.Publish("test-topic", fmt.Sprintf("Message %d", i))
		time.Sleep(100 * time.Millisecond) // Publish every 100ms (10 messages/sec)
	}

	// Allow time for the rate limiter to apply
	time.Sleep(3 * time.Second)

	// Check received messages
	close(received)

	// Count the received messages
	count := 0
	for range received {
		count++
	}

	// Expected number of messages received should not exceed the rate limit
	// With a 2 messages/sec limit and 3 seconds, the max messages should be 6
	if count > 6 {
		t.Errorf("Rate limiter failed: received %d messages, expected <= 6", count)
	} else {
		t.Logf("Rate limiter passed: received %d messages within limit", count)
	}

	// Cleanup
	ps.Unsubscribe("test-topic", sub)
	ps.Shutdown()
}
