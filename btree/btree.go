package btree

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

type items[E any] []E

func (self items[E]) first() *E {
	if len(self) == 0 {
		return nil
	}
	return &self[0]
}
func (self items[E]) last() *E {
	if len(self) == 0 {
		return nil
	}
	return &self[len(self)-1]
}

func (self *items[E]) add(elem E) {
	*self = append(*self, elem)
}

func (self *items[E]) pop() E {
	elem := *self.last()
	*self = (*self)[:len(*self)-1]
	return elem
}

func (self items[E]) sort(less func(a E, b E) bool) {
	slices.SortFunc(self, less)
}

func (src *items[E]) splitAt(fromIndex int) items[E] {
	ret := (*src)[fromIndex:]
	*src = (*src)[:fromIndex:fromIndex]
	return ret
}

type entry[K any, V any] struct {
	key   K
	value V
}

type page[K constraints.Ordered, V any] struct {
	parent   *page[K, V]
	entries  items[entry[K, V]]
	children items[*page[K, V]]
}

func newPage[K constraints.Ordered, V any](parent *page[K, V]) *page[K, V] {
	return &page[K, V]{parent: parent}
}

func (p *page[K, V]) isRoot() bool {
	return p.parent == nil
}

func (p *page[K, V]) add(key K, value V) {
	newEntry := entry[K, V]{key: key, value: value}
	p.entries.add(newEntry)
	p.sort()
}

func (p *page[K, V]) setChildren(incoming ...*page[K, V]) {
	p.children = incoming
	for _, child := range p.children {
		child.parent = p
	}
}

func (p *page[K, V]) sort() {
	p.entries.sort(func(e1, e2 entry[K, V]) bool { return e1.key < e2.key })

	p.children.sort(func(c1, c2 *page[K, V]) bool {
		return c1.entries.last().key < c2.entries.first().key
	})
}

func (p *page[K, V]) traverseSubtree(f func(K, V)) {
	i := 0
	for ; i < len(p.entries); i += 1 {
		if len(p.children) > i {
			p.children[i].traverseSubtree(f)
		}
		f(p.entries[i].key, p.entries[i].value)
	}

	if len(p.children) > i {
		p.children[len(p.entries)].traverseSubtree(f)
	}
}

func (p *page[K, V]) find(key K, parent *page[K, V]) (*entry[K, V], *page[K, V]) {
	if p == nil {
		return nil, parent
	}

	i := 0
	for ; i < len(p.entries); i += 1 {
		if key == p.entries[i].key {
			return &p.entries[i], p
		}

		if key < p.entries[i].key {
			break
		}
	}

	if i < len(p.children) {
		return p.children[i].find(key, p)
	}

	return nil, p
}
