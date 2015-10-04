package reader

import (
	"fmt"
	"time"

	rss "github.com/jteeuwen/go-pkg-rss"
)

// An Item is a stripped-down RSS item.
type Item struct {
	Channel string
	Title   string
	GUID    string
}

// A Fetcher fetches Items and returns the time when the next fetch should be
// attempted.  On failure, Fetch returns a non-nil error.
type Fetcher interface {
	Fetch() (items []Item, next time.Time, err error)
}

// Fetch returns a Fetcher for Items from domain.
func Fetch(domain string) Fetcher {
	return NewFetcher(fmt.Sprintf("http://%s", domain))
}

// fetcher implements Fetcher interface
type fetcher struct {
	uri   string
	feed  *rss.Feed
	items []Item
}

// NewFetcher returns a Fetcher for uri.
func NewFetcher(uri string) Fetcher {
	f := &fetcher{
		uri: uri,
	}
	newChans := func(feed *rss.Feed, chans []*rss.Channel) {}
	newItems := func(feed *rss.Feed, ch *rss.Channel, items []*rss.Item) {
		for _, item := range items {
			// Initialize guid to Atom Id
			// item.Id is empty string if item comes from RSS feed
			guid := item.Id
			// if item comes from RSS feed Guid must be provided
			if item.Guid != nil {
				guid = *(item.Guid)
			}
			f.items = append(f.items, Item{
				Channel: ch.Title,
				Title:   item.Title,
				GUID:    guid,
			})
		}
	}
	// min interval 1min, respect limits
	f.feed = rss.New(1, true, newChans, newItems)
	return f
}

// Fetch fetches articles from provided feed
func (f *fetcher) Fetch() (items []Item, next time.Time, err error) {
	fmt.Println("fetching", f.uri)
	if err = f.feed.Fetch(f.uri, nil); err != nil {
		return
	}
	items = f.items
	f.items = nil
	next = time.Now().Add(time.Duration(f.feed.SecondsTillUpdate()) * time.Second)
	return
}
