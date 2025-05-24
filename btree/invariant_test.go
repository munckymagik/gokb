package btree

import (
	"fmt"

	"golang.org/x/exp/slices"
)

type invariantState struct {
	leafLevel int
}

func (t *Tree[K, V]) AssertInvariantsHold() {
	var state invariantState
	t.root.validateSubtree(1, &state)
}

func (p *page[K, V]) validateSubtree(level int, state *invariantState) {
	if p == nil {
		return
	}

	for _, child := range p.children {
		child.validateSubtree(level+1, state)
	}

	p.validate(level, state)
}

func (p *page[K, V]) validate(level int, state *invariantState) {
	isLeaf := len(p.children) == 0
	isRoot := p.parent == nil

	// 1. Every page has at most 2k+1 children.
	assert("too many children", len(p.children) <= maxChildren)
	assert("children cap too large", cap(p.children) <= maxChildren+1)

	// 1.1 Every page has at most 2k entries.
	assert("too many entries", len(p.entries) <= maxEntries)
	assert("entries cap too large", cap(p.entries) <= maxEntries+1)

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
		if state.leafLevel == 0 {
			state.leafLevel = level
		} else {
			assert(fmt.Sprintf("Leaf at wrong level. Got %d, expected %d", level, state.leafLevel), level == state.leafLevel)
		}

		return
	}

	for i, entry := range p.entries {
		assert(fmt.Sprintf("NOT child[%d].lastEntry.key < entries[%d].key", i, i), p.children[i].entries.last().key < entry.key)
		assert(fmt.Sprintf("NOT entries[%d].key < child[%d].firstEntry.key", i, i+1), entry.key < p.children[i+1].entries.first().key)
	}
}

func assert(msg string, value bool) {
	if !value {
		panic(msg)
	}
}
