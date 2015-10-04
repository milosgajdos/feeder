package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/milosgajdos83/feeder/reader"
)

func main() {
	// registers signal handler
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
	reader, err := reader.NewReader()
	if err != nil {
		fmt.Println("Error creating reader:", err)
		os.Exit(1)
	}
	// Signal handler
	go func() {
		// Wait for a SIGINT or SIGKILL:
		sig := <-sigc
		fmt.Println("Shutting down reader. Got signal:", sig)
		// Stop listening (and unlink the socket if unix type):
		err := reader.Close()
		if err != nil {
			fmt.Println("Error closing reader: ", err)
		}
		os.Exit(1)
	}()
	// Goal.com RSS feed
	if err := reader.Subscribe("http://www.goal.com/en-gb/feeds/news?fmt=rss"); err != nil {
		fmt.Println("Error subscribing:", err)
		fmt.Println("Closing reader:", reader.Close())
		os.Exit(1)
	}
	// BBC RSS feed
	if err := reader.Subscribe("http://feeds.bbci.co.uk/sport/0/football/rss.xml?edition=uk"); err != nil {
		fmt.Println("Error subscribing:", err)
		fmt.Println("Closing reader:", reader.Close())
		os.Exit(1)
	}

	for it := range reader.Updates() {
		fmt.Println(it.Channel, it.Title)
	}
}
