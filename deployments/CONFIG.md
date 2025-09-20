# Docker Compose Configuration Guide

## Overview

The `docker-compose-local-dev.yaml` file has been refactored to be highly configurable using YAML anchors, environment variables, and Docker Compose features.

## Key Improvements

### 1. YAML Anchors for Reusability
- **`&common-configs`**: Shared configuration across all services
- **`&healthcheck-defaults`**: Standardized health check settings
- **`&environment`**: Centralized environment variable definitions

### 2. Environment Variable Support
All configuration values can be overridden via environment variables or `.env` file:

```bash
# Using environment variables
MONGO_PORT=27018 docker-compose -f docker-compose-local-dev.yaml up

# Using .env file (recommended)
cp .env.example .env
# Edit .env file with your values
docker-compose -f docker-compose-local-dev.yaml up
```

### 3. Enhanced Features
- **Health checks**: All services have proper health checks
- **Logging**: Structured logging with rotation
- **Networks**: Dedicated network with configurable subnet
- **Volumes**: Named volumes for better data management
- **Resource limits**: Configurable memory limits for services

## Configuration Sections

### Database (MongoDB)
```yaml
MONGO_VERSION=latest          # MongoDB version
MONGO_HOST=openusp-db        # Container hostname
MONGO_PORT=27017             # Internal port
MONGO_HOST_PORT=27017        # Host port mapping
MONGO_DATABASE=openusp       # Database name
MONGO_ROOT_USERNAME=admin    # Root username
MONGO_ROOT_PASSWORD=admin    # Root password
```

### Message Broker (ActiveMQ)
```yaml
ACTIVEMQ_VERSION=mariadb_max_packet  # ActiveMQ version
ACTIVEMQ_HOST=openusp-broker         # Container hostname
ACTIVEMQ_ADMIN_USER=admin            # Admin username
ACTIVEMQ_ADMIN_PASSWORD=admin        # Admin password
ACTIVEMQ_MIN_MEMORY=512              # Min memory (MB)
ACTIVEMQ_MAX_MEMORY=2048             # Max memory (MB)

# Port configurations
ACTIVEMQ_STOMP_PORT=61613            # STOMP protocol
ACTIVEMQ_OPENWIRE_PORT=61616         # OpenWire protocol
ACTIVEMQ_WEB_PORT=8161               # Web console
ACTIVEMQ_MQTT_PORT=1883              # MQTT protocol
```

### Cache (Redis)
```yaml
REDIS_VERSION=latest         # Redis version
REDIS_HOST=openusp-cache    # Container hostname
REDIS_PORT=6379             # Internal port
REDIS_HOST_PORT=6379        # Host port mapping
REDIS_PASSWORD=             # Optional password
REDIS_MAX_MEMORY=256mb      # Memory limit
```

### Network
```yaml
NETWORK_SUBNET=172.20.0.0/16  # Docker network subnet
```

## Usage Examples

### Basic Development Setup
```bash
# Use defaults
docker-compose -f docker-compose-local-dev.yaml up -d
```

### Custom Port Configuration
```bash
# Create .env file
cat > .env << EOF
MONGO_HOST_PORT=27018
REDIS_HOST_PORT=6380
ACTIVEMQ_WEB_HOST_PORT=8162
EOF

docker-compose -f docker-compose-local-dev.yaml up -d
```

### Production-like Setup
```bash
# Create .env file with stronger configuration
cat > .env << EOF
MONGO_ROOT_PASSWORD=secure_password_123
REDIS_PASSWORD=redis_secure_pass
ACTIVEMQ_ADMIN_PASSWORD=activemq_secure_pass
ACTIVEMQ_MAX_MEMORY=4096
REDIS_MAX_MEMORY=512mb
EOF

docker-compose -f docker-compose-local-dev.yaml up -d
```

## Service Health Monitoring

All services include health checks. Check status:

```bash
# View service health
docker-compose -f docker-compose-local-dev.yaml ps

# Check specific service logs
docker-compose -f docker-compose-local-dev.yaml logs openusp-db
docker-compose -f docker-compose-local-dev.yaml logs openusp-broker
docker-compose -f docker-compose-local-dev.yaml logs openusp-cache
```

## Data Persistence

Named volumes ensure data persistence:
- `openusp_mongodb_data`: MongoDB data
- `openusp_mongodb_config`: MongoDB configuration
- `openusp_activemq_data`: ActiveMQ data
- `openusp_activemq_conf`: ActiveMQ configuration
- `openusp_redis_data`: Redis data
- `openusp_redis_conf`: Redis configuration

## Network Configuration

Services communicate via the `openusp-network` bridge network with configurable subnet.

## Best Practices

1. **Always use `.env` file** for local development
2. **Keep `.env` out of version control** (add to `.gitignore`)
3. **Use `.env.example`** as template for team members
4. **Monitor service health** using `docker-compose ps`
5. **Use named volumes** for data that should persist
6. **Configure logging** to prevent disk space issues

## Troubleshooting

### Service Won't Start
```bash
# Check service logs
docker-compose -f docker-compose-local-dev.yaml logs [service-name]

# Check health status
docker-compose -f docker-compose-local-dev.yaml ps
```

### Port Conflicts
```bash
# Change ports in .env file
echo "MONGO_HOST_PORT=27018" >> .env
echo "REDIS_HOST_PORT=6380" >> .env

# Restart services
docker-compose -f docker-compose-local-dev.yaml down
docker-compose -f docker-compose-local-dev.yaml up -d
```

### Reset All Data
```bash
# Stop services and remove volumes
docker-compose -f docker-compose-local-dev.yaml down -v

# Start fresh
docker-compose -f docker-compose-local-dev.yaml up -d
```