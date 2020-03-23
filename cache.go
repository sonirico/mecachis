package mecachis

import (
	e "github.com/sonirico/mecachis/engines"
	lru "github.com/sonirico/mecachis/engines/lru"
	"sync"
)

const (
	basePath = "/mecachis/"
)

type Cache interface {
	Add(k string, v MemoryView) error
	Get(k string) (MemoryView, bool)
}

type cache struct {
	sync.RWMutex

	engine e.Engine
}

func NewCache(cap uint64, cType e.CacheType) *cache {
	return &cache{
		engine: newEngine(cType, cap),
	}
}

func (c *cache) Add(key string, value MemoryView) error {
	c.Lock()
	defer c.Unlock()
	res := c.engine.Insert(key, value)
	if !res {
		return NewDuplicatedKeyError(key)
	}
	return nil
}

func (c *cache) Get(key string) (MemoryView, bool) {
	c.RLock()
	defer c.RUnlock()
	res, ok := c.engine.Access(key)
	if !ok {
		return nil, false
	}
	data := res.(MemoryView)
	return data, ok
}

func newEngine(cType e.CacheType, capacity uint64) e.Engine {
	switch cType {
	case e.LRU:
		return lru.New(capacity)
	}
	return nil
}
