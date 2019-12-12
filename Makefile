.PHONY: test clean

test:
	go test ./...

format:
	go fmt ./...

clean:
	go clean -modcache
	go clean -testcache
