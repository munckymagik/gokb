package cache

import "container/list"

type Cache interface {
	Get(key string) string
	Set(key, val string)
	Len() int
}

type elem struct {
	listElem *list.Element
	val      string
}

type lruCache struct {
	maxSize int
	data    map[string]elem
	history *list.List
}

func NewLRUCache(maxSize int) Cache {
	return &lruCache{maxSize, make(map[string]elem), list.New()}
}

func (cache *lruCache) Get(key string) string {
	elem, ok := cache.data[key]

	if ok {
		cache.history.MoveToFront(elem.listElem)
		return elem.val
	}

	return ""
}
func (cache *lruCache) Set(key, val string) {
	if cache.Len() == cache.maxSize {
		evictee := cache.history.Back()
		delete(cache.data, evictee.Value.(string))
		cache.history.Remove(evictee)
	}

	listElem := cache.history.PushFront(key)
	elem := elem{listElem, val}
	cache.data[key] = elem
}
func (cache *lruCache) Len() int {
	return len(cache.data)
}
