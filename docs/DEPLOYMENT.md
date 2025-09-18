# Deployment Guide

This guide covers deployment strategies, environment setup, and operational procedures for the OpenUSP platform.

## Deployment Strategies

### Local Development
- Docker Compose for local testing
- Hot-reload development environment
- Local database instances
- Mock external services

### Staging Environment
- Kubernetes cluster deployment
- Production-like configuration
- Integration testing automation
- Performance testing

### Production Deployment
- Multi-region Kubernetes
- High availability setup
- Auto-scaling configuration
- Disaster recovery planning

## Prerequisites

### Infrastructure Requirements
- **Kubernetes**: Version 1.24 or later
- **Docker**: Version 20.10 or later
- **Helm**: Version 3.8 or later
- **Ingress Controller**: nginx-ingress or Traefik
- **Certificate Manager**: cert-manager for TLS

### Resource Requirements

#### Minimum (Development)
- **CPU**: 2 cores total
- **Memory**: 4GB RAM total
- **Storage**: 20GB persistent storage
- **Network**: 1Gbps bandwidth

#### Recommended (Production)
- **CPU**: 8 cores per service
- **Memory**: 16GB RAM per service
- **Storage**: 100GB SSD persistent storage
- **Network**: 10Gbps bandwidth with redundancy

## Docker Deployment

### Local Development
```bash
# Start all services
docker-compose -f deployments/docker-compose_local.yaml up -d

# Scale specific service
docker-compose up -d --scale apiserver=3

# View service logs
docker-compose logs -f apiserver

# Stop all services
docker-compose down
```

### Production Docker
```bash
# Build production images
make docker-build-prod

# Deploy with production compose
docker-compose -f deployments/docker-compose.yaml up -d

# Update single service
docker-compose up -d --no-deps apiserver
```

### Docker Configuration
```yaml
# docker-compose.yaml
version: '3.8'

services:
  apiserver:
    image: openusp/apiserver:latest
    ports:
      - "8080:8080"
    environment:
      - MONGODB_URI=mongodb://mongodb:27017/openusp
      - REDIS_URI=redis://redis:6379
    depends_on:
      - mongodb
      - redis
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  controller:
    image: openusp/controller:latest
    ports:
      - "9090:9090"
    environment:
      - MONGODB_URI=mongodb://mongodb:27017/openusp
      - GRPC_PORT=9090
    depends_on:
      - mongodb
    restart: unless-stopped

  mongodb:
    image: mongo:6.0
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db
    environment:
      - MONGO_INITDB_ROOT_USERNAME=admin
      - MONGO_INITDB_ROOT_PASSWORD=password

  redis:
    image: redis:7.0-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  mongodb_data:
  redis_data:
```

## Kubernetes Deployment

### Helm Chart Installation
```bash
# Add Helm repository
helm repo add openusp https://charts.openusp.io
helm repo update

# Install with default values
helm install openusp openusp/openusp

# Install with custom values
helm install openusp openusp/openusp -f values-production.yaml

# Upgrade deployment
helm upgrade openusp openusp/openusp --reuse-values
```

### Manual Kubernetes Deployment
```yaml
# namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: openusp
---
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: openusp-config
  namespace: openusp
data:
  MONGODB_URI: "mongodb://mongodb:27017/openusp"
  REDIS_URI: "redis://redis:6379"
  LOG_LEVEL: "info"
---
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: apiserver
  namespace: openusp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: apiserver
  template:
    metadata:
      labels:
        app: apiserver
    spec:
      containers:
      - name: apiserver
        image: openusp/apiserver:latest
        ports:
        - containerPort: 8080
        envFrom:
        - configMapRef:
            name: openusp-config
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
# service.yaml
apiVersion: v1
kind: Service
metadata:
  name: apiserver-service
  namespace: openusp
spec:
  selector:
    app: apiserver
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
---
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: openusp-ingress
  namespace: openusp
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
spec:
  tls:
  - hosts:
    - api.openusp.com
    secretName: openusp-tls
  rules:
  - host: api.openusp.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: apiserver-service
            port:
              number: 80
```

### Apply Kubernetes Manifests
```bash
# Apply all manifests
kubectl apply -f k8s/

# Check deployment status
kubectl get pods -n openusp
kubectl get services -n openusp
kubectl get ingress -n openusp

# View logs
kubectl logs -f deployment/apiserver -n openusp

# Scale deployment
kubectl scale deployment apiserver --replicas=5 -n openusp
```

## Environment Configuration

