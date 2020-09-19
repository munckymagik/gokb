// Demonstrates using functions as first class data, method values and function
// signatures as "interfaces of one".
//
// Together they let you inject dependencies while still abstracting whether
// the resource is a function or a stateful type.
//
// Run as cd funcs && go run ./polynomial
package main

import (
	"fmt"
	"math"
)

// Function is a function signature that defines the "interface" type we are
// going to be using to abstract the dependencies Polynomial relies on to do
// its thing.
type Function func(float64) float64

// Polynomial models a polynomial function like P(x) = ax^3 + bx^2 + cx + d
// Imagine this is a type in your application domain model that relies on being
// passed resources of type Function to do its job.
type Polynomial struct {
	terms []Function
}

// Apply applies the polynomial to the variable x.
// It is a method that happens to satisify the Function type signature.
func (p Polynomial) Apply(x float64) float64 {
	var accumulator float64 = 0.0
	for _, term := range p.terms {
		accumulator += term(x)
	}

	return accumulator
}

// NewPolynomial creates a Polynomial from the given set of terms.
func NewPolynomial(terms ...Function) Polynomial {
	return Polynomial{terms}
}

// PowerTerm models a term that raises x to a power then multiplies it by a
// constant.
// Imagine this is a type too complex for our tests, but we do want to inject
// as a dependency for our Polynomial at runtime.
type PowerTerm struct {
	constant float64
	exponent float64
}

// Apply applies the term function to the variable x.
// It is a method that happens to satisify the Function type signature.
func (t PowerTerm) Apply(x float64) float64 {
	return t.constant * math.Pow(x, t.exponent)
}

// testPolynomial is a pretend test case
func testPolynomial( /* t *testing.T */ ) {
	// In tests you can pass simple functions with planted values
	identity := func(x float64) float64 { return x }
	five := func(x float64) float64 { return 5 }

	fixtureP := NewPolynomial(identity, identity, five)

	result := fixtureP.Apply(1)
	if result != 7 {
		panic(fmt.Sprintf("Expected P(1) = 7, but was %f", result))
	}
}

func runApp() {
	// In app code you can pass more complicated types with their own state
	// by passing "method values" that satisfy the function signature
	actualP := NewPolynomial(
		PowerTerm{1.5, 2}.Apply,                // a bound method value
		PowerTerm{2, 1}.Apply,                  // a bound method value
		func(x float64) float64 { return -.5 }, // an anonymous function
	)

	for i := -5; i <= 5; i++ {
		x := float64(i)
		fmt.Printf("P(%.1f) = %.1f\n", x, actualP.Apply(x))
	}
}

func main() {
	testPolynomial()
	runApp()
}
