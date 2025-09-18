# CWMP ACS Docker Configuration

This directory contains Docker configuration for the OpenUSP CWMP ACS (Auto Configuration Server).

## Files

- `Dockerfile` - Standard production Dockerfile for CWMP ACS
- `Dockerfile.enhanced` - Enhanced version with entrypoint script and more flexibility
- `entrypoint.sh` - Configurable entrypoint script for the CWMP ACS
- `README.md` - This documentation file

## CWMP ACS Overview

The CWMP ACS is a TR-069 Auto Configuration Server that manages CPE (Customer Premises Equipment) devices using the TR-069 protocol. It provides:

- Device auto-configuration and management
- Firmware upgrade management  
- Configuration parameter management
- Device monitoring and diagnostics
- Bulk configuration operations

## Standard Usage

### Building
```bash
# Build using standard Dockerfile
docker build -t openusp-cwmpacs -f build/cwmpacs/Dockerfile .

# Build using enhanced Dockerfile
docker build -t openusp-cwmpacs:enhanced -f build/cwmpacs/Dockerfile.enhanced .
```

### Running
```bash
# Basic usage
docker run -p 7547:7547 -p 7548:7548 openusp-cwmpacs

# With custom ports
docker run -p 8547:7547 -p 8548:7548 openusp-cwmpacs

# With configuration volume
docker run -v ./configs:/opt/openusp/configs -p 7547:7547 -p 7548:7548 openusp-cwmpacs

# With persistent data
docker run -v cwmpacs-data:/var/lib/openusp -p 7547:7547 -p 7548:7548 openusp-cwmpacs
```

## Enhanced Version Usage

The enhanced version includes a configurable entrypoint script:

```bash
# With custom configuration
docker run -p 7547:7547 -p 7548:7548 openusp-cwmpacs:enhanced --log-level debug

# With custom ports via arguments
docker run -p 8547:8547 -p 8548:8548 openusp-cwmpacs:enhanced --http-port 8547 --https-port 8548

# With configuration file
docker run -v ./cwmpacs.yaml:/opt/openusp/configs/cwmpacs.yaml \
           -p 7547:7547 -p 7548:7548 \
           openusp-cwmpacs:enhanced --config-file /opt/openusp/configs/cwmpacs.yaml
```

## Environment Variables

- `CWMP_ACS_HTTP_PORT` - HTTP port (default: 7547)
- `CWMP_ACS_HTTPS_PORT` - HTTPS port (default: 7548) 
- `CWMP_ACS_LOG_LEVEL` - Log level: debug, info, warn, error (default: info)
- `CWMP_ACS_CONFIG_FILE` - Configuration file path
- `CWMP_ACS_DB_URL` - Database URL override

## Docker Compose

The CWMP ACS is included in the local development docker-compose configuration:

```bash
cd deployments
docker-compose -f docker-compose_local.yaml up openusp-cwmpacs
```

## Port Information

- **7547/TCP** - Standard TR-069 ACS HTTP port
- **7548/TCP** - Standard TR-069 ACS HTTPS port (if SSL/TLS enabled)

These are the standard ports defined in the TR-069 specification for CWMP communication between ACS and CPE devices.

## Directory Structure

```
/opt/openusp/
├── bin/           # CWMP ACS binary
├── configs/       # Configuration files (mountable)
├── scripts/       # Helper scripts (entrypoint.sh)
/var/lib/openusp/  # Application data (persistent)
/var/log/openusp/  # Log files (mountable)
```

## Security Notes

- Runs as non-root user (`openusp:openusp`)
- Minimal runtime dependencies (debian:stable + ca-certificates)
- Follows Linux FHS conventions
- Supports configuration and data volume mounting