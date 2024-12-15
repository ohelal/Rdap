# RDAP Service

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/ohelal/rdap)](https://goreportcard.com/report/github.com/ohelal/rdap)
[![Go Reference](https://pkg.go.dev/badge/github.com/ohelal/rdap.svg)](https://pkg.go.dev/github.com/ohelal/rdap)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ohelal/rdap)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/ohelal/rdap)](https://github.com/ohelal/rdap/releases)
[![Build Status](https://github.com/ohelal/rdap/workflows/Build%20and%20Test/badge.svg)](https://github.com/ohelal/rdap/actions)
[![codecov](https://codecov.io/gh/ohelal/rdap/branch/main/graph/badge.svg)](https://codecov.io/gh/ohelal/rdap)

A high-performance Registration Data Access Protocol (RDAP) service implementation in Go, with support for IP addresses, ASNs, and domain lookups. The service includes caching with Redis and event streaming with Apache Kafka.

## Features

- Full RDAP protocol implementation
- High-performance Go implementation
- Redis caching for improved response times
- Kafka integration for event streaming
- Automatic RDAP bootstrap file updates
- Docker and Kubernetes ready
- CLI tool for easy querying

## Prerequisites

- Go 1.22 or later
- Docker
- Kubernetes (optional, for deployment)
- Redis
- Apache Kafka

## Quick Start

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

3. Run Redis and Kafka (using Docker):
```bash
docker-compose up -d redis kafka
```

4. Build and run the service:
```bash
go run ./cmd/rdap/main.go
```

### Running Tests

Run all tests (including integration tests):
```bash
go test ./... -v
```

Run only unit tests (skips integration tests):
```bash
go test ./... -v -short
```

### Using Docker

```bash
# Build the image
docker build -t rdap:latest .

# Run the container
docker run -p 8080:8080 rdap:latest
```

### Kubernetes Deployment

1. Apply the Kubernetes manifests:
```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap.yaml -n rdap
kubectl apply -f k8s/kafka -n rdap
kubectl apply -f k8s/redis.yaml -n rdap
kubectl apply -f k8s/zookeeper -n rdap
kubectl apply -f k8s/rdap.yaml -n rdap
```

2. Get the service URL:
```bash
minikube service rdap-service -n rdap --url
```

## Using the CLI

The RDAP service includes a CLI tool for easy querying. For local Minikube deployment, use:

```bash
# Get Minikube IP
MINIKUBE_IP=$(minikube ip)

# IP lookup
go run ./cmd/rdap/main.go --base-url http://$MINIKUBE_IP:31080 ip 8.8.8.8

# Domain lookup
go run ./cmd/rdap/main.go --base-url http://$MINIKUBE_IP:31080 domain google.com

# ASN lookup
go run ./cmd/rdap/main.go --base-url http://$MINIKUBE_IP:31080 asn 15169
```

The service exposes the following ports:
- HTTP API: 31080
- Metrics: 31090

## Configuration

The service can be configured using environment variables or a configuration file. See the [configuration documentation](docs/configuration.md) for details.

## API Documentation

- [API Reference](docs/api.md)
- [Kubernetes Deployment](docs/kubernetes.md)
- [Development Guide](docs/development.md)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). See the following files for details:

- [LICENSE](LICENSE) - Full AGPL-3.0 license text
- [COPYING](COPYING) - Additional terms and conditions
- [AUTHORS](AUTHORS) - List of contributors
- [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) - Community guidelines
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contribution guidelines

### Commercial Use
Commercial use requires explicit written permission from the author. Please contact mohamed@helal.me for licensing inquiries.

### License Terms
This software is free to use for:
- Personal projects
- Open source projects
- Non-commercial use
- Educational purposes
- Research and development

All modifications and distributions must comply with AGPL-3.0 requirements, including:
- Making source code available
- Maintaining copyright notices
- Sharing modifications under AGPL-3.0

## Acknowledgments

- [IANA RDAP Bootstrap Registry](https://data.iana.org/rdap/)
- [RDAP Protocol Specification](https://tools.ietf.org/html/rfc7482)
