// From https://go.dev/doc/tutorial/generics
package main

import "fmt"

func SumInts(m map[string]int64) int64 {
	var s int64
	for _, v := range m {
		s += v
	}
	return s
}

func SumFloats(m map[string]float64) float64 {
	var s float64
	for _, v := range m {
		s += v
	}
	return s
}

func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

type Number interface {
	int64 | float64
}

func SumNumbers[K comparable, V Number](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

func tutorialMain() {
	ints := map[string]int64{
		"first":  34,
		"second": 12,
	}

	floats := map[string]float64{
		"first":  35.98,
		"second": 26.99,
	}

	fmt.Printf("Non-generic Sums: %v and %v\n",
		SumInts(ints),
		SumFloats(floats))

	// Explicitly passing type arguments
	fmt.Printf("Generic Sums: %v and %v\n",
		SumIntsOrFloats[string, int64](ints),
		SumIntsOrFloats[string, float64](floats))

	// Using type inference
	fmt.Printf("Generic Sums with inferred types: %v and %v\n",
		SumIntsOrFloats(ints),
		SumIntsOrFloats(floats))

	// Using type inference
	fmt.Printf("Generic Sums with a constraint: %v and %v\n",
		SumNumbers(ints),
		SumNumbers(floats))
}