### Production Values (values-production.yaml)
```yaml
# Helm values for production
global:
  imageRegistry: "your-registry.com"
  imagePullSecrets:
    - name: regcred

apiserver:
  replicaCount: 3
  image:
    repository: openusp/apiserver
    tag: "v1.0.0"
  resources:
    requests:
      memory: "512Mi"
      cpu: "500m"
    limits:
      memory: "1Gi"
      cpu: "1000m"
  autoscaling:
    enabled: true
    minReplicas: 3
    maxReplicas: 10
    targetCPUUtilizationPercentage: 70

controller:
  replicaCount: 2
  image:
    repository: openusp/controller
    tag: "v1.0.0"
  resources:
    requests:
      memory: "256Mi"
      cpu: "250m"
    limits:
      memory: "512Mi"
      cpu: "500m"

mongodb:
  enabled: true
  replicaCount: 3
  persistence:
    enabled: true
    size: "100Gi"
    storageClass: "fast-ssd"

redis:
  enabled: true
  cluster:
    enabled: true
    slaveCount: 2
  persistence:
    enabled: true
    size: "10Gi"

ingress:
  enabled: true
  className: "nginx"
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/rate-limit: "100"
  hosts:
    - host: api.openusp.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: openusp-tls
      hosts:
        - api.openusp.com

monitoring:
  enabled: true
  prometheus:
    enabled: true
  grafana:
    enabled: true
```

### Environment Variables
```bash
# Required Environment Variables
export MONGODB_URI="mongodb://user:pass@host:27017/openusp"
export REDIS_URI="redis://host:6379/0"
export JWT_SECRET="your-jwt-secret-key"

# Optional Environment Variables
export LOG_LEVEL="info"                    # debug, info, warn, error
export LOG_FORMAT="json"                   # json, text
export API_SERVER_PORT="8080"
export CONTROLLER_GRPC_PORT="9090"
export METRICS_PORT="9091"

# TLS Configuration
export TLS_ENABLED="true"
export TLS_CERT_PATH="/certs/tls.crt"
export TLS_KEY_PATH="/certs/tls.key"

# Database Configuration
export DB_MAX_CONNECTIONS="100"
export DB_CONNECTION_TIMEOUT="30s"
export REDIS_MAX_CONNECTIONS="50"

# Performance Tuning
export GOMAXPROCS="4"
export GOGC="100"
```

## Database Setup

### MongoDB Setup
```bash
# Create MongoDB user
mongo admin --eval "
  db.createUser({
    user: 'openusp',
    pwd: 'secure-password',
    roles: [
      {role: 'readWrite', db: 'openusp'},
      {role: 'dbAdmin', db: 'openusp'}
    ]
  })
"

# Create indexes
mongo openusp --eval "
  db.devices.createIndex({'endpoint_id': 1}, {unique: true})
  db.agents.createIndex({'agent_id': 1}, {unique: true})
  db.parameters.createIndex({'device_id': 1, 'path': 1})
  db.events.createIndex({'timestamp': 1}, {expireAfterSeconds: 2592000})
"

# Enable replica set for production
mongod --replSet rs0 --bind_ip localhost,<hostname(s)>
mongo --eval "rs.initiate()"
```

### Redis Setup
```bash
# Redis configuration for production
cat > redis.conf << EOF
bind 0.0.0.0
port 6379
requirepass your-redis-password
maxmemory 2gb
maxmemory-policy allkeys-lru
save 900 1
save 300 10
save 60 10000
EOF

# Start Redis with configuration
redis-server redis.conf

# Test Redis connection
redis-cli -a your-redis-password ping
```

## Load Balancing

### Nginx Configuration
```nginx
upstream apiserver {
    server apiserver-1:8080 max_fails=3 fail_timeout=30s;
    server apiserver-2:8080 max_fails=3 fail_timeout=30s;
    server apiserver-3:8080 max_fails=3 fail_timeout=30s;
}

server {
    listen 80;
    server_name api.openusp.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.openusp.com;

    ssl_certificate /etc/certs/tls.crt;
    ssl_certificate_key /etc/certs/tls.key;

    location / {
        proxy_pass http://apiserver;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Health check
        proxy_next_upstream error timeout invalid_header http_500 http_502 http_503;
        
        # Timeouts
        proxy_connect_timeout 5s;
        proxy_send_timeout 30s;
        proxy_read_timeout 30s;
    }

    location /health {
        access_log off;
        return 200 "healthy\n";
        add_header Content-Type text/plain;
    }
}
```

## SSL/TLS Configuration

