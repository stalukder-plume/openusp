# OpenUSP TR-069 CWMP ACS Server

This document describes the TR-069 CWMP (Customer Premises Equipment WAN Management Protocol) Auto Configuration Server (ACS) implementation in OpenUSP.

## Overview

The TR-069 CWMP ACS server provides device management capabilities for TR-069 compliant devices alongside the existing USP (User Services Platform) functionality in OpenUSP. This implementation allows managing both USP and CWMP devices from a unified platform.

## Features

### TR-069 Protocol Support
- CWMP protocol version 1.0, 1.1, 1.2, and 1.3 compatibility
- SOAP/HTTP and SOAP/HTTPS communication
- SSL/TLS support for secure communication
- Session management with proper state handling
- Connection request support

### Device Management
- Device discovery and registration via Inform messages
- Device inventory with manufacturer, model, and version tracking
- Online/offline status monitoring
- Device grouping and filtering capabilities

### Parameter Management
- Get/Set parameter values with full path support
- Parameter type validation and conversion
- Bulk parameter operations for efficiency
- Parameter change notifications and events

### File Transfer Operations
- Firmware upgrades via Download RPC
- Configuration backup/restore via Upload RPC
- File transfer status monitoring and notifications
- Scheduled file transfers with retry mechanisms

### Fault Handling
- Comprehensive fault code support per TR-069 specification
- Detailed error reporting and logging
- Automatic fault recovery mechanisms
- Fault history tracking for troubleshooting

## Architecture

### Components
1. **CWMP Protocol Layer** (`pkg/cwmp/soap.go`) - SOAP message handling and CWMP RPC implementations
2. **ACS Server** (`pkg/cwmp/acs.go`) - HTTP server with session management and device communication
3. **Controller Integration** (`pkg/cntlr/cwmp.go`) - Integration with existing OpenUSP controller
4. **Database Layer** (`pkg/db/cwmpdb.go`) - MongoDB collections for CWMP device and session data
5. **REST API** (`pkg/apiserver/cwmp_handlers.go`) - HTTP API endpoints for device management
6. **CLI Interface** (`pkg/cli/cwmp.go`) - Command-line tools for CWMP operations

### Database Schema
- **cwmpdevices** - Device inventory and metadata
- **cwmpsessions** - Active session tracking
- **cwmpparams** - Device parameter values and history
- **cwmpfiles** - File transfer operations and status

## Configuration

### Environment Variables
```bash
# Enable CWMP ACS Server
CWMP_ACS_ENABLE=true

# ACS Server Ports
CWMP_ACS_PORT=7547              # HTTP port (standard TR-069)
CWMP_ACS_TLS_PORT=7548          # HTTPS port

# Authentication
CWMP_ACS_USERNAME=acs           # ACS authentication username
CWMP_ACS_PASSWORD=admin         # ACS authentication password

# Connection Request Authentication  
CWMP_ACS_CONNECTION_REQUEST_USERNAME=cr
CWMP_ACS_CONNECTION_REQUEST_PASSWORD=admin

# Session Management
CWMP_ACS_SESSION_TIMEOUT=30     # Session timeout in minutes
CWMP_ACS_INFORM_INTERVAL=300    # Periodic inform interval in seconds
CWMP_ACS_PERIODIC_INFORM_ENABLE=true
```

### Docker Deployment
The CWMP ACS server runs as a separate container in the OpenUSP Docker Compose stack:

```yaml
openusp-cwmpacs:
  build:
    context: ..
    dockerfile: ./build/controller/Dockerfile
    target: cwmpacs
  image: n4networks/openusp-cwmpacs:latest
  container_name: openusp-cwmpacs
  ports:
    - 7547:7547   # HTTP TR-069 ACS port
    - 7548:7548   # HTTPS TR-069 ACS port
```

## Usage

### Starting the Services
```bash
# Start all OpenUSP services including CWMP ACS
docker-compose up -d

# Start only CWMP ACS server
docker-compose up openusp-cwmpacs
```

### Device Registration
TR-069 devices automatically register with the ACS when they send their initial Inform message. The device information is stored in the database and becomes available for management.

### CLI Operations
```bash
# List all CWMP devices
./openusp-cli cwmp device list

# Get device information
./openusp-cli cwmp device info <device-id>

# Get parameter values
./openusp-cli cwmp param get <device-id> <parameter-path>

# Set parameter values  
./openusp-cli cwmp param set <device-id> <parameter-path> <value>

# Reboot device
./openusp-cli cwmp device reboot <device-id>

# Initiate file transfer
./openusp-cli cwmp file download <device-id> <file-url> <target-filename>
```

### REST API
The CWMP functionality is exposed via REST API endpoints:

```
GET    /api/v1/cwmp/devices           - List all CWMP devices
GET    /api/v1/cwmp/devices/{id}      - Get device details
POST   /api/v1/cwmp/devices/{id}/params/get    - Get parameter values
POST   /api/v1/cwmp/devices/{id}/params/set    - Set parameter values
POST   /api/v1/cwmp/devices/{id}/reboot        - Reboot device
POST   /api/v1/cwmp/devices/{id}/files/download - Download file to device
POST   /api/v1/cwmp/devices/{id}/files/upload   - Upload file from device
```

## Security Considerations

### Authentication
- HTTP Basic Authentication for ACS server access
- Separate credentials for connection request operations
- Configurable username/password combinations

### Transport Security
- HTTPS/SSL support for encrypted communication
- Certificate-based device authentication (optional)
- Configurable SSL/TLS settings

### Access Control
- Role-based access control via API authentication
- Device-level access restrictions
- Audit logging for all operations

## Monitoring and Troubleshooting

### Logging
- Comprehensive logging at INFO, DEBUG, and ERROR levels
- SOAP message logging for protocol debugging
- Session state transition logging
- File transfer progress and status logging

### Health Checks
- ACS server health endpoint
- Database connectivity monitoring  
- Device connectivity status tracking
- Session timeout monitoring

### Performance Metrics
- Device count and online status
- Session duration statistics
- Parameter operation performance
- File transfer success rates

## Integration with USP

The TR-069 CWMP implementation runs alongside the existing USP functionality:

- **Unified Device Management** - Both USP and CWMP devices in single interface
- **Shared Database** - Common MongoDB instance for all device data  
- **Common API Server** - Single REST API for both protocols
- **Integrated CLI** - Unified command-line interface for all operations
- **Docker Deployment** - All services in single Docker Compose stack

## Compliance

This implementation follows the TR-069 specification standards:
- TR-069 Amendment 6 (2018) compliance
- CWMP Data Model support per TR-098, TR-181, and device-specific models
- Standard RPC method implementations
- Proper fault code handling and reporting
- Session management per specification requirements

## Development

### Building from Source
```bash
# Build CWMP ACS server
make cwmpacs

# Build Docker image
docker build -t openusp-cwmpacs --target=cwmpacs -f build/controller/Dockerfile .

# Run tests
cd cmd/cwmpacs && go test ./...
```

### Adding Custom Features
The modular architecture allows easy extension:
1. Add new RPC methods in `pkg/cwmp/soap.go`
2. Implement handlers in `pkg/cwmp/acs.go`
3. Add database operations in `pkg/db/cwmpdb.go`
4. Expose via API in `pkg/apiserver/cwmp_handlers.go`
5. Add CLI commands in `pkg/cli/cwmp.go`