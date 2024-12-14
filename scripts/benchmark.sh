#!/bin/bash

# Benchmark script for RDAP service
# Requires: wrk (HTTP benchmarking tool)
# Usage: ./benchmark.sh [duration] [connections] [threads]

DURATION=${1:-30}  # Default 30 seconds
CONNECTIONS=${2:-1000}  # Default 1000 connections
THREADS=${3:-$(nproc)}  # Default to number of CPU cores
URL="http://localhost:8080/lookup?q=8.8.8.8"

# Check if wrk is installed
if ! command -v wrk &> /dev/null; then
    echo "wrk is not installed. Please install it first:"
    echo "Ubuntu/Debian: sudo apt-get install wrk"
    echo "CentOS/RHEL: sudo yum install wrk"
    exit 1
fi

# Print test parameters
echo "=== RDAP Service Benchmark ==="
echo "Duration: ${DURATION}s"
echo "Connections: ${CONNECTIONS}"
echo "Threads: ${THREADS}"
echo "URL: ${URL}"
echo "=========================="

# Run warmup
echo "Warming up for 5 seconds..."
wrk -t${THREADS} -c${CONNECTIONS} -d5s ${URL} > /dev/null 2>&1

# Run actual benchmark
echo "Running benchmark..."
wrk -t${THREADS} -c${CONNECTIONS} -d${DURATION}s --latency ${URL}

# Run with different payload sizes
echo -e "\nTesting with different query types..."
QUERIES=("8.8.8.8" "2001:db8::" "AS15169" "example.com")
for query in "${QUERIES[@]}"; do
    echo -e "\nTesting with query: ${query}"
    wrk -t${THREADS} -c${CONNECTIONS} -d10s "http://localhost:8080/lookup?q=${query}"
done

# Run with increasing concurrency
echo -e "\nTesting with increasing concurrency..."
for conn in 100 1000 10000 50000 100000; do
    echo -e "\nTesting with ${conn} connections"
    wrk -t${THREADS} -c${conn} -d10s ${URL}
done

echo -e "\nBenchmark complete!"
