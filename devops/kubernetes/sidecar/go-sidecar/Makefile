.PHONY: build test run clean docker

# Build the Go binary
build:
	go build -o bin/sidecar .

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run local test (connects to real K8s cluster using kubeconfig)
local-test:
	PROPERTIES_FILE=/tmp/application.properties POD_NAMESPACE=default CONFIGMAP_NAME=properties go run ./cmd/local-test
	kubectl get configmap properties -n default -o yaml

# Run the application locally (requires Kubernetes cluster)
run: build
	PROPERTIES_FILE=/etc/config/application.properties POD_NAMESPACE=default CONFIGMAP_NAME=properties ./bin/sidecar

# Clean build artifacts
clean:
	rm -rf bin/ coverage.out coverage.html

# Docker operations
docker-build:
	docker build -t go-sidecar:latest .

# Lint code
lint:
	go fmt ./...
	go vet ./...

# Download dependencies
deps:
	go mod download
	go mod tidy

# Show help
help:
	@echo "Available targets:"
	@echo "  make build         - Build the Go binary"
	@echo "  make test          - Run tests with coverage"
	@echo "  make local-test    - Run sidecar locally (connects to K8s cluster)"
	@echo "  make run           - Run locally (requires K8s cluster)"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make lint          - Run linters"
	@echo "  make deps          - Download and tidy dependencies"
