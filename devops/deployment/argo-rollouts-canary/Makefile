CLUSTER_NAME  := argo-rollouts-demo
NAMESPACE     := demo
MONITORING_NS := monitoring
ROLLOUT_NAME  := my-app
RELEASE_NAME  := my-app
APP_IMAGE     := my-app:latest

setup: clean build-image
	@echo "\n==> Creating Kind cluster '$(CLUSTER_NAME)'..."
	kind create cluster --name $(CLUSTER_NAME)

	@echo "\n==> Installing Argo Rollouts..."
	kubectl create namespace argo-rollouts
	kubectl apply -n argo-rollouts \
		-f https://github.com/argoproj/argo-rollouts/releases/latest/download/install.yaml
	
	sleep 30
	kubectl wait --namespace argo-rollouts \
		--for=condition=ready pod \
		--selector=app.kubernetes.io/name=argo-rollouts \
		--timeout=180s

	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo update

	kubectl create namespace $(MONITORING_NS)
	helm upgrade --install kube-prometheus-stack prometheus-community/kube-prometheus-stack \
		--namespace $(MONITORING_NS) \
		--values prometheus/kube-prometheus-stack-values.yaml \
		--wait --timeout 5m

	@echo "\n==> Loading app image into Kind cluster..."
	kind load docker-image $(APP_IMAGE) --name $(CLUSTER_NAME)

	@echo "\n  Setup complete. Run 'make deploy' to deploy the app."

build-image:
	docker build -t $(APP_IMAGE) ./app

build-load-image: build-image
	kind load docker-image $(APP_IMAGE) --name $(CLUSTER_NAME)

clean:
	kind delete cluster --name $(CLUSTER_NAME)

deploy:
	helm upgrade --install $(RELEASE_NAME) ./chart \
		--namespace $(NAMESPACE) \
		--create-namespace \
		--force-conflicts \
		--wait
	
	helm list -n $(NAMESPACE)
	kubectl get pods -n $(NAMESPACE)

canary-bad:
	helm upgrade $(RELEASE_NAME) ./chart \
		--namespace $(NAMESPACE) \
		--set dbUrl="" \
		--force-conflicts
	@echo "\n▶  Watch progress: make status"

canary-good:
	helm upgrade $(RELEASE_NAME) ./chart \
		--namespace $(NAMESPACE) \
		--set dbUrl="postgres://user:pass@db:5432/mydb" \
		--force-conflicts
	@echo "\n▶  Watch progress: make status"

promote:
	kubectl argo rollouts promote $(ROLLOUT_NAME) -n $(NAMESPACE)

curl:
	kubectl run curl-debug -it --rm --image=curlimages/curl --restart=Never -- \
		sh -c 'while true; do curl -s http://my-app.demo.svc.cluster.local:8080/api; echo; sleep 1; done'

status:
	kubectl argo rollouts get rollout $(ROLLOUT_NAME) -n $(NAMESPACE) --watch

dashboard:
	kubectl argo rollouts dashboard &
	sleep 2
	open http://localhost:3100/rollouts 2>/dev/null || true

