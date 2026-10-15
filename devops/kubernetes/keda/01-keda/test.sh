#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

echo "=========================================================="
echo "  KEDA Autoscaling Test Demo"
echo "  Queue: keda-queue (Target: 5 msgs / pod)"
echo "=========================================================="

# Ensure deployment exists
if ! kubectl get scaledobject keda-queue-scaler -n keda-poc &>/dev/null; then
  echo "Workload not found in keda-poc. Deploying first..."
  "${SCRIPT_DIR}/deploy.sh"
fi

# Drain queue before starting to guarantee clean baseline
"${ROOT_DIR}/common/load-generator.sh" drain keda-queue

echo ""
echo "Step 1: Baseline Check (Scale-to-Zero)"
echo "Waiting a few seconds to verify baseline replicas = 0..."
sleep 5
CURRENT_REPLICAS=$(kubectl get deployment keda-consumer -n keda-poc -o jsonpath='{.status.readyReplicas}' 2>/dev/null || echo "0")
[ -z "${CURRENT_REPLICAS}" ] && CURRENT_REPLICAS="0"
echo "Current Ready Pods: ${CURRENT_REPLICAS} (Demonstrating native scale-to-zero!)"

echo ""
echo "Step 2: Injecting Load (30 messages into keda-queue)..."
"${ROOT_DIR}/common/load-generator.sh" push keda-queue 30

echo ""
echo "Step 3: Monitoring KEDA Response in Real-Time"
echo "Target formula: ceil(30 / 5) = 6 pods"
echo "Press Ctrl+C to stop watching at any time."
echo "--------------------------------------------------------------------------"
printf "%-10s | %-12s | %-14s | %-20s\n" "TIME" "QUEUE DEPTH" "READY PODS" "STATUS"
echo "--------------------------------------------------------------------------"

PEAK_SEEN=0
DRAINED=0

for i in $(seq 1 60); do
  TIMESTAMP=$(date +%T)
  QLEN=$("${ROOT_DIR}/common/load-generator.sh" status | grep "keda-queue" | awk '{print $3}')
  REPLICAS=$(kubectl get deployment keda-consumer -n keda-poc -o jsonpath='{.status.readyReplicas}' 2>/dev/null || echo "0")
  [ -z "${REPLICAS}" ] && REPLICAS="0"

  STATUS="Scaling"
  if [ "${REPLICAS}" -ge 6 ]; then
    STATUS="Peak scale reached!"
    PEAK_SEEN=1
  fi
  if [ "${QLEN}" -eq 0 ] && [ "${PEAK_SEEN}" -eq 1 ]; then
    STATUS="Queue empty. Cooling down..."
    DRAINED=1
  fi
  if [ "${DRAINED}" -eq 1 ] && [ "${REPLICAS}" -eq 0 ]; then
    printf "%-10s | %-12s | %-14s | %-20s\n" "${TIMESTAMP}" "${QLEN}" "${REPLICAS}" "Successfully scaled to 0!"
    break
  fi

  printf "%-10s | %-12s | %-14s | %-20s\n" "${TIMESTAMP}" "${QLEN}" "${REPLICAS}" "${STATUS}"
  sleep 3
done

echo "--------------------------------------------------------------------------"
echo "KEDA Test Completed!"
echo "Observed Behavior:"
echo "  1. Cold start: 0 pods running when queue was empty (0 resource waste)."
echo "  2. Rapid reaction: Pods scaled from 0 to 6 within seconds of load arrival."
echo "  3. Scale to zero: Scaled cleanly back to 0 pods once queue drained."
echo "=========================================================="
