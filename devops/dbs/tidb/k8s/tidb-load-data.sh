#!/usr/bin/env bash
set -euo pipefail

DB_NAME="demo"
TABLE_NAME="orders"
ROW_COUNT="10000"
PORT_FORWARD_PID=""


mysql_exec() {
  local sql="$1"
  local -a mysql_cmd=(mysql -h 127.0.0.1 -P 4000 -u root -N -s)

  "${mysql_cmd[@]}" -e "${sql}"
}

wait_tiflash_replica() {
  local available
  local i
  for i in $(seq 1 60); do
    available="$(mysql_exec "SELECT IFNULL(AVAILABLE, 0) FROM information_schema.tiflash_replica WHERE TABLE_SCHEMA='${DB_NAME}' AND TABLE_NAME='${TABLE_NAME}' LIMIT 1;" || true)"
    if [ "${available:-0}" = "1" ]; then
      return 0
    fi
    sleep 2
  done
  return 1
}


echo "Creating database/table (idempotent)..."
mysql_exec "DROP DATABASE IF EXISTS ${DB_NAME};"
mysql_exec "CREATE DATABASE IF NOT EXISTS ${DB_NAME};"

mysql_exec "
  CREATE TABLE IF NOT EXISTS ${DB_NAME}.${TABLE_NAME} (
    id BIGINT PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    region VARCHAR(16) NOT NULL,
    status VARCHAR(16) NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    created_at DATETIME NOT NULL,
    INDEX idx_region (region)
  );
"

echo ""
echo "Enabling TiFlash replica for aggregation query..."
mysql_exec "ALTER TABLE ${DB_NAME}.${TABLE_NAME} SET TIFLASH REPLICA 1;"


echo "Loading ${ROW_COUNT} rows..."
mysql_exec "
  INSERT INTO ${DB_NAME}.${TABLE_NAME} (id, customer_id, region, status, amount, created_at)
  SELECT n,
         (n % 1000) + 1,
         ELT((n % 5) + 1, 'us', 'eu', 'apac', 'latam', 'mea'),
         ELT((n % 3) + 1, 'new', 'paid', 'cancelled'),
         ROUND(((n * 17) % 10000) / 100 + 10, 2),
         DATE_ADD(
           DATE_ADD(TIMESTAMP('2026-01-01 00:00:00'), INTERVAL (n % 60) DAY),
           INTERVAL (n % 24) HOUR
         )
  FROM (
    SELECT ones.n + tens.n * 10 + hundreds.n * 100 + thousands.n * 1000 + 2 AS n
    FROM (SELECT 0 n UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) ones
    CROSS JOIN (SELECT 0 n UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) tens
    CROSS JOIN (SELECT 0 n UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) hundreds
    CROSS JOIN (SELECT 0 n UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) thousands
  ) t
  WHERE n <= ${ROW_COUNT};
"

# insert a known row for testing
mysql_exec "INSERT INTO ${DB_NAME}.${TABLE_NAME} (id, customer_id, region, status, amount, created_at) VALUES (1, 1, 'us', 'new', 100.00, '2026-01-01 00:00:00');"

echo "Row count check:"
mysql_exec "SELECT COUNT(*) AS total_rows FROM ${DB_NAME}.${TABLE_NAME};"

wait_tiflash_replica && echo "TiFlash replica is available."