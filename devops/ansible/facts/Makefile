build:
	docker-compose build

up: build
	docker rm -f ansible-master node1 node2 2>/dev/null || true
	docker-compose up -d

down:
	docker-compose down

run:
	docker exec ansible-master ansible-playbook /etc/ansible/playbooks/facts.yml
	
facts:
	docker exec ansible-master ansible all -m setup