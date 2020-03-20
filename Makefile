.PHONY: test clean format build

PORT ?= 8000

test:
	go test ./server ./engines ./container

format:
	go fmt ./...

clean:
	go clean -modcache
	go clean -testcache
	rm -rf ./bin/*

build:
	mkdir -p ./bin
	go build -o ./bin/server.bin ./build
	./bin/server.bin -http $(PORT)
