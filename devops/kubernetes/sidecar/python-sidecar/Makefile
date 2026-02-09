.PHONY: test run clean docker-build deploy kind-load help

test:
	python3 -m pytest tests/ -v --cov=sidecar --cov-report=html

run:
	WATCH_DIR=/etc/config POD_NAMESPACE=default CONFIGMAP_NAME=properties python3 sidecar.py

clean:
	rm -rf __pycache__ .pytest_cache htmlcov .coverage
	find . -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null || true
	find . -type f -name "*.pyc" -delete

docker-build:
	docker build -t python-sidecar:latest .

deploy: kind-load
	kubectl delete -f k8s-deployment.yaml --ignore-not-found
	kubectl apply -f k8s-deployment.yaml

kind-load: docker-build
	kind load docker-image configmap-printer:latest
	kind load docker-image python-sidecar:latest

deps:
	pip install kubernetes watchdog

lint:
	python3 -m black sidecar.py
	python3 -m flake8 sidecar.py --max-line-length=100

help:
	@echo "Available targets:"
	@echo "  make test          - Run tests with coverage"
	@echo "  make run           - Run sidecar locally (connects to K8s cluster)"
	@echo "  make clean         - Remove build artifacts and cache"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make deploy        - Deploy to Kubernetes cluster"
	@echo "  make kind-load     - Build and load image to kind cluster"
	@echo "  make deps          - Install Python dependencies"
	@echo "  make lint          - Format and lint code"
