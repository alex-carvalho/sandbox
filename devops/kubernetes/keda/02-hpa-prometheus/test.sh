#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

echo "=========================================================="
echo "  Standard HPA + Prometheus Autoscaling Test Demo"
echo "  Queue: hpa-prometheus-queue (Target: 5 msgs / pod)"
echo "=========================================================="

# Ensure deployment exists
if ! kubectl get hpa hpa-prometheus-consumer -n hpa-poc &>/dev/null; then
  echo "Workload not found in hpa-poc. Deploying first..."
  "${SCRIPT_DIR}/deploy.sh"
fi

# Drain queue before starting
"${ROOT_DIR}/common/load-generator.sh" drain hpa-prometheus-queue

echo ""
echo "Step 1: Baseline Check (Standard HPA Baseline Pod)"
CURRENT_REPLICAS=$(kubectl get deployment hpa-consumer -n hpa-poc -o jsonpath='{.status.readyReplicas}' 2>/dev/null || echo "1")
[ -z "${CURRENT_REPLICAS}" ] && CURRENT_REPLICAS="1"
echo "Current Ready Pods: ${CURRENT_REPLICAS} (Standard HPA keeps minReplicas >= 1; 1 pod always consumes idle compute)."

echo ""
echo "Step 2: Injecting Load (30 messages into hpa-prometheus-queue)..."
"${ROOT_DIR}/common/load-generator.sh" push hpa-prometheus-queue 30

echo ""
echo "Step 3: Monitoring HPA Response in Real-Time"
echo "Notice the scrape pipeline latency: Redis -> redis_exporter -> Prometheus -> Adapter -> HPA"
echo "Target formula: ceil(30 / 5) = 6 pods"
echo "Press Ctrl+C to stop watching at any time."
echo "--------------------------------------------------------------------------------------------------"
printf "%-10s | %-12s | %-20s | %-12s | %-25s\n" "TIME" "QUEUE DEPTH" "HPA METRIC VALUE" "READY PODS" "STATUS"
echo "--------------------------------------------------------------------------------------------------"

PEAK_SEEN=0
DRAINED=0

for i in $(seq 1 60); do
  TIMESTAMP=$(date +%T)
  QLEN=$("${ROOT_DIR}/common/load-generator.sh" status | grep "hpa-prometheus-queue" | awk '{print $3}')
  HPA_METRIC=$(kubectl get hpa hpa-prometheus-consumer -n hpa-poc -o jsonpath='{.status.currentMetrics[0].object.current.averageValue}' 2>/dev/null || echo "<unknown>")
  [ -z "${HPA_METRIC}" ] && HPA_METRIC="<unknown>"
  REPLICAS=$(kubectl get deployment hpa-consumer -n hpa-poc -o jsonpath='{.status.readyReplicas}' 2>/dev/null || echo "0")
  [ -z "${REPLICAS}" ] && REPLICAS="0"

  STATUS="Waiting for scrape pipeline"
  if [ "${HPA_METRIC}" != "<unknown>" ] && [ "${HPA_METRIC}" != "0" ]; then
    STATUS="Metric arrived at HPA"
  fi
  if [ "${REPLICAS}" -ge 6 ]; then
    STATUS="Peak scale reached (6 pods)!"
    PEAK_SEEN=1
  fi
  if [ "${QLEN}" -eq 0 ] && [ "${PEAK_SEEN}" -eq 1 ]; then
    STATUS="Queue empty. Cooling down..."
    DRAINED=1
  fi
  if [ "${DRAINED}" -eq 1 ] && [ "${REPLICAS}" -eq 1 ]; then
    printf "%-10s | %-12s | %-20s | %-12s | %-25s\n" "${TIMESTAMP}" "${QLEN}" "${HPA_METRIC}" "${REPLICAS}" "Scaled down (min 1 pod)"
    break
  fi

  printf "%-10s | %-12s | %-20s | %-12s | %-25s\n" "${TIMESTAMP}" "${QLEN}" "${HPA_METRIC}" "${REPLICAS}" "${STATUS}"
  sleep 3
done

echo "--------------------------------------------------------------------------------------------------"
echo "HPA Test Completed!"
echo "Observed Behavior:"
echo "  1. Idle resource usage: 1 pod was always running even when queue had 0 messages."
echo "  2. Scrape Pipeline Latency: Took longer to scale up due to Prometheus scrape & Adapter sync."
echo "  3. Cannot scale to zero: Scaled down and stopped at 1 pod (never drops to 0)."
echo "=========================================================="
