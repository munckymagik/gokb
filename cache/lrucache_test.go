package cache_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/munckymagik/gokb/cache"
)

func assert(t *testing.T, condition bool, msg string) {
	if !condition {
		t.Fatal(msg)
	}
}

func assertEq(t *testing.T, a, b interface{}, msg string) {
	assert(t, reflect.DeepEqual(a, b), fmt.Sprintf("%s:\n\tgot: %v\n\texpected: %v", msg, a, b))
}

func TestLRUCache1(t *testing.T) {
	c := cache.NewLRUCache(1)
	assertEq(t, c.Len(), 0, "length wasn't 0")
	assertEq(t, c.Get("a"), "", "a didn't map to ''")

	c.Set("a", "b")
	assertEq(t, c.Len(), 1, "length wasn't 1a")
	assertEq(t, c.Get("a"), "b", "a didn't map to b")

	c.Set("a", "c")
	assertEq(t, c.Get("a"), "c", "a didn't map to c")

	c.Set("b", "a")
	assertEq(t, c.Len(), 1, "length wasn't 1b")
	assertEq(t, c.Get("b"), "a", "b didn't map to a")
	assertEq(t, c.Get("a"), "", "a didn't map to ''")
}

func TestLRUCache2(t *testing.T) {
	c := cache.NewLRUCache(2)
	assertEq(t, c.Len(), 0, "length wasn't 0")

	c.Set("a", "1")
	c.Set("b", "2")
	assertEq(t, c.Len(), 2, "length wasn't 2")

	c.Set("c", "3")
	assertEq(t, c.Get("a"), "", "Adding c should have removed a")
	assertEq(t, c.Get("b"), "2", "Adding c should not have removed b")

	// Make b the most recently used
	c.Get("b")

	// Adding d should cause eviction of c
	c.Set("d", "4")
	assertEq(t, c.Get("a"), "", "a should be long gone")
	assertEq(t, c.Get("b"), "2", "Adding d should not have removed b")
	assertEq(t, c.Get("c"), "", "Adding d should have evicted c")
}
