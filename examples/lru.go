package main

import (
	"fmt"
	"io/ioutil"
	"mecachis/lru"
	"net/http"
	"time"
)

type CacheableFunc func(k interface{}) interface{}

func withLRUCache(capacity uint, fn CacheableFunc) (CacheableFunc, func()) {
	cache := lru.NewCache(capacity)
	return func(key interface{}) interface{} {
			value, err := cache.Access(key)
			if err != nil {
				value = fn(key)
				cache.Insert(key, value)
				return value
			}
			return value
		}, func() {
			cache.Free()
		}
}

func main() {
	requestCached, free := withLRUCache(2, func(web interface{}) interface{} {
		w, _ := web.(string)
		resp, _ := http.Get(w)
		bytes, _ := ioutil.ReadAll(resp.Body)
		return len(bytes)
	})
	webs := []string{
		"https://github.com/sonirico/mecachis",
		"https://github.com/sonirico/node.go",
		"https://github.com/sonirico/mecachis",
		"https://github.com/sonirico/node.go",
		"https://github.com/sonirico/wpoke",
	}
	for _, web := range webs {
		t := time.Now()
		fmt.Println(fmt.Sprintf("web: %s, content-length: %d", web, requestCached(web)))
		fmt.Println(fmt.Sprintf("took=%v", time.Since(t)))
	}
	free()
}
