# Development Guide

This guide covers setting up a development environment and contributing to the OpenUSP platform following Go standard layout principles.

## Prerequisites

### Required Tools
- **Go**: Version 1.21 or later
- **Docker**: Version 20.10 or later  
- **Docker Compose**: Version 2.0 or later
- **Git**: Version 2.30 or later
- **Make**: GNU Make 4.0 or later

### Recommended Tools
- **golangci-lint**: Code linting and quality checks
- **air**: Live reload during development
- **delve**: Go debugger
- **grpcurl**: gRPC testing tool
- **yq**: YAML processing tool

## Project Structure (Go Standard Layout)

The project follows Go standard layout principles:

```
openusp/
├── cmd/                  # Application entry points
│   ├── apiserver/       # REST API server
│   ├── controller/      # USP message processor  
│   ├── cli/            # Command-line interface
│   └── cwmpacs/        # CWMP ACS server
├── internal/            # Private application code
│   ├── apiserver/      # API server implementation
│   ├── controller/     # Controller business logic
│   ├── cli/           # CLI implementation
│   ├── cwmp/          # CWMP protocol handlers
│   ├── db/            # Database access layer
│   ├── mtp/           # Message Transport Protocol
│   └── parser/        # Protocol parsing logic
├── pkg/                # Public libraries
│   ├── config/        # Configuration management
│   └── pb/            # Protocol Buffer definitions
├── configs/            # YAML configuration files
├── deployments/        # Docker Compose manifests
└── docs/              # Documentation
```

## Environment Setup

### 1. Clone Repository
```bash
git clone https://github.com/stalukder-plume/openusp.git
cd openusp
```

### 2. Install Development Dependencies
```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install air for live reload
go install github.com/cosmtrek/air@latest

# Install delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest

# Install yq for YAML processing
go install github.com/mikefarah/yq/v4@latest

# Verify installation
golangci-lint --version
air -v
dlv version
yq --version
```

### 3. Configuration Setup
OpenUSP uses YAML-based configuration. The default configuration files are in `configs/`:

```bash
# Review configuration files
ls configs/
# apiserver.yaml  cli.yaml  controller.yaml  cwmpacs.yaml

# Edit configuration
vim configs/openusp.env
```

Example configuration:
```env
# Database
MONGODB_URI=mongodb://localhost:27017/openusp
REDIS_URI=redis://localhost:6379

# API Server
API_SERVER_PORT=8080
API_SERVER_HOST=localhost

# Controller
CONTROLLER_GRPC_PORT=9090
CONTROLLER_HTTP_PORT=8081

# Logging
LOG_LEVEL=debug
LOG_FORMAT=json

# Security
JWT_SECRET=your-secret-key
TLS_CERT_PATH=/path/to/cert.pem
TLS_KEY_PATH=/path/to/key.pem
```

## Development Workflow

### 1. Start Development Environment
```bash
# Start infrastructure services
docker-compose -f deployments/docker-compose_local.yaml up -d

# Verify services are running
docker-compose ps
```

### 2. Build and Run Services

#### API Server
```bash
cd cmd/apiserver
make build
./apiserver

# Or with live reload
air
```

#### Controller
```bash
cd cmd/controller
make build
./controller

# Or with live reload
air -c .air.controller.toml
```

#### CLI Tool
```bash
cd cmd/cli
make build
./cli --help
```

### 3. Development with Docker
```bash
# Build all services
make docker-build

# Run development stack
docker-compose up -d

# View logs
docker-compose logs -f apiserver
docker-compose logs -f controller
```

## Code Quality

### Linting
```bash
# Run linter on entire codebase
golangci-lint run

# Run linter on specific package
golangci-lint run ./pkg/apiserver/...

# Fix auto-fixable issues
golangci-lint run --fix
```

### Code Formatting
```bash
# Format code
go fmt ./...

# Or use gofumpt for stricter formatting
go install mvdan.cc/gofumpt@latest
gofumpt -w .
```

### Import Organization
```bash
# Install goimports
go install golang.org/x/tools/cmd/goimports@latest

# Organize imports
goimports -w .
```

## Testing

### Unit Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Integration Tests
```bash
# Start test dependencies
docker-compose -f deployments/docker-compose_local.yaml up -d mongodb redis

# Run integration tests
go test -tags=integration ./...

# Run specific test suite
go test -tags=integration ./pkg/apiserver/...
```

### End-to-End Tests
```bash
# Start full environment
docker-compose up -d

# Run E2E tests
cd test
go test -tags=e2e ./...
```

## Debugging

### Local Debugging
```bash
# Debug API server
cd cmd/apiserver
dlv debug

# Debug with arguments
dlv debug -- --config=/path/to/config.yaml

# Attach to running process
dlv attach $(pgrep apiserver)
```

### Remote Debugging
```bash
# Start service with debug server
dlv exec ./apiserver --listen=:2345 --headless=true --api-version=2

# Connect from IDE or command line
dlv connect :2345
```

