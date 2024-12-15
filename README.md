# RDAP Service

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/ohelal/rdap)](https://goreportcard.com/report/github.com/ohelal/rdap)
[![Go Reference](https://pkg.go.dev/badge/github.com/ohelal/rdap.svg)](https://pkg.go.dev/github.com/ohelal/rdap)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ohelal/rdap)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/ohelal/rdap)](https://github.com/ohelal/rdap/releases)
[![Build Status](https://github.com/ohelal/rdap/workflows/Build%20and%20Test/badge.svg)](https://github.com/ohelal/rdap/actions)

A high-performance Registration Data Access Protocol (RDAP) service implementation in Go. This project provides both a server implementation and a CLI tool for querying RDAP data about IP addresses, ASNs, and domains. The service includes advanced features like Redis caching for improved performance and Kafka event streaming for real-time data processing.

## 🚀 Features

- **Full RDAP Protocol Support**
  - Domain name queries
  - IP address queries (IPv4 and IPv6)
  - Autonomous System Number (ASN) queries
  - Bootstrap registry integration
  
- **High Performance**
  - Redis caching for faster responses
  - Connection pooling and request coalescing
  - Efficient memory management
  
- **Advanced Features**
  - Kafka event streaming for real-time updates
  - Metrics and monitoring support
  - Rate limiting and circuit breaking
  - Automatic RDAP bootstrap file updates
  
- **Developer Friendly**
  - Easy-to-use CLI tool
  - Docker and Kubernetes support
  - Comprehensive API documentation
  - Extensive test coverage

## 📦 Installation

### CLI Tool

Install the RDAP CLI tool using Go:

```bash
go install github.com/ohelal/rdap/cmd/rdap@latest
```

Or download pre-built binaries from the [releases page](https://github.com/ohelal/rdap/releases).

### Server

Install the RDAP server using Go:

```bash
go install github.com/ohelal/rdap/cmd/server@latest
```

Or use Docker:

```bash
docker pull ohelal/rdap
```

## 🔧 Prerequisites

- Go 1.22 or later (for building from source)
- Docker (for containerized deployment)
- Redis (for caching)
- Apache Kafka (for event streaming)
- Kubernetes (optional, for orchestration)

## 🏃 Quick Start

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/ohelal/rdap.git
cd rdap
```

2. Install dependencies:
```bash
go mod download
```

3. Start required services:
```bash
docker-compose up -d redis kafka
```

4. Run the server:
```bash
go run ./cmd/server
```

### Using the CLI

Query domain information:
```bash
rdap domain google.com
```

Query IP address information:
```bash
rdap ip 8.8.8.8
```

Query ASN information:
```bash
rdap asn AS15169
```

### Running Tests

Run the full test suite:
```bash
go test ./... -v
```

Run only unit tests (skip integration tests):
```bash
go test ./... -v -short
```

## ⚙️ Configuration

The service can be configured using environment variables or a configuration file. For detailed configuration options, see our [Configuration Guide](docs/configuration.md).

### Key Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `RDAP_SERVER_PORT` | Server port | 8080 |
| `RDAP_REDIS_URL` | Redis connection URL | redis://localhost:6379 |
| `RDAP_KAFKA_BROKERS` | Kafka broker list | localhost:9092 |
| `RDAP_LOG_LEVEL` | Logging level | info |

### Docker Configuration

```bash
docker run -p 8080:8080 \
  -e RDAP_REDIS_URL=redis://redis:6379 \
  -e RDAP_KAFKA_BROKERS=kafka:9092 \
  ohelal/rdap
```

## 📚 Documentation

- [API Documentation](docs/api.md) - Detailed API reference
- [Configuration Guide](docs/configuration.md) - Configuration options
- [Development Guide](docs/development.md) - Development setup and guidelines
- [Kubernetes Guide](docs/kubernetes.md) - Kubernetes deployment instructions
- [Project Structure](docs/project_structure.md) - Codebase organization

## 🤝 Contributing

We welcome contributions! Please read our [Contributing Guide](CONTRIBUTING.md) for details on:
- Code of Conduct
- Development process
- How to submit pull requests
- Coding standards

## 🔒 Security

For security issues, please see our [Security Policy](SECURITY.md).

## 📄 License

This project is licensed under the GNU Affero General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## 📝 Changelog

See [CHANGELOG.md](CHANGELOG.md) for a list of changes in each release.
