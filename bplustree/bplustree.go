package bplustree

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

type entry[K any, V any] struct {
	key   K
	value V
}

type page[K constraints.Ordered, V any] struct {
	entries []entry[K, V]
}

func (p *page[K, V]) add(key K, value V) {
	newEntry := entry[K, V]{key: key, value: value}
	p.entries = append(p.entries, newEntry)
	slices.SortFunc(p.entries, func(e1, e2 entry[K, V]) bool { return e1.key < e2.key })
}

type Tree[K constraints.Ordered, V any] struct {
	root *page[K, V]
}

func New[K constraints.Ordered, V any]() Tree[K, V] {
	return Tree[K, V]{}
}

func (t *Tree[K, V]) Insert(key K, value V) {
	if t.root == nil {
		t.root = &page[K, V]{}
	}
	t.root.add(key, value)
}

func (t *Tree[K, V]) Find(key K) *V {
	return &t.root.entries[0].value
}

func (t *Tree[K, V]) Each(f func(K, V)) {
	t.traverseSubtree(t.root, f)
	// f(t.root.entries[0].key, t.root.entries[0].value)
}

func (t *Tree[K, V]) traverseSubtree(page *page[K, V], f func(K, V)) {
	for _, entry := range page.entries {
		f(entry.key, entry.value)
	}
}
