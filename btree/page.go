package btree

import "golang.org/x/exp/constraints"

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
