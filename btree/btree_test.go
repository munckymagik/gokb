package btree_test

import (
	"testing"
	"testing/quick"

	"github.com/munckymagik/gokb/btree"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func TestBTreeInsert(t *testing.T) {
	t.Run("when the key already exists it updates the value", func(t *testing.T) {
		// Given
		bt := btree.New[int, int]()
		bt.Insert(123, 456)

		// When
		bt.Insert(123, 789)

		// Then
		value, _ := bt.Find(123)
		require.Equal(t, 789, value)
	})

	t.Run("it maintains the invariant as new keys are added", func(t *testing.T) {
		propFunc := func(keys []int) bool {
			// Given
			bt := btree.New[int, emptyValue]()

			for _, key := range keys {
				// When
				bt.Insert(key, emptyValue{})

				// Then
				if !assert.NotPanics(t, bt.AssertInvariantsHold, "key=%d", key) {
					return false
				}
			}

			return true
		}

		require.NoError(t, quick.Check(propFunc, &quick.Config{MaxCount: 1000}))
	})
}

func TestBTreeFind(t *testing.T) {
	t.Run("when the key exists it returns the value", func(t *testing.T) {
		// Given
		bt := btree.New[int, int]()
		bt.Insert(123, 456)

		// When
		value, found := bt.Find(123)

		// Then
		require.True(t, found)
		require.Equal(t, 456, value)
	})

	t.Run("when the key does not exist it returns the zero value of the type", func(t *testing.T) {
		// Given
		bt := btree.New[int, int]()
		bt.Insert(456, 789)

		// When
		value, found := bt.Find(123)

		// Then
		require.False(t, found)
		require.Equal(t, 0, value)
	})
}

func TestBTreeEach(t *testing.T) {
	t.Run("it does not crash or execute when the tree is empty", func(t *testing.T) {
		// Given
		bt := btree.New[int, emptyValue]()
		executions := 0

		// When
		bt.Each(func(key int, _ emptyValue) {
			executions += 1
		})

		// Then
		require.Zero(t, executions)
	})

	t.Run("it performs in-order traversal of entries based on keys", func(t *testing.T) {
		propFunc := func(keys []int8) bool {
			// Given
			bt := btree.New[int8, emptyValue]()
			sortedSetKeys := cloneSortAndCompact(keys)

			for _, key := range keys {
				bt.Insert(key, emptyValue{})
			}

			// When
			collected := make([]int8, 0)
			bt.Each(func(key int8, _ emptyValue) {
				collected = append(collected, key)
			})

			// Then
			return assert.Equal(t, sortedSetKeys, collected)
		}

		require.NoError(t, quick.Check(propFunc, &quick.Config{MaxCount: 1000}))
	})
}
func TestBTreeLen(t *testing.T) {
	t.Run("it returns the number of unique entries in the tree", func(t *testing.T) {
		propFunc := func(keys []int8) bool {
			// Given
			bt := btree.New[int8, emptyValue]()
			sortedSetKeys := cloneSortAndCompact(keys)
			for _, key := range keys {
				bt.Insert(key, emptyValue{})
			}

			// Then
			return assert.Equal(t, len(sortedSetKeys), bt.Len())
		}

		require.NoError(t, quick.Check(propFunc, &quick.Config{MaxCount: 1000}))
	})
}

type emptyValue = struct{}

func cloneSortAndCompact[S ~[]E, E constraints.Ordered](s S) S {
	sortedKeys := slices.Clone(s)
	slices.Sort(sortedKeys)
	return slices.Compact(sortedKeys)
}
