// Demonstrates using functions as first class data to model a Polynomial.
//
// Run as cd funcs && go run ./polynomialfp
package main

import (
	"fmt"
	"math"
)

type Function func(float64) float64

func Polynomial(terms ...Function) Function {
	return func(x float64) float64 {
		var accumulator float64 = 0.0
		for _, term := range terms {
			accumulator += term(x)
		}

		return accumulator
	}
}

func Term(constant, exponent float64) Function {
	return func(x float64) float64 {
		return constant * math.Pow(x, exponent)
	}
}

func Linear(constant float64) Function {
	return func(x float64) float64 {
		return constant * x
	}
}

func Constant(constant float64) Function {
	return func(x float64) float64 {
		return constant
	}
}

func main() {
	// Models 1.5x^2 + 2x - 0.5
	actualP := Polynomial(
		Term(1.5, 2),
		Linear(2),
		Constant(-.5),
	)

	for i := -5; i <= 5; i++ {
		x := float64(i)
		fmt.Printf("P(%.1f) = %.1f\n", x, actualP(x))
	}
}
