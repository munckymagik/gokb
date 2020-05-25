package mapsync

import (
	"log"
	"sync"
	"time"
)

type mutexrepo struct {
	db    RecordMap
	mutex sync.Mutex
}

func NewMutexRepo() Repo {
	repo := mutexrepo{db: make(RecordMap)}
	repo.updateDB(1)

	go repo.updateLoop()

	return &repo
}

func (repo *mutexrepo) updateDB(value int) {
	db := make(RecordMap)
	db["x"] = &Record{value}
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	repo.db = db
}

func (repo *mutexrepo) Get(key string) *Record {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	return repo.db[key]
}

func (repo *mutexrepo) updateLoop() {
	tick := time.Tick(time.Second * 3)

	for range tick {
		newValue := repo.Get("x").value + 1
		repo.updateDB(newValue)
		log.Println("Updated to:", newValue)
	}
}
