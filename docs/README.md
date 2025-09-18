# OpenUSP Documentation Index

This folder contains detailed documentation referenced from the root README. Each document is intentionally scoped and should stay concise and task‑focused.

## Index

| Document | Purpose |
|----------|---------|
| [ARCHITECTURE.md](ARCHITECTURE.md) | High-level system design, component boundaries, data flows |
| [COMPONENTS.md](COMPONENTS.md) | Detailed description of each service/component and its responsibilities |
| [DEVELOPMENT.md](DEVELOPMENT.md) | Local dev environment setup, build, test, lint, debug workflows |
| [DEPLOYMENT.md](DEPLOYMENT.md) | Deployment options (Docker Compose, Kubernetes, production hardening) |
| [CONFIGURATION.md](CONFIGURATION.md) | Environment variables, configuration files, secrets, TLS |
| [SECURITY.md](SECURITY.md) | Security model, authn/z, certificates, vulnerability scanning |
| [OPERATIONS.md](OPERATIONS.md) | Monitoring, metrics, logging, tracing, troubleshooting |
| [RELEASES.md](RELEASES.md) | Semantic versioning, tagging, Docker image strategy, upgrade notes |
| [PROTOCOLS.md](PROTOCOLS.md) | USP (TR-369) & CWMP (TR-069) protocol behavior, mapping, extensions |
| [API.md](API.md) | REST/CLI usage pointers and linkouts to generated references |

## Contribution Notes
- Keep root README brief—deep explanations belong here.
- Prefer linking between docs instead of duplicating content.
- Add diagrams as `.drawio` or `.svg` alongside the doc that references them.

## TODO Placeholders
Files are created as stubs to be iteratively filled. Open an issue before introducing major structural changes.