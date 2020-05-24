package rss

import "time"

func ExtFetch(uri string) (items []Item, next time.Time, err error)

type Item struct {
	Title, Channel, GUID string // a subset of RSS fields
}
