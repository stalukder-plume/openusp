# Components (Stub)

Purpose: Provide a concise catalog of all runtime services, their boundaries, ports, scaling traits, and key dependencies.

## Table of Contents
1. Summary Matrix
2. Service Descriptions
3. Inter-Service Communication
4. Scaling & HA Notes
5. Future / Planned Components

## 1. Summary Matrix
| Component | Binary | Role | Primary Ports | Depends On | Scale Pattern |
|-----------|--------|------|---------------|------------|---------------|
| API Server | apiserver | REST + auth gateway | 8081 (HTTP) | MongoDB, Redis, ActiveMQ | Stateless horizontal |
| Controller | controller | USP orchestration & MTP mgmt | 8082 (HTTP), gRPC? | MongoDB, ActiveMQ | Partitioned horizontal |
| CWMP ACS | cwmpacs | TR-069 ACS endpoint | 7547 (HTTP) | MongoDB, ActiveMQ | Session-aware horizontal |
| CLI | cli | Administrative / ops tooling | n/a | API, Controller | Local client |
| Message Broker | (external) | Event + MTP transport (STOMP/MQTT) | 61613/61614 | N/A | Cluster (broker native) |
| MongoDB | (external) | Persistent state store | 27017 | Disk | Replica set / shard |
| Redis | (external) | Cache / ephemeral state | 6379 | Memory | Primary + replicas |

## 2. Service Descriptions
### API Server
Responsibilities:
- HTTP REST endpoints
- AuthN/AuthZ (e.g., JWT)
- Basic request validation
Excludes business orchestration logic (delegated to Controller).

### Controller
Responsibilities:
- USP message handling
- Device lifecycle orchestration
- Subscription & event dispatch
Excludes HTTP user-facing APIs.

### CWMP ACS
Responsibilities:
- TR-069 session management
- Parameter set/get mediation
- Firmware/diagnostics bridging

### CLI
Responsibilities:
- Operator productivity
- Scriptable fleet actions
- Diagnostics

## 3. Inter-Service Communication (Initial View)
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
