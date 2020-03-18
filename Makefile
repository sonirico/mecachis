.PHONY: test clean format

test:
	go test ./...

format:
	go fmt ./...

clean:
	go clean -modcache
	go clean -testcache
