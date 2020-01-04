package lfu

import (
	"bytes"
	"fmt"
	"testing"
)

type testNode struct {
	Key   interface{}
	Value interface{}
}

type maybeNode struct {
	one   *testNode
	other *testNode
}

func (tn *testNode) String() string {
	return fmt.Sprintf("<k: %v, v: %v>", tn.Key, tn.Value)
}

type expectedNode struct {
	Key  interface{}
	Freq uint
}

func (m *maybeNode) String() string {
	var buf bytes.Buffer
	buf.WriteString("{One:")
	buf.WriteString(m.one.String())
	buf.WriteString("}")
	buf.WriteString("\n")
	buf.WriteString("{Other:")
	buf.WriteString(m.one.String())
	buf.WriteString("}")
	return buf.String()
}
func (tn *expectedNode) String() string {
	return fmt.Sprintf("<k: %v, f: %v>", tn.Key, tn.Freq)
}

func testCacheSizeEquals(t *testing.T, c *Cache, expectedSize uint) bool {
	t.Helper()

	if c.Size() != uint(expectedSize) {
		t.Errorf("wrong cache size. want %d. have %d", expectedSize, c.Size())
		return false
	}

	return true
}

func testCacheFrequencyEquals(t *testing.T, c *Cache, size uint, elements []expectedNode) bool {
	t.Helper()

	if !testCacheSizeEquals(t, c, size) {
		t.FailNow()
		return false
	}

	for _, expectedNode := range elements {
		ok := c.Has(expectedNode.Key)
		if !ok {
			t.Errorf("expected node to be in the cache: %s", expectedNode.String())
			return false
		}
		frequency, _ := c.FreqKey(expectedNode.Key)
		if frequency != expectedNode.Freq {
			t.Errorf("unexpected frequency. want %d. have %d", expectedNode.Freq, frequency)
			return false
		}
	}
	return true
}

func testCacheFrequencyEqualsMaybe(t *testing.T, c *Cache, size uint, elements []expectedNode, maybe maybeNode) {
	t.Helper()

	testCacheFrequencyEquals(t, c, size, elements)

	if !(c.Has(maybe.one.Key) || c.Has(maybe.other.Key)) {
		t.Errorf("expected cache to have any element of %s", maybe.String())
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
	cache := newCache(3, payload)
	testCacheSizeEquals(t, cache, 3)
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
	cache.Insert("e", 5) // should evict b or c

	expectedFreqs := []expectedNode{
		{"a", 3},
		{"d", 3},
		{"e", 1},
	}

	maybeNode := maybeNode{
		one:   &testNode{"b", 2},
		other: &testNode{"c", 3},
	}

	testCacheFrequencyEqualsMaybe(t, cache, 4, expectedFreqs, maybeNode)
}
