# Development Guide

This guide will help you set up your development environment and understand the project structure.

## Project Structure

```
.
├── cmd/
│   ├── rdap/       # CLI tool
│   └── server/     # RDAP server
├── internal/
│   ├── api/        # API handlers
│   ├── cache/      # Caching layer
│   ├── config/     # Configuration
│   ├── models/     # Data models
│   └── service/    # Business logic
├── pkg/
│   └── rdap/       # Public API package
├── k8s/            # Kubernetes manifests
├── docs/           # Documentation
└── config/         # Configuration files
```

## Development Setup

1. Install prerequisites:
   - Go 1.22 or later
   - Docker and Docker Compose
   - Make
   - golangci-lint

2. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/rdap.git
   cd rdap
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Set up development environment:
   ```bash
   make dev-setup
   ```

## Development Workflow

1. Create a new branch:
   ```bash
   git checkout -b feature/your-feature
   ```

2. Make your changes and ensure they follow our coding standards:
   ```bash
   make lint
   make test
   ```

3. Build and run locally:
   ```bash
   make build
   make run
   ```

4. Submit a pull request:
   - Write clear commit messages
   - Include tests for new features
   - Update documentation as needed
   - Ensure CI passes

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with race detection
make test-race

# Run specific test
go test ./... -run TestName

# Generate coverage report
make coverage
```

### Test Structure

- Unit tests: `*_test.go` files alongside the code they test
- Integration tests: `tests/integration/` directory
- Benchmarks: `*_bench_test.go` files

## Code Style

We follow standard Go code style guidelines:

1. Use `gofmt` for formatting
2. Follow [Effective Go](https://golang.org/doc/effective_go.html)
3. Document all exported functions and types
4. Write clear, concise commit messages

## Debugging

1. Use delve for debugging:
   ```bash
   dlv debug ./cmd/server/main.go
   ```

2. Enable debug logging:
   ```bash
   LOG_LEVEL=debug go run ./cmd/server/main.go
   ```

## Common Tasks

### Adding a New API Endpoint

1. Define the request/response models in `internal/models`
2. Add the handler in `internal/api`
3. Register the route in `internal/api/router.go`
4. Add tests in `*_test.go`
5. Update API documentation

### Adding a New CLI Command

1. Create command in `cmd/rdap/commands`
2. Register in `cmd/rdap/main.go`
3. Add tests and documentation

### Updating Dependencies

1. Update Go modules:
   ```bash
   go get -u ./...
   go mod tidy
   ```

2. Verify everything works:
   ```bash
   make test
   make build
   ```

## Release Process

1. Update version in `version.go`
2. Update CHANGELOG.md
3. Create a new tag:
   ```bash
   git tag v1.2.3
   git push origin v1.2.3
   ```

4. GitHub Actions will:
   - Run tests
   - Build binaries
   - Create GitHub release
   - Push Docker image

## Getting Help

- Check existing issues and documentation
- Join our community chat
- Ask questions in pull requests
- Contact the maintainers
