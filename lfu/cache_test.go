package lfu

import (
	"fmt"
	"testing"
)

type testNode struct {
	Key   interface{}
	Value interface{}
}

func (tn *testNode) String() string {
	return fmt.Sprintf("<k: %v, v: %v>", tn.Key, tn.Value)
}

type expectedNode struct {
	Key  interface{}
	Freq uint
}

func (tn *expectedNode) String() string {
	return fmt.Sprintf("<k: %v, f: %v>", tn.Key, tn.Freq)
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

	for _, expectedNode := range elements {
		actualValue, err := c.Access(expectedNode.Key)
		if err != nil {
			t.Errorf("expected node to be in the cache: %s", expectedNode.String())
		}
		if actualValue != expectedNode.Value {
			t.Errorf("values missmatch. want %v, have %v.", expectedNode.Value, actualValue)
		}
	}
}

func testCacheFrequencyEquals(t *testing.T, c *Cache, elements []expectedNode) {
	t.Helper()

	if !testCacheSizeEquals(t, c, len(elements)) {
		t.FailNow()
	}

	for _, expectedNode := range elements {
		ok := c.Has(expectedNode.Key)
		if !ok {
			t.Errorf("expected node to be in the cache: %s", expectedNode.String())
		}
		frequency, _ := c.FreqKey(expectedNode.Key)
		if frequency != expectedNode.Freq {
			t.Errorf("unexpected frequency. want %d. have %d", expectedNode.Freq, frequency)
		}
	}
}

func newCache(cap uint, initialState []testNode) *Cache {
	cache := NewCache(cap)
	for _, item := range initialState {
		cache.Insert(item.Key, item.Value)
	}
	return cache
}

func TestCacheEvictsLFUifExceedingCapacity_Insert(t *testing.T) {
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

func TestCache_Access(t *testing.T) {
	payload := []testNode{
		{"a", 1},
		{"b", 2},
		{"c", 3},
		{"d", 4},
	}

	cache := newCache(4, payload)
	cache.Access("a")
	cache.Access("a")
	cache.Access("d")
	cache.Access("d")
	cache.Insert("e", 5) // should evict b

	expectedFreqs := []expectedNode{
		{"a", 3},
		{"d", 3},
		{"e", 1},
		{"c", 1},
	}
	testCacheFrequencyEquals(t, cache, expectedFreqs)
}
