package bplustree

import (
	"fmt"

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
	root *page[K, V]
}

func New[K constraints.Ordered, V any]() Tree[K, V] {
	return Tree[K, V]{}
}

func (t *Tree[K, V]) Insert(key K, value V) {
	_, page := t.root.doFind(key, nil)
	if page == nil {
		t.root = newPage[K, V](nil)
		page = t.root
	}

	page.add(key, value)

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

func (t *Tree[K, V]) Find(key K) *V {
	found, _ := t.root.doFind(key, nil)
	if found == nil {
		return nil
	}

	return &found.value
}

func (t *Tree[K, V]) Each(f func(K, V)) {
	if t.root != nil {
		t.root.traverseSubtree(f)
	}
}

func (t *Tree[K, V]) AssertInvariantsHold() {
	t.root.validateSubtree()
}

func (t *Tree[K, V]) CheckInvariantsHold() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	t.AssertInvariantsHold()
	return nil
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

func (p *page[K, V]) doFind(key K, parent *page[K, V]) (*entry[K, V], *page[K, V]) {
	if p == nil {
		return nil, parent
	}

	i := 0
	for ; i < len(p.entries); i = i + 1 {
		if key < p.entries[i].key {
			// recurse into the left child
			if len(p.children) > i {
				return p.children[i].doFind(key, p)
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
		return p.children[i].doFind(key, p)
	}

	return nil, p
}

func (p *page[K, V]) validateSubtree() {
	if p == nil {
		return
	}

	for _, child := range p.children {
		child.validateSubtree()
	}

	p.validate()
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

func (p *page[K, V]) validate() {
	isLeaf := len(p.children) == 0
	isRoot := p.parent == nil

	// 1. Every page has at most 2k+1 children.
	assert("too many children", len(p.children) <= maxChildren)

	// 1.1 Every page has at most 2k entries.
	assert("too many entries", len(p.entries) <= maxEntries)

	// 2. Every non-leaf page (except root) has at least ⌈(2k+1)/2⌉ child pages.
	if !isLeaf && !isRoot {
		assert("too few children", minChildren <= len(p.children))
	}

	// 3. The root has at least two children if it is not a leaf page.
	if isRoot && !isLeaf {
		assert("root has too few children", len(p.children) >= 2)
	}

	// 4. A non-leaf page with n children contains n−1 entries.
	if !isLeaf {
		assert("n children but not n-1 entries", len(p.entries) == len(p.children)-1)
	}

	if !isRoot {
		assert("child has wrong parent", slices.Contains(p.parent.children[:], p))
	}

	if isLeaf {
		return
	}

	for i, entry := range p.entries {
		assert(fmt.Sprintf("NOT child[%d].lastEntry.key < entries[%d].key", i, i), p.children[i].lastEntry().key < entry.key)
		// if len(p.children) > i+1 {
		assert(fmt.Sprintf("NOT entries[%d].key < child[%d].firstEntry.key", i, i+1), entry.key < p.children[i+1].firstEntry().key)
		// }
	}
}

func assert(msg string, value bool) {
	if !value {
		panic(msg)
	}
}
