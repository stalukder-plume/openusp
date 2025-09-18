# API & CLI (Stub)

Purpose: Quick orientation for interacting with OpenUSP programmatically or via the CLI. Deep reference should link to generated Swagger / future gRPC docs.

## 1. REST API
Base URL (default local): `http://localhost:8081`

Example: Health
```bash
curl -f http://localhost:8081/api/v1/health
```

Example: List Devices (JWT omitted for brevity)
```bash
curl http://localhost:8081/api/v1/devices
```

## 2. Common Endpoints (Draft)
| Purpose | Method | Path | Notes |
|---------|--------|------|-------|
| Health | GET | /api/v1/health | Liveness summary |
| List Devices | GET | /api/v1/devices | Pagination TBD |
| Get Device Params | GET | /api/v1/devices/{id}/params | Filter options |
| Set Device Param | POST | /api/v1/devices/{id}/params | JSON body |

## 3. CLI
Binary: `./build/bin/openusp-cli`

```bash
./build/bin/openusp-cli --help
./build/bin/openusp-cli devices list
./build/bin/openusp-cli params get --device <id> --path Device.DeviceInfo.ModelName
```

## 4. Authentication (Planned Outline)
| Mode | Description | Notes |
|------|-------------|-------|
| JWT | Bearer tokens | Default dev mode |
| mTLS | Mutual TLS | Production recommended |
| Basic | Simple user/pass | Transitional only |

## 5. Error Format (Draft)
```json
{
  "error": "invalid_parameter",
  "message": "Parameter path not found",
  "request_id": "abc123"
}
```

## 6. Versioning Strategy
- REST endpoints versioned under `/api/v1/`
- Backward compatible additions do not bump major API prefix

## 7. Future
- gRPC interface for internal services
- GraphQL gateway (planned)
- Async command/task status endpoints

## TODO
- Replace placeholder endpoints with generated list.
- Document pagination, filtering, and auth header examples.
- Provide curl + CLI parity examples.
