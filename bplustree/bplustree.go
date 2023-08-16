package bplustree

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
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
	if p.parent == nil {
		p.parent = newPage[K, V](nil)
		p.parent.children = append(p.parent.children, p)
		t.root = p.parent
	}

	newRight := newPage[K, V](p.parent)
	newRight.entries = append(newRight.entries, p.entries[minEntries+1:]...)
	newRight.entries = slices.Clone(p.entries[minEntries+1:])
	p.entries = p.entries[:minEntries+1]
	if len(p.children) > 0 {
		newRight.children = slices.Clone(p.children[minChildren+1:])
		p.children = p.children[:minChildren+1]
		for _, child := range newRight.children {
			child.parent = newRight
		}
	}

	p.parent.entries = append(p.parent.entries, p.entries[minEntries])
	p.entries = p.entries[:minEntries]
	p.parent.children = append(p.parent.children, newRight)
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

type entry[K any, V any] struct {
	key   K
	value V
}

type page[K constraints.Ordered, V any] struct {
	parent   *page[K, V]
	entries  []entry[K, V]
	children []*page[K, V]
}

func newPage[K constraints.Ordered, V any](parent *page[K, V]) *page[K, V] {
	return &page[K, V]{parent: parent}
}

func (p *page[K, V]) add(key K, value V) {
	newEntry := entry[K, V]{key: key, value: value}
	p.entries = append(p.entries, newEntry)
	p.sort()
}

func (p *page[K, V]) firstEntry() *entry[K, V] {
	if len(p.entries) == 0 {
		return nil
	}
	return &p.entries[0]
}
func (p *page[K, V]) lastEntry() *entry[K, V] {
	if len(p.entries) == 0 {
		return nil
	}
	return &p.entries[len(p.entries)-1]
}

func (p *page[K, V]) sort() {
	slices.SortFunc(p.entries, func(e1, e2 entry[K, V]) bool { return e1.key < e2.key })
	slices.SortFunc(p.children, func(c1, c2 *page[K, V]) bool {
		return c1.lastEntry().key < c2.firstEntry().key
	})
}

func (page *page[K, V]) traverseSubtree(f func(K, V)) {
	i := 0
	for ; i < len(page.entries); i += 1 {
		if len(page.children) > i {
			page.children[i].traverseSubtree(f)
		}
		f(page.entries[i].key, page.entries[i].value)
	}

	if len(page.children) > i {
		page.children[len(page.entries)].traverseSubtree(f)
	}
}

func (p *page[K, V]) find(key K, parent *page[K, V]) (*entry[K, V], *page[K, V]) {
	if p == nil {
		return nil, parent
	}

	i := 0
	for ; i < len(p.entries); i = i + 1 {
		if key < p.entries[i].key {
			// recurse into the left child
			if len(p.children) > i {
				return p.children[i].find(key, p)
			}
			return nil, p
		}

		if key == p.entries[i].key {
			// We've found it, return
			return &p.entries[i], p
		}

		// Continue search to next entry
	}

	// Final option is to check the remaining child
	if i < len(p.children) {
		return p.children[i].find(key, p)
	}

	return nil, p
}
