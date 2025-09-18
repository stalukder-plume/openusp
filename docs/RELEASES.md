# Releases & Versioning (Stub)

Purpose: Define semantic versioning policy, artifact publishing workflow, and upgrade guidance.

## 1. Versioning Policy
- Format: `vMAJOR.MINOR.PATCH` (SemVer)
- Increment rules:
  - MAJOR: Backward incompatible API/behavior changes
  - MINOR: Backward compatible feature additions
  - PATCH: Backward compatible fixes / perf / docs
- Pre-releases: `-rc.N`, `-beta.N` supported (plan usage)

## 2. Tagging & Automation
```bash
make version-check      # show current version
make version-suggest    # patch/minor/major hints
make release VERSION=vX.Y.Z  # builds + tags + pushes images
```

## 3. Docker Image Tagging Strategy
Each component image gets:
- `vX.Y.Z`
- `<git-sha>`
- Optional `latest` (only on stable releases)

## 4. Supported Upgrade Windows (Draft)
| Series | Status | Notes |
|--------|--------|-------|
| v0.x   | Development | Rapid iteration |
| v1.x   | Planned GA  | Contract stabilization |

## 5. Breaking Change Handling
- Document in `CHANGELOG.md`
- Provide migration snippets (DB schema / env var changes)

## 6. Release Checklist (Draft)
1. All CI checks green
2. Security scan clean (critical/high addressed)
3. Docs updated (API + config changes)
4. `make version-suggest` reviewed
5. `make release VERSION=...`
6. Verify images on registry

## 7. Post-Release Validation
| Area | Action |
|------|--------|
| Deploy | Spin ephemeral env using new tag |
| Smoke | Health + simple device command |
| Rollback | Confirm prior tag deployment procedure |

## 8. Future Enhancements
- Signed container images (cosign)
- SBOM publication
- Automated changelog generation

## TODO
- Add actual supported version matrix once GA.
- Introduce CHANGELOG automation (e.g., git-chglog or conventional commits parser).
- Define formal EOL policy per minor series.
