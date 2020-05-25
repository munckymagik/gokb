package mapsync

import (
	"log"
	"sync"
	"time"
)

type rwmutexrepo struct {
	db    RecordMap
	mutex sync.RWMutex
}

func NewRWMutexRepo() Repo {
	repo := rwmutexrepo{db: make(RecordMap)}
	repo.updateDB(1)

	go repo.updateLoop()

	return &repo
}

func (repo *rwmutexrepo) updateDB(value int) {
	db := make(RecordMap)
	db["x"] = &Record{value}
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	repo.db = db
}

func (repo *rwmutexrepo) Get(key string) *Record {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	return repo.db[key]
}

func (repo *rwmutexrepo) updateLoop() {
	tick := time.Tick(time.Second * 3)

	for range tick {
		newValue := repo.Get("x").value + 1
		repo.updateDB(newValue)
		log.Println("Updated to:", newValue)
	}
}
