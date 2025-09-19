# OpenUSP API Documentation

This folder contains the OpenAPI specifications and documentation for the OpenUSP platform APIs.

## Files

- **openusp.yaml** - Complete OpenAPI 3.0.3 specification for all OpenUSP APIs

## API Overview

The OpenUSP platform provides comprehensive REST APIs for managing both USP (TR-369) and TR-069 CWMP devices:

### USP (TR-369) APIs
- **Agent Management** - List and monitor USP agents
- **Data Model Operations** - Explore and manage device data models
- **Parameter Management** - Get and set device parameters
- **Instance Operations** - Create, read, update, delete object instances
- **Command Execution** - Execute operations and commands

### TR-069 CWMP APIs
- **Device Management** - Manage TR-069 devices and connections
- **Parameter Operations** - Get and set device parameters
- **Device Control** - Reboot, factory reset, connection requests
- **File Transfer** - Firmware updates, configuration management

### Administrative APIs
- **System Health** - Health checks and status monitoring
- **MTP Management** - Message Transport Protocol operations
- **Database Operations** - Administrative database functions

## Authentication

All APIs use HTTP Basic Authentication. Default credentials:
- Username: `admin`
- Password: `admin`

## Base URLs

- **Development**: http://localhost:8081
- **Swagger UI**: http://localhost:8080/swagger

## Usage with Docker Compose

When running the full OpenUSP stack with Docker Compose:

```bash
# Start all services including Swagger UI
docker-compose up -d

# Access Swagger UI
open http://localhost:8080

# Access API directly
curl -u admin:admin http://localhost:8081/health
```

## API Examples

### Get USP Agents
```bash
curl -u admin:admin http://localhost:8081/get/agents/
```

### Get Device Parameters
```bash
# USP device
curl -u admin:admin http://localhost:8081/get/params/os::012345-000000000000/Device.WiFi.

# TR-069 device
curl -u admin:admin http://localhost:8081/cwmp/device/acs-device-001/params
```

### Set Parameters
```bash
# USP parameter
curl -X POST -u admin:admin \
  -H "Content-Type: application/json" \
  -d '{"parameters":[{"name":"Device.WiFi.Radio.1.Enable","value":"true","type":"boolean"}]}' \
  http://localhost:8081/set/params/os::012345-000000000000/Device.WiFi.Radio.1.Enable

# TR-069 parameter
curl -X POST -u admin:admin \
  -H "Content-Type: application/json" \
  -d '{"parameters":[{"name":"Device.WiFi.Radio.1.Enable","value":"true","type":"boolean"}]}' \
  http://localhost:8081/cwmp/device/acs-device-001/params
```

## Interactive Documentation

The Swagger UI provides interactive documentation where you can:
- Explore all available endpoints
- Test API calls directly from the browser
- View request/response schemas
- Authenticate and execute real API calls

Access it at: http://localhost:8080 when running with Docker Compose.

## API Development

For API development and testing:

1. **Local Development**: Use the OpenAPI spec with tools like Postman or Insomnia
2. **Code Generation**: Generate client SDKs using OpenAPI Generator
3. **Testing**: Use the interactive Swagger UI for manual testing
4. **Integration**: Use automated testing tools or frameworks for API validation

## OpenAPI Specification

The `openusp.yaml` file is a complete OpenAPI 3.0.3 specification that includes:
- All endpoint definitions with parameters and responses
- Authentication schemes and security requirements
- Request/response schemas and examples
- Comprehensive documentation and descriptions
- Tagged organization for easy navigation

This specification can be used with any OpenAPI-compatible tool for documentation, testing, or client generation.
