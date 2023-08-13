package bplustree_test

import (
	"testing"
	"testing/quick"

	"github.com/munckymagik/gokb/bplustree"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
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

	t.Run("it maintains the invariant as new keys are added", func(t *testing.T) {
		propFunc := func(keys []int) bool {
			// Given
			bt := bplustree.New[int, emptyValue]()

			// When
			for _, key := range keys {
				bt.Insert(key, emptyValue{})
			}

			// Then
			return assert.NoError(t, bt.CheckInvariantsHold())
		}

		require.NoError(t, quick.Check(propFunc, &quick.Config{MaxCount: 1000}))
	})
}

func TestBPlusTreeFind(t *testing.T) {
	t.Run("when the key exists it returns a pointer to the value", func(t *testing.T) {
		// Given
		bt := bplustree.New[int, int]()
		bt.Insert(123, 456)

		// When
		found := bt.Find(123)

		// Then
		require.Equal(t, ptr(456), found)
	})

	t.Run("when the key does not exist it returns nil", func(t *testing.T) {
		// Given
		bt := bplustree.New[int, int]()
		bt.Insert(456, 789)

		// When
		found := bt.Find(123)

		// Then
		require.Nil(t, found)
	})
}

func TestBPlusTreeEach(t *testing.T) {
	t.Run("it does not crash or execute when the tree is empty", func(t *testing.T) {
		// Given
		bt := bplustree.New[int, emptyValue]()
		executions := 0

		// When
		bt.Each(func(key int, _ emptyValue) {
			executions += 1
		})

		// Then
		require.Zero(t, executions)
	})

	t.Run("it performs in-order traversal of entries based on keys", func(t *testing.T) {
		propFunc := func(keys []int) bool {
			// Given
			bt := bplustree.New[int, emptyValue]()
			sortedKeys := slices.Clone(keys)
			slices.Sort(sortedKeys)

			for _, key := range keys {
				bt.Insert(key, emptyValue{})
			}

			// When
			collected := make([]int, 0)
			bt.Each(func(key int, _ emptyValue) {
				collected = append(collected, key)
			})

			// Then
			return assert.Equal(t, sortedKeys, collected)
		}

		require.NoError(t, quick.Check(propFunc, &quick.Config{MaxCount: 1000}))
	})
}

type emptyValue = struct{}

func ptr[T any](v T) *T {
	return &v
}
