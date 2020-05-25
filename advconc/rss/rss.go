package rss

import "time"

func Fetch(uri string) (items []Item, next time.Time, err error) {
	items = append(
		items,
		Item{"title1", "chan1", "GUID1"},
		Item{"title2", "chan2", "GUID2"},
	)
	next = time.Now().Add(time.Second * 3)
	err = nil
	return
}

type Item struct {
	Title, Channel, GUID string // a subset of RSS fields
}
