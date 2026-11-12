FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum* ./
COPY main.go ./
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o /loan-backend main.go

FROM alpine:3.20
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /loan-backend /loan-backend
EXPOSE 8081
ENTRYPOINT ["/loan-backend"]
