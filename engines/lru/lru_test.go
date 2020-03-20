package engines

import (
	"github.com/sonirico/mecachis/engines"
	"reflect"
	"testing"
)

type cachevalue string

func (v cachevalue) Value() interface{} {
	return v
}

func (v cachevalue) Len() uint64 {
	return uint64(len(v))
}

type testNode struct {
	Key   string
	Value cachevalue
}

type expectedState struct {
	Nodes     []testNode
	CacheSize uint64
}

func testCacheSizeEquals(t *testing.T, c *lru, expectedSize uint64) bool {
	t.Helper()

	if c.Size() != expectedSize {
		t.Errorf("wrong cache size. want %d. have %d", expectedSize, c.Size())
		return false
	}

	return true
}

func testNodeEquals(t *testing.T, en testNode, cn engines.Entry) bool {
	t.Helper()

	if en.Key != cn.Key() {
		t.Errorf("keys missmatch. want %v, have %v.", en.Key, cn.Key())
		return false
	}

	if en.Value != cn.Value().Value() {
		t.Errorf("values missmatch. want %v, have %v.", en.Value, cn.Value().Value())
		return false
	}

	return true
}

func testCacheStateEquals(t *testing.T, c *lru, eState *expectedState) {
	t.Helper()

	if !testCacheSizeEquals(t, c, eState.CacheSize) {
		t.FailNow()
	}

	for position, actualNode := range c.Dump() {
		expectedNode := eState.Nodes[position]
		testNodeEquals(t, expectedNode, actualNode)
	}
}

func newCache(cap uint64, initialState []testNode) *lru {
	cache := New(cap)
	for _, item := range initialState {
		cache.Insert(item.Key, item.Value)
	}
	return cache
}

func TestCacheLRU_EvictsLRUIfExceedingCapacity_Insert(t *testing.T) {
	payload := []testNode{
		{"a", cachevalue("1")}, // +2
		{"b", cachevalue("2")}, // +2
		{"c", cachevalue("3")}, // +2
		{"d", cachevalue("4")}, // +2
	}
	expectedState := &expectedState{
		Nodes: []testNode{
			{"d", cachevalue("4")},
			{"c", cachevalue("3")},
			{"b", cachevalue("2")},
		},
		CacheSize: 6,
	}
	cache := newCache(6, payload)
	testCacheStateEquals(t, cache, expectedState)
}

func TestCacheLRUReturnsErrorIfDuplicated_Insert(t *testing.T) {
	var payload []testNode
	cache := newCache(3, payload)
	ok := cache.Insert("a", cachevalue("1"))
	if !ok {
		t.Errorf("expected successful insertion. want %t, have %t", true, ok)
	}
	ok = cache.Insert("a", cachevalue("1"))
	if ok {
		t.Errorf("expected no insertion. want %t, have %t", false, ok)
	}
}

func TestCacheLRU_Access_UpgradesToHead(t *testing.T) {
	payload := []testNode{
		{"a", cachevalue("1")},
		{"b", cachevalue("2")},
		{"c", cachevalue("3")},
	}
	expectedState := &expectedState{
		Nodes: []testNode{
			{"a", cachevalue("1")},
			{"c", cachevalue("3")},
			{"b", cachevalue("2")},
		},
		CacheSize: 6,
	}
	cache := newCache(32, payload)
	value, _ := cache.Access("a") // "a" should be put on top, leaving "b" at the bottom
	testCacheStateEquals(t, cache, expectedState)
	cached := value.(cachevalue)
	if cached != "1" {
		t.Errorf("wrong cachevalue returned. want '%s', have '%v'", "1", cached)
	}
}

func TestCacheLRU_Access_UpgradesToHead_OneElement(t *testing.T) {
	payload := []testNode{
		{"a", cachevalue(1)},
	}
	expectedState := &expectedState{
		Nodes: []testNode{
			{"a", cachevalue(1)},
		},
		CacheSize: 2,
	}
	cache := newCache(3, payload)
	_, _ = cache.Access("a")
	testCacheStateEquals(t, cache, expectedState)
}

func TestCacheLRU_OnEvicted(t *testing.T) {
	payload := []testNode{
		{"a", cachevalue(1)}, // +2
	}
	keys := make([]string, 0)
	onEvicted := func(v engines.Entry) {
		keys = append(keys, v.Key())
	}
	cache := newCache(4, payload)
	cache.OnEvict(onEvicted)
	cache.Insert("b", cachevalue("2")) // +2
	cache.Insert("c", cachevalue("3")) // +2, one element should have been evicted
	cache.Insert("d", cachevalue("4")) // +2, one element should have been evicted
	if !reflect.DeepEqual(keys, []string{"a", "b"}) {
		t.Fatalf("wrong set of elements have been evicted. instead have %v", keys)
	}
}
