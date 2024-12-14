.PHONY: all build test clean benchmark tune

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
BINARY_NAME=rdap_service

# Build flags
BUILD_FLAGS=-ldflags="-s -w" -trimpath

all: clean build test

build:
	$(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME) ./cmd/server

test:
	$(GOTEST) -v -race -cover ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

benchmark:
	./scripts/benchmark.sh

tune:
	sudo ./scripts/tune_system.sh

# Development targets
dev: build
	./$(BINARY_NAME)

lint:
	golangci-lint run

# Production targets
prod-build:
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME) ./cmd/server

docker-build:
	docker build -t rdap_service:latest .

# Test targets
test-coverage:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-bench:
	$(GOTEST) -v -bench=. -benchmem ./...

# Maintenance targets
deps:
	$(GOCMD) mod tidy
	$(GOCMD) mod verify

update-deps:
	$(GOCMD) get -u ./...
	$(GOCMD) mod tidy

# Documentation targets
docs:
	godoc -http=:6060

# Help target
help:
	@echo "Available targets:"
	@echo "  make          - Build and test the project"
	@echo "  make build    - Build the binary"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean build files"
	@echo "  make benchmark- Run performance benchmarks"
	@echo "  make tune     - Tune system for performance"
	@echo "  make dev      - Build and run for development"
	@echo "  make prod-build- Build for production"
	@echo "  make docker-build- Build Docker image"
	@echo "  make deps     - Verify and tidy dependencies"
	@echo "  make docs     - Start godoc server"
