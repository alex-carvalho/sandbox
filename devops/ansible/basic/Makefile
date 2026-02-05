check: run
	@echo "=== Checking nginx installation ==="
	curl http://localhost:8086
	curl http://localhost:8087
	$(MAKE) clean

run: setup
	@echo "Running Ansible playbook..."
	ansible-playbook -i inventory.ini playbook.yml

setup: clean
	./setup.sh

clean:
	@echo "Removing Docker containers..."
	@docker rm -f node1 node2 2>/dev/null || true

