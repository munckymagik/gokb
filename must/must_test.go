package must_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMust(t *testing.T) {
	// This doesn't compile
	// _ = must(t)(newThing())

	// Sadly not sure this feels simple enough for my purposes
	_ = must[*Thing](t)(newThing())

	mustThing := must[*Thing](t)
	_ = mustThing(newThing())
}

type Thing struct{}

func newThing() (*Thing, error) {
	return nil, fmt.Errorf("some error")
}

func must[T any](t *testing.T) func(returnValue T, err error) T {
	return func(returnValue T, err error) T {
		t.Helper()
		require.NoError(t, err)
		return returnValue
	}
}
