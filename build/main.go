package main

import (
	"flag"
	"fmt"
	"github.com/sonirico/mecachis"
	"net/http"
)

func main() {
	var port int
	flag.IntVar(&port, "http", 8000, "http port")
	flag.Parse()

	hub := mecachis.NewHub()
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), hub)
	if err != nil {
		panic(err)
	}
}
