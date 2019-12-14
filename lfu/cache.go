package lfu

import (
	"errors"
	"fmt"
)

type cacheKey interface{}
type cacheValue interface{}

// Cache represents the LRU cache
type Cache struct {
	// Maximum capacity
	Capacity uint
	// How many elements are cached
	itemsCount uint
	// cache nodes hash-map
	items map[cacheKey]*cacheNode
	// The head of the frequencies dll
	freqHeadNode *freqNode
}

// NewCache initializes a new Cache by providing the maximum
// capacity which, once reached, will provoke to evict the LRU
// element
func NewCache(capacity uint) *Cache {
	cache := &Cache{
		Capacity:     capacity,
		itemsCount:   0,
		items:        make(map[cacheKey]*cacheNode, capacity),
		freqHeadNode: newHeadFreqNode(),
	}
	return cache
}

func (c *Cache) evict() *cacheNode {
	// Get node with lowest frequency
	lfuNode := c.freqHeadNode.next
	// remove a random node (for now) from it
	node := lfuNode.Pop()
	// remove it from cache registry
	delete(c.items, node)
	// update counter accordingly
	c.itemsCount--
	// Remove the frequency node if it has run out of items
	if lfuNode.Size() < 1 {
		c.removeNode(lfuNode)
	}
	return node
}

func (c *Cache) removeNode(node *freqNode) {
	node.prev.next = node.next
	if node.next != nil {
		node.next.prev = node.prev
	}
}

func (c *Cache) Insert(key, value interface{}) bool {
	if _, ok := c.items[key]; ok {
		// The key is already in the cache
		return false
	}

	if c.itemsCount >= c.Capacity {
		c.evict()
	}

	freq := c.freqHeadNode.next
	// frequency list is empty. First insertion
	if freq.value != 1 {
		freq = newFreqNode(1, c.freqHeadNode, freq)
		c.freqHeadNode.next = freq
	}

	node := newCacheNode(key, value, freq)
	freq.Add(node)
	c.items[key] = node
	c.itemsCount++

	return true
}

func (c *Cache) Has(key interface{}) bool {
	_, ok := c.items[key]
	return ok
}

func (c *Cache) FreqKey(key interface{}) (uint, error) {
	node, ok := c.items[key]
	if !ok {
		return 0, nil
	}
	return node.parent.value, nil
}

func (c *Cache) Access(key interface{}) (interface{}, error) {
	node, ok := c.items[key]
	if !ok {
		return nil, errors.New(fmt.Sprintf("element %v is not cached", key))
	}

	freq := node.parent
	nextFreq := freq.next
	if nextFreq == nil || nextFreq == c.freqHeadNode || nextFreq.value != freq.value+1 {
		nextFreq = newFreqNode(freq.value+1, freq, nextFreq)
		freq.next = nextFreq
	}
	node.parent = nextFreq
	nextFreq.Add(node)
	freq.Remove(node)
	if freq.Size() < 1 {
		c.removeNode(freq)
	}
	return node.value, nil
}

func (c *Cache) Size() uint {
	return c.itemsCount
}
