build:
	docker-compose build

up: build
	docker rm -f ansible-master proxy app-blue app-green 2>/dev/null || true
	docker-compose up -d

down:
	docker-compose down

setup: up
	docker exec ansible-master ansible-playbook /etc/ansible/playbooks/setup.yml

deploy:
	docker exec ansible-master ansible-playbook /etc/ansible/playbooks/deploy.yml

check:
	@echo "=== Active slot (X-Active-Slot header) ==="
	@curl -sI http://localhost:8080 | grep -i x-active-slot
	@echo ""
	@echo "=== HTTP status ==="
	@curl -s -o /dev/null -w "HTTP status: %{http_code}\n" http://localhost:8080

