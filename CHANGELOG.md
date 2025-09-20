# Changelog

All notable changes to the OpenUSP project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Complete YAML-based configuration system for all services
- Unified configuration loader with environment variable support
- Swagger UI integration for API documentation at http://localhost:8080
- Enhanced development scripts for infrastructure and application management
- Comprehensive CWMP/TR-069 ACS server implementation
- Docker Compose setup for local development infrastructure
- Service health checks and monitoring endpoints
- Enhanced logging configuration with structured output
- API specification file (`api/openusp.yaml`) for Swagger UI

### Enhanced
- **Development Environment**: Complete overhaul of local development setup
  - Infrastructure services (MongoDB, ActiveMQ, Redis, Swagger UI) in containers
  - Application services run locally for faster development iteration
  - Centralized environment configuration
  - Automated service startup and shutdown scripts
- **Project Structure**: Migrated to Go standard layout
  - `cmd/` for application entry points
  - `internal/` for private application code
  - `pkg/` for public libraries
  - `configs/` for YAML configuration files
- **Configuration Management**: 
  - Service-specific YAML configs (apiserver.yaml, controller.yaml, etc.)
  - Environment variable fallback support
  - Configuration validation and error handling
- **API Server**: Enhanced with health checks, structured logging, and YAML config
- **Controller**: Improved message handling and database integration
- **CLI Tool**: Better configuration management and API server integration
- **CWMP ACS**: Full TR-069 implementation with SOAP message handling

### Changed
- **Configuration**: Migrated from environment variables to YAML-based configuration
- **Database Connection**: Enhanced with connection pooling and timeout configuration
- **Service Discovery**: Updated to use YAML configuration instead of environment variables
- **Development Workflow**: Simplified with automated scripts and containerized infrastructure
- **Documentation**: Updated `DEVELOPMENT.md` with comprehensive setup instructions
- **Build System**: Enhanced Makefile with better dependency management

### Removed
- Legacy environment variable configuration files
- Disabled test files and empty test directories
- Deprecated configuration patterns
- Test files with skip directives:
  - `internal/cli/testutils.go.disabled`
  - `internal/cli/wifi_test.go.disabled` 
  - `internal/cli/main_test.go.disabled`
  - `internal/cli/ip_test.go.disabled`
  - `internal/controller/msg_test.go.disabled`
  - `internal/controller/mtp_test.go` (contained t.Skip())
  - `internal/cli/test/` (empty directory)

### Fixed
- Service startup dependencies and ordering
- Configuration loading errors and fallback mechanisms
- Database connection reliability with retry logic
- Message broker connectivity issues
- Service health check implementations

### Security
- Enhanced TLS configuration support for all protocols
- Improved authentication mechanisms
- Secure credential handling in configuration

### Infrastructure
- **MongoDB**: Latest version with authentication and health checks
- **ActiveMQ**: Multi-protocol support (STOMP, MQTT, OpenWire) with web console
- **Redis**: Caching layer with persistence
- **Swagger UI**: API documentation with interactive testing

### Development Tools
- **Scripts**: Automated infrastructure and application management
  - `scripts/start-infrastructure.sh` - Start containerized services
  - `scripts/start-applications.sh` - Start local applications
  - `scripts/stop-applications.sh` - Gracefully stop applications
  - `scripts/load-env.sh` - Load development environment variables
- **Docker Compose**: Complete local development environment
- **Logging**: Centralized log management in `logs/` directory
- **Monitoring**: Service health checks and status monitoring

### API Documentation
- **Swagger UI**: Interactive API documentation at http://localhost:8080
- **OpenAPI Specification**: Complete API specification in `api/openusp.yaml`
- **Testing Interface**: Direct API testing from documentation interface

---

## Version History

### Development Setup Migration (September 2025)
This release represents a major overhaul of the development environment and project structure, transitioning from environment variable-based configuration to a modern YAML-based system with containerized infrastructure services and local application development.

### Key Benefits
- **Faster Development**: Local applications with quick rebuild cycles
- **Consistent Environment**: Containerized infrastructure ensures consistency
- **Better Documentation**: Comprehensive setup guides and API documentation
- **Improved Debugging**: Local processes with direct debugger access
- **Modern Architecture**: Go standard layout with clean separation of concerns

### Migration Notes
Developers should:
1. Update local development setup using new scripts
2. Review new YAML configuration files
3. Use Swagger UI for API testing and documentation
4. Follow updated development workflow in `DEVELOPMENT.md`