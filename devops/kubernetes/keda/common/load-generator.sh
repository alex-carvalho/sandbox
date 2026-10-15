#!/usr/bin/env bash
set -euo pipefail

usage() {
  echo "Usage: $0 {push|drain|status} [queue_name] [count]"
  echo ""
  echo "Commands:"
  echo "  push <queue_name> <count>  Push <count> messages into <queue_name>"
  echo "                             (e.g., $0 push keda-queue 30)"
  echo "                             (e.g., $0 push hpa-prometheus-queue 30)"
  echo "  drain <queue_name>         Empty the specified queue"
  echo "  status                     Display length of both queues and active replica counts"
  exit 1
}

REDIS_POD=$(kubectl get pod -n common -l app=redis -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || true)
if [ -z "${REDIS_POD}" ]; then
  echo "Error: Redis pod not found in namespace 'common'."
  exit 1
fi

CMD="${1:-status}"

case "${CMD}" in
  push)
    QUEUE="${2:-keda-queue}"
    COUNT="${3:-25}"
    echo "Pushing ${COUNT} messages to '${QUEUE}'..."
    # Generate bulk payload for redis-cli pipe to be fast
    # Each message: "item-<timestamp>-<i>"
    kubectl exec -i -n common "${REDIS_POD}" -c redis -- /bin/sh -c \
      "for i in \$(seq 1 ${COUNT}); do redis-cli rpush ${QUEUE} \"job-\$(date +%s%N)-\$i\" >/dev/null; done"
    
    LEN=$(kubectl exec -n common "${REDIS_POD}" -c redis -- redis-cli llen "${QUEUE}" | tr -d '\r')
    echo "Done! Current length of '${QUEUE}': ${LEN}"
    ;;

  drain)
    QUEUE="${2:-keda-queue}"
    echo "Draining '${QUEUE}'..."
    kubectl exec -n common "${REDIS_POD}" -c redis -- redis-cli del "${QUEUE}" >/dev/null
    echo "'${QUEUE}' is now empty."
    ;;

  status)
    KEDA_LEN=$(kubectl exec -n common "${REDIS_POD}" -c redis -- redis-cli llen keda-queue 2>/dev/null | tr -d '\r' || echo "0")
    HPA_LEN=$(kubectl exec -n common "${REDIS_POD}" -c redis -- redis-cli llen hpa-prometheus-queue 2>/dev/null | tr -d '\r' || echo "0")
    
    KEDA_REPLICAS=$(kubectl get deployment keda-consumer -n keda-poc -o jsonpath='{.status.readyReplicas}' 2>/dev/null || echo "0")
    [ -z "${KEDA_REPLICAS}" ] && KEDA_REPLICAS="0"

    HPA_REPLICAS=$(kubectl get deployment hpa-consumer -n hpa-poc -o jsonpath='{.status.readyReplicas}' 2>/dev/null || echo "0")
    [ -z "${HPA_REPLICAS}" ] && HPA_REPLICAS="0"

    echo "=========================================================="
    echo "  Queue & Replica Status"
    echo "=========================================================="
    printf "%-25s : %s items | %s ready pods (Scale-to-0)\n" "keda-queue" "${KEDA_LEN}" "${KEDA_REPLICAS}"
    printf "%-25s : %s items | %s ready pods (Min-1 pod)\n" "hpa-prometheus-queue" "${HPA_LEN}" "${HPA_REPLICAS}"
    echo "=========================================================="
    ;;

  *)
    usage
    ;;
esac
