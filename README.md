# RDAP Service & CLI Tool

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/ohelal/rdap)](https://goreportcard.com/report/github.com/ohelal/rdap)
[![Go Reference](https://pkg.go.dev/badge/github.com/ohelal/rdap.svg)](https://pkg.go.dev/github.com/ohelal/rdap)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ohelal/rdap)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/ohelal/rdap)](https://github.com/ohelal/rdap/releases)
[![Build Status](https://github.com/ohelal/rdap/workflows/build/badge.svg)](https://github.com/ohelal/rdap/actions)
[![Coverage](https://img.shields.io/codecov/c/github/ohelal/rdap)](https://codecov.io/gh/ohelal/rdap)

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

### Configuration
Key configuration options (via environment variables or .env file):
- `REDIS_URL`: Redis connection string (default: redis:6379)
- `CACHE_TTL`: Cache duration in seconds (default: 3600)
- `MAX_CONCURRENT_REQUESTS`: Maximum concurrent requests (default: 5000)
- `LOG_LEVEL`: Logging level (default: info)
- `GOMAXPROCS`: Number of processors to use (default: 2)

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
