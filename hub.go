package mecachis

import (
	"sync"
)

type Hub struct {
	mx     sync.RWMutex
	groups map[string]*group
}

func NewHub() *Hub {
	return &Hub{
		groups: make(map[string]*group),
	}
}

func (h *Hub) group(name string) (*group, bool) {
	h.mx.RLock()
	defer h.mx.RUnlock()
	g, ok := h.groups[name]
	return g, ok
}

func (h *Hub) getOrCreateGroup(name string) (*group, bool) {
	if g, ok := h.group(name); ok {
		return g, false
	}
	g := newGroup(name)
	h.mx.Lock()
	h.groups[name] = g
	h.mx.Unlock()
	return g, true
}
