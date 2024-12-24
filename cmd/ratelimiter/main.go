package main

import (
	"fmt"
	"time"

	ps "github.com/ganeshskudva/Golang-Concurrency-Playground/pubsub/ratelimiter"
	"golang.org/x/time/rate"
)

func main() {
	// Create a PubSub system with a rate limit of 5 messages per second and a burst size of 10
	ps := ps.NewPubSub[string](rate.Limit(5), 10)

	// Subscribe to a topic
	sub := ps.Subscribe("news")

	// Goroutine to process messages
	go func() {
		for msg := range sub {
			fmt.Println("Received:", msg)
		}
	}()

	// Publish messages
	for i := 0; i < 20; i++ {
		ps.Publish("news", fmt.Sprintf("Message %d", i))
		time.Sleep(100 * time.Millisecond) // Simulate a delay between messages
	}

	// Shutdown the system
	ps.Unsubscribe("news", sub)
	ps.Shutdown()
}
