CLUSTER_NAME = ansible-poc
KUBECONFIG   = $(PWD)/kubeconfig

cluster:
	kind create cluster --config cluster-config.yaml
	kind export kubeconfig --name $(CLUSTER_NAME) --kubeconfig $(KUBECONFIG)
	@echo "Cluster ready. Kubeconfig: $(KUBECONFIG)"

init:
	ansible-galaxy collection install -r requirements.yml
	/opt/homebrew/bin/python3.14 -m pip install kubernetes --break-system-packages

run:
	ANSIBLE_CONFIG=$(PWD)/ansible.cfg ansible-playbook playbooks/site.yml

down:
	kind delete cluster --name $(CLUSTER_NAME)
	rm -f $(KUBECONFIG)

check:
	@echo "=== Cluster nodes ==="
	kubectl --kubeconfig $(KUBECONFIG) get nodes
	@echo ""
	@echo "=== Pods ==="
	kubectl --kubeconfig $(KUBECONFIG) get pods -n default
	@echo ""
	@echo "=== Service ==="
	kubectl --kubeconfig $(KUBECONFIG) get svc nginx-poc
	@echo ""
	@echo "=== HTTP check (http://localhost:8080) ==="
	@curl -s -o /dev/null -w "HTTP status: %{http_code}\n" http://localhost:8080

status:
	kubectl --kubeconfig $(KUBECONFIG) get all -n default
