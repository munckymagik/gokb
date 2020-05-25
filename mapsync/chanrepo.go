package mapsync

import (
	"log"
	"time"
)

type request struct {
	key   string
	reply chan *Record
}

type chanrepo struct {
	requests chan<- request
}

func NewChanRepo() Repo {
	requests := make(chan request)
	repo := chanrepo{requests}
	innerRepo := innerchanrepo{make(RecordMap)}
	innerRepo.updateDB(1)

	go innerRepo.updateLoop(requests)

	return &repo
}

func (repo *chanrepo) Get(key string) *Record {
	req := request{key, make(chan *Record, 1)}
	repo.requests <- req
	return <-req.reply
}

type innerchanrepo struct {
	db RecordMap
}

func (repo *innerchanrepo) updateDB(value int) {
	db := make(RecordMap)
	db["x"] = &Record{value}
	repo.db = db
}

func (repo *innerchanrepo) updateLoop(requests <-chan request) {
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			newValue := repo.db["x"].value + 1
			repo.updateDB(newValue)
			log.Println("Updated to:", newValue)
		case req := <-requests:
			record := repo.db[req.key]
			req.reply <- record
		}
	}
}
