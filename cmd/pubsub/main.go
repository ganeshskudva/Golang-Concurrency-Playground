package main

import (
	"fmt"
	"sync"

	ps "github.com/ganeshskudva/Golang-Concurrency-Playground/pubsub"
)

func main() {
	ps := ps.NewPubSub()

	// Subscribe to a topic
	sub1 := ps.Subscribe("news")
	sub2 := ps.Subscribe("news")

	// Read messages from subscribers
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for msg := range sub1 {
			fmt.Println("Subscriber 1 received:", msg)
		}
	}()

	go func() {
		defer wg.Done()
		for msg := range sub2 {
			fmt.Println("Subscriber 2 received:", msg)
		}
	}()

	// Publish messages
	ps.Publish("news", "Breaking News: Go is awesome!")
	ps.Publish("news", "Another update!")

	// Unsubscribe after some delay
	ps.Unsubscribe("news", sub1)
	ps.Unsubscribe("news", sub2)

	// Shutdown and wait for goroutines to finish
	ps.Shutdown()
	wg.Wait()
	fmt.Println("PubSub system shut down gracefully.")
}
