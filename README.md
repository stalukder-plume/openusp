# OpenUSP

<div align="center">

[![CI](https://github.com/stalukder-plume/openusp/actions/workflows/ci.yml/badge.svg)](https://github.com/stalukder-plume/openusp/actions/workflows/ci.yml)
[![Security](https://github.com/stalukder-plume/openusp/actions/workflows/security.yml/badge.svg)](https://github.com/stalukder-plume/openusp/actions/workflows/security.yml)
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
docker-compose -f deployments/docker-compose.yaml up -d

# Health
curl -f http://localhost:8081/api/v1/health

# Build CLI (optional)
make build-cli
./build/bin/openusp-cli devices list || true
```

Endpoints:
- REST API: `http://localhost:8081`
- Swagger: `http://localhost:8080/swagger/` (if enabled)
- CWMP ACS: `http://localhost:7547`

Stop stack: `docker-compose -f deployments/docker-compose.yaml down -v`

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

## 4. Minimal Configuration
Most evaluation scenarios work with defaults from `deployments/docker-compose.yaml`.

Useful environment overrides:
```bash
export OPENUSP_MONGO_URI=mongodb://mongo:27017
export OPENUSP_REDIS_ADDR=redis:6379
export OPENUSP_AMQ_URI=stomp://activemq:61613
export OPENUSP_LOG_LEVEL=info
```
See `docs/CONFIGURATION.md` for the full matrix (to be filled).

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

## 7. Project Layout (Selected)
```
cmd/            # Entry points (apiserver, controller, cli, cwmpacs)
pkg/            # Internal packages (protocols, db, mtp, etc.)
deployments/    # Compose / (future) Helm manifests
scripts/        # Release + utility scripts
docs/           # Extended documentation (stubs)
api/            # API spec / swagger helpers
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
