#!/bin/bash
# filepath: /workspaces/sandbox/devops/challenges/docker-compose-prometheus-elk/create_index_pattern.sh

KIBANA_URL="http://localhost:5601"
INDEX_PATTERN="java-logs-*"
PATTERN_NAME="java-logs-*"

# Wait for Kibana to be up
until curl -s "$KIBANA_URL/api/status" | grep -q '"All services are available"'; do
  echo "Waiting for Kibana..."
  sleep 5
done

# Create the index pattern
curl -X POST "$KIBANA_URL/api/saved_objects/index-pattern" \
  -H 'kbn-xsrf: true' \
  -H 'Content-Type: application/json' \
  -d "{\"attributes\":{\"title\":\"$INDEX_PATTERN\",\"timeFieldName\":\"@timestamp\"}}"