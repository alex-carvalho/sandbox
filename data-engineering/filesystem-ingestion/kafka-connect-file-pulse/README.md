

```shell
docker-compose up

# check connector is loaded
 curl -sX GET http://localhost:8083/connector-plugins | grep FilePulseSourceConnector

# load connector configuration
curl -X POST -H "Content-Type: application/json" \
     --data @file-pulse-csv.json \
     http://localhost:8083/connectors



open http://localhost:8080/

```