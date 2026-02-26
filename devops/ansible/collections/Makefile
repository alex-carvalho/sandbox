build:
	docker-compose build

up: build
	docker rm -f ansible-master node1 node2 2>/dev/null || true
	docker-compose up -d

down:
	docker-compose down

init:
	docker exec ansible-master ansible-galaxy collection install -r /etc/ansible/requirements.yml

run:
	docker exec ansible-master ansible-playbook /etc/ansible/playbooks/site.yml

list:
	docker exec ansible-master ansible-galaxy collection list

check:
	@echo "=== node1 (https://localhost:8443) ==="
	@curl -sk -o /dev/null -w "HTTP status: %{http_code}\n" https://localhost:8443
	@curl -sk --cert-status -w "SSL: %{ssl_verify_result} | Cert expires: %{ssl_cert_expire}\n" https://localhost:8443 -o /dev/null 2>/dev/null || \
		openssl s_client -connect localhost:8443 -servername node1 </dev/null 2>/dev/null | openssl x509 -noout -subject -dates
	@echo "=== node2 (https://localhost:8444) ==="
	@curl -sk -o /dev/null -w "HTTP status: %{http_code}\n" https://localhost:8444
	@curl -sk --cert-status -w "SSL: %{ssl_verify_result} | Cert expires: %{ssl_cert_expire}\n" https://localhost:8444 -o /dev/null 2>/dev/null || \
		openssl s_client -connect localhost:8444 -servername node2 </dev/null 2>/dev/null | openssl x509 -noout -subject -dates
