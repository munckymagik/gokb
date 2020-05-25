package mapsync

import (
	"log"
	"time"
)

type naiverepo struct {
	db RecordMap
}

func NewNaiveRepo() Repo {
	repo := naiverepo{make(RecordMap)}
	repo.db["x"] = &Record{1}
	go repo.updateLoop()
	return &repo
}

func (repo *naiverepo) Get(key string) *Record {
	return repo.db[key]
}

func (repo *naiverepo) updateLoop() {
	tick := time.Tick(time.Second * 3)

	for range tick {
		newValue := repo.db["x"].value + 1
		repo.db = make(RecordMap)
		repo.db["x"] = &Record{newValue}
		log.Println("Updated to:", newValue)
	}
}
