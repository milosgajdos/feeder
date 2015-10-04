package feed

// NewDeduper returns new Deduper that converts
// Subscription that may send duplicate Items into one that doesn't.
func NewDeduper(s Subscription) Subscription {
	d := &deduper{
		s:       s,
		updates: make(chan Item),
		closing: make(chan chan error),
	}
	go d.loop()
	return d
}

// deduper implements Subscription interface
type deduper struct {
	s       Subscription
	updates chan Item
	closing chan chan error
}

// loop() implements deduplication logic
func (d *deduper) loop() {
	in := d.s.Updates() // enable receive
	var pending Item
	var out chan Item // disable send
	seen := make(map[string]bool)
	for {
		select {
		case item := <-in:
			if !seen[item.GUID] {
				pending = item
				in = nil        // disable receive
				out = d.updates // enable send
				seen[item.GUID] = true
			}
		case out <- pending:
			in = d.s.Updates() // enable receive
			out = nil          // disable send
		case errc := <-d.closing:
			err := d.s.Close()
			errc <- err
			close(d.updates)
			return
		}
	}
}

// Closes deduplicator which in turn closes underlying Subscription
func (d *deduper) Close() error {
	errc := make(chan error)
	d.closing <- errc
	return <-errc
}

// Returns deduplicated stream of Items
func (d *deduper) Updates() <-chan Item {
	return d.updates
}
