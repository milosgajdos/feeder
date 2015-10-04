package reader

import (
	"fmt"
	"sync"

	"github.com/milosgajdos83/feeder/feed"
)

// Reader adds or removes subscriptions and provides updates stream
// Close closes updates stream and returns error if any
type Reader interface {
	Subscribe(string) error
	Unsubscribe(string) error
	Updates() <-chan feed.Item
	Close() error
}

// reader implements Reader interface
type reader struct {
	cache   Cache
	mu      sync.RWMutex
	subs    chan feed.Subscription
	updates chan feed.Item
	quit    chan struct{}
	errs    chan error
}

// NewReader creates new RSS reader or returns error
func NewReader() (Reader, error) {
	cache, err := NewCache()
	if err != nil {
		return nil, err
	}
	r := &reader{
		cache:   cache,
		subs:    make(chan feed.Subscription),
		updates: make(chan feed.Item),
		quit:    make(chan struct{}),
		errs:    make(chan error),
	}
	go r.run()
	return r, nil
}

func (r *reader) run() {
	for sub := range r.subs {
		go func(s feed.Subscription) {
			for {
				var it feed.Item
				select {
				// receive Items from particular Subscription
				case it = <-s.Updates():
				case <-r.quit:
					r.errs <- s.Close()
					return
				}
				select {
				// write to the merge stream of Items
				case r.updates <- it:
				case <-r.quit:
					r.errs <- s.Close()
					return
				}
			}
		}(sub)
	}
}

// Subscribe adds a new RSS subscription or returns error if it cant be added
func (r *reader) Subscribe(uri string) error {
	s := feed.NewDeduper(feed.NewSubscription(feed.NewFetcher(uri)))
	if err := r.cache.Insert(uri, s); err != nil {
		s.Close()
		return err
	}
	r.subs <- s
	return nil
}

// Unsubscribe removes an existing RSS subscription or fails with error
func (r *reader) Unsubscribe(uri string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, err := r.cache.Find(uri)
	if err != nil {
		return err
	}
	if err := s.Close(); err != nil {
		fmt.Printf("Could not stop subscription: %s", err)
	}
	return r.cache.Delete(uri)
}

// Updates returns deduplicated merged Item stream
func (r *reader) Updates() <-chan feed.Item {
	return r.updates
}

// Close stops all existing subscriptions and returns error if any
func (r *reader) Close() (err error) {
	close(r.quit)
	subs := r.cache.Close()
	for range subs {
		if e := <-r.errs; e != nil {
			// Collects return values of terminated
			// subscriptions which unblocks goroutines in r.run()
			err = e
		}
	}
	close(r.subs)
	close(r.updates)
	return
}
