.PHONY: build run clean test deps

build:
	go build -o bin/kafka-postgres-consumer ./cmd/main

run: build
	./bin/kafka-postgres-consumer

clean:
	rm -f bin/kafka-postgres-consumer

test:
	go test ./...

deps:
	go mod download
	go mod tidy
