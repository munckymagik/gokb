package mapsync

type Record struct {
	value int
}

type Repo interface {
	Get(key string) *Record
}
