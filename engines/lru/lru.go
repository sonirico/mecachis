package engines

import (
	"github.com/sonirico/mecachis/engines"

	"container/list"
)

// cache represents the lru cache
type lru struct {
	// how much capacity in bytes
	capacity  uint64
	size      uint64
	list      *list.List
	cache     map[string]*list.Element
	onEvicted engines.EvictionFn
}

// New initializes a new cache by providing the maximum
// capacity which, once reached, will provoke to evict the lru
// element
func New(capacity uint64) *lru {
	return &lru{
		capacity: capacity,
		size:     0,
		list:     list.New(),
		cache:    make(map[string]*list.Element),
	}
}

func (c *lru) OnEvict(onEvicted engines.EvictionFn) {
	c.onEvicted = onEvicted
}

func (c *lru) evict() {
	el := c.list.Back()
	if el == nil {
		return
	}
	c.list.Remove(el)
	entry := el.Value.(engines.Entry)
	delete(c.cache, entry.Key())
	c.size -= entry.Len()
	if c.onEvicted != nil {
		c.onEvicted(entry)
	}
	return
}

// Insert puts a key-value pair into the cache. Returns whether the pair
// was inserted. `false` means that the element was cached already
func (c *lru) Insert(key string, value engines.Value) bool {
	if el, ok := c.cache[key]; ok {
		c.list.MoveToFront(el)
		return false
	}
	if c.capacity > 0 {
		// Limit configured
		if c.size == 0 || c.size >= c.capacity {
			c.evict()
		}
	}
	entry := engines.NewEntry(key, value)
	el := c.list.PushFront(entry)
	c.cache[key] = el
	c.size += entry.Len()
	return true
}

// Access returns an element by key if it is within the cache already. Otherwise
// it returns an error
func (c *lru) Access(key string) (engines.Value, bool) {
	el, ok := c.cache[key]
	if !ok {
		return nil, ok
	}
	c.list.MoveToFront(el)
	entry := el.Value.(engines.Entry)
	return entry.Value(), true
}

// Size returns the current length of the cache
func (c *lru) Size() uint64 {
	return c.size
}

// Dump returns the current state of the cache
func (c *lru) Dump() []engines.Entry {
	var result []engines.Entry
	el := c.list.Front()
	for el != nil {
		entry := el.Value.(engines.Entry)
		result = append(result, entry)
		el = el.Next()
	}
	return result
}

// Free empties the cache, leaving it with the initial state
func (c *lru) Free() {
	for k, _ := range c.cache {
		delete(c.cache, k)
	}
	c.list.Init()
	c.size = 0
}
