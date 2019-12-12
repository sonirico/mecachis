package lru

import (
	"bytes"
	"fmt"
)

type cacheKey interface{}
type cacheValue interface{}

type CacheNode struct {
	key   cacheKey
	value cacheValue

	next *CacheNode
	prev *CacheNode
}

func NewCacheNode() *CacheNode {
	node := &CacheNode{}
	node.prev = node
	node.next = node
	return node
}

func (cn *CacheNode) String() string {
	return fmt.Sprintf("<key: %v, value: %v>", cn.key, cn.value)
}

type Cache struct {
	capacity int

	itemsCount int

	items map[cacheKey]*CacheNode

	head *CacheNode
	foot *CacheNode
}

func NewCache(capacity int) *Cache {
	head := NewCacheNode()
	cache := &Cache{
		capacity:   capacity,
		itemsCount: 0,
		items:      make(map[cacheKey]*CacheNode),
		head:       head,
		foot:       head,
	}
	return cache
}

func (c *Cache) getHead() *CacheNode {
	return c.head.next
}

func (c *Cache) insertHead(node *CacheNode) {
	c.head.prev = node
	node.prev = node
	node.next = c.head
	c.head = node
}

func (c *Cache) Insert(key, value interface{}) bool {
	if _, ok := c.items[key]; ok {
		return false
	}
	node := NewCacheNode()
	node.key = key
	node.value = value
	if c.itemsCount >= c.capacity {
		delete(c.items, c.foot.key)
		foot := c.foot
		c.foot = foot.prev
		foot.prev = nil
		foot.next = nil
		c.itemsCount--
	}
	c.insertHead(node)
	if c.itemsCount == 0 {
		c.foot = node
	}
	c.items[key] = node
	c.itemsCount++

	return true
}

func (c *Cache) Access(key interface{}) (interface{}, error) {
	return nil, nil
}

func (c *Cache) Size() int {
	return c.itemsCount
}

func (c *Cache) Dump() string {
	var buf bytes.Buffer

	item := c.head
	for item != nil {
		buf.WriteString(item.String())
		buf.WriteString("\n")
		item = item.next
	}

	return buf.String()
}
