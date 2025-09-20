# OpenUSP Local Development Setup

This guide helps you set up a local development environment using the new YAML-based configuration system and Go standard layout structure.

## Prerequisites

- **Go**: Version 1.21+ installed
- **Docker & Docker Compose**: For infrastructure services
- **Make**: For build automation

## Project Structure

The project follows Go standard layout:
- `cmd/` - Application entry points (apiserver, controller, cli, cwmpacs)  
- `internal/` - Private application code (not importable externally)
- `pkg/` - Public libraries (config, protobuf definitions)
- `configs/` - YAML configuration files for each service

## Quick Start

### 1. Start Infrastructure Services

Start MongoDB, ActiveMQ, and Redis in containers:

```bash
./scripts/start-infrastructure.sh
```

This will:
- Start containerized MongoDB, ActiveMQ, Redis, and Swagger UI
- Use environment variables for service configuration
- Display connection details and status

### 2. Build Applications

Build all OpenUSP applications using the new structure:

```bash
make build-all
./scripts/start-applications.sh
```

This will:
- Load local development environment variables
- Start Controller, API Server, and CWMP ACS as background processes
- Display service endpoints and log locations

### 3. Use the CLI

The CLI can now connect to the local API server:

```bash
# Export environment variables for CLI use
export $(cat configs/openusp-local-dev.env | xargs)

# Use the CLI
./build/bin/openusp-cli --help
./build/bin/openusp-cli agent list
```

### 4. Stop Everything

Stop local applications:
```bash
./scripts/stop-applications.sh
```

Stop infrastructure services:
```bash
docker-compose -f deployments/docker-compose-local-dev.yaml down
```

## Service Details

### Infrastructure Services (Containerized)
- **MongoDB**: `localhost:27017` (admin/admin)
- **ActiveMQ**: 
  - STOMP: `localhost:61613`
  - MQTT: `localhost:1883`
  - Web Console: http://localhost:8161/admin (admin/admin)
- **Redis**: `localhost:6379`
- **Swagger UI**: http://localhost:8080 (OpenUSP API Documentation)

### Application Services (Local)
- **Controller**: gRPC on `localhost:9001`
- **API Server**: HTTP on `localhost:8081`
  - Health Check: http://localhost:8081/health
- **CWMP ACS**: HTTP on `localhost:7547`, HTTPS on `localhost:7548`
- **CLI**: Command-line tool (connects to API Server)

## Configuration

### API Documentation
The OpenUSP API documentation is available via Swagger UI at http://localhost:8080 when infrastructure services are running. This provides:
- Interactive API exploration and testing
- Complete endpoint documentation  
- Request/response schemas
- Authentication details

The API specification file is located at `api/openusp.yaml` and is automatically loaded into Swagger UI.

### Environment File
Local development configuration is in `configs/openusp-local-dev.env`:
- Database connection to localhost MongoDB
- STOMP/MQTT connections to localhost ActiveMQ
- Redis connection to localhost
- Service binding addresses

### Logging
Application logs are written to `logs/` directory:
- `logs/controller.log`
- `logs/apiserver.log`
- `logs/cwmpacs.log`

Process IDs are stored in corresponding `.pid` files.

## Development Workflow

### Building Applications
```bash
make clean
make build-all
```

### Monitoring Services
```bash
# Check running processes
ps aux | grep -E '(openusp-controller|openusp-apiserver|openusp-cwmpacs)'

# Follow logs
tail -f logs/controller.log
tail -f logs/apiserver.log
tail -f logs/cwmpacs.log

# Check API health
curl http://localhost:8081/health
```

### Debugging
Since applications run locally, you can:
1. Attach debuggers directly to processes
2. Add debug prints and rebuild quickly
3. Monitor logs in real-time
4. Use IDE debugging features

### Testing Infrastructure Connection
```bash
# Test MongoDB connection
docker exec -it openusp-mongodb mongosh -u admin -p admin

# Test ActiveMQ (Web Console)
open http://localhost:8161/admin

# Test Redis connection  
docker exec -it openusp-redis redis-cli

# View API Documentation
open http://localhost:8080  # Swagger UI with OpenUSP API docs
```

## Troubleshooting

### Port Conflicts
If you get port binding errors:
```bash
# Check what's using the port
lsof -i :27017  # or other ports
netstat -tulpn | grep :27017
```

### Service Connection Issues
1. Ensure infrastructure services are running:
   ```bash
   docker-compose -f deployments/docker-compose-local-dev.yaml ps
   ```

2. Check environment variables are loaded:
   ```bash
   echo $DB_ADDR  # Should show localhost:27017
   ```

3. Verify network connectivity:
   ```bash
   telnet localhost 27017  # Test MongoDB
   telnet localhost 61613  # Test ActiveMQ STOMP
   ```

### Application Startup Issues
1. Check logs: `tail -f logs/[service].log`
2. Verify binaries exist: `ls -la bin/`
3. Check environment file: `cat configs/openusp-local-dev.env`

### Clean Restart
```bash
# Stop everything
./scripts/stop-applications.sh
docker-compose -f deployments/docker-compose-local-dev.yaml down

# Clean build
make clean
make build-all

# Restart
./scripts/start-infrastructure.sh
./scripts/start-applications.sh
```

## Development Tips

1. **Fast Iteration**: Since apps run locally, you can rebuild and restart quickly
2. **Debugging**: Use your favorite Go debugger (dlv, IDE debuggers)
3. **Log Monitoring**: Use `tail -f logs/*.log` to monitor all services
4. **Configuration Changes**: Edit `configs/openusp-local-dev.env` and restart apps
5. **Infrastructure Persistence**: Container data persists between restarts (MongoDB data, etc.)
6. **API Testing**: Use Swagger UI at http://localhost:8080 to test API endpoints interactively

## Architecture Benefits

This hybrid setup provides:
- **Fast Development**: Local apps start quickly, no container build time
- **Easy Debugging**: Direct process access, IDE integration
- **Isolated Infrastructure**: Consistent database/messaging environment
- **Resource Efficiency**: Only infrastructure services in containers
- **Network Simplicity**: All services accessible via localhost