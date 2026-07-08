.PHONY: build up down logs clean test deps

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down --volumes

logs:
	docker-compose logs -f

logs-api:
	docker-compose logs -f voting-api

logs-postgres:
	docker-compose logs -f kafka-postgres-consumer

logs-redis:
	docker-compose logs -f kafka-redis-consumer

clean:
	docker-compose down -v
	find . -name "*.o" -o -name "*_test" | xargs rm -f

deps:
	cd voting-api && go mod download && go mod tidy
	cd kafka-postgres-consumer && go mod download && go mod tidy
	cd kafka-redis-consumer && go mod download && go mod tidy

test:
	cd voting-api && go test -v ./...

status:
	docker-compose ps

shell-api:
	docker-compose exec voting-api sh

shell-postgres:
	docker-compose exec postgres psql -U postgres -d voting

shell-redis:
	docker-compose exec redis redis-cli

total_db:
	docker-compose exec postgres psql -U postgres -d voting -c "SELECT COUNT(*) FROM votes;"

total_api:
	curl -s http://localhost:8081/results?voting_id=100  | jq