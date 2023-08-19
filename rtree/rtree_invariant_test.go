package rtree

import (
	"fmt"
)

type invariantState struct {
	leafLevel int
}

func (t *RTree[T]) AssertInvariantsHold() {
	var state invariantState
	t.root.validateSubtree(1, &state)
}

func (t *RTree[T]) CheckInvariantsHold() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	t.AssertInvariantsHold()
	return nil
}

func (n *node[T]) validateSubtree(level int, state *invariantState) {
	if n == nil {
		return
	}

	for _, child := range n.children {
		child.validateSubtree(level+1, state)
	}

	n.validate(level, state)
}

func (n *node[T]) validate(level int, state *invariantState) {
	isLeaf := len(n.children) == 0
	isRoot := level == 1

	// (1) Every leaf node contalns between m and M index records unless it the root
	if isLeaf {
		assert("too many entries", len(n.bbox) <= maxEntries)
		assert("too few entries", len(n.bbox) >= minEntries)
		assert("len(bbox) != len(data)", len(n.bbox) == len(n.data))
	}

	// (2) For each index record (I, tuple-identifier) in a leaf node, I is the smallest
	// rectangle that spatially contains the n-dimensional data object represented
	// by the indicated tuple
	// if isLeaf {
	// }

	// (3) Every non-leaf node has between m and M children unless it is the root
	if !isLeaf {
		assert("too many children", len(n.children) <= maxEntries)
		assert("too few children", len(n.children) >= minEntries)
		assert("len(bbox) != len(children)", len(n.bbox) == len(n.children))
	}

	// (4) For each entry (I, child-pointer) in a non-leaf node, I is the
	// smallest rectangle that spatially contains the rectangles in the child node
	if !isLeaf {
		for i := 0; i < len(n.bbox); i += 1 {
			bbox := n.bbox[i]
			bboxes := rect{}
			for _, childBbox := range n.children[0].bbox {
				bboxes = bboxes.add(childBbox)
			}
			assert("bbox does not tightly contain child bboxes", bbox == bboxes) // TODO floating point comparison
		}
	}

	// (5) The root node has at least two children unless it is a leaf
	if isRoot && !isLeaf {
		assert("too few children in non-leaf root", len(n.children) >= 2)
	}

	// (6) All leaves appear on the same level
	if isLeaf {
		if state.leafLevel == 0 {
			state.leafLevel = level
		} else {
			assert(fmt.Sprintf("Leaf at wrong level. Got %d, expected %d", level, state.leafLevel), level == state.leafLevel)
		}

		return
	}
}

func assert(msg string, value bool) {
	if !value {
		panic(msg)
	}
}
