build:
	docker-compose build

up: build
	docker-compose up -d

down:
	docker-compose down

list:
	docker exec ansible-master ansible-inventory --list

run:
	docker exec ansible-master ansible-playbook /etc/ansible/playbooks/webservers.yml

check:
	@echo "=== Checking nginx installation ==="
	@echo "Testing node1 (http://localhost:8001)..."
	@curl -s -o /dev/null -w "Status: %{http_code}\n" http://localhost:8001
	@echo "Testing node2 (http://localhost:8002)..."
	@curl -s -o /dev/null -w "Status: %{http_code}\n" http://localhost:8002
