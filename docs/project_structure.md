# Project Structure

```
rdap_service/
├── api/                    # API documentation and OpenAPI/Swagger specs
│   └── swagger/
├── cmd/                    # Main applications
│   └── server/            # The RDAP server application
├── configs/               # Configuration files
│   ├── ipv4.json         # IPv4 bootstrap configuration
│   ├── ipv6.json         # IPv6 bootstrap configuration
│   ├── asn.json          # ASN bootstrap configuration
│   └── dns.json          # DNS bootstrap configuration
├── deployments/          # Deployment configurations
│   ├── docker/          # Docker related files
│   ├── k8s/             # Kubernetes manifests
│   └── terraform/       # Infrastructure as code
├── docs/                 # Documentation files
│   ├── architecture.md  # System architecture
│   ├── performance.md   # Performance tuning guide
│   └── api.md          # API documentation
├── internal/            # Private application code
│   ├── config/         # Configuration handling
│   ├── errors/         # Custom error types
│   ├── handlers/       # HTTP handlers
│   ├── logger/         # Logging package
│   ├── metrics/        # Metrics collection
│   ├── middleware/     # HTTP middleware
│   ├── models/         # Data models
│   └── service/        # Business logic
├── pkg/                # Public libraries that can be used by external applications
│   └── rdap/          # RDAP protocol implementation
├── scripts/           # Scripts for various tasks
│   ├── tune_system.sh # System tuning script
│   └── benchmark.sh   # Performance benchmarking script
├── tests/            # Additional test suites
│   ├── e2e/         # End-to-end tests
│   ├── integration/ # Integration tests
│   └── load/        # Load tests
├── tools/           # Tools and utilities
│   └── codegen/    # Code generation tools
├── vendor/         # Vendored dependencies
├── .dockerignore
├── .gitignore
├── Dockerfile
├── Makefile
├── README.md
└── go.mod
```

## Directory Descriptions

### `/api`
Contains API specifications and documentation. OpenAPI/Swagger specs live here.

### `/cmd`
Main applications for this project. Each subdirectory is a separate executable.

### `/configs`
Configuration file templates or default configs. RDAP bootstrap configurations.

### `/deployments`
IaaC, Docker, and orchestration files for various deployment environments.

### `/docs`
Design documents, user guides, and other documentation.

### `/internal`
Private application code. This is the code you don't want others importing.

### `/pkg`
Library code that's safe to use by external applications.

### `/scripts`
Scripts for build, install, analysis, etc.

### `/tests`
Additional external test apps and test data.

### `/tools`
Supporting tools for this project.

## Best Practices

1. Use `/internal` for private code that you don't want others importing
2. Put your main applications in `/cmd`
3. Use small packages and a flat structure
4. Group by context, not by type
5. Keep test files next to the code they test
6. Use `/pkg` for code you want others to import
