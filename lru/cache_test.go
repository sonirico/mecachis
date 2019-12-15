package lru

import (
	"testing"
)

type testNode struct {
	Key   interface{}
	Value interface{}
}

func testCacheSizeEquals(t *testing.T, c *Cache, expectedSize int) bool {
	t.Helper()

	if c.Size() != uint(expectedSize) {
		t.Errorf("wrong cache size. want %d. have %d", expectedSize, c.Size())
		return false
	}

	return true
}

func testNodeEquals(t *testing.T, en testNode, cn cacheNode) bool {
	t.Helper()

	if en.Key != cn.key {
		t.Errorf("keys missmatch. want %v, have %v.", en.Key, cn.key)
		return false
	}

	if en.Value != cn.value {
		t.Errorf("values missmatch. want %v, have %v.", en.Value, cn.value)
		return false
	}

	return true
}

func testCacheStateEquals(t *testing.T, c *Cache, elements []testNode) {
	t.Helper()

	if !testCacheSizeEquals(t, c, len(elements)) {
		t.FailNow()
	}

	for position, actualNode := range c.Nodes() {
		expectedNode := elements[position]

		testNodeEquals(t, expectedNode, actualNode)
	}
}

func newCache(cap uint, initialState []testNode) *Cache {
	cache := NewCache(cap)
	for _, item := range initialState {
		cache.Insert(item.Key, item.Value)
	}
	return cache
}

func TestCacheEvictsLRUifExceedingCapacity_Insert(t *testing.T) {
	payload := []testNode{
		{"a", 1},
		{"b", 2},
		{"c", 3},
		{"d", 4},
	}
	expectedState := []testNode{
		{"d", 4},
		{"c", 3},
		{"b", 2},
	}
	cache := newCache(3, payload)
	testCacheStateEquals(t, cache, expectedState)
}

func TestCacheReturnsErrorIfDuplicated_Insert(t *testing.T) {
	var payload []testNode
	cache := newCache(3, payload)
	ok := cache.Insert("a", 1)
	if !ok {
		t.Errorf("expected successful insertion. want %t, have %t", true, ok)
	}
	ok = cache.Insert("a", 1)
	if ok {
		t.Errorf("expected no insertion. want %t, have %t", false, ok)
	}
}

func TestCache_Access_UpgradesToHead(t *testing.T) {
	payload := []testNode{
		{"a", 1},
		{"b", 2},
		{"c", 3},
	}
	expectedState := []testNode{
		{"a", 1},
		{"c", 3},
		{"b", 2},
	}
	cache := newCache(3, payload)
	value, _ := cache.Access("a") // "a" should be put on top, leaving "b" at the bottom
	testCacheStateEquals(t, cache, expectedState)
	if value != 1 {
		t.Errorf("wrong value returned. want %d, have %v", 1, value)
	}
}

func TestCache_Access_UpgradesToHead_OneElement(t *testing.T) {
	payload := []testNode{
		{"a", 1},
	}
	expectedState := []testNode{
		{"a", 1},
	}
	cache := newCache(3, payload)
	_, _ = cache.Access("a")
	testCacheStateEquals(t, cache, expectedState)
}
