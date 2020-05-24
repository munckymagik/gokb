package advconc

import (
	"time"

	"github.com/munckymagik/gokb/advconc/rss"
)

type naivesub struct {
	fetcher Fetcher
	updates chan rss.Item
	closed  bool
	err     error
}

func (s *naivesub) loop() {
	for {
		// BUG unsync access to s.closed
		if s.closed {
			close(s.updates)
			return
		}
		items, next, err := s.fetcher.Fetch()
		if err != nil {
			// BUG unsync access to s.err
			s.err = err
			// BUG may keep loop running
			time.Sleep(10 * time.Second)
			continue
		}
		for _, item := range items {
			// BUG may block forever
			s.updates <- item
		}
		if now := time.Now(); next.After(now) {
			// BUG may keep loop running
			time.Sleep((next.Sub(now)))
		}
	}
}

func (s *naivesub) Updates() <-chan rss.Item {
	return s.updates
}

func (s *naivesub) Close() error {
	s.closed = true
	return s.err
}
