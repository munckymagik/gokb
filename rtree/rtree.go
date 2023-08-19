package rtree

import "math"

const (
	maxEntries = 3
	minEntries = maxEntries / 2
)

type RTree[T any] struct {
	root *node[T]
}

type rect struct {
	xMin, xMax float64
	yMin, yMax float64
}

func (r rect) add(other rect) rect {
	return rect{
		xMin: math.Min(r.xMin, other.xMin),
		xMax: math.Max(r.xMax, other.xMax),
		yMin: math.Min(r.yMin, other.yMin),
		yMax: math.Max(r.yMax, other.yMax),
	}
}

type node[T any] struct {
	bbox     []rect
	children []*node[T]
	data     []*T
}
