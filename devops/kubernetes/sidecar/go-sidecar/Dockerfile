FROM golang:1.23-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix crunched -o sidecar .

# Final image
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/sidecar .

# Create non-root user
RUN adduser -D -u 1000 sidecar
USER sidecar

ENTRYPOINT ["./sidecar"]
