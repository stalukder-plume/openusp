# Configuration

OpenUSP uses YAML-based configuration with environment variable substitution. Each service has its own configuration file with shared configuration patterns.

## 1. Configuration Files

| Service | Configuration File | Purpose |
|---------|-------------------|---------|
| API Server | `configs/apiserver.yaml` | REST API, authentication, database connections |
| Controller | `configs/controller.yaml` | USP processing, MTP settings, device management |
| CLI | `configs/cli.yaml` | Command-line tool settings and defaults |
| CWMP ACS | `configs/cwmpacs.yaml` | TR-069/CWMP protocol settings |

## 2. Configuration Structure

### Common Configuration Patterns
All services share common configuration sections:

```yaml
# Database configuration
database:
  uri: ${OPENUSP_MONGO_URI:mongodb://localhost:27017}
  name: ${OPENUSP_DB_NAME:openusp}
  timeout: ${OPENUSP_DB_TIMEOUT:10s}

# Redis configuration
redis:
  addr: ${OPENUSP_REDIS_ADDR:localhost:6379}
  password: ${OPENUSP_REDIS_PASSWORD:}
  db: ${OPENUSP_REDIS_DB:0}

# Logging configuration
log:
  level: ${OPENUSP_LOG_LEVEL:info}
  format: ${OPENUSP_LOG_FORMAT:json}
  output: ${OPENUSP_LOG_OUTPUT:stdout}

# Server configuration
server:
  host: ${OPENUSP_HOST:0.0.0.0}
  port: ${OPENUSP_PORT:8080}
  read_timeout: ${OPENUSP_READ_TIMEOUT:30s}
  write_timeout: ${OPENUSP_WRITE_TIMEOUT:30s}
```

### Environment Variable Substitution
OpenUSP supports environment variable substitution with default values using the syntax:
- `${VAR_NAME}` - Required environment variable
- `${VAR_NAME:default_value}` - Optional with default value

## 3. Service-Specific Configuration

### API Server (`configs/apiserver.yaml`)
```yaml
server:
  port: ${OPENUSP_API_PORT:8081}

auth:
  method: ${OPENUSP_AUTH_METHOD:jwt}
  jwt_secret: ${OPENUSP_JWT_SECRET:your-secret-key}
  
swagger:
  enabled: ${OPENUSP_SWAGGER_ENABLED:true}
  host: ${OPENUSP_SWAGGER_HOST:localhost:8081}
```

### Controller (`configs/controller.yaml`)
```yaml
server:
  port: ${OPENUSP_CONTROLLER_PORT:8082}

usp:
  endpoint_id: ${OPENUSP_USP_ENDPOINT_ID:controller.openusp.org}
  
mtp:
  stomp:
    broker_url: ${OPENUSP_STOMP_URL:stomp://localhost:61613}
  mqtt:
    broker_url: ${OPENUSP_MQTT_URL:mqtt://localhost:1883}
```

### CWMP ACS (`configs/cwmpacs.yaml`)
```yaml
server:
  port: ${OPENUSP_CWMP_PORT:7547}

cwmp:
  acs_url: ${OPENUSP_ACS_URL:http://localhost:7547}
  connection_request_auth: ${OPENUSP_CWMP_AUTH:digest}
```

## 4. Environment Variables Reference

| Variable | Default | Purpose | Used By |
|----------|---------|---------|---------|
| OPENUSP_MONGO_URI | mongodb://localhost:27017 | MongoDB connection string | All services |
| OPENUSP_DB_NAME | openusp | Database name | All services |
| OPENUSP_REDIS_ADDR | localhost:6379 | Redis server address | All services |
| OPENUSP_REDIS_PASSWORD | (empty) | Redis password | All services |
| OPENUSP_LOG_LEVEL | info | Log level (debug,info,warn,error) | All services |
| OPENUSP_LOG_FORMAT | json | Log format (json,text) | All services |
| OPENUSP_API_PORT | 8081 | API server HTTP port | apiserver |
| OPENUSP_CONTROLLER_PORT | 8082 | Controller HTTP port | controller |
| OPENUSP_CWMP_PORT | 7547 | CWMP ACS port | cwmpacs |
| OPENUSP_AUTH_METHOD | jwt | Authentication method | apiserver |
| OPENUSP_JWT_SECRET | your-secret-key | JWT signing secret | apiserver |
| OPENUSP_SWAGGER_ENABLED | true | Enable Swagger UI | apiserver |
| OPENUSP_USP_ENDPOINT_ID | controller.openusp.org | USP endpoint identifier | controller |
| OPENUSP_STOMP_URL | stomp://localhost:61613 | STOMP broker URL | controller |
| OPENUSP_MQTT_URL | mqtt://localhost:1883 | MQTT broker URL | controller |
| OPENUSP_ACS_URL | http://localhost:7547 | CWMP ACS URL | cwmpacs |

## 5. Configuration Loading

Configuration is loaded by the `pkg/config` package using this precedence (highest to lowest):

1. **Environment Variables**: Direct environment variable values
2. **YAML Files**: Values from service-specific configuration files  
3. **Default Values**: Defaults specified in environment variable substitution syntax

### Loading Process
```go
// Example usage in service initialization
config, err := config.LoadConfig("configs/apiserver.yaml")
if err != nil {
    log.Fatal("Failed to load configuration:", err)
}
```

## 6. Production Configuration

### Security Considerations
- Use external secret management (Kubernetes Secrets, HashiCorp Vault)
- Never commit secrets to version control
- Use TLS for all inter-service communication in production
- Rotate JWT secrets regularly

### Performance Tuning
| Component | Parameter | Recommendation | Purpose |
|-----------|-----------|----------------|---------|
| MongoDB | Connection Pool | 10-100 connections | Optimize for concurrent load |
| Redis | Max Memory | 2GB+ | Cache performance |
| Go Runtime | GOGC | 100-200% | Memory vs CPU trade-off |
| HTTP Server | Timeouts | 30s read, 30s write | Prevent resource exhaustion |

### Example Production Configuration
```yaml
# Production apiserver.yaml
database:
  uri: ${OPENUSP_MONGO_URI:mongodb://mongo-cluster:27017/openusp?replicaSet=rs0}
  timeout: 10s
  
redis:
  addr: ${OPENUSP_REDIS_ADDR:redis-cluster:6379}
  password: ${OPENUSP_REDIS_PASSWORD}
  
server:
  host: 0.0.0.0
  port: ${OPENUSP_API_PORT:8081}
  read_timeout: 30s
  write_timeout: 30s
  
auth:
  method: jwt
  jwt_secret: ${OPENUSP_JWT_SECRET}
  
log:
  level: info
  format: json
  output: stdout
```

## 7. Migration from Environment Variables

The configuration system has been migrated from pure environment variables to YAML with environment variable substitution. This provides:

- **Better structure**: Hierarchical configuration organization
- **Type safety**: Proper typing for numbers, booleans, durations
- **Documentation**: Self-documenting configuration files
- **Flexibility**: Easy to override specific values via environment variables
- **Validation**: Built-in configuration validation

See `docs/YAML_CONFIGURATION.md` for detailed information about the YAML configuration system.
