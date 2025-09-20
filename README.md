# OpenUSP

<div align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/stalukder-plume/openusp)](https://goreportcard.com/report/github.com/stalukder-plume/openusp)
[![Release](https://img.shields.io/github/v/release/stalukder-plume/openusp?include_prereleases)](https://github.com/stalukder-plume/openusp/releases)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

<strong>Unified USP (TR-369) + CWMP (TR-069) device management platform</strong>

<sub>Brief overview only. Full documentation lives in <a href="docs/README.md">/docs</a>.</sub>

</div>

---

## 1. What Is It?
OpenUSP is a cloud-friendly control plane for managing heterogeneous broadband CPE fleets via modern USP and legacy CWMP protocols. It provides a small core you can deploy with Docker Compose for evaluation or scale out in Kubernetes for production.

### Core Capabilities
- Dual protocol support (USP + CWMP)
- API + CLI for fleet operations
- Message-driven design (ActiveMQ STOMP/MQTT)
- Semantic versioned, multi-arch Docker images

---

## 2. Try It (5‑Minute Demo)

```bash
git clone https://github.com/stalukder-plume/openusp.git
cd openusp

# Start infrastructure services (MongoDB, ActiveMQ, Redis, Swagger UI)
./scripts/start-infrastructure.sh

# Health check (no authentication required)
curl -f http://localhost:8081/health

# Build CLI (optional)
make build-cli
./build/bin/openusp-cli devices list || true
```

Endpoints:
- REST API: `http://localhost:8081`
- Swagger UI: `http://localhost:8080` (interactive API documentation)
- CWMP ACS: `http://localhost:7547`
- ActiveMQ Console: `http://localhost:8161/admin` (admin/admin)

Stop stack: `./scripts/stop-applications.sh && docker-compose -f deployments/docker-compose-local-dev.yaml down`

---

## 3. Build From Source
Prerequisites: Go 1.21+, Make, Docker (optional for images)

```bash
make deps        # (idempotent)
make build-all   # apiserver, controller, cli, cwmpacs
make test        # run unit tests
```

Artifacts land in `./build/bin`.

---

## 4. Configuration
OpenUSP uses YAML-based configuration with environment variable substitution.

Configuration files:
```bash
configs/apiserver.yaml   # API server settings
configs/controller.yaml  # Controller settings  
configs/cli.yaml        # CLI configuration
configs/cwmpacs.yaml    # CWMP ACS settings
```

Environment variable substitution syntax:
```yaml
database:
  uri: ${OPENUSP_MONGO_URI:mongodb://localhost:27017}
  name: ${OPENUSP_DB_NAME:openusp}
redis:
  addr: ${OPENUSP_REDIS_ADDR:localhost:6379}
```

See `docs/CONFIGURATION.md` and `docs/YAML_CONFIGURATION.md` for detailed configuration options.

---

## 5. CLI Quick Glance
```bash
./build/bin/openusp-cli --help
./build/bin/openusp-cli devices list
./build/bin/openusp-cli params get --device <id> --path Device.DeviceInfo.ModelName
```
More examples: `docs/API.md` (stub).

---

## 6. Releases & Images
```bash
make version-check      # current tag
make version-suggest    # next semver hints
make release VERSION=vX.Y.Z
```
Images (multi-arch):
```
stalukder-plume/openusp-controller:<tag>
stalukder-plume/openusp-apiserver:<tag>
stalukder-plume/openusp-cli:<tag>
stalukder-plume/openusp-cwmpacs:<tag>
```
Details: `docs/RELEASES.md`.

---

## 7. Project Layout (Go Standard Layout)
```
cmd/            # Application entry points (apiserver, controller, cli, cwmpacs)
internal/       # Private application code (not importable by other projects)
  ├─ apiserver/ # REST API server implementation
  ├─ controller/# USP/CWMP controller logic
  ├─ cli/       # Command-line interface
  ├─ cwmp/      # CWMP protocol handlers
  ├─ db/        # Database access layer
  ├─ mtp/       # Message Transport Protocol
  └─ parser/    # Protocol message parsing
pkg/            # Public libraries (importable by external projects)
  ├─ config/    # YAML configuration management
  └─ pb/        # Protocol Buffer definitions
configs/        # YAML configuration files
deployments/    # Docker Compose manifests
scripts/        # Release + utility scripts
docs/           # Extended documentation
api/            # OpenAPI/Swagger specifications
```

---

## 8. Contributing (Fast Path)
```bash
git checkout -b feat/my-change
make build test lint
git commit -m "feat: my change"
git push origin HEAD
```
Please follow conventional commits & add tests for behavior changes.
See `docs/DEVELOPMENT.md` (stub) for full workflow.

---

## 9. Where Next?
| Need | Start Here |
|------|------------|
| Architecture overview | `docs/ARCHITECTURE.md` |
| Component responsibilities | `docs/COMPONENTS.md` |
| Config & env vars | `docs/CONFIGURATION.md` |
| Deployment patterns | `docs/DEPLOYMENT.md` |
| Protocol specifics | `docs/PROTOCOLS.md` |
| Operations & monitoring | `docs/OPERATIONS.md` |
| Release process | `docs/RELEASES.md` |
| Security model | `docs/SECURITY.md` |

---

## 10. License
Apache 2.0 — see `LICENSE`.

---
<sub>Intentionally concise. Expand only via linked docs to avoid drift.</sub>
