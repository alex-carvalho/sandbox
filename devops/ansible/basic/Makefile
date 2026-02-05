check: run
	@echo "=== Checking nginx installation ==="
	@echo "Testing node1 (http://localhost:8086)..."
	@curl -s -o /dev/null -w "Status: %{http_code}\n" http://localhost:8086
	@echo "Testing node2 (http://localhost:8087)..."
	@curl -s -o /dev/null -w "Status: %{http_code}\n" http://localhost:8087
	$(MAKE) clean

run: setup
	@echo "Running Ansible playbook..."
	ansible-playbook -i inventory.ini playbook.yml

setup: clean
	./setup.sh

clean:
	@echo "Removing Docker containers..."
	@docker rm -f node1 node2 2>/dev/null || true