### Debug in Docker
```dockerfile
# Add to Dockerfile for debug builds
RUN go install github.com/go-delve/delve/cmd/dlv@latest
EXPOSE 2345
CMD ["dlv", "exec", "/app/apiserver", "--listen=:2345", "--headless=true", "--api-version=2"]
```

## Database Development

### Schema Migrations
```bash
# Run database setup
cd pkg/db/cmd
go run main.go --migrate

# Reset database
go run main.go --reset

# Seed test data
go run main.go --seed
```

### Database Queries
```bash
# Connect to MongoDB
mongo openusp

# Sample queries
db.devices.find()
db.agents.find({status: "online"})
db.parameters.createIndex({device_id: 1, path: 1})
```

## API Development

### Generate API Documentation
```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger docs
cd cmd/apiserver
swag init

# View API documentation
open http://localhost:8080/swagger/index.html
```

### Test API Endpoints
```bash
# Using curl
curl -X GET "http://localhost:8080/api/v1/agents" -H "Authorization: Bearer $TOKEN"

# Using HTTPie
http GET localhost:8080/api/v1/agents Authorization:"Bearer $TOKEN"

# Load testing with hey
hey -n 1000 -c 10 http://localhost:8080/api/v1/health
```

## gRPC Development

### Generate Protobuf Code
```bash
cd pkg/pb/cntlrgrpc
./gengo.sh
```

### Test gRPC Services
```bash
# List services
grpcurl -plaintext localhost:9090 list

# Call method
grpcurl -plaintext -d '{"agent_id": "123"}' localhost:9090 cntlr.ControllerService/GetAgent
```

## Performance Profiling

### CPU Profiling
```bash
# Add pprof endpoint to main.go
import _ "net/http/pprof"
go func() {
    http.ListenAndServe(":6060", nil)
}()

# Generate CPU profile
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Analyze profile
(pprof) top
(pprof) web
```

### Memory Profiling
```bash
# Generate memory profile
go tool pprof http://localhost:6060/debug/pprof/heap

# Memory allocation profiling
go test -memprofile mem.prof -bench .
go tool pprof mem.prof
```

## Contributing Guidelines

### Branch Strategy
- `main`: Production-ready code
- `develop`: Integration branch for features
- `feature/*`: Feature development branches
- `hotfix/*`: Critical bug fixes
- `release/*`: Release preparation branches

### Commit Message Format
```
type(scope): description

[optional body]

[optional footer]
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

Example:
```
feat(apiserver): add device filtering endpoint

Add new endpoint to filter devices by status and type.
Includes pagination support and query validation.

Closes #123
```

### Pull Request Process
1. Fork repository and create feature branch
2. Implement changes with tests
3. Run linting and tests locally
4. Update documentation if needed
5. Create PR with clear description
6. Address review feedback
7. Squash and merge after approval

### Code Review Checklist
- [ ] Code follows project style guidelines
- [ ] All tests pass locally
- [ ] Documentation updated for API changes
- [ ] Error handling implemented
- [ ] Logging added for important operations
- [ ] Security considerations addressed
- [ ] Performance impact assessed

## Troubleshooting

### Common Issues

#### Port Conflicts
```bash
# Check port usage
lsof -i :8080
netstat -tulpn | grep 8080

# Kill process using port
kill -9 $(lsof -t -i:8080)
```

#### Database Connection Issues
```bash
# Test MongoDB connection
mongo --eval "db.runCommand('ping')"

# Test Redis connection
redis-cli ping
```

#### Build Issues
```bash
# Clean module cache
go clean -modcache

# Update dependencies
go mod download
go mod tidy

# Rebuild from scratch
make clean
make build
```

### IDE Setup

#### VS Code
Recommended extensions:
- Go (official)
- golangci-lint
- Docker
- REST Client
- GitLens

Settings:
```json
{
    "go.lintTool": "golangci-lint",
    "go.formatTool": "goimports",
    "go.useLanguageServer": true,
    "go.testFlags": ["-v"],
    "files.exclude": {
        "**/.git": true,
        "**/vendor": true
    }
}
```

#### GoLand/IntelliJ
- Install Go plugin
- Configure golangci-lint as external tool
- Set up run configurations for services
- Configure database connections

## Performance Guidelines

### Database Optimization
- Use appropriate indexes
- Implement connection pooling
- Monitor query performance
- Use aggregation pipelines for complex queries

### API Performance
- Implement response caching
- Use pagination for large datasets
- Optimize JSON marshaling
- Monitor response times

### Memory Management
- Avoid memory leaks in goroutines
- Use object pools for frequently allocated objects
- Profile memory usage regularly
- Implement proper cleanup in defer statements

## Security Considerations

### Code Security
- Validate all input data
- Use parameterized queries
- Implement rate limiting
- Sanitize log output (no secrets)

### Authentication
- Use strong JWT secrets
- Implement token refresh
- Support multiple auth providers
- Log authentication events

### Transport Security
- Use TLS for all external communication
- Validate certificates
- Implement certificate rotation
- Use secure random number generation