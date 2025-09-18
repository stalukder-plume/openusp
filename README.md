# OpenUSP - Universal Device Management Platform

<div align="center">

[![CI/CD Pipeline](https://github.com/stalukder-plume/openusp/actions/workflows/ci.yml/badge.svg)](https://github.com/stalukder-plume/openusp/actions/workflows/ci.yml)
[![Security Scan](https://github.com/stalukder-plume/openusp/actions/workflows/security.yml/badge.svg)](https://github.com/stalukder-plume/openusp/actions/workflows/security.yml)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![Docker](https://img.shields.io/badge/docker-ready-blue.svg)](https://docker.com)
[![Release](https://img.shields.io/github/v/release/stalukder-plume/openusp?include_prereleases)](https://github.com/stalukder-plume/openusp/releases)
[![Docker Pulls](https://img.shields.io/docker/pulls/stalukder-plume/openusp-controller)](https://hub.docker.com/r/stalukder-plume/openusp-controller)
[![USP TR-369](https://img.shields.io/badge/USP-TR--369-green.svg)](https://usp.technology)
[![TR-069](https://img.shields.io/badge/CWMP-TR--069-orange.svg)](https://www.broadband-forum.org)

**Production-ready device management platform supporting both modern USP (TR-369) and legacy TR-069 CWMP protocols**

[🚀 Quick Start](#-quick-start) • [🏗️ Architecture](#-architecture) • [📖 Documentation](#-documentation) • [🤝 Contributing](#-contributing)

</div>

## ✨ What is OpenUSP?

OpenUSP is a **unified device management platform** that bridges the gap between legacy TR-069 CWMP infrastructure and modern USP (User Services Platform) implementations. Built with **Go** and designed for **cloud-native deployments**.

### 🎯 Key Benefits

- **🔄 Unified Management** - Single platform for both USP and TR-069 devices
- **🏗️ Production Ready** - Battle-tested with comprehensive monitoring
- **📋 Standards Compliant** - Full TR-369 (USP) and TR-069 (CWMP) compliance
- **☁️ Cloud Native** - Containerized microservices with Docker/Kubernetes
- **🔌 Multi-Protocol** - STOMP, CoAP, MQTT, WebSocket, and SOAP support

## 🚀 Quick Start

Get OpenUSP running in under 5 minutes:

```bash
# Clone and start
git clone https://github.com/n4-networks/openusp.git
cd openusp
docker-compose up -d

# Verify
curl http://localhost:8081/api/v1/health
```

**Access Points:**
- 🌐 **REST API**: http://localhost:8081
- 📚 **API Documentation**: http://localhost:8080  
- 🔧 **TR-069 ACS**: http://localhost:7547
- 💻 **CLI Tools**: `./openusp-cli`

## 🏗️ Architecture

OpenUSP follows a cloud-native microservices architecture with support for both modern USP (TR-369) and legacy TR-069 CWMP protocols.

**Key architectural principles:**
- 🏗️ **Microservices Design** - Loosely coupled, independently deployable services
- ☁️ **Cloud Native** - Containerized with Kubernetes orchestration support  
- 🔄 **Protocol Agnostic** - Unified management for USP and CWMP devices
- 📊 **Event Driven** - Real-time device state changes and notifications
- 🛡️ **Security First** - TLS encryption, certificate-based authentication

For detailed architecture information, see [📋 Architecture Documentation](docs/ARCHITECTURE.md).

### Core Components

| Component | Purpose | Technology | Ports |
|-----------|---------|------------|-------|
| **API Server** | REST gateway with OpenAPI 3.0 | Go, Gorilla Mux | 8081, 8443 |
| **Controller** | USP device management engine | Go, gRPC | 8082 |
| **CWMP ACS** | TR-069 Auto Configuration Server | Go, SOAP/HTTP | 7547, 7548 |
| **CLI** | Interactive command-line interface | Go, Ishell | N/A |

## 🌟 Features

### USP (TR-369) Support
- Multi-protocol transport (STOMP, CoAP, MQTT, WebSocket)
- Device lifecycle management and real-time configuration
- Parameter operations (Get/Set/Add/Delete) with bulk operations
- Event subscriptions and notifications
- Security with end-to-end encryption

### TR-069 CWMP Support  
- Complete SOAP/HTTP protocol implementation
- Device operations (parameters, file transfers, reboot, factory reset)
- Session management with proper timeout handling
- Bidirectional communication and connection requests
- HTTP Basic/Digest authentication

### Platform Features
- Microservices architecture with Docker containerization
- MongoDB database with Redis caching
- Health monitoring and audit logging
- Interactive Swagger UI documentation
- Comprehensive CLI with multiple output formats

## 📦 Installation

### Docker Deployment (Recommended)

```bash
git clone https://github.com/n4-networks/openusp.git
cd openusp
docker-compose up -d
```

### Building from Source

```bash
# Prerequisites: Go 1.21+, MongoDB, Redis, ActiveMQ
git clone https://github.com/n4-networks/openusp.git
cd openusp
make build
```

## 💻 Usage

### CLI Interface
```bash
# Interactive mode
./openusp-cli

# Direct commands
./openusp-cli show agents
./openusp-cli get Device.WiFi.SSID
./openusp-cli cwmp list-devices
```

### REST API Examples
```bash
# Health check
curl http://localhost:8081/api/v1/health

# List USP agents
curl -u admin:admin http://localhost:8081/api/v1/agents

# List CWMP devices
curl -u admin:admin http://localhost:8081/api/v1/cwmp/devices
```

## ⚙️ Configuration

Key environment variables:

```bash
# Database
DB_ADDR=localhost:27017
DB_USER=admin
DB_PASSWD=admin

# API Server
HTTP_PORT=8081
API_SERVER_AUTH_NAME=admin
API_SERVER_AUTH_PASSWD=admin

# CWMP ACS
CWMP_ACS_ENABLE=true
CWMP_ACS_PORT=7547
CWMP_ACS_USERNAME=acs
CWMP_ACS_PASSWORD=admin

# Protocol Settings
STOMP_ADDR=localhost:61613
MQTT_ADDR=localhost:1883
```

## 📖 Documentation

- [🏗️ Architecture Guide](docs/ARCHITECTURE.md) - Comprehensive technical architecture
- [🚀 Deployment Guide](docs/DEPLOYMENT.md) - Production deployment guide  
- [📡 TR-069 CWMP Guide](docs/TR069_CWMP_ACS.md) - CWMP ACS setup and usage
- [🔧 API Documentation](api/README.md) - REST API reference
- [💻 CLI Reference](pkg/cli/doc/) - Command-line interface docs

## 🚦 Network Ports

| Port | Service | Protocol | Description |
|------|---------|----------|-------------|
| 8081 | API Server | HTTP | REST API endpoints |
| 8082 | Controller | gRPC | Internal service communication |
| 7547 | CWMP ACS | HTTP | TR-069 ACS server |
| 7548 | CWMP ACS | HTTPS | Secure TR-069 ACS |
| 8080 | Swagger UI | HTTP | API documentation |
| 27017 | MongoDB | TCP | Database server |
| 6379 | Redis | TCP | Cache server |
| 61613 | ActiveMQ | TCP | Message broker |

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md).

### Development Workflow
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make changes and add tests
4. Ensure tests pass (`make test`)
5. Commit changes (`git commit -m 'Add amazing feature'`)
6. Push and open a Pull Request

## 📄 License

Licensed under the Apache License 2.0 - see [LICENSE](LICENSE) for details.

## 📞 Support

- 🌐 **Website**: [openusp.org](https://openusp.org)
- 📚 **Documentation**: [docs.openusp.org](https://docs.openusp.org)
- 🐛 **Issues**: [GitHub Issues](https://github.com/n4-networks/openusp/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/n4-networks/openusp/discussions)

---

<div align="center">

**⭐ Star us on GitHub if OpenUSP helps you manage your devices! ⭐**

**OpenUSP - Bridging legacy TR-069 and modern USP device management**

</div>
