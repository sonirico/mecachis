package mecachis

import (
	e "github.com/sonirico/mecachis/engines"
	lru "github.com/sonirico/mecachis/engines/lru"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	basePath = "/mecachis/"
)

type Cache interface {
	Add(k string, v *MemoryView) error
	Get(k string) (*MemoryView, bool)
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

func (c *cache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	log.Printf("%s: %s\n", r.Method, uri)
	if !strings.HasPrefix(uri, basePath) {
		http.NotFound(w, r)
		return
	}
	uriParts := strings.SplitN(r.URL.Path[len(basePath):], "/", 2)
	log.Println(uriParts)
	if len(uriParts) < 1 {
		http.NotFound(w, r)
		return
	}
	if r.Method == http.MethodGet {
		key := uriParts[0]
		value, ok := c.Get(key)
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/octet-stream")
		if _, err := w.Write(value.Clone()); err != nil {
			log.Printf(err.Error())
			http.Error(w, "error when writing to response buffer", http.StatusInternalServerError)
		}
		return
	}
	if r.Method == http.MethodPost {
		key := uriParts[0]
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf(err.Error())
			http.Error(w, "error when reading request buffer", http.StatusInternalServerError)
			return
		}
		if err := c.Add(key, content); err != nil {
			log.Printf(err.Error())
			http.Error(w, "error when writing to response buffer", http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	}
	http.NotFound(w, r)
}

func newEngine(cType e.CacheType, capacity uint64) e.Engine {
	switch cType {
	case e.LRU:
		return lru.New(capacity)
	}
	return nil
}
