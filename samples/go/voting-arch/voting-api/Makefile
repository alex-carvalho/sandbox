.PHONY: build run clean test deps

build:
	mkdir -p bin
	go build -o bin/voting-api ./cmd/main

run: build
	./bin/voting-api

clean:
	rm -f bin/voting-api

test:
	go test ./...

deps:
	go mod download
	go mod tidy
