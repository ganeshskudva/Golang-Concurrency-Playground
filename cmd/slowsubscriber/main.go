package main

import (
	"fmt"
	"time"

	ps "github.com/ganeshskudva/Golang-Concurrency-Playground/pubsub/slowsubscriber"
)

func main() {
	ps := ps.NewPubSub[string]()

	// Fast subscriber
	fastSub := ps.Subscribe("news")
	go func() {
		for msg := range fastSub {
			fmt.Println("Fast Subscriber received:", msg)
		}
	}()

	// Slow subscriber
	slowSub := ps.Subscribe("news")
	go func() {
		for msg := range slowSub {
			fmt.Println("Slow Subscriber received:", msg)
			time.Sleep(200 * time.Millisecond) // Simulate slow processing
		}
	}()

	// Publish messages
	for i := 0; i < 20; i++ {
		ps.Publish("news", fmt.Sprintf("Message %d", i))
		time.Sleep(50 * time.Millisecond) // Simulate publishing rate
	}

	// Unsubscribe and shutdown
	ps.Unsubscribe("news", fastSub)
	ps.Unsubscribe("news", slowSub)
	ps.Shutdown()
}
