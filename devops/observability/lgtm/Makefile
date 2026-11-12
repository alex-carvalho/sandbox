CLUSTER_NAME ?= lgtm-poc
NAMESPACE_OBS ?= observability
NAMESPACE_APP ?= loan-app

.PHONY: all cluster-up helm-repos deploy-lgtm build-load-apps deploy-apps status verify clean

all: cluster-up helm-repos deploy-lgtm build-load-apps deploy-apps status

cluster-up:
	@if kind get clusters | grep -q "^$(CLUSTER_NAME)$$"; then \
		echo "Kind cluster '$(CLUSTER_NAME)' already exists."; \
	else \
		kind create cluster --name $(CLUSTER_NAME) --config kind/cluster-config.yaml; \
	fi
	kubectl cluster-info --context kind-$(CLUSTER_NAME)

helm-repos:
	helm repo add grafana https://grafana.github.io/helm-charts
	helm repo add grafana-community https://grafana-community.github.io/helm-charts
	helm repo update

deploy-lgtm:
	helm repo add grafana https://grafana.github.io/helm-charts 2>/dev/null || true
	helm repo add grafana-community https://grafana-community.github.io/helm-charts 2>/dev/null || true
	helm repo update
	kubectl create namespace $(NAMESPACE_OBS) --dry-run=client -o yaml | kubectl apply -f -
	# 1. Deploy Mimir (Single-Binary Monolithic with Exemplars)
	helm upgrade --install mimir ./helm/mimir-monolithic \
		-n $(NAMESPACE_OBS) \
		-f ./helm/mimir-values.yaml \
		--wait --timeout 5m
	# 2. Deploy Loki (SingleBinary with TSDB v13 & Structured Metadata)
	helm upgrade --install loki grafana/loki \
		-n $(NAMESPACE_OBS) \
		-f ./helm/loki-values.yaml \
		--wait --timeout 5m
	# 3. Deploy Tempo (SingleBinary via grafana-community/tempo with Metrics Generator -> Mimir)
	helm upgrade --install tempo grafana-community/tempo \
		-n $(NAMESPACE_OBS) \
		-f ./helm/tempo-values.yaml \
		--wait --timeout 5m
	# 4. Deploy Grafana Alloy (OTLP Collector -> Loki, Tempo, Mimir)
	helm upgrade --install alloy grafana/alloy \
		-n $(NAMESPACE_OBS) \
		-f ./helm/alloy-values.yaml \
		--wait --timeout 5m
	# 5. Provision Dashboard ConfigMap & Deploy Grafana (via grafana-community/grafana)
	kubectl create configmap loan-observability-dashboards \
		-n $(NAMESPACE_OBS) \
		--from-file=loan-observability-dashboard.json=./dashboards/loan-observability-dashboard.json \
		--dry-run=client -o yaml | kubectl apply -f -
	helm upgrade --install grafana grafana-community/grafana \
		-n $(NAMESPACE_OBS) \
		-f ./helm/grafana-values.yaml \
		--wait --timeout 5m
	kubectl rollout restart deployment/alloy deployment/grafana -n $(NAMESPACE_OBS)
	kubectl rollout status deployment/alloy -n $(NAMESPACE_OBS) --timeout=2m
	kubectl rollout status deployment/grafana -n $(NAMESPACE_OBS) --timeout=2m

build-load-apps:
	docker build -t loan-backend:local ./apps/loan-backend
	docker build -t loan-underwriting:local ./apps/loan-underwriting
	docker build -t loan-frontend:local ./apps/loan-frontend
	kind load docker-image loan-backend:local --name $(CLUSTER_NAME)
	kind load docker-image loan-underwriting:local --name $(CLUSTER_NAME)
	kind load docker-image loan-frontend:local --name $(CLUSTER_NAME)

deploy-apps:
	kubectl apply -k ./apps/k8s
	kubectl rollout status deployment/postgres -n $(NAMESPACE_APP) --timeout=3m
	kubectl rollout restart deployment/loan-underwriting deployment/loan-backend deployment/loan-frontend deployment/loan-traffic-gen -n $(NAMESPACE_APP)
	kubectl rollout status deployment/loan-underwriting -n $(NAMESPACE_APP) --timeout=3m
	kubectl rollout status deployment/loan-backend -n $(NAMESPACE_APP) --timeout=3m
	kubectl rollout status deployment/loan-frontend -n $(NAMESPACE_APP) --timeout=3m
	kubectl rollout status deployment/loan-traffic-gen -n $(NAMESPACE_APP) --timeout=3m

status:
	@echo "=== Observability Namespace Pods ==="
	kubectl get pods -n $(NAMESPACE_OBS)
	@echo ""
	@echo "=== Loan App Namespace Pods ==="
	kubectl get pods -n $(NAMESPACE_APP)
	@echo ""
	@echo "Grafana UI:        http://localhost:3000 (admin / admin)"
	@echo "Loan Dashboard:    http://localhost:3000/d/loan-lgtm-poc"
	@echo "Loan Frontend UI:  http://localhost:8080"

verify:
	@echo "1. Testing Normal Approval Scenario..."
	curl -s -X POST http://localhost:8080/api/loans/apply \
		-H "Content-Type: application/json" \
		-d '{"applicant_name":"Alex Verification","loan_type":"home","amount":300000,"term_months":360,"annual_income":160000,"credit_score":775,"simulate_mode":"normal"}' | python3 -m json.tool
	@echo "2. Testing Slow PostgreSQL Query Scenario (amount > 900k)..."
	curl -s -X POST http://localhost:8080/api/loans/apply \
		-H "Content-Type: application/json" \
		-d '{"applicant_name":"Sofia HighValue","loan_type":"home","amount":960000,"term_months":360,"annual_income":320000,"credit_score":780,"simulate_mode":"slow_db"}' | python3 -m json.tool
	@echo "3. Testing Fraud Blacklist Error Scenario (403)..."
	curl -s -X POST http://localhost:8080/api/loans/apply \
		-H "Content-Type: application/json" \
		-d '{"applicant_name":"Victor Fraud Blacklist","loan_type":"personal","amount":50000,"term_months":36,"annual_income":90000,"credit_score":690,"simulate_mode":"fraud_error"}' | python3 -m json.tool

clean:
	kind delete cluster --name $(CLUSTER_NAME)
