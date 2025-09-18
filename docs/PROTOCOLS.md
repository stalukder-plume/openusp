# Protocols: USP & CWMP (Stub)

Purpose: Describe how OpenUSP implements / bridges USP (TR-369) and CWMP (TR-069), including mapping, extensions, and operational nuances.

## 1. USP (TR-369) Overview
| Aspect | Notes (Draft) |
|--------|---------------|
| Version Support | Targeting 1.x (confirm) |
| Transports | STOMP, CoAP, MQTT (validate list) |
| Security | TLS + endpoint identity (cert-based) |
| Data Model | TR-181 mapping subset |

### USP Message Flow (High-Level)
```
Agent → Broker → Controller → Persistence / Events → API Consumer
```

## 2. CWMP (TR-069) Overview
| Aspect | Notes |
|--------|-------|
| Amendment | 6 (planned) |
| Transport | HTTP/S SOAP |
| Sessions | Connection-based RPC style |
| Mapping | Parameter bridging into unified model |

## 3. Bridging Concepts
| Unified Concept | USP Source | CWMP Source | Notes |
|-----------------|-----------|------------|-------|
| Device Identity | EndpointID | SerialNumber | Normalized key TBD |
| Params | Path-based | Hierarchical | Provide translation layer |
| Events | Notify | Inform | Event normalization |

## 4. Parameter Mapping Strategy (Draft)
- Maintain translation catalog (JSON/YAML) (planned)
- Cache resolved mappings in Redis
- Fallback logging for unmapped parameters

## 5. Firmware / Operations
| Operation | USP Method | CWMP Method | Notes |
|-----------|-----------|------------|-------|
| Firmware Download | Download() | Download | Align progress reporting |
| Reboot | Reboot() | Reboot | Uniform audit log entry |
| Factory Reset | FactoryReset() | FactoryReset | Confirm device capability |

## 6. Diagnostics (Examples)
| Diagnostic | USP Path | CWMP RPC | Normalization Plan |
|------------|----------|----------|--------------------|
| Ping | Device.IP.Diagnostics.Ping | Ping | Standard result schema |
| Traceroute | Device.IP.Diagnostics.TraceRoute | TraceRoute | TBD |

## 7. Security Considerations
- Ensure transport-level encryption for both protocols
- Validate agent/device identity consistency across migrations

## 8. Migration Path (CWMP → USP)
Stages:
1. Discovery & inventory
2. Dual registration
3. Progressive feature shift
4. Protocol retirement (legacy subset)

## 9. Future Enhancements
- WebSocket USP transport
- Bulk parameter operations optimization
- Policy-driven parameter sync

## TODO
- Confirm supported USP version + transports.
- Add real mapping examples.
- Provide sequence diagrams for multi-step operations.
