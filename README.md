[![CI Build Status](https://github.com/n4-networks/openusp/actions/workflows/build.yml/badge.svg)](https://github.com/n4-networks/openusp/actions/workflows/build.yml)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![Docker](https://img.shields.io/badge/docker-ready-blue.svg)](https://docker.com)

# OpenUSP - Universal Device Management Platform

Open source implementation of **USP (User Services Platform)** controller based on Broadband Forum's [USP specification](https://usp.technology), now enhanced with **TR-069 CWMP ACS** server capabilities for comprehensive legacy and next-generation device management.

> 🎉 **Latest Update**: OpenUSP now provides unified management for both USP (TR-369) and TR-069 CWMP devices in a single, production-ready platform.

## 📋 Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [Installation](#installation)
- [Usage](#usage)
- [API Reference](#api-reference)
- [Configuration](#configuration)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

## 🎯 Overview

OpenUSP bridges the gap between legacy TR-069 device management and modern USP protocol implementations, providing a unified platform for managing broadband devices, IoT devices, and network equipment. Built with Go and containerized deployment, it offers scalable, production-ready device management capabilities.

### Why OpenUSP?

- **Unified Management**: Single platform for USP and TR-069 devices
- **Production Ready**: Battle-tested with comprehensive error handling and monitoring
- **Standards Compliant**: Full TR-369 (USP) and TR-069 (CWMP) specification compliance
- **Cloud Native**: Docker-based deployment with microservices architecture
- **Developer Friendly**: Complete CLI tools, REST APIs, and extensive documentation

## ⭐ Key Features

### 🚀 USP (TR-369) Support
- **Multi-Protocol Transport**: STOMP, CoAP, MQTT, WebSocket message protocols
- **Device Lifecycle Management**: Discovery, onboarding, configuration, and monitoring
- **Parameter Operations**: Get/Set/Add/Delete with full data model support
- **Object Management**: Create, update, and delete managed objects
- **Event Subscriptions**: Real-time notifications and event handling
- **Bulk Operations**: Efficient batch operations for large device deployments
- **Security**: End-to-end encryption, certificate-based authentication
- **Standards Compliance**: Full TR-369 specification implementation

### 🔧 TR-069 CWMP ACS Support
- **Legacy Device Support**: Seamless management of existing TR-069 infrastructure
- **SOAP/HTTP Protocol**: Complete CWMP protocol implementation with session management
- **Device Operations**: Parameter management, file transfers, factory resets, reboots
- **Session Management**: Stateful CWMP sessions with proper timeout handling
- **Connection Requests**: Bidirectional communication with CPE devices
- **File Transfer**: Firmware upgrades, configuration backup/restore
- **Fault Handling**: Comprehensive error reporting and recovery mechanisms
- **Authentication**: HTTP Basic authentication with configurable credentials

### 🏗️ Platform Features
- **Microservices Architecture**: Scalable, containerized deployment
- **MongoDB Integration**: High-performance document database with indexing
- **Redis Caching**: Session management and high-speed data access
- **REST API**: Comprehensive HTTP API for all operations
- **CLI Interface**: Interactive command-line tools for device management
- **Docker Deployment**: Complete containerization with Docker Compose
- **Health Monitoring**: Built-in health checks and monitoring endpoints
- **Audit Logging**: Comprehensive logging for compliance and troubleshooting

## 🏛️ Architecture

OpenUSP follows a microservices architecture with clear separation of concerns:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web UI        │    │   CLI Tools     │    │  External Apps  │
│   (Future)      │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   API Server    │
                    │   (REST API)    │
                    └─────────────────┘
                                 │
                    ┌─────────────────┐
                    │   Controller    │
                    │  (USP + CWMP)   │
                    └─────────────────┘
                                 │
        ┌────────────────────────┼────────────────────────┐
        │                       │                        │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CWMP ACS      │    │   MongoDB       │    │   Redis Cache   │
│   Server        │    │   Database      │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
        │
┌─────────────────┐
│   TR-069        │
│   Devices       │
└─────────────────┘

         MTP Layer (STOMP, CoAP, MQTT, WebSocket)
                           │
                ┌─────────────────┐
                │   USP Devices   │
                │   (TR-369)      │
                └─────────────────┘
```

### Core Components

| Component | Description | Technology |
|-----------|-------------|------------|
| **Controller** | Core device management engine | Go, gRPC |
| **API Server** | REST API gateway | Go, Gorilla Mux |
| **CWMP ACS** | TR-069 Auto Configuration Server | Go, SOAP/HTTP |
| **CLI** | Command-line interface | Go, Ishell |
| **Database** | Device data and configuration storage | MongoDB |
| **Cache** | Session and temporary data | Redis |
| **Message Broker** | USP protocol transport | ActiveMQ |

## 🚀 Quick Start

Get OpenUSP running in minutes with Docker Compose:

```bash
# Clone the repository
git clone https://github.com/n4-networks/openusp.git
cd openusp

# Set up aliases (optional)
source scripts/bash/aliases.sh

# Start all services
docker-compose up -d

# Verify services are running
docker-compose ps

# Access the CLI
./scripts/docker/cli.sh
```

### Service Status Check
```bash
$ docker-compose ps
NAME                    IMAGE                              STATUS
openusp-apiserver       n4networks/openusp-apiserver       Up 2 minutes
openusp-broker          islandora/activemq                 Up 2 minutes  
openusp-cache           redis:latest                       Up 2 minutes
openusp-controller      n4networks/openusp-controller      Up 2 minutes
openusp-cwmpacs         n4networks/openusp-cwmpacs         Up 2 minutes
openusp-db              mongo:latest                       Up 2 minutes (healthy)
```

## 📦 Installation

### Prerequisites

- **Docker & Docker Compose**: For containerized deployment
- **Go 1.21+**: For building from source  
- **MongoDB**: Database backend
- **Redis**: Caching layer
- **ActiveMQ**: Message broker for USP protocols

### Docker Deployment (Recommended)

1. **Clone and Start Services**:
   ```bash
   git clone https://github.com/n4-networks/openusp.git
   cd openusp
   docker-compose up -d
   ```

2. **Verify Installation**:
   ```bash
   # Check service health
   curl http://localhost:8081/api/v1/health
   
   # Check CWMP ACS
   curl http://localhost:7547/cwmp
   ```

### Building from Source

1. **Prerequisites**:
   ```bash
   # Install Go 1.21+
   # Install MongoDB, Redis, ActiveMQ
   ```

2. **Build Components**:
   ```bash
   git clone https://github.com/n4-networks/openusp.git
   cd openusp
   
   # Build all components
   make build
   
   # Or build individually
   make controller
   make apiserver  
   make cli
   make cwmpacs
   ```

3. **Configuration**:
   ```bash
   # Copy and edit configuration
   cp configs/openusp.env.example configs/openusp.env
   # Edit database and service URLs
   ```

## 💻 Usage

### CLI Interface

OpenUSP provides a comprehensive command-line interface for both USP and TR-069 device management:

```bash
# Start interactive CLI
$ ./openusp-cli
OpenUsp-Cli>> 
**************************************************************
                          OpenUSP CLI
**************************************************************

# USP Device Management
OpenUsp-Cli>> show agent all-ids
Agent Number              : 1           
EndpointId                : os::012345-000000000000
-------------------------------------------------
Agent Number              : 2           
EndpointId                : os::012345-02420A050007
-------------------------------------------------

# Get USP device details
OpenUsp-Cli>> show agent os::012345-000000000000
Object Path                 : Device.LocalAgent.
 UpTime                     : 3628        
 SupportedProtocols         : STOMP, CoAP, MQTT, WebSocket
 SoftwareVersion            : 7.0.2       
 EndpointID                 : os::012345-000000000000
 ControllerNumberOfEntries  : 1           
 MTPNumberOfEntries         : 4           

# USP Parameter Operations
OpenUsp-Cli>> get Device.DeviceInfo.SoftwareVersion
Parameter Value: 2.1.0-release

OpenUsp-Cli>> set Device.WiFi.Radio.1.Enable true
Operation successful

# TR-069 CWMP Device Management  
OpenUsp-Cli>> cwmp device list
Device ID                           Manufacturer    Model           Status
acs-device-001                      Broadcom        BCM63138        Online
acs-device-002                      Qualcomm        QCA9563         Offline

# CWMP Parameter Operations
OpenUsp-Cli>> cwmp param get acs-device-001 Device.DeviceInfo.SoftwareVersion
Parameter: Device.DeviceInfo.SoftwareVersion
Value: 2.1.0-beta
Type: string

OpenUsp-Cli>> cwmp param set acs-device-001 Device.WiFi.Radio.1.Enable true
Setting parameter Device.WiFi.Radio.1.Enable to true... Success

# CWMP File Operations
OpenUsp-Cli>> cwmp file download acs-device-001 http://server.com/firmware.bin firmware.bin
File transfer initiated successfully
Command Key: fw_upgrade_20240915_001

# CWMP Device Control
OpenUsp-Cli>> cwmp device reboot acs-device-001
Reboot command sent successfully
```

### REST API Examples

#### USP Device Management
```bash
# List all USP agents
curl -X GET http://localhost:8081/api/v1/agents

# Get agent information
curl -X GET http://localhost:8081/api/v1/agents/os::012345-000000000000

# Parameter operations
curl -X POST http://localhost:8081/api/v1/agents/os::012345-000000000000/get \
  -H "Content-Type: application/json" \
  -d '{"parameters": ["Device.DeviceInfo.SoftwareVersion"]}'

curl -X POST http://localhost:8081/api/v1/agents/os::012345-000000000000/set \
  -H "Content-Type: application/json" \
  -d '{"parameters": [{"name": "Device.WiFi.Radio.1.Enable", "value": "true"}]}'
```

#### TR-069 CWMP Device Management
```bash
# List all CWMP devices
curl -X GET http://localhost:8081/api/v1/cwmp/devices

# Get CWMP device details
curl -X GET http://localhost:8081/api/v1/cwmp/devices/acs-device-001

# Get parameter values
curl -X POST http://localhost:8081/api/v1/cwmp/devices/acs-device-001/params/get \
  -H "Content-Type: application/json" \
  -d '{"parameters": ["Device.DeviceInfo.SoftwareVersion", "Device.WiFi.Radio.1.Enable"]}'

# Set parameter values
curl -X POST http://localhost:8081/api/v1/cwmp/devices/acs-device-001/params/set \
  -H "Content-Type: application/json" \
  -d '{"parameters": [{"name": "Device.WiFi.Radio.1.Enable", "value": "true", "type": "boolean"}]}'

# Initiate file download
curl -X POST http://localhost:8081/api/v1/cwmp/devices/acs-device-001/files/download \
  -H "Content-Type: application/json" \
  -d '{"url": "http://server.com/firmware.bin", "filename": "firmware.bin", "type": "firmware"}'

# Reboot device
curl -X POST http://localhost:8081/api/v1/cwmp/devices/acs-device-001/reboot \
  -H "Content-Type: application/json" \
  -d '{"command_key": "reboot_maintenance_001"}'
```

## 📚 API Reference

### USP API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/agents` | List all USP agents |
| GET | `/api/v1/agents/{id}` | Get agent details |
| POST | `/api/v1/agents/{id}/get` | Get parameter values |
| POST | `/api/v1/agents/{id}/set` | Set parameter values |
| POST | `/api/v1/agents/{id}/add` | Add objects |
| POST | `/api/v1/agents/{id}/delete` | Delete objects |
| POST | `/api/v1/agents/{id}/operate` | Execute operations |
| GET | `/api/v1/agents/{id}/subscriptions` | List subscriptions |
| POST | `/api/v1/agents/{id}/subscribe` | Create subscription |

### TR-069 CWMP API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/cwmp/devices` | List all CWMP devices |
| GET | `/api/v1/cwmp/devices/{id}` | Get device details |
| POST | `/api/v1/cwmp/devices/{id}/params/get` | Get parameter values |
| POST | `/api/v1/cwmp/devices/{id}/params/set` | Set parameter values |
| POST | `/api/v1/cwmp/devices/{id}/files/download` | Download file to device |
| POST | `/api/v1/cwmp/devices/{id}/files/upload` | Upload file from device |
| POST | `/api/v1/cwmp/devices/{id}/reboot` | Reboot device |
| POST | `/api/v1/cwmp/devices/{id}/factory-reset` | Factory reset device |
| GET | `/api/v1/cwmp/devices/{id}/sessions` | Get active sessions |
| GET | `/api/v1/cwmp/stats` | Get ACS statistics |

## ⚙️ Configuration

### Environment Variables

OpenUSP uses environment variables for configuration. Key settings include:

```bash
# Database Configuration
DB_ADDR=localhost:27017
DB_USER=admin
DB_PASSWD=admin
DB_NAME=usp

# USP Controller Settings  
CNTLR_GRPC_PORT=9001
CNTLR_EPID=self::openusp-controller

# API Server Settings
HTTP_PORT=8081
API_SERVER_AUTH_NAME=admin
API_SERVER_AUTH_PASSWD=admin

# CWMP ACS Settings
CWMP_ACS_ENABLE=true
CWMP_ACS_PORT=7547
CWMP_ACS_TLS_PORT=7548
CWMP_ACS_USERNAME=acs
CWMP_ACS_PASSWORD=admin

# Protocol Settings
STOMP_ADDR=localhost:61613
MQTT_ADDR=localhost:1883
COAP_SERVER_PORT=5683
WS_SERVER_PORT=8080
```

### Docker Compose Configuration

```yaml
version: '3.7'
services:
  openusp-controller:
    image: n4networks/openusp-controller:latest
    ports:
      - "9001:9001"
    environment:
      - DB_ADDR=openusp-db:27017
      - STOMP_ADDR=openusp-broker:61613
      
  openusp-apiserver:
    image: n4networks/openusp-apiserver:latest
    ports:
      - "8081:8081"
    depends_on:
      - openusp-controller
      
  openusp-cwmpacs:
    image: n4networks/openusp-cwmpacs:latest
    ports:
      - "7547:7547"
      - "7548:7548"
    depends_on:
      - openusp-db
```

## 🔧 Development

### Building from Source

```bash
# Clone repository
git clone https://github.com/n4-networks/openusp.git
cd openusp

# Install dependencies
go mod download

# Build all components
make build

# Build specific components
make controller    # Build USP controller
make apiserver     # Build API server  
make cli          # Build CLI tools
make cwmpacs      # Build CWMP ACS server

# Run tests
make test

# Build Docker images
make images
```

### Project Structure

```
openusp/
├── cmd/                    # Application entry points
│   ├── controller/         # USP Controller main
│   ├── apiserver/         # API Server main  
│   ├── cli/               # CLI application main
│   └── cwmpacs/           # CWMP ACS Server main
├── pkg/                   # Core packages
│   ├── cntlr/            # Controller logic
│   ├── apiserver/        # API server handlers
│   ├── cli/              # CLI implementation
│   ├── cwmp/             # TR-069 CWMP implementation
│   ├── db/               # Database layer
│   ├── mtp/              # Message Transport Protocols
│   └── pb/               # Protocol buffer definitions
├── configs/               # Configuration files
├── deployments/          # Docker compose and deployment files
├── docs/                 # Documentation
├── scripts/              # Utility scripts
└── test/                 # Test files and test data
```

### Adding New Features

1. **Protocol Extensions**: Add new MTP protocols in `pkg/mtp/`
2. **API Endpoints**: Extend handlers in `pkg/apiserver/`
3. **CLI Commands**: Add commands in `pkg/cli/`
4. **Database Models**: Extend models in `pkg/db/`
5. **CWMP Methods**: Add TR-069 methods in `pkg/cwmp/`

## 🌐 Network Ports

| Port | Service | Protocol | Description |
|------|---------|----------|-------------|
| 8081 | API Server | HTTP | REST API endpoints |
| 8443 | API Server | HTTPS | Secure REST API |
| 9001 | Controller | gRPC | Internal service communication |
| 7547 | CWMP ACS | HTTP | TR-069 ACS server |
| 7548 | CWMP ACS | HTTPS | Secure TR-069 ACS |
| 5683 | CoAP | UDP | CoAP message transport |
| 5684 | CoAP | UDP | Secure CoAP (DTLS) |
| 8080 | WebSocket | TCP | WebSocket transport |
| 8443 | WebSocket | TCP | Secure WebSocket |
| 61613 | STOMP | TCP | STOMP message broker |
| 61614 | STOMP | TCP | Secure STOMP (TLS) |
| 1883 | MQTT | TCP | MQTT message broker |
| 8883 | MQTT | TCP | Secure MQTT (TLS) |
| 27017 | MongoDB | TCP | Database |
| 6379 | Redis | TCP | Cache server |

## 📖 Documentation

- **[TR-069 CWMP ACS Guide](docs/TR069_CWMP_ACS.md)**: Complete guide for TR-069 CWMP ACS setup and usage
- **[API Documentation](api/README.md)**: Detailed REST API reference
- **[CLI Reference](pkg/cli/doc/)**: Command-line interface documentation  
- **[Developer Guide](docs/DEVELOPMENT.md)**: Development setup and contribution guidelines
- **[Deployment Guide](docs/DEPLOYMENT.md)**: Production deployment recommendations
- **[Protocol Support](docs/PROTOCOLS.md)**: USP and TR-069 protocol implementation details

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`make test`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Broadband Forum](https://www.broadband-forum.org/) for USP and TR-069 specifications
- [USP Technology](https://usp.technology/) community
- Open source contributors and maintainers

## 📞 Support

- **Website**: [https://openusp.org](https://openusp.org)
- **Documentation**: [https://docs.openusp.org](https://docs.openusp.org)  
- **Issues**: [GitHub Issues](https://github.com/n4-networks/openusp/issues)
- **Discussions**: [GitHub Discussions](https://github.com/n4-networks/openusp/discussions)

---

**OpenUSP** - Bridging legacy TR-069 and modern USP device management in a unified, production-ready platform.


