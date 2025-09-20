# OpenUSP Documentation Index

This directory contains comprehensive documentation for the OpenUSP platform following the current Go standard layout and YAML-based configuration system.

## Index

| Document | Purpose | Status |
|----------|---------|--------|
| [ARCHITECTURE.md](ARCHITECTURE.md) | System design, microservices architecture, Go standard layout | ✅ Updated |
| [COMPONENTS.md](COMPONENTS.md) | Service descriptions, package structure, inter-service communication | ✅ Updated |
| [CONFIGURATION.md](CONFIGURATION.md) | YAML configuration, environment variables, production setup | ✅ Updated |
| [YAML_CONFIGURATION.md](YAML_CONFIGURATION.md) | Detailed YAML configuration system documentation | ✅ Current |
| [DEVELOPMENT.md](DEVELOPMENT.md) | Development environment, Go standard layout, build workflows | ✅ Updated |
| [DEPLOYMENT.md](DEPLOYMENT.md) | Docker Compose, Kubernetes deployment with YAML configs | ✅ Updated |
| [API.md](API.md) | REST API documentation and usage examples | 📝 Needs Update |
| [DATABASE.md](DATABASE.md) | Database schema, migrations, access patterns | 📝 Needs Update |
| [SECURITY.md](SECURITY.md) | Authentication, authorization, TLS configuration | 📝 Needs Update |
| [OPERATIONS.md](OPERATIONS.md) | Monitoring, logging, troubleshooting | 📝 Needs Update |
| [RELEASES.md](RELEASES.md) | Release process, versioning, Docker images | 📝 Needs Update |
| [PROTOCOLS.md](PROTOCOLS.md) | USP (TR-369) & CWMP (TR-069) protocol implementation | 📝 Needs Update |
| [CWMP.md](CWMP.md) | CWMP/TR-069 specific documentation | 📝 Needs Update |
| [USP.md](USP.md) | USP/TR-369 specific documentation | 📝 Needs Update |

## Recent Updates (September 2025)

### ✅ Completed Migrations
- **Go Standard Layout**: Migrated from `pkg/` to `internal/` + `pkg/` structure
- **YAML Configuration**: Replaced environment-only config with structured YAML files
- **Docker Integration**: Updated Dockerfiles to use YAML configurations
- **Service Architecture**: Refactored microservices following Go best practices

### 📁 Current Project Structure
```
openusp/
├── cmd/                  # Application entry points
├── internal/            # Private application code  
├── pkg/                # Public libraries (config, pb)
├── configs/            # YAML configuration files
├── deployments/        # Docker deployment manifests
└── docs/              # This documentation
```

### 🔧 Configuration System
- **Service Configs**: `configs/{apiserver,controller,cli,cwmpacs}.yaml`
- **Environment Substitution**: `${VAR_NAME:default_value}` syntax
- **Type Safety**: Structured YAML with validation
- **Migration**: Complete removal of `.env` file dependencies

## Navigation Guide

### 🚀 Getting Started
1. [DEVELOPMENT.md](DEVELOPMENT.md) - Set up your development environment
2. [CONFIGURATION.md](CONFIGURATION.md) - Understand the YAML configuration system
3. [DEPLOYMENT.md](DEPLOYMENT.md) - Deploy with Docker Compose

### 🏗️ Architecture & Design
1. [ARCHITECTURE.md](ARCHITECTURE.md) - High-level system overview
2. [COMPONENTS.md](COMPONENTS.md) - Detailed component breakdown
3. [YAML_CONFIGURATION.md](YAML_CONFIGURATION.md) - Configuration deep dive

### 🔧 Operations
1. [OPERATIONS.md](OPERATIONS.md) - Monitoring and troubleshooting
2. [SECURITY.md](SECURITY.md) - Security implementation
3. [RELEASES.md](RELEASES.md) - Release and upgrade procedures

## Contribution Guidelines
- Keep root README.md brief—detailed explanations belong in this docs/ directory
- Update documentation when making architectural changes
- Use consistent formatting and cross-reference related documents
- Include code examples in configuration and development guides

## Status Legend
- ✅ **Updated**: Recently updated to reflect current codebase
- 📝 **Needs Update**: Contains outdated information requiring revision  
- 🆕 **New**: Recently created documentation
- ❌ **Deprecated**: No longer relevant to current architecture