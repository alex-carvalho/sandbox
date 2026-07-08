.PHONY: build run clean test deps

build:
	mkdir -p bin
	go build -o bin/kafka-redis-consumer ./cmd/main

run: build
	./bin/kafka-redis-consumer

clean:
	rm -f bin/kafka-redis-consumer

test:
	go test ./...

deps:
	go mod download
	go mod tidy
