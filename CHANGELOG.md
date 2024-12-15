# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-12-15

### Added
- Initial release of the RDAP service
- Full RDAP protocol implementation for domains, IP addresses, and ASNs
- Redis caching support
- Kafka event streaming integration
- CLI tool for easy querying
- Docker and Kubernetes support
- Comprehensive test suite with unit and integration tests

### Changed
- Simplified GitHub Actions workflow to focus on essential checks
- Added support for skipping integration tests in CI using -short flag

### Fixed
- Fixed type assertions in integration tests
- Resolved CDN package dependencies
- Fixed import issues in coalesced handler

## [0.1.0] - 2024-12-14

### Added
- Initial project setup
- Basic RDAP protocol implementation
- Redis caching integration
- Kafka event streaming setup
- Docker support
- GitHub Actions workflow setup
