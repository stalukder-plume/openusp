# Configuration (Stub)

Purpose: Central reference for environment variables, config files, and operational tuning parameters.

## 1. Environment Variables (Initial Set)
| Name | Default | Purpose | Notes |
|------|---------|---------|-------|
| OPENUSP_MONGO_URI | mongodb://mongo:27017 | MongoDB connection | Add auth in prod |
| OPENUSP_REDIS_ADDR | redis:6379 | Redis address | Password via secret |
| OPENUSP_AMQ_URI | stomp://activemq:61613 | Broker URI | Supports STOMP/MQTT |
| OPENUSP_LOG_LEVEL | info | Log verbosity | debug, info, warn, error |
| OPENUSP_AUTH_METHOD | jwt | Authentication mode | jwt|mtls|basic (planned) |
| OPENUSP_TLS_ENABLED | false | Enable TLS endpoints | Prod recommended |
| OPENUSP_USP_ENDPOINT_ID | controller.openusp.org | USP endpoint identifier | Unique per deployment |
| OPENUSP_USP_MTP_ENABLED | stomp,coap,mqtt | Allowed USP transports | Validate final list |

(Expand with authoritative list once codebase enumeration is done.)

## 2. Configuration Files
Provide sample `controller.yaml`, `apiserver.yaml`, and mapping to env overrides.

## 3. Secrets Management
- Recommend external secret store (K8s Secrets / Vault)
- Never bake credentials in images

## 4. TLS & Certificates
- Device cert trust store
- Controller/service mutual TLS (future doc)

## 5. Performance Tuning (Placeholders)
| Area | Knob | Description |
|------|------|-------------|
| MongoDB | Pool size | Concurrency tuning |
| Redis | TTL policies | Cache eviction |
| Broker | Heartbeats | Connection liveness |
| Go runtime | GOGC | Memory pressure control |

## 6. Configuration Precedence
1. Explicit CLI flags (if/when implemented)
2. Environment variables
3. Config file defaults
4. Internal hard-coded defaults

## TODO
- Auto-generate env var table from code.
- Add validation matrix (required vs optional).
- Document safe production baseline values.