### Certificate Management with cert-manager
```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: admin@openusp.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: nginx
```

### Manual Certificate Setup
```bash
# Generate self-signed certificate for development
openssl req -x509 -newkey rsa:4096 -keyout tls.key -out tls.crt -days 365 -nodes \
  -subj "/C=US/ST=State/L=City/O=Org/CN=api.openusp.com"

# Create Kubernetes secret
kubectl create secret tls openusp-tls --cert=tls.crt --key=tls.key -n openusp
```

## Monitoring Setup

### Prometheus Configuration
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: prometheus-config
data:
  prometheus.yml: |
    global:
      scrape_interval: 15s
    scrape_configs:
    - job_name: 'apiserver'
      static_configs:
      - targets: ['apiserver:9091']
    - job_name: 'controller'
      static_configs:
      - targets: ['controller:9091']
    - job_name: 'mongodb'
      static_configs:
      - targets: ['mongodb-exporter:9216']
```

### Grafana Dashboard
```bash
# Import OpenUSP dashboard
curl -X POST \
  http://admin:password@grafana:3000/api/dashboards/db \
  -H 'Content-Type: application/json' \
  -d @grafana-dashboard.json
```

## Backup and Recovery

### Database Backup
```bash
# MongoDB backup
mongodump --uri="mongodb://user:pass@host:27017/openusp" --out=/backup/$(date +%Y%m%d)

# Automated backup script
#!/bin/bash
BACKUP_DIR="/backup/mongodb"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR

mongodump --uri="$MONGODB_URI" --out="$BACKUP_DIR/$DATE"
tar -czf "$BACKUP_DIR/$DATE.tar.gz" -C "$BACKUP_DIR" "$DATE"
rm -rf "$BACKUP_DIR/$DATE"

# Keep only last 7 days
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete
```

### Disaster Recovery
```bash
# Restore MongoDB from backup
mongorestore --uri="mongodb://user:pass@host:27017/openusp" /backup/20231201/openusp/

# Kubernetes disaster recovery
# 1. Backup etcd
etcdctl snapshot save /backup/etcd-snapshot.db

# 2. Backup persistent volumes
kubectl get pv -o yaml > pv-backup.yaml

# 3. Export configurations
kubectl get all -n openusp -o yaml > openusp-backup.yaml
```

## Performance Tuning

### Go Application Tuning
```bash
# Set Go runtime variables
export GOMAXPROCS=4
export GOGC=100
export GODEBUG=gctrace=1

# Memory ballast (for stable GC)
export BALLAST_SIZE=1GB
```

### Database Performance
```javascript
// MongoDB indexes for performance
db.devices.createIndex({"endpoint_id": 1}, {unique: true})
db.devices.createIndex({"last_seen": -1})
db.agents.createIndex({"status": 1, "last_heartbeat": -1})
db.parameters.createIndex({"device_id": 1, "path": 1})
db.events.createIndex({"timestamp": -1, "device_id": 1})

// Redis performance settings
CONFIG SET maxmemory-policy allkeys-lru
CONFIG SET timeout 300
CONFIG SET tcp-keepalive 300
```

## Security Configuration

### Network Security
```bash
# Firewall rules
ufw allow 22/tcp    # SSH
ufw allow 80/tcp    # HTTP
ufw allow 443/tcp   # HTTPS
ufw deny 27017/tcp  # MongoDB (internal only)
ufw deny 6379/tcp   # Redis (internal only)
ufw enable
```

### Kubernetes Security
```yaml
apiVersion: v1
kind: NetworkPolicy
metadata:
  name: openusp-network-policy
  namespace: openusp
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
    ports:
    - protocol: TCP
      port: 8080
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: mongodb
    ports:
    - protocol: TCP
      port: 27017
```

## Troubleshooting Deployment

### Common Issues
```bash
# Check pod status
kubectl describe pod apiserver-xxx -n openusp

# View container logs
kubectl logs -f pod/apiserver-xxx -c apiserver -n openusp

# Debug networking
kubectl exec -it apiserver-xxx -n openusp -- nslookup mongodb

# Check resource usage
kubectl top pods -n openusp
kubectl top nodes
```

### Health Checks
```bash
# API server health
curl -f http://localhost:8080/health

# Database connectivity
mongo --eval "db.runCommand('ping')"
redis-cli ping

# Service discovery
kubectl get endpoints -n openusp
```

This deployment guide provides comprehensive instructions for deploying OpenUSP in various environments, from local development to production Kubernetes clusters.