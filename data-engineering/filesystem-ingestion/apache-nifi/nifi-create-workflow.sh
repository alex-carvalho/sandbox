#!/usr/bin/env bash
set -euo pipefail

echo "🔐 Getting NiFi auth token..."
TOKEN=$(curl -k -s -X POST "https://localhost:8443/nifi-api/access/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=admin&password=admin1234567")

if [[ -z "$TOKEN" ]]; then
  echo "❌ Failed to get auth token"
  exit 1
fi
echo "✅ Token obtained."

# -------------------------------------------------------------------------
# Get Root Process Group ID
# -------------------------------------------------------------------------
ROOT_PG_ID=$(curl -k -s "https://localhost:8443/nifi-api/flow/process-groups/root" \
  -H "Authorization: Bearer ${TOKEN}" | jq -r '.processGroupFlow.id')

echo "📦 Root Process Group ID: ${ROOT_PG_ID}"

# -------------------------------------------------------------------------
# Create GetFile Processor
# -------------------------------------------------------------------------
GET_FILE_RESPONSE=$(curl -k -s -X POST "https://localhost:8443/nifi-api/process-groups/${ROOT_PG_ID}/processors" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "revision": {"version": 0},
    "component": {
      "type": "org.apache.nifi.processors.standard.GetFile",
      "name": "GetFile - CSV",
      "position": {"x": 100, "y": 100},
      "config": {
        "properties": {
          "Input Directory": "/data/input",
          "Keep Source File": "false",
          "File Filter": ".*\\.csv"
        }
      }
    }
  }')

GET_FILE_ID=$(echo "$GET_FILE_RESPONSE" | jq -r '.id')
echo "📂 GetFile created with ID: ${GET_FILE_ID}"

# -------------------------------------------------------------------------
# Create Kafka Connection Service
# -------------------------------------------------------------------------
KAFKA_RESPONSE=$(curl -k -s -X POST "https://localhost:8443/nifi-api/process-groups/${ROOT_PG_ID}/controller-services" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "revision": {"version": 0},
    "component": {
      "type": "org.apache.nifi.kafka.service.Kafka3ConnectionService",
      "name": "KafkaConnectionService 2",
      
        "properties": {
          "bootstrap.servers": "kafka:9092",
          "security.protocol": "PLAINTEXT"
        }
      
    }
  }')

KAFKA_SERVICE_ID=$(echo "$KAFKA_RESPONSE" | jq -r '.id')
echo "🔌 KafkaConnectionService created: ${KAFKA_SERVICE_ID}"


# curl 'https://localhost:8443/nifi-api/controller-services/356c0d62-019a-1000-dad9-ed24c30a55eb' \
#   -X 'PUT' \\
#   --data-raw '{"revision":{"clientId":"f96791ef-aeed-49a3-9463-77270045a379","version":3},"disconnectedNodeAcknowledged":false,"component":{"id":"356c0d62-019a-1000-dad9-ed24c30a55eb","name":"KafkaConnectionService","bulletinLevel":"WARN","comments":"","properties":{"bootstrap.servers":"kafka:9092"},"sensitiveDynamicPropertyNames":[]}}' \
#   --insecure

# -------------------------------------------------------------------------
# Enable Kafka Connection Service
# -------------------------------------------------------------------------
echo "⚙️  Enabling Kafka service..."
curl -k -s -X PUT "https://localhost:8443/nifi-api/controller-services/${KAFKA_SERVICE_ID}/run-status" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d "{\"revision\":{\"version\":1},\"state\":\"ENABLED\"}"
echo "✅ Kafka service enabled."

sleep 2 # allow service to register

# -------------------------------------------------------------------------
# Create PublishKafka Processor
# -------------------------------------------------------------------------


KAFKA_PROCESSOR=$(curl -k -s -X POST "https://localhost:8443/nifi-api/process-groups/${ROOT_PG_ID}/processors" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d "{
    \"revision\": {\"version\": 0},
    \"component\": {
      \"type\": \"org.apache.nifi.kafka.processors.PublishKafka\",
      \"name\": \"Publish to Kafka\",
      \"position\": {\"x\": 400, \"y\": 200},
      \"bundle\": {
          \"group\": \"org.apache.nifi\",
          \"artifact\": \"nifi-kafka-nar\",
          \"version\": \"2.6.0\"\
        }
    }
  }")

KAFKA_ID=$(echo "$KAFKA_PROCESSOR" | jq -r '.id')
echo "🪶 PublishKafka created with ID: ${KAFKA_ID}"

# -------------------------------------------------------------------------
# Configure PublishKafka Processor
# -------------------------------------------------------------------------

curl -k -s -X PUT "https://localhost:8443/nifi-api/processors/${KAFKA_ID}" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d "{
    \"revision\": {\"version\": 1},
    \"component\": {
      \"id\": \"${KAFKA_ID}\",
      \"name\": \"Publish to Kafka\",
      \"config\": {
        \"properties\": {
          \"Kafka Connection Service\": \"${KAFKA_SERVICE_ID}\",
          \"Topic Name\": \"test-topic\"
        },
        \"schedulingPeriod\": \"1 sec\",
        \"schedulingStrategy\": \"TIMER_DRIVEN\"
      }
    }
  }" 

# -------------------------------------------------------------------------
# Create Connection (GetFile → PublishKafka)
# -------------------------------------------------------------------------
echo "🔗 Creating connection between GetFile and PublishKafka..."

curl -k -s -X POST "https://localhost:8443/nifi-api/process-groups/${ROOT_PG_ID}/connections" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d "{
    \"revision\": {\"version\": 0, \"clientId\": \"auto-flow-script\"},
    \"component\": {
      \"name\": \"GetFile to Kafka\",
      \"backPressureDataSizeThreshold\": \"1 GB\",
      \"backPressureObjectThreshold\": 10000,
      \"flowFileExpiration\": \"0 sec\",
      \"loadBalanceStrategy\": \"DO_NOT_LOAD_BALANCE\",
      \"loadBalanceCompression\": \"DO_NOT_COMPRESS\",
      \"selectedRelationships\": [\"success\"],
      \"source\": {
        \"groupId\": \"${ROOT_PG_ID}\",
        \"id\": \"${GET_FILE_ID}\",
        \"type\": \"PROCESSOR\"
      },
      \"destination\": {
        \"groupId\": \"${ROOT_PG_ID}\",
        \"id\": \"${KAFKA_ID}\",
        \"type\": \"PROCESSOR\"
      },
      \"bends\": [],
      \"zIndex\": 1
    }
  }" > /dev/null

echo "✅ Connection created successfully."

# -------------------------------------------------------------------------
# Start Processors
# -------------------------------------------------------------------------
echo "🚀 Starting processors..."

curl -k -s -X PUT "https://localhost:8443/nifi-api/processors/${GET_FILE_ID}/run-status" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"revision": {"version": 1}, "state": "RUNNING"}' > /dev/null

curl -k -s -X PUT "https://localhost:8443/nifi-api/processors/${KAFKA_ID}/run-status" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"revision": {"version": 1}, "state": "RUNNING"}' > /dev/null

echo "✅ Flow deployed and started successfully."
