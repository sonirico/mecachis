package main

import (
	"fmt"
	"io/ioutil"
	"mecachis/lfu"
	"net/http"
	"time"
)

type cacheableFuncLFU func(k string) int

func withLFUCache(capacity uint, fn cacheableFuncLFU) (cacheableFuncLFU, func()) {
	cache := lfu.NewCache(capacity)
	free := func() { cache.Free() }
	decorator := func(key string) int {
		value, err := cache.Access(key)
		if err != nil {
			value = fn(key)
			cache.Insert(key, value)
		}
		intValue, _ := value.(int)
		return intValue
	}
	return decorator, free
}

func main() {
	requestCached, free := withLFUCache(2, func(web string) int {
		resp, _ := http.Get(web)
		bytes, _ := ioutil.ReadAll(resp.Body)
		return len(bytes)
	})
	webs := []string{
		"https://github.com/sonirico/mecachis",
		"https://github.com/sonirico/node.go",
		"https://github.com/sonirico/node.go",
		"https://github.com/sonirico/node.go",
		"https://github.com/sonirico/datetoken",
		"https://github.com/sonirico/paranoid",
		"https://github.com/sonirico/go-fist",
		// Note that node.go should not be evicted yet
		"https://github.com/sonirico/node.go",
	}
	for _, web := range webs {
		t := time.Now()
		fmt.Println(fmt.Sprintf("web: %s, content-length: %d", web, requestCached(web)))
		fmt.Println(fmt.Sprintf("took=%v", time.Since(t)))
	}
	free()
}
