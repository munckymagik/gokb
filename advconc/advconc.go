package advconc

import (
	"fmt"
	"time"

	"github.com/munckymagik/gokb/advconc/rss"
)

type Fetcher interface {
	Fetch() (items []rss.Item, next time.Time, err error)
}

func Fetch(domain string) Fetcher {} // fetches Items from domain

type Subscription interface {
	Updates() <-chan rss.Item // a stream of items
	Close() error             // shuts down the stream
}

func Subscribe(fetcher Fetcher) Subscription {
	s := &selectsub{
		fetcher: fetcher,
		updates: make(chan rss.Item),
	}
	go s.loop()
	return s
}

func Merge(subs ...Subscription) Subscription {}

func main() {
	merged := Merge(
		Subscribe(Fetch("blog.golang.org")),
		Subscribe(Fetch("googleblog.blogspot.org")),
		Subscribe(Fetch("googledevelopers.blogspot.com")),
	)

	time.AfterFunc(3*time.Second, func() {
		fmt.Println("Closed:", merged.Close())
	})

	for it := range merged.Updates() {
		fmt.Println(it.Channel, it.Title)
	}

	panic("show me stacks")
}
