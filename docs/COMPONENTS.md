# Components

This document provides a comprehensive catalog of all OpenUSP runtime services following the Go standard layout, their boundaries, ports, scaling traits, and key dependencies.

## Table of Contents
1. Summary Matrix
2. Service Descriptions
3. Package Structure
4. Inter-Service Communication
5. Scaling & HA Notes
6. Future / Planned Components

## 1. Summary Matrix
| Component | Entry Point | Implementation | Role | Primary Ports | Configuration | Depends On | Scale Pattern |
|-----------|-------------|----------------|------|---------------|---------------|------------|---------------|
| API Server | `cmd/apiserver/` | `internal/apiserver/` | REST + auth gateway | 8081 (HTTP) | `configs/apiserver.yaml` | MongoDB, Redis, ActiveMQ | Stateless horizontal |
| Controller | `cmd/controller/` | `internal/controller/` | USP orchestration & MTP mgmt | 8082 (HTTP), gRPC | `configs/controller.yaml` | MongoDB, ActiveMQ | Partitioned horizontal |
| CWMP ACS | `cmd/cwmpacs/` | `internal/cwmp/` | TR-069 ACS endpoint | 7547 (HTTP) | `configs/cwmpacs.yaml` | MongoDB, ActiveMQ | Session-aware horizontal |
| CLI | `cmd/cli/` | `internal/cli/` | Administrative / ops tooling | n/a | `configs/cli.yaml` | API, Controller | Local client |

### Shared Internal Packages
| Package | Location | Purpose | Used By |
|---------|----------|---------|---------|
| Database | `internal/db/` | Database access layer | apiserver, controller, cwmp |
| Parser | `internal/parser/` | Protocol message parsing | controller |
| MTP | `internal/mtp/` | Message Transport Protocol | controller |

### Public Packages
| Package | Location | Purpose | Accessibility |
|---------|----------|---------|---------------|
| Config | `pkg/config/` | YAML configuration management | Public library |
| Protocol Buffers | `pkg/pb/` | gRPC/Protobuf definitions | Public API |

## 2. Service Descriptions
### API Server (`internal/apiserver/`)
**Entry Point**: `cmd/apiserver/main.go`
**Configuration**: `configs/apiserver.yaml`

Responsibilities:
- HTTP REST endpoints with OpenAPI/Swagger documentation
- Authentication/Authorization (JWT-based)
- Request validation and response formatting
- WebSocket connections for real-time updates
- Integration with database layer via `internal/db/`

Excludes business orchestration logic (delegated to Controller).

### Controller (`internal/controller/`)
**Entry Point**: `cmd/controller/main.go`
**Configuration**: `configs/controller.yaml`

Responsibilities:
- USP message handling and processing
- Device lifecycle orchestration
- Subscription & event dispatch
- Message routing via `internal/mtp/`
- Protocol parsing via `internal/parser/`

Excludes HTTP user-facing APIs.

### CWMP ACS (`internal/cwmp/`)
**Entry Point**: `cmd/cwmpacs/main.go`  
**Configuration**: `configs/cwmpacs.yaml`

Responsibilities:
- TR-069 session management
- SOAP envelope processing
- Parameter set/get mediation
- Firmware/diagnostics bridging
- Legacy device support

### CLI (`internal/cli/`)
**Entry Point**: `cmd/cli/main.go`
**Configuration**: `configs/cli.yaml`

Responsibilities:
- Operator productivity tools
- Scriptable fleet actions
- System diagnostics and monitoring
- Integration testing utilities

## 3. Package Structure
The project follows Go standard layout principles:

### Internal Packages (`internal/`)
These packages are private to the OpenUSP project and cannot be imported by external projects:
- **apiserver/**: REST API implementation
- **controller/**: USP message processing engine
- **cli/**: Command-line interface implementation  
- **cwmp/**: CWMP/TR-069 protocol handlers
- **db/**: Database access layer (MongoDB, Redis)
- **mtp/**: Message Transport Protocol implementations
- **parser/**: USP/CWMP message parsing logic

### Public Packages (`pkg/`)
These packages can be imported by external projects:
- **config/**: YAML configuration management with environment variable substitution
- **pb/**: Protocol Buffer definitions for gRPC APIs

## 4. Inter-Service Communication
- API Server → Controller: internal RPC or DB indirection (to clarify)
- Controller ↔ Broker: STOMP/MQTT frames for USP MTP
- CWMP ACS ↔ Broker: event publication (future?)
- All services → MongoDB / Redis for state/persistence.

## 4. Scaling & HA Notes (Draft)
| Concern | Current Thought | TODO |
|---------|-----------------|------|
| Stateless API | Yes | Confirm session mgmt approach |
| Controller Partitioning | Device-domain partitioning | Define sharding key |
| CWMP Session Stickiness | Required per session | Document load-balancer strategy |
| Broker HA | Rely on external cluster | Provide HA config sample |
| MongoDB | Replica set baseline | Add sharding guidance |

## 5. Future / Planned Components
- Web UI (management console)
- GraphQL gateway
- Policy engine
- Edge/agent satellite controllers

## TODO
- Replace guesses with confirmed port & transport matrix.
- Add sequence diagrams for device provisioning & parameter update flows.
- Clarify gRPC usage if present.
