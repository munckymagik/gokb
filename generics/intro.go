// From https://go.dev/blog/intro-generics
package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

func Min(x, y float64) float64 {
	if x < y {
		return x
	}

	return y
}

func GMin[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}

	return y
}

type Tree[T any] struct {
	left, right *Tree[T]
	value       T
}

func NewNode[T any](x T) *Tree[T] {
	return &Tree[T]{value: x}
}

func Scale[E constraints.Integer](s []E, c E) []E {
	r := make([]E, len(s))
	for i, v := range s {
		r[i] = v * c
	}
	return r
}

// S is any type where the underlying type is []E
func ScaleWithInf[S ~[]E, E constraints.Integer](s S, c E) S {
	r := make(S, len(s))
	for i, v := range s {
		r[i] = v * c
	}
	return r
}

type Point []int32

func NewPoint(i ...int32) Point {
	return Point(i)
}

func (p Point) String() string {
	return fmt.Sprintf("Point(%v)", []int32(p))
}

func introMain() {
	fmt.Printf("Non-generic Min: %v\n", Min(1.0, 2.0))
	fmt.Printf("Generic Min: %v and %v\n", GMin(1.0, 2.0), GMin(1, 2))
	fmt.Printf("Generic Tree: %v\n", NewNode("hello"))
	fmt.Printf("Generic Scale: %v and %v\n",
		Scale([]byte{2, 3}, 3),
		Scale([]int16{3, 9}, 7))

	// Type constraint inference
	p := NewPoint(2, 4)
	r := Scale(p, 3)        // returns an []int32 not point
	s := ScaleWithInf(p, 4) // returns a Point
	fmt.Printf("Generic Scale: %v, %v and %v\n", p, r, s)
}
