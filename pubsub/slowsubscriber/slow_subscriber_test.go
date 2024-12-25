package pubsub

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSlowSubscriberHandling(t *testing.T) {
	ps := NewPubSub[string]()

	// Subscriber 1: Fast subscriber
	fastSub := ps.Subscribe("news")
	var fastMessages []string
	fastWg := sync.WaitGroup{}
	fastWg.Add(1)

	go func() {
		defer fastWg.Done()
		for msg := range fastSub {
			fastMessages = append(fastMessages, msg)
		}
	}()

	// Subscriber 2: Slow subscriber
	slowSub := ps.Subscribe("news")
	var slowMessages []string
	slowWg := sync.WaitGroup{}
	slowWg.Add(1)

	go func() {
		defer slowWg.Done()
		for msg := range slowSub {
			slowMessages = append(slowMessages, msg)
			time.Sleep(200 * time.Millisecond) // Simulate slow processing
		}
	}()

	// Publish messages
	for i := 0; i < 10; i++ {
		ps.Publish("news", fmt.Sprintf("Message %d", i))
		time.Sleep(50 * time.Millisecond) // Simulate publishing rate
	}

	// Unsubscribe and shutdown
	ps.Unsubscribe("news", fastSub)
	ps.Unsubscribe("news", slowSub)
	ps.Shutdown()

	// Wait for all subscribers to finish processing
	fastWg.Wait()
	slowWg.Wait()

	// Check results
	if len(fastMessages) != 10 {
		t.Errorf("Fast subscriber missed messages, received: %d, expected: 10", len(fastMessages))
	}

	if len(slowMessages) < 5 {
		t.Errorf("Slow subscriber missed too many messages, received: %d, expected at least: 5", len(slowMessages))
	} else {
		t.Logf("Slow subscriber received %d messages (some may have been dropped).", len(slowMessages))
	}
}
