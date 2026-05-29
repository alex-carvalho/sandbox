#!/usr/bin/env bash

set -e

LOKI_URL="http://localhost:3101"
JOB=${1:-"bash-test"}
LEVEL=${2:-"INFO"}
MESSAGE=${3:-"Hello, Grafana Loki! This is a test log message."}

TIMESTAMP=$(date +%s)000000000

PAYLOAD=$(cat <<EOF
{
  "streams": [
    {
      "stream": {
        "job": "${JOB}",
        "level": "${LEVEL}",
        "host": "$(hostname)",
        "source": "shell-script"
      },
      "values": [
        [ "${TIMESTAMP}", "${MESSAGE}" ]
      ]
    }
  ]
}
EOF
)

echo "Pushing log to Loki at ${LOKI_URL}/loki/api/v1/push..."
echo ""

RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" \
  -H "Content-Type: application/json" \
  -X POST \
  -d "${PAYLOAD}" \
  "${LOKI_URL}/loki/api/v1/push")

HTTP_STATUS=$(echo "${RESPONSE}" | tr -d '\r' | tail -n 1 | cut -d':' -f2)
BODY=$(echo "${RESPONSE}" | sed '$d')

if [ "${HTTP_STATUS}" -eq 204 ]; then
  echo "✅ Log successfully pushed (HTTP 204 No Content)"
else
  echo "❌ Failed to push log (HTTP ${HTTP_STATUS})"
  echo "Response body: ${BODY}"
  exit 1
fi
