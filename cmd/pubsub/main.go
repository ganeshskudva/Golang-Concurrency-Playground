package main

import (
	"fmt"
	"sync"

	ps "github.com/ganeshskudva/Golang-Concurrency-Playground/pubsub"
)

type NewsUpdate struct {
	Headline string
	Details  string
}

func main() {
	// PubSub for string messages
	stringPubSub := ps.NewPubSub[string]()
	// PubSub for custom NewsUpdate struct
	newsPubSub := ps.NewPubSub[NewsUpdate]()

	// Subscribe to topics
	stringSub := stringPubSub.Subscribe("general")
	newsSub := newsPubSub.Subscribe("news")

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Goroutine to process string messages
	go func() {
		defer wg.Done()
		for msg := range stringSub {
			fmt.Println("String Subscriber received:", msg)
		}
	}()

	// Goroutine to process NewsUpdate messages
	go func() {
		defer wg.Done()
		for news := range newsSub {
			fmt.Printf("News Subscriber received: Headline - %s, Details - %s\n", news.Headline, news.Details)
		}
	}()

	// Publish heterogeneous messages
	stringPubSub.Publish("general", "Hello, World!")
	stringPubSub.Publish("general", "Go generics are powerful!")

	newsPubSub.Publish("news", NewsUpdate{
		Headline: "Breaking News",
		Details:  "Go supports generics from version 1.18",
	})
	newsPubSub.Publish("news", NewsUpdate{
		Headline: "Another Update",
		Details:  "PubSub system now supports heterogeneous types.",
	})

	// Unsubscribe and shutdown
	stringPubSub.Unsubscribe("general", stringSub)
	newsPubSub.Unsubscribe("news", newsSub)

	stringPubSub.Shutdown()
	newsPubSub.Shutdown()

	wg.Wait()
	fmt.Println("PubSub system shut down gracefully.")
}
