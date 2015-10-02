package reader

import (
	"sync"

	"github.com/milosgajdos83/feeder/feed"
)

// Database stores subscriptions in memory
// and provides a basic operations on top of it
type Database interface {
	Insert(feed.Subscription) error
	Delete(string) error
	Find(string) (feed.Subscription, error)
	Close() error
}

// database implements Database interface
type database struct {
	// in-memory subscription data store
	subs map[string]feed.Subscription
	mu   sync.RWMutex
}

// NewDatabase returns new database or error
func NewDatabase() (Database, error) {
	subs := make(map[string]feed.Subscription)
	return &database{
		subs: subs,
	}, nil
}

// Insert adds a new sibscription to database
func (db *database) Insert(s feed.Subscription) error {
	return nil
}

// Delete removes a susbcription from database
func (db *database) Delete(uri string) error {
	return nil
}

// Finds searches for a given subscription in databases and returns is
// if it can't find it in the database it returns error
func (db *database) Find(uri string) (feed.Subscription, error) {
	return nil, nil
}

func (db *database) Close() error {
	return nil
}
