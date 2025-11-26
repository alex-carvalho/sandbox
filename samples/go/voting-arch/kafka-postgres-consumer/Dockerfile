FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kafka-postgres-consumer ./cmd/main

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/kafka-postgres-consumer .

CMD ["./kafka-postgres-consumer"]
