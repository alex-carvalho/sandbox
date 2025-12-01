# voting system

- voting service - POST and GET api
- postgres consumer
- redis consumer


## Project Structure

This project follows Go standard project layout conventions.

## Directory Structure

```
voting-arch/
│
├── voting-api/
│   ├── cmd/
│   │   └── main/
│   │       └── main.go                 # Application entry point
│   ├── internal/
│   │   ├── handlers/
│   │   │   ├── vote.go                 # Vote submission handler
│   │   │   └── results.go              # Results retrieval handler
│   │   ├── kafka/
│   │   │   └── producer.go             # Kafka producer implementation
│   │   ├── middleware/
│   │   │   └── middleware.go           # HTTP middleware (panic recovery)
│   │   └── redis/
│   │       └── client.go               # Redis client wrapper
│   ├── pkg/
│   │   └── models/
│   │       └── models.go               # Data models (VoteRequest, VoteResponse)
│   ├── go.mod                          # Go module definition
│   ├── go.sum                          # Go module checksums
│   ├── Dockerfile                      # Docker image definition
│   └── Makefile                        # Build and run targets
│
├── kafka-postgres-consumer/
│   ├── cmd/
│   │   └── main/
│   │       └── main.go                 # Consumer entry point
│   ├── internal/
│   │   ├── consumer/
│   │   │   └── consumer.go             # Kafka consumer logic
│   │   └── database/
│   │       └── postgres.go             # PostgreSQL operations
│   ├── pkg/
│   │   └── models/
│   │       └── models.go               # Data models (Vote)
│   ├── go.mod                          # Go module definition
│   ├── go.sum                          # Go module checksums
│   ├── Dockerfile                      # Docker image definition
│   └── Makefile                        # Build and run targets
│
├── kafka-redis-consumer/
│   ├── cmd/
│   │   └── main/
│   │       └── main.go                 # Consumer entry point
│   ├── internal/
│   │   ├── cache/
│   │   │   └── redis.go                # Redis cache operations
│   │   └── consumer/
│   │       └── consumer.go             # Kafka consumer logic
│   ├── pkg/
│   │   └── models/
│   │       └── models.go               # Data models (Vote)
│   ├── go.mod                          # Go module definition
│   ├── go.sum                          # Go module checksums
│   ├── Dockerfile                      # Docker image definition
│   └── Makefile                        # Build and run targets
│
├── docker-compose.yml                  # Docker Compose configuration
├── init.sql                            # PostgreSQL initialization script
├── .env.example                        # Environment variables template
├── README.md                           # Project documentation
└── Makefile                            # Root orchestration targets

```

## Package Organization

### `cmd/` - Command/Application Entry Points
- Contains main packages for standalone applications
- Minimal logic - delegates to internal packages
- Specific to the application binary

### `internal/` - Private Packages
- Packages that should not be imported by external projects
- Go tooling enforces this restriction
- Organized by functionality:
  - `handlers/` - HTTP handlers
  - `kafka/` - Kafka producer/consumer
  - `middleware/` - HTTP middleware
  - `redis/` - Redis client wrapper
  - `cache/` - Cache operations
  - `consumer/` - Message consumer logic
  - `database/` - Database operations

### `pkg/` - Public Packages
- Can be imported by external projects
- Contains reusable code
- `models/` - Data structures and types

## Naming Conventions

- Package names are lowercase, single word when possible
- File names use snake_case for clarity (e.g., `producer.go`, `middleware.go`)
- Keep package names short and descriptive
- Avoid unnecessary abbreviations

## Import Paths

```go
// Correct - using module path
import "github.com/alex-carvalho/voting-api/internal/handlers"
import "github.com/alex-carvalho/voting-api/pkg/models"

// Correct - standard library
import "net/http"
import "encoding/json"
```

## Building

Each service can be built independently:

```bash
cd voting-api && make build
cd kafka-postgres-consumer && make build
cd kafka-redis-consumer && make build
```

Or from root:

```bash
make deps    # Download all dependencies
make build   # Build all Docker images
```
