package lru

import (
	"bytes"
	"errors"
	"fmt"
)

type cacheKey interface{}
type cacheValue interface{}

type cacheNode struct {
	key   cacheKey
	value cacheValue

	next *cacheNode
	prev *cacheNode
}

func newCacheNode(key cacheKey, value cacheValue) *cacheNode {
	return &cacheNode{key: key, value: value}
}

func (cn *cacheNode) String() string {
	return fmt.Sprintf("<key: %v, value: %v>", cn.key, cn.value)
}

// Cache represents the LRU cache
type Cache struct {
	// Maximum capacity
	Capacity uint
	// How many elements are cached
	itemsCount uint
	// cache nodes hash-map
	items map[cacheKey]*cacheNode
	// pointer to the last inserted/access node
	head *cacheNode
	// pointer to the least recently updated node
	foot *cacheNode
}

// NewCache initializes a new Cache by providing the maximum
// capacity which, once reached, will provoke to evict the LRU
// element
func NewCache(capacity uint) *Cache {
	cache := &Cache{
		Capacity:   capacity,
		itemsCount: 0,
		items:      make(map[cacheKey]*cacheNode, capacity),
	}
	return cache
}

func (c *Cache) getHead() *cacheNode {
	return c.head.next
}

// Inserts a node at the top
func (c *Cache) insertNode(node *cacheNode) {
	if c.head == nil {
		c.head = node
		c.foot = node
	} else {
		c.head.prev = node
		node.next = c.head
		c.head = node
	}
	c.items[node.key] = node
	c.itemsCount++
}

func (c *Cache) getNodes() []cacheNode {
	var nodes []cacheNode
	item := c.head
	for item != nil {
		nodes = append(nodes, *item)
		item = item.next
	}
	return nodes
}

func (c *Cache) evict() {
	k := c.foot.key
	nextFoot := c.foot.prev
	c.removeNode(c.foot)
	c.foot = nextFoot
	c.itemsCount--
	delete(c.items, k)
}

func (c *Cache) removeNode(node *cacheNode) {
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
}

// Insert puts a key-value pair into the cache. Returns whether the pair
// was inserted. `false` means that the element was cached already
func (c *Cache) Insert(key, value interface{}) bool {
	if _, ok := c.items[key]; ok {
		return false
	}
	if c.itemsCount >= c.Capacity {
		c.evict()
	}
	node := newCacheNode(key, value)
	c.insertNode(node)
	return true
}

// Access returns an element by key if it is within the cache already. Otherwise
// it returns an error
func (c *Cache) Access(key interface{}) (interface{}, error) {
	node, ok := c.items[key]
	if !ok {
		return nil, errors.New(fmt.Sprintf("non existent key \"%v\"", key))
	}
	if node != c.head {
		// Skip reindex if requested key was the last one requested as well
		c.removeNode(node)
		c.insertNode(node)
	}
	return node.value, nil
}

// Size returns the current length of the cache
func (c *Cache) Size() uint {
	return c.itemsCount
}

// Nodes returns the list of nodes ordered. Highest prioritized elements first
func (c *Cache) Nodes() []cacheNode {
	return c.getNodes()
}

// Dump returns the current state of the cache as string
func (c *Cache) Dump() string {
	var buf bytes.Buffer

	for _, item := range c.getNodes() {
		buf.WriteString(item.String())
		buf.WriteString("\n")
	}

	return buf.String()
}

// Free empties the cache, leaving it with the initial state
func (c *Cache) Free() {
	c.itemsCount = 0
	c.items = make(map[cacheKey]*cacheNode, c.Capacity)
	item := c.head
	for item != nil {
		next := item.next
		item.prev = nil
		item.next = nil
		item = next
	}
	c.head = nil
	c.foot = nil
}
