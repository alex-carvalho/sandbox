# Bitonic Sequence Generator

## Overview
Generates bitonic sequences of length n from integers in a given range. A bitonic sequence first increases monotonically, then decreases monotonically.

## Algorithm
The implementation divides the sequence into two parts:
- Increasing part: from min to max (first half)
- Decreasing part: from max back down (second half)

## Build & Run

### Prerequisites
```bash
opam install dune dream yojson ounit2
```

### Build
```bash
dune build
```

### Run Tests
```bash
dune exec test/test_bitonic.exe
```

### Run API Server
```bash
dune exec bitonic_sequence
```
Server runs on `http://localhost:8080`

## API Endpoints

### POST /bitonic
Generate a bitonic sequence.

**Request:**
```json
{
  "n": 10,
  "min": 1,
  "max": 100
}
```

**Response:**
```json
{
  "sequence": [1, 23, 45, 67, 89, 100, 78, 56, 34, 12],
  "length": 10,
  "is_bitonic": true
}
```

### GET /health
Health check endpoint.

**Response:**
```json
{
  "status": "ok"
}
```

## Performance Testing

### Run k6 Test
```bash
k6 run k6/load_test.js
```

## Redis Setup

### Start Redis
```bash
docker-compose up -d
```

### Stop Redis
```bash
docker-compose down
```

## Examples

```bash
# Generate sequence
curl -X POST http://localhost:8080/bitonic \
  -H "Content-Type: application/json" \
  -d '{"n": 7, "min": 5, "max": 50}'
```
