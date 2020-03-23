package mecachis

import (
	"github.com/sonirico/mecachis/engines"
	"sync"
)

type group struct {
	mx sync.RWMutex

	Ns    string
	Cap   uint64
	Ct    engines.CacheType
	cache *cache
}

func newGroup(name string) *group {
	g := &group{Ns: name}
	return g
}

func (g *group) Add(k string, v MemoryView) error {
	if g.cache == nil {
		g.cache = NewCache(g.Cap, g.Ct)
	}
	return g.cache.Add(k, v)
}

func (g *group) Get(k string) (MemoryView, bool) {
	if g.cache == nil {
		g.cache = NewCache(g.Cap, g.Ct)
	}
	return g.cache.Get(k)
}
