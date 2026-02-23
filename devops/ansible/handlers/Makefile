build:
	docker-compose build

up: build
	docker rm -f ansible-master node1 node2 2>/dev/null || true
	docker-compose up -d

down:
	docker-compose down

run: up
	docker exec ansible-master ansible-playbook /etc/ansible/playbooks/handlers.yml

check:
	@echo "=== Checking apache installation ==="
	@echo "Testing node1 (http://localhost:8001)..."
	@curl -s -o /dev/null -w "Status: %{http_code}\n" http://localhost:8001
	@echo "Testing node2 (http://localhost:8002)..."
	@curl -s -o /dev/null -w "Status: %{http_code}\n" http://localhost:8002

logs:
	docker-compose logs -f
