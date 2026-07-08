# Task Manager HTTP API

A simple OCaml HTTP API for task management with CRUD operations.

## Features

- Add new tasks
- Get all tasks
- Mark tasks as complete
- Delete tasks
- JSON API responses

## Prerequisites

- OCaml (>= 4.14)
- Dune (>= 3.0)
- opam

## Installation

```bash
# Install dependencies
opam install dune cohttp-lwt-unix yojson alcotest

# Build the project
dune build

# Run tests
dune test --verbose
```

## Running the Server

```bash
dune exec bin/main.exe
```

The server will start on `http://localhost:8080`

## API Endpoints

## Example Usage

1. Start the server:
```bash
dune exec bin/main.exe
```

2. Add some tasks:
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn OCaml"}'

curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Build HTTP API"}'
```

3. Get all tasks:
```bash
curl -X GET http://localhost:8080/tasks
```

4. Complete a task:
```bash
curl -X PUT http://localhost:8080/tasks/1
```

5. Delete a task:
```bash
curl -X DELETE http://localhost:8080/tasks/2
```

## Response Format

All responses are in JSON format:

**Task Object:**
```json
{
  "id": 1,
  "title": "Task title",
  "completed": false
}
```

**Task List:**
```json
[
  {
    "id": 1,
    "title": "First task",
    "completed": true
  },
  {
    "id": 2,
    "title": "Second task", 
    "completed": false
  }
]
```

## Testing

Run the unit tests:
```bash
dune test
```