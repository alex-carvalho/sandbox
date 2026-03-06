#!/usr/bin/env bash
set -euo pipefail

DB_NAME="demo"
TABLE_NAME="orders"
ROW_COUNT="10000"

mysql_exec() {
  local sql="$1"
  local -a mysql_cmd=(mysql -h 127.0.0.1 -P 4000 -u root -N -s)

  "${mysql_cmd[@]}" -e "${sql}"
}


echo "Row count check:"
mysql_exec "SELECT COUNT(*) AS total_rows FROM ${DB_NAME}.${TABLE_NAME};"

echo ""
echo "Query 1: by id"
mysql_exec "SELECT * FROM ${DB_NAME}.${TABLE_NAME} WHERE id = 1;" 


echo ""
echo "Query 2: by index field (region)"
mysql_exec "SELECT count(*) FROM ${DB_NAME}.${TABLE_NAME} WHERE region = 'us';" 


echo ""
echo "Query 3: by non-index field (status)"
mysql_exec "SELECT count(*) FROM ${DB_NAME}.${TABLE_NAME} WHERE status = 'new';" 


echo ""
echo "Query 4: aggregation routed to TiFlash"
mysql_exec "
  SET SESSION tidb_allow_mpp = 1;
  SET SESSION tidb_enforce_mpp = 1;
  SELECT region, COUNT(*) AS total_rows, ROUND(SUM(amount), 2) AS total_amount
  FROM ${DB_NAME}.${TABLE_NAME}
  GROUP BY region
  ORDER BY total_amount DESC;
"