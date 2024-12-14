# RDAP Service & CLI Tool

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/ohelal/rdap)](https://goreportcard.com/report/github.com/ohelal/rdap)
[![Go Reference](https://pkg.go.dev/badge/github.com/ohelal/rdap.svg)](https://pkg.go.dev/github.com/ohelal/rdap)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ohelal/rdap)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/ohelal/rdap)](https://github.com/ohelal/rdap/releases)
[![Build Status](https://github.com/ohelal/rdap/workflows/build/badge.svg)](https://github.com/ohelal/rdap/actions)
[![codecov](https://codecov.io/gh/ohelal/rdap/branch/main/graph/badge.svg)](https://codecov.io/gh/ohelal/rdap)

A high-performance Registration Data Access Protocol (RDAP) service and command-line tool implemented in Go. This project provides fast and reliable lookups for IP addresses, ASNs, and domain names with built-in distributed caching, message queuing, and scalable architecture support.

🚀 **Key Features**:
- High-performance RDAP lookups
- Built-in distributed caching
- Message queue integration
- CLI and Service modes
- AGPL v3 Licensed - Free for non-commercial use

📧 **Contact**: mohamed@helal.me

## Quick Links
- [Installation](#installation)
- [Service Usage](#service-usage)
- [CLI Usage](#cli-usage)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [Security](#security)
- [License](#license)

## Table of Contents
- [Service Features](#service-features)
- [CLI Features](#cli-features)
- [Installation](#installation)
- [Service Usage](#service-usage)
- [CLI Usage](#cli-usage)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [Security](#security)
- [License](#license)
- [Legal Notices](#legal-notices)
- [Acknowledgments](#acknowledgments)

## Service Features

### Complete RDAP Implementation
- IP address lookups (IPv4 and IPv6)
- Autonomous System Number (ASN) lookups
- Domain name lookups
- Standards-compliant RDAP responses
- Bootstrap configuration support

### High Performance Architecture
- Distributed caching with Redis
- Request coalescing to prevent cache stampede
- Efficient connection pooling
- Concurrent request handling
- Response compression
- Horizontal scalability
- Rate limiting with sliding window algorithm

### Distributed System Components
- Redis for high-speed caching
  - Configurable TTL
  - Distributed cache invalidation
  - Support for cache clusters
  - Connection pooling with automatic retry
  - Cache metrics and monitoring

### Reliability & Monitoring
- Prometheus metrics integration
- Structured logging with configurable levels
- Health checks for all components
- Graceful shutdown
- Error handling and recovery
- Circuit breakers for external services
- Rate limiting per endpoint
- Request tracing

## CLI Features
- Interactive command-line interface
- Query information for:
  - Domains
  - IP addresses (IPv4 and IPv6)
  - ASNs (supports both formats: "AS15169" or "15169")
- Multiple output formats:
  - Pretty (default)
  - JSON
  - Table
  - Box
- Result caching
- Progress indicators
- Configurable timeout and base URL

## Requirements
- Go 1.22 or higher
- Redis 7.0 or higher
- Docker and Docker Compose (optional)

## Installation

### Using as a Library
```bash
go get github.com/ohelal/rdap@latest
```

### Using the CLI Tool
```bash
go install github.com/ohelal/rdap/cmd/rdap@latest
```

### Library Usage Example
```go
package main

import (
    "context"
    "fmt"
    "github.com/ohelal/rdap"
)

func main() {
    client := rdap.NewClient()
    
    // IP Lookup
    ipResult, err := client.LookupIP(context.Background(), "8.8.8.8")
    if err != nil {
        panic(err)
    }
    fmt.Printf("IP Owner: %s\n", ipResult.Name)
    
    // ASN Lookup
    asnResult, err := client.LookupASN(context.Background(), "15169")
    if err != nil {
        panic(err)
    }
    fmt.Printf("ASN Name: %s\n", asnResult.Name)
    
    // Domain Lookup
    domainResult, err := client.LookupDomain(context.Background(), "example.com")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Domain Status: %v\n", domainResult.Status)
}
```

### Prerequisites
- Docker Engine 20.10+
- Docker Compose v2+
- Minimum 4GB RAM
- 2 CPU cores recommended

### Using Docker
```bash
docker pull ghcr.io/ohelal/rdap
docker run -p 8080:8080 ghcr.io/ohelal/rdap
```

### Using Docker Compose (Recommended for Production)
```yaml
version: '3.8'
services:
  rdap:
    image: ghcr.io/ohelal/rdap
    ports:
      - "8080:8080"
    environment:
      - REDIS_URL=redis:6379
      - CACHE_TTL=3600
    depends_on:
      - redis
  redis:
    image: redis:alpine
    ports:
      - "6379:6379"
```

### From Source
```bash
git clone https://github.com/ohelal/rdap.git
cd rdap
go mod download
go build ./...
```

## Service Usage

### As a Library
```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/ohelal/rdap/pkg/rdap"
)

func main() {
    // Create a new RDAP client
    client := rdap.NewClient(
        rdap.WithBaseURL("https://your-rdap-server.com"),
        rdap.WithTimeout(5 * time.Second),
    )

    // Query domain information
    domain, err := client.QueryDomain(context.Background(), "example.com")
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Domain info: %+v", domain)
}
```

### HTTP API Endpoints
```bash
# IP lookup
curl http://localhost:8080/ip/8.8.8.8

# ASN lookup
curl http://localhost:8080/autnum/15169

# Domain lookup
curl http://localhost:8080/domain/google.com
```

### Rate Limits
- `/ip`: 100 requests/minute
- `/domain`: 200 requests/minute
- `/autnum`: 150 requests/minute
- `/nameserver`: 100 requests/minute
- Other endpoints: 50 requests/minute

## Package Documentation

### Available Packages

```go
github.com/ohelal/rdap              // Main package with client interface
github.com/ohelal/rdap/cmd/rdap     // CLI tool
github.com/ohelal/rdap/cmd/server   // RDAP server
github.com/ohelal/rdap/pkg/rdap     // Public API package
```

### Using as a Library

```go
import "github.com/ohelal/rdap/pkg/rdap"

// Create a new client with default configuration
client := rdap.NewClient()

// Create a client with custom configuration
client := rdap.NewClient(rdap.Config{
    BaseURL: "https://rdap.example.com",
    Timeout: 30 * time.Second,
    Cache: rdap.CacheConfig{
        Enabled: true,
        TTL:     time.Hour,
    },
})
```

### Running the Server

1. Using Docker:
```bash
docker run -p 8080:8080 \
    -e REDIS_URL=redis:6379 \
    -e KAFKA_BROKERS=kafka:9092 \
    ohelal/rdap server
```

2. Using Kubernetes:
```bash
kubectl apply -f deployments/k8s/
```

3. From source:
```bash
go run cmd/server/main.go
```

### Environment Variables

| Category | Variable | Description | Default |
|----------|----------|-------------|---------|
| **Server** |
| | `PORT` | Server port | `8080` |
| | `MAX_CONCURRENT_REQUESTS` | Max concurrent requests | `5000` |
| **Redis** |
| | `REDIS_URL` | Redis connection URL | `redis:6379` |
| | `REDIS_PASSWORD` | Redis password | `` |
| | `CACHE_TTL` | Cache TTL in seconds | `3600` |
| **Kafka** |
| | `KAFKA_BROKERS` | Kafka brokers (comma-separated) | `` |
| | `KAFKA_TOPIC` | Kafka topic | `rdap-events` |
| **Metrics** |
| | `METRICS_PORT` | Prometheus metrics port | `9090` |
| **Logging** |
| | `LOG_LEVEL` | Log level (debug,info,warn,error) | `info` |
| | `LOG_FORMAT` | Log format (json,text) | `json` |

### API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/ip/:ip` | GET | IP address lookup |
| `/asn/:asn` | GET | ASN lookup |
| `/domain/:domain` | GET | Domain lookup |
| `/health` | GET | Health check |
| `/metrics` | GET | Prometheus metrics |

### Rate Limits

| Endpoint | Limit |
|----------|-------|
| `/ip/*` | 100/min |
| `/asn/*` | 100/min |
| `/domain/*` | 100/min |

## Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `REDIS_URL` | Redis connection string | `redis:6379` | Yes |
| `CACHE_TTL` | Cache duration in seconds | `3600` | No |
| `MAX_CONCURRENT_REQUESTS` | Maximum concurrent requests | `5000` | No |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` | No |
| `KAFKA_BROKERS` | Comma-separated list of Kafka brokers | | No |
| `KAFKA_TOPIC` | Kafka topic for event streaming | | No |
| `METRICS_PORT` | Port for Prometheus metrics | `9090` | No |
| `HEALTH_CHECK_PORT` | Port for health checks | `8081` | No |
| `TLS_CERT_FILE` | Path to TLS certificate file | | No |
| `TLS_KEY_FILE` | Path to TLS key file | | No |

### Configuration File
You can also use a configuration file. Create `config.yaml`:

```yaml
redis:
  url: redis:6379
  ttl: 3600
kafka:
  brokers:
    - kafka-1:9092
    - kafka-2:9092
  topic: rdap-events
server:
  port: 8080
  metrics_port: 9090
  health_check_port: 8081
  max_concurrent_requests: 5000
logging:
  level: info
  format: json
```

## CLI Usage

### Domain Lookup
```bash
# Basic query
rdap domain example.com

# With different output formats
rdap domain example.com --format json
rdap domain example.com --style table
rdap domain example.com --style box
```

### IP Lookup
```bash
# Basic query
rdap ip 8.8.8.8

# With different output formats
rdap ip 8.8.8.8 --format json
rdap ip 8.8.8.8 --style table
rdap ip 8.8.8.8 --style box
```

### ASN Lookup
```bash
# Both formats supported
rdap asn 15169
rdap asn AS15169

# With different output formats
rdap asn 15169 --format json
rdap asn 15169 --style table
rdap asn 15169 --style box
```

### Cache Management
```bash
# View statistics
rdap cache stats

# Clear cache
rdap cache clear
```

### Global Flags
```bash
  --base-url string    Base URL for RDAP server
  -f, --format string  Output format: pretty, json, or compact (default "pretty")
  -s, --style string   Output style: default, table, box (default "default")
  -t, --timeout duration   Query timeout (default 10s)
  -v, --verbose       Enable verbose output
```

## Kubernetes Deployment

For Kubernetes deployment, you'll need:
1. Kubernetes cluster (minikube, kind, or cloud provider)
2. kubectl installed
3. Helm (optional, for using charts)

To set up minikube:
```bash
# Install minikube from official release
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube

# Start minikube
minikube start

# Enable required addons
minikube addons enable ingress
minikube addons enable metrics-server
```

For deployment instructions, see [Kubernetes Deployment Guide](docs/kubernetes.md)

## Contributing
1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

Please ensure your PR:
- Includes tests
- Updates documentation
- Follows Go best practices
- Passes all CI checks

## Security
- [Security Policy](SECURITY.md)

## License
- [License Terms](#license-terms)

## Legal Notices

### Copyright Notice
Copyright 2024 Helal <mohamed@helal.me>
All rights reserved.

### License Terms
This project is licensed under the GNU Affero General Public License v3 (AGPL-3.0) with additional terms.

✅ **Free to Use For**:
- Personal projects
- Open source projects
- Non-commercial use
- Educational purposes
- Research and development

⚠️ **Commercial Use Requirements**:
- Requires explicit written permission from the author
- Contact mohamed@helal.me for licensing
- Must comply with AGPL-3.0 requirements
- Must make source code available to network users
- Must share all modifications under AGPL-3.0
- Must maintain all copyright notices and attributions

For complete terms, see:
- [LICENSE](LICENSE) - Full AGPL-3.0 license text
- [COPYING](COPYING) - Additional terms and conditions

## Documentation
- [Contributing Guide](CONTRIBUTING.md)
- [Security Policy](SECURITY.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Authors](AUTHORS)

## Acknowledgments
- IANA for RDAP standards and bootstrap files
- Go community for excellent libraries
- All [contributors](AUTHORS) to this project
