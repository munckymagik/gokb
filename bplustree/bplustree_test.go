package bplustree_test

import (
	"testing"

	"github.com/munckymagik/gokb/bplustree"
	"github.com/stretchr/testify/require"
)

func TestBPlusTree(t *testing.T) {
	// Given
	bt := bplustree.New[int, string]()
	bt2 := bplustree.New[string, int]()

	// When
	bt.Insert(123, "some value")
	bt2.Insert("some key", 123)

	// Then
	require.Equal(t, bt.Find(123), "some value")
	require.Equal(t, bt2.Find("some key"), 123)
}
