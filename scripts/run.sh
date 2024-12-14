#!/bin/bash

# Check prerequisites
command -v docker >/dev/null 2>&1 || { echo "Docker is required but not installed. Aborting." >&2; exit 1; }
command -v docker-compose >/dev/null 2>&1 || { echo "Docker Compose is required but not installed. Aborting." >&2; exit 1; }
command -v go >/dev/null 2>&1 || { echo "Go is required but not installed. Aborting." >&2; exit 1; }

# Build the project
echo "Building RDAP service..."
go build -o rdap_service ./cmd/server

# Run system tuning if root
if [ "$EUID" -eq 0 ]; then
    echo "Running system tuning..."
    ./scripts/tune_system.sh
else
    echo "Skipping system tuning (requires root)"
fi

# Start dependencies with Docker Compose
echo "Starting dependencies..."
docker-compose -f docker-compose.dev.yml up -d redis kafka zookeeper

# Wait for dependencies
echo "Waiting for dependencies to be ready..."
sleep 10

# Start the RDAP service
echo "Starting RDAP service..."
./rdap_service 