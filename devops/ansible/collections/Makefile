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
