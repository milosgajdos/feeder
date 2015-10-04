package feed

import "time"

// A Subscription delivers Items over a channel. Close cancels the
// subscription, closes the Updates channel, and returns the last fetch error
type Subscription interface {
	Updates() <-chan Item // stream of updates
	Close() error         // close subscription
}

// NewSubscription returns a new Subscription that uses fetcher to fetch Items.
func NewSubscription(fetcher Fetcher) Subscription {
	s := &sub{
		fetcher: fetcher,
		updates: make(chan Item),       // for Updates
		closing: make(chan chan error), // for Close
	}
	go s.loop()
	return s
}

// sub implements the Subscription interface.
type sub struct {
	fetcher Fetcher         // fetches items
	updates chan Item       // sends items to the user
	closing chan chan error // for Close
}

// Updates returns channel to read Items
func (s *sub) Updates() <-chan Item {
	return s.updates
}

// Close closes the Item stream and reads latest error
func (s *sub) Close() error {
	errc := make(chan error)
	s.closing <- errc // Asks subscr to close the stream and return error if any
	return <-errc
}

// Implements Subscription feed logic
// Item stream is controlled via updaes, closing and startFetch channels
func (s *sub) loop() {
	const maxPending = 10
	type fetchResult struct {
		fetched []Item
		next    time.Time
		err     error
	}
	var fetchDone chan fetchResult // if non-nil, Fetch is in progress
	var pending []Item
	var next time.Time
	var err error
	for {
		// start with no delay - i.e. fetch immediately
		// next fetch is set by fetcher.Fetch()
		var fetchDelay time.Duration
		if now := time.Now(); next.After(now) {
			fetchDelay = next.Sub(now)
		}
		// controls fetching - disable fetch (nil channel)
		var startFetch <-chan time.Time
		// if no fetch is in progress (async goroutine in startFetch case)
		// or if number of pending Items doesnt outgrow maxPending limit
		if fetchDone == nil && len(pending) < maxPending {
			// enable fetch after fetchDelay
			startFetch = time.After(fetchDelay)
		}
		var first Item
		var updates chan Item
		if len(pending) > 0 {
			first = pending[0]
			// enable send case if at least one Item has been read
			updates = s.updates
		}
		select {
		case <-startFetch:
			fetchDone = make(chan fetchResult, 1)
			go func() {
				fetched, next, err := s.fetcher.Fetch()
				fetchDone <- fetchResult{fetched, next, err}
			}()
		case result := <-fetchDone:
			fetchDone = nil
			next, err = result.next, result.err
			if err != nil {
				next = time.Now().Add(10 * time.Second)
				break
			}
			pending = append(pending, result.fetched...)
		case errc := <-s.closing:
			errc <- err
			close(s.updates)
			return
		case updates <- first:
			pending = pending[1:]
		}
	}
}
