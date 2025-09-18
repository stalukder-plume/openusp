# GitHub Actions Setup Guide

This document explains how to set up the required secrets and configurations for the OpenUSP GitHub Actions workflows.

## 🔐 Required Secrets

Navigate to your repository's **Settings > Secrets and variables > Actions** and add the following secrets:

### Docker Registry Secrets

| Secret Name | Description | Example |
|-------------|-------------|---------|
| `DOCKER_USERNAME` | Your Docker Hub username | `stalukder-plume` |
| `DOCKER_PASSWORD` | Your Docker Hub access token or password | `dckr_pat_abc123...` |

**How to create Docker Hub access token:**
1. Go to [Docker Hub Account Settings](https://hub.docker.com/settings/security)
2. Click **New Access Token**
3. Give it a name like "GitHub Actions OpenUSP"
4. Select **Read, Write, Delete** permissions
5. Copy the generated token and add it as `DOCKER_PASSWORD` secret

### Optional Secrets (for advanced features)

| Secret Name | Description | When Needed |
|-------------|-------------|-------------|
| `SLACK_WEBHOOK_URL` | Slack webhook for notifications | If using Slack notifications |
| `TELEGRAM_BOT_TOKEN` | Telegram bot token | If using Telegram notifications |
| `TELEGRAM_CHAT_ID` | Telegram chat ID | If using Telegram notifications |

## 🚀 Workflow Overview

### 1. **CI/CD Pipeline** (`.github/workflows/ci.yml`)

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches  
- Release published

**What it does:**
- ✅ Runs tests and linting
- 🔨 Builds binaries for multiple platforms
- 🐳 Builds and tests Docker images
- 📤 Pushes Docker images to registry (on main/develop/releases)
- 📦 Uploads release assets (on GitHub releases)

### 2. **Security Scan** (`.github/workflows/security.yml`)

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches
- Weekly schedule (Mondays at 6 AM UTC)

**What it does:**
- 🔍 Runs Gosec security scanner
- 🔍 Scans for dependency vulnerabilities  
- 🔍 Scans Docker images with Trivy
- 📊 Uploads results to GitHub Security tab

### 3. **Release Management** (`.github/workflows/release.yml`)

**Triggers:**
- Manual workflow dispatch only

**What it does:**
- ✅ Validates semantic version format
- 🏷️ Creates and pushes git tag
- 🔨 Builds and tests release
- 🐳 Builds and pushes Docker images with release tags
- 📋 Creates GitHub release with assets and changelog

### 4. **Dependency Updates** (`.github/workflows/dependencies.yml`)

**Triggers:**
- Weekly schedule (Mondays at 2 AM UTC)
- Manual workflow dispatch

**What it does:**
- 📦 Updates Go dependencies
- 🔄 Updates GitHub Actions versions
- 📝 Creates PRs for review

## 📋 How to Use

### Normal Development Workflow

1. **Create feature branch**
   ```bash
   git checkout -b feature/my-feature
   ```

2. **Make changes and push**
   ```bash
   git add .
   git commit -m "feat: add new feature"
   git push origin feature/my-feature
   ```

3. **Create Pull Request**
   - CI will run tests, linting, and security scans
   - Docker images will be built and tested
   - Must pass all checks before merging

4. **Merge to main**
   - Triggers full CI/CD pipeline
   - Builds and pushes Docker images with `dev-YYYYMMDD` tags

### Creating Releases

#### Option 1: Using GitHub UI (Recommended)

1. Go to **Actions** tab in your repository
2. Click **Release Management** workflow
3. Click **Run workflow**
4. Enter version (e.g., `v1.2.3`)
5. Select release type and options
6. Click **Run workflow**

The workflow will:
- Create git tag
- Build and test everything
- Push Docker images with version tags
- Create GitHub release with binaries

#### Option 2: Using Command Line

```bash
# Create and push tag
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3

# Then create GitHub release manually or use gh CLI
gh release create v1.2.3 --title "OpenUSP v1.2.3" --generate-notes
```

### Docker Images

After successful workflows, your images will be available:

```bash
# Latest development build
docker pull stalukder-plume/openusp-controller:dev-20240918

# Specific release
docker pull stalukder-plume/openusp-controller:v1.2.3

# Latest stable (only for stable releases)
docker pull stalukder-plume/openusp-controller:latest
```

## 🐛 Troubleshooting

### Common Issues

#### 1. **Docker push fails with "unauthorized"**
- Check `DOCKER_USERNAME` and `DOCKER_PASSWORD` secrets
- Ensure Docker Hub access token has correct permissions
- Verify repository name in `DOCKER_REPO_PREFIX`

#### 2. **Tests fail in CI but pass locally**
- Check if tests depend on local environment
- Ensure all dependencies are in `go.mod`
- Check for race conditions (CI runs with `-race` flag)

#### 3. **Release workflow fails**
- Ensure version follows semantic versioning (e.g., `v1.2.3`)
- Check if tag already exists
- Verify Docker secrets are set up

#### 4. **Security scan failures**
- Check GitHub Security tab for details
- Update vulnerable dependencies
- Scan reports are informational by default

### Viewing Workflow Results

1. **Actions Tab**: See all workflow runs and their status
2. **Security Tab**: View security scan results and advisories  
3. **Releases**: See all releases with binaries and Docker images
4. **Packages**: View published Docker images (if using GitHub Container Registry)

## 🔧 Customization

### Changing Docker Registry

Edit these files to use a different registry:
- `.github/workflows/ci.yml`: Update `REGISTRY` and `DOCKER_REPO_PREFIX`
- `scripts/docker-release.sh`: Update `REPO_PREFIX`
- `Makefile`: Update `DOCKER_REGISTRY`

### Adding Notifications

Add notification steps to workflows:

```yaml
- name: Notify Slack
  if: failure()
  uses: 8398a7/action-slack@v3
  with:
    status: failure
    webhook_url: ${{ secrets.SLACK_WEBHOOK_URL }}
```

### Changing Build Targets

Edit `build` job matrix in `.github/workflows/ci.yml`:
```yaml
strategy:
  matrix:
    goos: [linux, darwin, windows]
    goarch: [amd64, arm64]
```

## 📚 Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Hub Access Tokens](https://docs.docker.com/docker-hub/access-tokens/)
- [Semantic Versioning](https://semver.org/)
- [Go Module Reference](https://golang.org/ref/mod)