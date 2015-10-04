package reader

import (
	"fmt"
	"sync"

	"github.com/milosgajdos83/feeder/feed"
)

const (
	MaxCache = 100
)

// Cache stores subscriptions in memory keyed by uri
// and provides a basic operations on top of it
type Cache interface {
	Insert(string, feed.Subscription) error
	Delete(string) error
	Find(string) (feed.Subscription, error)
	Close() map[string]feed.Subscription
}

// cache implements Cache interface
type cache struct {
	subs   map[string]feed.Subscription
	closed bool
	mu     sync.RWMutex
}

// Newcache returns new cache or error
func NewCache() (Cache, error) {
	subs := make(map[string]feed.Subscription)
	return &cache{
		subs: subs,
	}, nil
}

// Insert adds a new sibscription to cache
func (c *cache) Insert(key string, s feed.Subscription) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return fmt.Errorf("Can not store %s: Cache read only!", key)
	}

	if len(c.subs)+1 > MaxCache {
		return fmt.Errorf("Can not store %s: MaxCache hit!", key)
	}

	if _, ok := c.subs[key]; ok {
		return fmt.Errorf("Can not write %s: Item already exists", key)
	}
	c.subs[key] = s
	return nil
}

// Delete removes a susbcription from cache
func (c *cache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return fmt.Errorf("Can not delete %s: Cache read only!", key)
	}
	delete(c.subs, key)
	return nil
}

// Finds searches for a given subscription in caches and returns is
// if it can't find it in the cache it returns error
func (c *cache) Find(uri string) (feed.Subscription, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	s, ok := c.subs[uri]
	if !ok {
		return nil, fmt.Errorf("Could not find %s", uri)
	}
	return s, nil
}

// Closes the cache and returns underlying store
func (c *cache) Close() map[string]feed.Subscription {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
	return c.subs
}
