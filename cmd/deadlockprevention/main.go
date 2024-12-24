package main

import (
	"fmt"
	"sync"
	"time"

	ps "github.com/ganeshskudva/Golang-Concurrency-Playground/pubsub/deadlockprevention"
)

func main() {
	// Initialize the PubSub system
	pubsub := ps.NewPubSub[string]()

	const numPublishers = 3
	const numSubscribers = 5
	const numMessages = 20

	// Subscribe multiple subscribers to the "news" topic
	subscribers := make([]chan string, numSubscribers)
	for i := 0; i < numSubscribers; i++ {
		subscribers[i] = pubsub.Subscribe("news")
	}

	// Start goroutines to process messages for each subscriber
	var subscriberWg sync.WaitGroup
	for i, ch := range subscribers {
		subscriberWg.Add(1)
		go func(id int, c chan string) {
			defer subscriberWg.Done()
			for msg := range c {
				fmt.Printf("Subscriber %d received: %s\n", id, msg)
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
				message := fmt.Sprintf("Publisher %d: Message %d", id, j)
				pubsub.Publish("news", message)
				time.Sleep(50 * time.Millisecond) // Simulate publishing delay
			}
		}(i)
	}

	// Wait for all publishers to finish
	publisherWg.Wait()

	// Unsubscribe all subscribers after publishing is done
	for i, ch := range subscribers {
		fmt.Printf("Unsubscribing Subscriber %d\n", i)
		pubsub.Unsubscribe("news", ch)
	}

	// Shutdown the PubSub system
	pubsub.Shutdown()

	// Wait for all subscribers to finish processing
	subscriberWg.Wait()

	fmt.Println("PubSub system shut down gracefully.")
}
