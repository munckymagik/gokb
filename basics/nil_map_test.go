package basics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNilMap(t *testing.T) {
	var m map[int]int

	require.Nil(t, m)
	require.Zero(t, len(m))

	// m[2] = 3 // nil dereference
	require.Zero(t, m[2])
}

func TestNilMap2(t *testing.T) {
	m := map[string]struct{}{}
	require.NotNil(t, m)

	require.Equal(t, 0, len(m))

	m["x"] = struct{}{}
	require.Equal(t, 1, len(m))
}
