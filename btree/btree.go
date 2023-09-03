package btree

import (
	"golang.org/x/exp/constraints"
)

const (
	k           int = 1
	maxEntries  int = 2 * k           // 2
	minEntries  int = maxEntries / 2  // 1
	maxChildren int = maxEntries + 1  // 3
	minChildren int = maxChildren / 2 // 1
)

type Tree[K constraints.Ordered, V any] struct {
	root       *page[K, V]
	numEntries int
}

func New[K constraints.Ordered, V any]() Tree[K, V] {
	return Tree[K, V]{}
}

func (t *Tree[K, V]) Len() int {
	return t.numEntries
}
func (t *Tree[K, V]) Insert(key K, value V) {
	entry, page := t.root.find(key, nil)
	if page == nil {
		t.root = newPage[K, V](nil)
		page = t.root
	}

	if entry != nil {
		entry.value = value
		return
	}

	page.add(key, value)
	t.numEntries += 1

	if len(page.entries) > maxEntries {
		t.split(page)
	}
}

func (t *Tree[K, V]) split(p *page[K, V]) {
	if p.isRoot() {
		newRoot := newPage[K, V](nil)
		newRoot.setChildren(p)
		t.root = newRoot
	}

	newRight := newPage[K, V](p.parent)
	newRight.entries = p.entries.splitAt(minEntries + 1)

	if len(p.children) > 0 {
		newRight.setChildren(p.children.splitAt(minChildren + 1)...)
	}

	p.parent.entries.add(p.entries.pop())
	p.parent.children.add(newRight)
	p.parent.sort()

	if len(p.parent.entries) > maxEntries {
		t.split(p.parent)
	}
}

func (t *Tree[K, V]) Find(key K) (V, bool) {
	entry, _ := t.root.find(key, nil)
	if entry == nil {
		var zeroValue V
		return zeroValue, false
	}

	return entry.value, true
}

func (t *Tree[K, V]) Each(f func(K, V)) {
	if t.root != nil {
		t.root.traverseSubtree(f)
	}
}
