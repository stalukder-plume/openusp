# OpenUSP Standardized Entrypoint System

## Overview

All OpenUSP services now use a standardized entrypoint system that provides:
- Consistent command-line argument parsing
- Environment variable management
- Help system
- Service-specific configuration
- Common operational patterns

## Architecture

### Core Components

1. **`build/common/entrypoint-common.sh`** - Shared functions and patterns
2. **`build/[service]/entrypoint.sh`** - Service-specific entrypoint scripts
3. **Dockerfiles** - Updated to use entrypoint scripts instead of direct binary execution

### Services

- **apiserver** - REST API server (port 8081)
- **controller** - USP controller service (port 8082)
- **cli** - Command-line interface tool
- **cwmpacs** - CWMP ACS server (ports 7547/7548)
- **obuspa** - OB-USP-A implementation (ports 8080/8443)

## Usage Examples

### API Server
```bash
# Basic usage
docker run openusp-apiserver

# With custom port and debug logging
docker run openusp-apiserver --port 9081 --log-level debug

# With TLS enabled
docker run openusp-apiserver --tls-enabled --cert-file /certs/server.crt --key-file /certs/server.key

# With configuration file
docker run -v ./config.yaml:/opt/openusp/configs/apiserver.yaml \
           openusp-apiserver --config-file /opt/openusp/configs/apiserver.yaml
```

### Controller
```bash
# Basic usage
docker run openusp-controller

# With custom gRPC port
docker run openusp-controller --grpc-port 9083

# With custom MQTT broker
docker run openusp-controller --mqtt-broker mqtt://broker.example.com:1883
```

### CLI
```bash
# Interactive mode
docker run -it openusp-cli --interactive

# Execute specific command
docker run openusp-cli device list

# JSON output format
docker run openusp-cli --output json device get device-123

# Connect to different API endpoint
docker run openusp-cli --api-endpoint http://prod-api:8081 device list
```

### CWMP ACS
```bash
# Basic TR-069 ACS
docker run -p 7547:7547 -p 7548:7548 openusp-cwmpacs

# Custom ports
docker run -p 8547:8547 -p 8548:8548 \
           openusp-cwmpacs --http-port 8547 --https-port 8548

# With TLS enabled
docker run -p 7547:7547 -p 7548:7548 \
           -v ./certs:/opt/openusp/configs/certs \
           openusp-cwmpacs --tls-enabled --cert-file /opt/openusp/configs/certs/acs.crt
```

### OB-USP-A
```bash
# Basic USP agent
docker run -p 8080:8080 -p 8443:8443 openusp-obuspa

# With persistent database
docker run -v obuspa-data:/var/lib/obuspa \
           -p 8080:8080 -p 8443:8443 openusp-obuspa

# Custom database location
docker run openusp-obuspa --database /opt/obuspa/configs/custom.db
```

## Environment Variables

### Common Variables (All Services)
- `OPENUSP_PORT` - Service port
- `OPENUSP_LOG_LEVEL` - Log level (debug, info, warn, error)
- `OPENUSP_CONFIG_FILE` - Configuration file path
- `OPENUSP_DB_URL` - Database URL

### Service-Specific Variables

#### API Server
- `APISERVER_HOST` - Bind host address (default: 0.0.0.0)
- `APISERVER_TLS_ENABLED` - Enable TLS (default: false)
- `APISERVER_CERT_FILE` - TLS certificate file path
- `APISERVER_KEY_FILE` - TLS private key file path
- `APISERVER_CORS_ENABLED` - Enable CORS (default: true)

#### Controller
- `CONTROLLER_GRPC_PORT` - gRPC server port (default: 8083)
- `CONTROLLER_MQTT_BROKER` - MQTT broker URL
- `CONTROLLER_STOMP_BROKER` - STOMP broker URL
- `CONTROLLER_COAP_PORT` - CoAP server port (default: 5683)
- `CONTROLLER_WEBSOCKET_PORT` - WebSocket server port (default: 8084)

#### CLI
- `CLI_API_ENDPOINT` - API server endpoint
- `CLI_OUTPUT_FORMAT` - Output format (table, json, yaml)
- `CLI_TIMEOUT` - Request timeout duration
- `CLI_INTERACTIVE` - Interactive mode enabled

#### CWMP ACS
- `CWMP_ACS_HTTP_PORT` - HTTP port (default: 7547)
- `CWMP_ACS_HTTPS_PORT` - HTTPS port (default: 7548)
- `CWMP_ACS_TLS_ENABLED` - TLS enabled flag
- `CWMP_ACS_CERT_FILE` - TLS certificate file path
- `CWMP_ACS_KEY_FILE` - TLS private key file path

#### OB-USP-A
- `OBUSPA_HTTP_PORT` - HTTP server port (default: 8080)
- `OBUSPA_HTTPS_PORT` - HTTPS server port (default: 8443)
- `OBUSPA_DATABASE_FILE` - Database file path
- `OBUSPA_FACTORY_RESET_FILE` - Factory reset configuration file
- `OBUSPA_TRUST_STORE_DIR` - Trust store directory

## Help System

Every service provides comprehensive help:

```bash
# Get general help
docker run openusp-apiserver --help

# Service-specific help is included in the output
docker run openusp-cwmpacs --help
```

## Directory Structure

All services follow the same directory structure:

```
/opt/openusp/          # Main application directory (or /opt/obuspa/ for obuspa)
├── bin/              # Service binaries
├── configs/          # Configuration files (mountable)
└── scripts/          # Entrypoint and helper scripts
/var/lib/openusp/     # Application data (persistent)
/var/log/openusp/     # Log files (mountable)
```

## Benefits

1. **Consistency** - All services use the same argument and environment patterns
2. **Flexibility** - Easy to configure services for different environments
3. **Documentation** - Built-in help system for each service
4. **Maintainability** - Shared common functions reduce duplication
5. **Production Ready** - Proper logging, configuration, and data management
6. **Docker Best Practices** - Non-root execution, proper directory structure
7. **Development Friendly** - Easy to override settings for local development

## Migration from Direct Binary Execution

**Before:**
```dockerfile
ENTRYPOINT ["/opt/openusp/bin/apiserver"]
```

**After:**
```dockerfile
ENTRYPOINT ["/opt/openusp/scripts/entrypoint.sh"]
CMD []
```

This change provides much more flexibility while maintaining backward compatibility for most use cases.