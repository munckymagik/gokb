package mapsync

import (
	"log"
	"sync/atomic"
	"time"
)

type atomicrepo struct {
	db atomic.Value
}

func NewAtomicRepo() Repo {
	repo := atomicrepo{}
	repo.updateDB(1)

	go repo.updateLoop()

	return &repo
}

func (repo *atomicrepo) updateDB(value int) {
	db := make(map[string]*Record)
	db["x"] = &Record{value}
	repo.db.Store(db)
}

func (repo *atomicrepo) Get(key string) *Record {
	maybeDb := repo.db.Load()
	db := maybeDb.(map[string]*Record)

	return db[key]
}

func (repo *atomicrepo) updateLoop() {
	tick := time.Tick(time.Second * 3)

	for range tick {
		newValue := repo.Get("x").value + 1
		repo.updateDB(newValue)
		log.Println("Updated to:", newValue)
	}
}
