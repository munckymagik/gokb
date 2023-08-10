package bplustree

type Tree[K any, V any] struct {
	key   K
	value V
}

func New[K any, V any]() Tree[K, V] {
	return Tree[K, V]{}
}

func (t *Tree[K, V]) Insert(key K, value V) {
	t.key = key
	t.value = value
}

func (t *Tree[K, V]) Find(key K) V {
	return t.value
}
