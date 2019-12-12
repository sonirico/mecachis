package lru

import (
	"fmt"
	"testing"
)

func TestCacheEvictsLRUifExceedingCapacity_Insert(t *testing.T) {
	payload := []string{"a", "b", "c", "d"}
	capacity := 3
	cache := NewCache(capacity)
	for _, item := range payload {
		cache.Insert(item, true)
	}
	fmt.Println(cache.Dump())
	if cache.Size() != capacity {
		t.Errorf("wrong cache size. want %d. have %d", capacity, cache.Size())
	}
}

func TestCacheReturnsErrorIfDuplicated_Insert(t *testing.T) {
	capacity := 3
	cache := NewCache(capacity)
	ok := cache.Insert("a", 1)
	if !ok {
		t.Errorf("expected successful insertion. want %t, have %t", true, ok)
	}
	ok = cache.Insert("a", 1)
	if ok {
		t.Errorf("expected no insertion. want %t, have %t", false, ok)
	}
}

func TestCacheAccessAlreadyInsertedKeyRaisesAsHead_Access(t *testing.T) {
	capacity := 3
	cache := NewCache(capacity)
	payloads := []struct {
		Key   string
		Value interface{}
	}{
		{"a", 1},
		{"b", 2},
		{"c", 3},
	}
	for _, item := range payloads {
		cache.Insert(item.Key, item.Value)
	}
	value, err := cache.Access("a") // "a" should be put on top, leaving "b" at the bottom
	if value != 1 {
		t.Errorf("wrong value returned. want %d, have %v", 1, value)
	}
	cache.Insert("d", 4) // "b" should be evicted
	_, err = cache.Access("b")
	if err == nil {
		t.Error("accessing non-existent key should yield error")
	}
}
