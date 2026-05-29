#!/usr/bin/env bash

set -e

LOKI_URL="http://localhost:3101"
QUERY=${1:-'{job="bash-test"}'}
LIMIT=10

echo "Querying Loki at ${LOKI_URL}/loki/api/v1/query_range..."
echo "LogQL Query: ${QUERY}"
echo ""

RESPONSE=$(curl -s -G \
  --data-urlencode "query=${QUERY}" \
  --data-urlencode "limit=${LIMIT}" \
  "${LOKI_URL}/loki/api/v1/query_range")


echo "${RESPONSE}" | jq .
