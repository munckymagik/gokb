package btree

import (
	"fmt"

	"golang.org/x/exp/slices"
)

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

func (p *page[K, V]) validateSubtree() {
	if p == nil {
		return
	}

	for _, child := range p.children {
		child.validateSubtree()
	}

	p.validate()
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
		assert(fmt.Sprintf("NOT entries[%d].key < child[%d].firstEntry.key", i, i+1), entry.key < p.children[i+1].firstEntry().key)
	}
}

func assert(msg string, value bool) {
	if !value {
		panic(msg)
	}
}
