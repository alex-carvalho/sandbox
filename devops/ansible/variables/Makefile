build:
	docker-compose build

up: build
	docker rm -f ansible-master node1 node2 2>/dev/null || true
	docker-compose up -d

down:
	docker-compose down

run:
	docker exec ansible-master ansible-playbook /etc/ansible/playbooks/apache.yml

ping:
	docker exec ansible-master ansible webservers -m ping

shell:
	docker exec -it ansible-master /bin/bash

logs:
	docker-compose logs -f
