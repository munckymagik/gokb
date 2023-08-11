package bplustree_test

import (
	"testing"

	"github.com/munckymagik/gokb/bplustree"
	"github.com/stretchr/testify/require"
)

func TestBPlusTreeInsert(t *testing.T) {
	t.Run("it supports generic key/value pairs", func(t *testing.T) {
		// Given
		bt1 := bplustree.New[int, string]()
		bt2 := bplustree.New[string, int]()

		// When
		bt1.Insert(123, "some value")
		bt2.Insert("some key", 123)

		// Then
		require.Equal(t, bt1.Find(123), ptr("some value"))
		require.Equal(t, bt2.Find("some key"), ptr(123))
	})
}

func TestBPlusTreeEach(t *testing.T) {
	t.Run("it performs in-order traversal of entries based on keys", func(t *testing.T) {
		// Given
		type emptyValue = struct{}
		bt := bplustree.New[int, emptyValue]()
		bt.Insert(2, emptyValue{})
		bt.Insert(1, emptyValue{})
		bt.Insert(3, emptyValue{})

		// When
		var collected []int
		bt.Each(func(key int, _ emptyValue) {
			collected = append(collected, key)
		})

		// Then
		require.Equal(t, []int{1, 2, 3}, collected)
	})
}

func ptr[T any](v T) *T {
	return &v
}
