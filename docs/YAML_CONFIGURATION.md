# OpenUSP YAML Configuration Migration

This document describes the migration from environment variable-based configuration to YAML-based configuration for all OpenUSP microservices.

## Overview

All OpenUSP microservices have been migrated to use YAML configuration files instead of environment variables, following modern microservice configuration standards. This provides better structure, validation, and maintainability.

## Configuration Files

Each microservice has its own YAML configuration file:

- `configs/apiserver.yaml` - API Server configuration
- `configs/controller.yaml` - Controller configuration  
- `configs/cli.yaml` - CLI configuration
- `configs/cwmpacs.yaml` - CWMP ACS configuration

## Configuration Structure

All configuration files follow a consistent structure:

```yaml
service:
  name: "service-name"
  version: "${SERVICE_VERSION:1.0.0}"
  environment: "${ENVIRONMENT:development}"
  debug: ${DEBUG:false}

database:
  type: "mongodb"
  host: "${DB_HOST:localhost}"
  port: ${DB_PORT:27017}
  name: "${DB_NAME:database}"
  username: "${DB_USER:admin}"
  password: "${DB_PASSWD:admin}"
  pool:
    maxConnections: ${DB_MAX_CONNECTIONS:10}
    timeout: "${DB_TIMEOUT:30s}"

protocols:
  http:
    enabled: ${HTTP_ENABLED:true}
    host: "${HTTP_HOST:0.0.0.0}"
    port: ${HTTP_PORT:8080}
    enableTLS: ${HTTP_TLS_ENABLED:false}
    
  grpc:
    enabled: ${GRPC_ENABLED:true}
    host: "${GRPC_HOST:0.0.0.0}" 
    port: ${GRPC_PORT:9001}

security:
  auth:
    type: "basic"
    username: "${API_SERVER_AUTH_NAME:admin}"
    password: "${API_SERVER_AUTH_PASSWD:admin}"
    
logging:
  level: "${LOGGING:info}"
  format: "json"
  output: "stdout"
```

## Environment Variable Support

The YAML configuration files support environment variable substitution using the syntax:
`${VARIABLE_NAME:default_value}`

This ensures backward compatibility with existing environment-based deployments.

## Changes Made

### 1. New Configuration Package

Created `pkg/config/config.go` with:
- Comprehensive configuration structures
- YAML parsing with environment variable expansion
- Configuration validation
- Helper methods for common address formatting

### 2. Service Modifications

#### API Server (`pkg/apiserver/init.go`)
- Replaced `loadConfigFromEnv()` with `loadConfig()`
- Added YAML configuration loading
- Maintained backward compatibility with existing struct fields

#### Controller (`pkg/cntlr/cfg.go`)
- Replaced environment variable loading with YAML config
- Updated gRPC and USP configuration mapping
- Added config struct to main Cntlr type

#### CLI (`pkg/cli/cli.go`) 
- Migrated from godotenv to YAML configuration
- Updated REST client configuration
- Maintained authentication and history file settings

#### CWMP ACS (`pkg/cwmp/acs.go`)
- Replaced environment variable configuration with YAML
- Updated HTTP/HTTPS port configuration
- Added TLS certificate configuration from YAML

### 3. Docker File Updates

Updated all Dockerfiles to copy service-specific YAML configuration files:
- `build/apiserver/Dockerfile` - copies `apiserver.yaml`
- `build/controller/Dockerfile` - copies `controller.yaml`  
- `build/cli/Dockerfile` - copies `cli.yaml`
- `build/cwmpacs/Dockerfile` - copies `cwmpacs.yaml`

## Configuration File Locations

The configuration loader searches for files in the following order:
1. `./config.yaml`
2. `./configs/config.yaml`
3. `./openusp.yaml`
4. `./configs/openusp.yaml`
5. `/etc/openusp/config.yaml`
6. `/usr/local/etc/openusp/config.yaml`

Service-specific files (e.g., `apiserver.yaml`) are loaded directly from `./configs/`.

## Migration Benefits

1. **Structure**: YAML provides better organization and readability
2. **Validation**: Built-in configuration validation
3. **Environment Support**: Maintains environment variable compatibility  
4. **Documentation**: Self-documenting configuration format
5. **Microservice Standards**: Follows modern configuration patterns
6. **Maintainability**: Centralized configuration management

## Usage Examples

### Development Environment
```bash
# Use default configuration with environment overrides
export DB_HOST=localhost
export HTTP_PORT=8081
./apiserver
```

### Production Environment  
```bash
# Modify YAML files directly or use environment variables
export DB_HOST=prod-mongo.internal
export DB_USER=prod_user
export DB_PASSWD=secure_password
./apiserver
```

### Docker Deployment
Configuration files are automatically copied into containers and can be overridden with environment variables in docker-compose files.

## Testing

All services have been tested to ensure:
- ✅ Successful compilation
- ✅ Configuration loading from YAML files
- ✅ Environment variable substitution
- ✅ Docker container builds
- ✅ Backward compatibility with existing deployments

## Migration Path

1. **Phase 1** ✅ - Create YAML configuration files with environment variable support
2. **Phase 2** ✅ - Update Go applications to use YAML configuration  
3. **Phase 3** ✅ - Update Docker files to include configuration files
4. **Phase 4** - Deploy and validate in test environments
5. **Phase 5** - Gradual production rollout with monitoring

The migration maintains full backward compatibility, allowing for gradual rollout without service disruption.