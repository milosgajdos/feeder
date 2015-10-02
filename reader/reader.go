package reader

import "github.com/milosgajdos83/feeder/feed"

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
	db Database
}

// NewReader creates new RSS reader or returns error
func NewReader() (Reader, error) {
	db, err := NewDatabase()
	if err != nil {
		return nil, err
	}
	return &reader{
		db: db,
	}, nil
}

// Subscribe adds a new RSS subscription or returns error if it cant be added
func (r *reader) Subscribe(uri string) error {
	return nil
}

// Unsubscribe removes an existing RSS subscription or fails with error
func (r *reader) Unsubscribe(uri string) error {
	return nil
}

// Updates returns deduplicated merged Item stream
func (r *reader) Updates() <-chan feed.Item {
	return nil
}

// Close stops all existing subscriptions and returns error if any
func (r *reader) Close() error {
	return nil
}
