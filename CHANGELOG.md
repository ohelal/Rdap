# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-12-15

### Added
- Initial stable release of the RDAP service
- Full RDAP protocol implementation for IP, ASN, and domain queries
- Redis caching integration for improved performance
- Kafka event streaming support
- CLI tool for easy querying
- Docker and Kubernetes support
- Comprehensive API documentation
- GitHub Actions workflow for CI/CD
- Unit and integration tests with `-short` flag support

### Changed
- Simplified GitHub Actions workflow to focus on essential checks
- Improved documentation structure and content
- Enhanced error handling and response formats
- Updated dependencies to latest stable versions

### Fixed
- Integration tests now handle type assertions correctly
- CI pipeline skips integration tests with `-short` flag
- Various code quality improvements and bug fixes

## [0.1.0] - 2024-12-01

### Added
- Initial beta release
- Basic RDAP server implementation
- Simple CLI tool
- Basic documentation
