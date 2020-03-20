package main

import (
	"flag"
	"fmt"
	"github.com/sonirico/mecachis"
	"github.com/sonirico/mecachis/engines"
	"net/http"
)

func main() {
	var port int
	flag.IntVar(&port, "http", 8000, "http port")
	flag.Parse()

	cache := mecachis.NewCache(10, engines.LRU)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), cache)
	if err != nil {
		panic(err)
	}
}
