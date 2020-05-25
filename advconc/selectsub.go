package main

import (
	"time"

	"github.com/munckymagik/gokb/advconc/rss"
)

type selectsub struct {
	fetcher Fetcher
	updates chan rss.Item
	closing chan chan error
}

func (s *selectsub) loop() {
	var pending []rss.Item
	var next time.Time
	var err error

	for {
		var first rss.Item
		var updates chan rss.Item
		if len(pending) > 0 {
			first = pending[0]
			// enables send by making updates non-nil
			updates = s.updates
		}

		var fetchDelay time.Duration
		if now := time.Now(); next.After(now) {
			fetchDelay = next.Sub(now)
		}
		startFetch := time.After(fetchDelay)

		select {
		case errc := <-s.closing:
			errc <- err
			close(s.updates)
			return
		case <-startFetch:
			var fetched []rss.Item
			fetched, next, err = s.fetcher.Fetch()
			if err != nil {
				next = time.Now().Add(10 * time.Second)
				break
			}
			pending = append(pending, fetched...)
		case updates <- first:
			// updates may be nil in which case select ignores it
			pending = pending[1:]
		}
	}
}

func (s *selectsub) Updates() <-chan rss.Item {
	return s.updates
}

func (s *selectsub) Close() error {
	errChannel := make(chan error)
	s.closing <- errChannel
	return <-errChannel
}
