package main

import (
	"fmt"
	"time"

	"github.com/milosgajdos83/feeder/feed"
)

func main() {
	// Subscribe to some feeds, and create a merged update stream
	merged := feed.Merge(
		feed.Dedupe(feed.Subscribe(feed.Fetch("www.goal.com/en-gb/feeds/news?fmt=rss"))))

	// Close the subscriptions after some time.
	time.AfterFunc(3*time.Second, func() {
		fmt.Println("Closing fetch")
		if err := merged.Close(); err != nil {
			fmt.Println("Error: ", err)
		}
	})

	// Print the stream.
	for it := range merged.Updates() {
		fmt.Println(it.Channel, it.Title)
	}
}
