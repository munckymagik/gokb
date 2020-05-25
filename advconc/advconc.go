package main

import (
	"fmt"
	"time"

	"github.com/munckymagik/gokb/advconc/rss"
)

type Fetcher interface {
	Fetch() (items []rss.Item, next time.Time, err error)
}

type fetcher struct {
	uri string
}

func Fetch(domain string) Fetcher {
	return fetcher{fmt.Sprintf("https://%s/feed.rss", domain)}
}

func (f fetcher) Fetch() (items []rss.Item, next time.Time, err error) {
	return rss.Fetch(f.uri)
}

type Subscription interface {
	Updates() <-chan rss.Item // a stream of items
	Close() error             // shuts down the stream
}

func Subscribe(fetcher Fetcher) Subscription {
	s := &naivesub{
		fetcher: fetcher,
		updates: make(chan rss.Item),
	}
	go s.loop()
	return s
}

func main() {
	subs := Subscribe(Fetch("blog.golang.org"))

	time.AfterFunc(3*time.Second, func() {
		fmt.Println("Closed:", subs.Close())
	})

	for it := range subs.Updates() {
		fmt.Println(it.Channel, it.Title)
	}

	panic("show me stacks")
}
