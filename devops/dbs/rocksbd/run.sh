#!/bin/bash

IMAGE_NAME="rocksdb-poc"
CONTAINER_NAME="rocksdb-container"

echo "Building RocksDB Docker image..."
docker build -t $IMAGE_NAME .

echo ""
echo "Starting RocksDB container..."
docker run -d --name $CONTAINER_NAME $IMAGE_NAME

echo ""
echo "Container started! Running client..."
echo ""
docker exec -it $CONTAINER_NAME ./client

echo ""
echo "Cleanup: docker rm -f $CONTAINER_NAME"
