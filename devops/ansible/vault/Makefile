build:
	docker-compose build

up: build
	docker rm -f ansible-master node1 node2 2>/dev/null || true
	docker-compose up -d

down:
	docker-compose down

init:
	docker exec ansible-master bash -c "echo 'ansible' > /etc/ansible/.vault_pass && chmod 600 /etc/ansible/.vault_pass && echo '.vault_pass created'"

run: init
	docker exec ansible-master ansible-playbook /etc/ansible/playbooks/site.yml

view: init
	docker exec ansible-master ansible-vault view /etc/ansible/group_vars/all/vault.yml \
		--vault-password-file /etc/ansible/.vault_pass

check:
	@echo "=== Checking apache installation ==="
	@echo "Testing node1 (http://localhost:8001)..."
	@curl -s -o /dev/null -w "Status: %{http_code}\n" http://localhost:8001
	@echo "Testing node2 (http://localhost:8002)..."
	@curl -s -o /dev/null -w "Status: %{http_code}\n" http://localhost:8002
