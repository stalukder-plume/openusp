# Database Schema

This document describes the database schema and data models used by the OpenUSP platform.

## Database Technology

OpenUSP uses **MongoDB** as the primary database for storing device information, parameters, events, and operational data. MongoDB was chosen for its:
- Flexible schema design
- Horizontal scalability
- JSON-like document storage
- Rich query capabilities
- Built-in replication and sharding

## Collections Overview

```
openusp/
├── agents              # USP/CWMP device agents
├── parameters          # Device parameter values
├── instances          # Object instances
├── operations          # Pending/completed operations
├── events             # Device events and notifications
├── subscriptions      # Event subscriptions
├── datamodels         # Supported data model definitions
├── sessions           # Active device sessions
├── users              # User accounts and authentication
└── audit_logs         # System audit trail
```

## Core Collections

### agents
Stores information about managed devices (both USP and CWMP).

```javascript
{
  _id: ObjectId("507f1f77bcf86cd799439011"),
  endpoint_id: "os::012345-ABCDEFG",           // Unique device identifier
  agent_type: "USP",                           // USP, CWMP, or HYBRID
  manufacturer: "Acme Corp",
  model: "Router-X1",
  serial_number: "SN123456789",
  software_version: "1.2.3",
  hardware_version: "v1.0",
  
  // Connection information
  status: "online",                            // online, offline, unknown
  last_seen: ISODate("2023-12-01T10:30:00Z"),
  ip_address: "192.168.1.100",
  connection_info: {
    transport: "WebSocket",
    endpoint: "wss://device.example.com/usp",
    last_heartbeat: ISODate("2023-12-01T10:29:00Z")
  },
  
  // Protocol-specific fields
  usp_info: {
    supported_protocols: ["USP"],
    supported_operations: ["Get", "Set", "Add", "Delete"],
    max_message_size: 65536
  },
  
  cwmp_info: {
    oui: "ABCDEF",
    product_class: "Router",
    provisioning_code: "PROV123",
    connection_request_url: "http://192.168.1.100:7547/",
    connection_request_username: "admin",
    connection_request_password: "password",
    periodic_inform_enable: true,
    periodic_inform_interval: 300
  },
  
  // Metadata
  created_at: ISODate("2023-11-01T09:00:00Z"),
  updated_at: ISODate("2023-12-01T10:30:00Z"),
  tags: ["production", "router", "wifi"],
  
  // Indexing fields
  search_text: "os::012345-ABCDEFG Acme Corp Router-X1 SN123456789"
}
```

**Indexes:**
```javascript
db.agents.createIndex({ "endpoint_id": 1 }, { unique: true })
db.agents.createIndex({ "status": 1, "last_seen": -1 })
db.agents.createIndex({ "agent_type": 1 })
db.agents.createIndex({ "manufacturer": 1, "model": 1 })
db.agents.createIndex({ "search_text": "text" })
db.agents.createIndex({ "tags": 1 })
```

### parameters
Stores current parameter values for all devices.

```javascript
{
  _id: ObjectId("507f1f77bcf86cd799439012"),
  device_id: ObjectId("507f1f77bcf86cd799439011"), // Reference to agents collection
  path: "Device.WiFi.Radio.1.Enable",
  value: true,
  type: "boolean",
  writable: true,
  
  // Metadata
  last_updated: ISODate("2023-12-01T10:00:00Z"),
  source: "device",                           // device, controller, user
  operation_id: ObjectId("507f1f77bcf86cd799439013"),
  
  // Data model information
  data_model: "Device:2.14",
  description: "Enable or disable the radio",
  
  // Validation constraints
  constraints: {
    min_value: null,
    max_value: null,
    allowed_values: [true, false],
    pattern: null
  }
}
```

**Indexes:**
```javascript
db.parameters.createIndex({ "device_id": 1, "path": 1 }, { unique: true })
db.parameters.createIndex({ "path": 1 })
db.parameters.createIndex({ "last_updated": -1 })
db.parameters.createIndex({ "writable": 1 })
```

### instances
Stores object instances for multi-instance objects.

```javascript
{
  _id: ObjectId("507f1f77bcf86cd799439014"),
  device_id: ObjectId("507f1f77bcf86cd799439011"),
  object_path: "Device.WiFi.SSID.",
  instance_number: 1,
  instance_path: "Device.WiFi.SSID.1.",
  
  // Instance parameters
  parameters: {
    "Enable": true,
    "SSID": "MyNetwork",
    "BeaconType": "11i",
    "BasicAuthenticationMode": "WPA2-PSK",
    "WPAEncryptionModes": "AESEncryption",
    "IEEE11iEncryptionModes": "AESEncryption",
    "KeyPassphrase": "MyPassword"
  },
  
  // Metadata
  created_at: ISODate("2023-11-01T09:00:00Z"),
  updated_at: ISODate("2023-12-01T10:00:00Z"),
  created_by: "user",                         // user, controller, device
  operation_id: ObjectId("507f1f77bcf86cd799439015")
}
```

**Indexes:**
```javascript
db.instances.createIndex({ "device_id": 1, "instance_path": 1 }, { unique: true })
db.instances.createIndex({ "device_id": 1, "object_path": 1 })
db.instances.createIndex({ "instance_number": 1 })
```

### operations
Tracks all operations performed on devices.

```javascript
{
  _id: ObjectId("507f1f77bcf86cd799439016"),
  operation_id: "op-123456",                  // User-friendly ID
  device_id: ObjectId("507f1f77bcf86cd799439011"),
  type: "Set",                               // Get, Set, Add, Delete, Operate
  status: "completed",                       // pending, in_progress, completed, failed, cancelled
  
  // Operation details
  request: {
    paths: ["Device.WiFi.Radio.1.Enable"],
    parameters: [
      {
        path: "Device.WiFi.Radio.1.Enable",
        value: true
      }
    ]
  },
  
  response: {
    results: [
      {
        path: "Device.WiFi.Radio.1.Enable",
        status: "success",
        value: true,
        message: null
      }
    ]
  },
  
  // Error information
  error: {
    code: null,
    message: null,
    details: null
  },
  
  // Timing information
  created_at: ISODate("2023-12-01T10:00:00Z"),
  started_at: ISODate("2023-12-01T10:00:01Z"),
  completed_at: ISODate("2023-12-01T10:00:05Z"),
  timeout_at: ISODate("2023-12-01T10:05:00Z"),
  
  // Metadata
  initiated_by: ObjectId("507f1f77bcf86cd799439017"), // User ID
  source: "api",                             // api, scheduler, event
  priority: "normal",                        // low, normal, high, urgent
  retry_count: 0,
  max_retries: 3
}
```

**Indexes:**
```javascript
db.operations.createIndex({ "operation_id": 1 }, { unique: true })
db.operations.createIndex({ "device_id": 1, "created_at": -1 })
db.operations.createIndex({ "status": 1, "created_at": -1 })
db.operations.createIndex({ "type": 1 })
db.operations.createIndex({ "initiated_by": 1 })
```

### events
Stores device events and notifications.

```javascript
{
  _id: ObjectId("507f1f77bcf86cd799439018"),
  event_id: "evt-123456",
  device_id: ObjectId("507f1f77bcf86cd799439011"),
  
  // Event details
  event_type: "ValueChange",                 // ValueChange, ObjectCreation, Boot, etc.
  path: "Device.WiFi.Radio.1.Channel",
  old_value: 1,
  new_value: 6,
  
  // Additional event data
  event_data: {
    command_key: "change-channel-123",
    source: "user"
  },
  
  // Timing
  timestamp: ISODate("2023-12-01T10:00:00Z"),
  received_at: ISODate("2023-12-01T10:00:01Z"),
  
  // Processing status
  processed: true,
  notification_sent: true,
  subscriptions_matched: ["sub-123", "sub-456"],
  
  // Metadata
  severity: "info",                          // debug, info, warn, error, critical
  category: "configuration",                 // configuration, diagnostic, system, security
  tags: ["wifi", "channel", "user-initiated"]
}
```

**Indexes:**
```javascript
db.events.createIndex({ "device_id": 1, "timestamp": -1 })
db.events.createIndex({ "event_type": 1, "timestamp": -1 })
db.events.createIndex({ "path": 1 })
db.events.createIndex({ "timestamp": -1 }, { expireAfterSeconds: 2592000 }) // 30 days
```

### subscriptions
Manages event subscriptions.

```javascript
{
  _id: ObjectId("507f1f77bcf86cd799439019"),
  subscription_id: "sub-123456",
  device_id: ObjectId("507f1f77bcf86cd799439011"), // null for global subscriptions
  user_id: ObjectId("507f1f77bcf86cd799439020"),
  
  // Subscription criteria
  event_type: "ValueChange",
  path_pattern: "Device.WiFi.Radio.*.Channel",
  conditions: {
    severity: ["warn", "error", "critical"],
    tags: ["production"]
  },
  
  // Delivery configuration
  delivery: {
    method: "webhook",                       // webhook, email, websocket
    endpoint: "https://app.example.com/webhook",
    authentication: {
      type: "bearer",
      token: "webhook-token-123"
    },
    retry_policy: {
      max_attempts: 3,
      backoff: "exponential"
    }
  },
  
  // Status
  status: "active",                          // active, paused, error
  last_delivery: ISODate("2023-12-01T09:55:00Z"),
  delivery_count: 156,
  error_count: 2,
  last_error: {
    message: "Webhook endpoint unreachable",
    timestamp: ISODate("2023-12-01T08:30:00Z")
  },
  
  // Metadata
  created_at: ISODate("2023-11-01T09:00:00Z"),
  updated_at: ISODate("2023-12-01T10:00:00Z"),
  expires_at: ISODate("2024-12-01T09:00:00Z")
}
```

**Indexes:**
```javascript
db.subscriptions.createIndex({ "device_id": 1, "event_type": 1 })
db.subscriptions.createIndex({ "user_id": 1 })
db.subscriptions.createIndex({ "status": 1 })
db.subscriptions.createIndex({ "expires_at": 1 })
```

## Authentication & Authorization

### users
User accounts and authentication information.

```javascript
{
  _id: ObjectId("507f1f77bcf86cd799439021"),
  username: "admin@example.com",
  email: "admin@example.com",
  password_hash: "$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/NewPkfYXHNt0jYOz6",
  
  // Profile information
  profile: {
    first_name: "John",
    last_name: "Doe",
    display_name: "John Doe",
    avatar_url: "https://example.com/avatars/johndoe.png",
    timezone: "America/New_York",
    locale: "en_US"
  },
  
  // Authentication
  status: "active",                          // active, inactive, suspended
  email_verified: true,
  phone_verified: false,
  mfa_enabled: true,
  mfa_secret: "JBSWY3DPEHPK3PXP",
  
  // Authorization
  role: "admin",                             // admin, operator, viewer
  permissions: [
    "device:read", "device:write", "device:delete",
    "user:read", "user:write",
    "system:read", "system:write"
  ],
  
  // Device access control
  device_groups: ["production", "staging"],
  allowed_devices: [],                       // Empty means all devices
  
  // Session management
  active_sessions: ["sess-123", "sess-456"],
  last_login: ISODate("2023-12-01T09:00:00Z"),
  password_changed_at: ISODate("2023-11-15T14:30:00Z"),
  
  // Metadata
  created_at: ISODate("2023-10-01T10:00:00Z"),
  updated_at: ISODate("2023-12-01T09:00:00Z"),
  created_by: ObjectId("507f1f77bcf86cd799439022"),
  
  // API access
  api_keys: [
    {
      key_id: "key-123456",
      key_hash: "$2b$12$...",
      name: "Production API Key",
      permissions: ["device:read", "device:write"],
      expires_at: ISODate("2024-06-01T00:00:00Z"),
      created_at: ISODate("2023-06-01T00:00:00Z"),
      last_used: ISODate("2023-12-01T09:30:00Z")
    }
  ]
}
```

**Indexes:**
```javascript
db.users.createIndex({ "username": 1 }, { unique: true })
db.users.createIndex({ "email": 1 }, { unique: true })
db.users.createIndex({ "status": 1 })
db.users.createIndex({ "role": 1 })
```

### sessions
Active user and device sessions.

```javascript
{
  _id: ObjectId("507f1f77bcf86cd799439023"),
  session_id: "sess-123456",
  session_type: "user",                      // user, device, api
  
  // User sessions
  user_id: ObjectId("507f1f77bcf86cd799439021"),
  
  // Device sessions (CWMP/USP)
  device_id: ObjectId("507f1f77bcf86cd799439011"),
  
  // Session data
  status: "active",                          // active, expired, terminated
  ip_address: "203.0.113.1",
  user_agent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
  
  // Timing
  created_at: ISODate("2023-12-01T09:00:00Z"),
  last_activity: ISODate("2023-12-01T10:30:00Z"),
  expires_at: ISODate("2023-12-01T17:00:00Z"),
  
  // Security
  csrf_token: "csrf-token-123",
  jwt_token_hash: "$2b$12$...",
  
  // Device session specific
  pending_operations: [
    ObjectId("507f1f77bcf86cd799439016")
  ],
  max_envelopes: 1,
  connection_info: {
    remote_addr: "192.168.1.100:12345",
    transport: "HTTP",
    tls_version: "TLSv1.3",
    cipher_suite: "TLS_AES_256_GCM_SHA384"
  }
}
```

**Indexes:**
```javascript
db.sessions.createIndex({ "session_id": 1 }, { unique: true })
db.sessions.createIndex({ "user_id": 1 })
db.sessions.createIndex({ "device_id": 1 })
db.sessions.createIndex({ "expires_at": 1 })
```

## Data Models & Schema

### datamodels
Data model definitions and schemas.

```javascript
{
  _id: ObjectId("507f1f77bcf86cd799439024"),
  name: "Device:2.14",
  version: "2.14",
  type: "USP",                               // USP, CWMP, Custom
  
  // Source information
  source_url: "https://usp-data-models.broadband-forum.org/tr-181-2-14-0-usp.xml",
  checksum: "sha256:abc123...",
  
  // Schema definition
  objects: [
    {
      name: "Device.",
      description: "Top-level object for Device data model",
      parameters: [
        {
          name: "DeviceInfo",
          type: "object",
          description: "Device information object",
          access: "readOnly"
        }
      ],
      children: [
        {
          name: "DeviceInfo.",
          description: "Device information",
          parameters: [
            {
              name: "Manufacturer",
              type: "string",
              access: "readOnly",
              description: "Device manufacturer",
              max_length: 64
            }
          ]
        }
      ]
    }
  ],
  
  // Statistics
  total_objects: 156,
  total_parameters: 2344,
  
  // Metadata
  imported_at: ISODate("2023-11-01T10:00:00Z"),
  status: "active",                          // active, deprecated, draft
  supported_by: ["USP-1.2", "USP-1.3"]
}
```

**Indexes:**
```javascript
db.datamodels.createIndex({ "name": 1, "version": 1 }, { unique: true })
db.datamodels.createIndex({ "type": 1 })
db.datamodels.createIndex({ "status": 1 })
```

## Audit & Logging

### audit_logs
System audit trail for compliance and debugging.

```javascript
{
  _id: ObjectId("507f1f77bcf86cd799439025"),
  timestamp: ISODate("2023-12-01T10:00:00Z"),
  
  // Action information
  action: "parameter_set",
  resource_type: "parameter",
  resource_id: "Device.WiFi.Radio.1.Enable",
  device_id: ObjectId("507f1f77bcf86cd799439011"),
  
  // Actor information
  actor_type: "user",                        // user, system, device
  actor_id: ObjectId("507f1f77bcf86cd799439021"),
  session_id: "sess-123456",
  
  // Request context
  ip_address: "203.0.113.1",
  user_agent: "OpenUSP Web Dashboard/1.0",
  request_id: "req-123456",
  
  // Change details
  changes: {
    old_value: false,
    new_value: true,
    operation_id: "op-123456"
  },
  
  // Result
  status: "success",                         // success, failure, partial
  error_message: null,
  
  // Metadata
  severity: "info",                          // debug, info, warn, error, critical
  category: "data_change",                   // auth, data_change, system, security
  tags: ["wifi", "radio", "configuration"]
}
```

**Indexes:**
```javascript
db.audit_logs.createIndex({ "timestamp": -1 })
db.audit_logs.createIndex({ "device_id": 1, "timestamp": -1 })
db.audit_logs.createIndex({ "actor_id": 1, "timestamp": -1 })
db.audit_logs.createIndex({ "action": 1 })
db.audit_logs.createIndex({ "timestamp": -1 }, { expireAfterSeconds: 7776000 }) // 90 days
```

## Database Operations

### Connection Management
```go
type DatabaseConfig struct {
    URI             string        `yaml:"uri"`
    Database        string        `yaml:"database"`
    MaxPoolSize     uint64        `yaml:"max_pool_size"`
    MinPoolSize     uint64        `yaml:"min_pool_size"`
    MaxConnIdleTime time.Duration `yaml:"max_conn_idle_time"`
    ConnectTimeout  time.Duration `yaml:"connect_timeout"`
    SocketTimeout   time.Duration `yaml:"socket_timeout"`
    TLS             TLSConfig     `yaml:"tls"`
}

func NewDatabase(config *DatabaseConfig) (*Database, error) {
    clientOptions := options.Client().ApplyURI(config.URI)
    clientOptions.SetMaxPoolSize(config.MaxPoolSize)
    clientOptions.SetMinPoolSize(config.MinPoolSize)
    clientOptions.SetMaxConnIdleTime(config.MaxConnIdleTime)
    clientOptions.SetConnectTimeout(config.ConnectTimeout)
    clientOptions.SetSocketTimeout(config.SocketTimeout)
    
    if config.TLS.Enabled {
        tlsConfig := &tls.Config{
            ServerName: config.TLS.ServerName,
        }
        clientOptions.SetTLSConfig(tlsConfig)
    }
    
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
    }
    
    // Test the connection
    err = client.Ping(context.Background(), readpref.Primary())
    if err != nil {
        return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
    }
    
    return &Database{
        client: client,
        db:     client.Database(config.Database),
    }, nil
}
```

### Query Patterns
```go
// Agent queries
func (db *Database) FindAgents(filter AgentFilter) ([]*Agent, error) {
    query := bson.M{}
    
    if filter.Status != "" {
        query["status"] = filter.Status
    }
    
    if filter.AgentType != "" {
        query["agent_type"] = filter.AgentType
    }
    
    if filter.Search != "" {
        query["$text"] = bson.M{"$search": filter.Search}
    }
    
    if len(filter.Tags) > 0 {
        query["tags"] = bson.M{"$in": filter.Tags}
    }
    
    cursor, err := db.agents.Find(context.Background(), query, 
        options.Find().
            SetSort(bson.D{{"last_seen", -1}}).
            SetSkip(int64(filter.Offset)).
            SetLimit(int64(filter.Limit)))
    
    if err != nil {
        return nil, err
    }
    
    var agents []*Agent
    err = cursor.All(context.Background(), &agents)
    return agents, err
}

// Parameter aggregation
func (db *Database) GetParameterStatistics() (*ParameterStats, error) {
    pipeline := mongo.Pipeline{
        {{"$group", bson.D{
            {"_id", "$type"},
            {"count", bson.D{{"$sum", 1}}},
            {"writable_count", bson.D{{"$sum", bson.D{{"$cond", []interface{}{"$writable", 1, 0}}}}}},
        }}},
        {{"$sort", bson.D{{"count", -1}}}},
    }
    
    cursor, err := db.parameters.Aggregate(context.Background(), pipeline)
    if err != nil {
        return nil, err
    }
    
    var results []bson.M
    err = cursor.All(context.Background(), &results)
    if err != nil {
        return nil, err
    }
    
    stats := &ParameterStats{
        ByType: make(map[string]TypeStats),
    }
    
    for _, result := range results {
        paramType := result["_id"].(string)
        stats.ByType[paramType] = TypeStats{
            Count:         int(result["count"].(int32)),
            WritableCount: int(result["writable_count"].(int32)),
        }
    }
    
    return stats, nil
}
```

### Backup and Maintenance
```bash
#!/bin/bash
# Database backup script
BACKUP_DIR="/backup/mongodb"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_PATH="$BACKUP_DIR/$DATE"

# Create backup directory
mkdir -p "$BACKUP_PATH"

# Dump database
mongodump --uri="$MONGODB_URI" --out="$BACKUP_PATH"

# Compress backup
tar -czf "$BACKUP_PATH.tar.gz" -C "$BACKUP_DIR" "$DATE"
rm -rf "$BACKUP_PATH"

# Cleanup old backups (keep 30 days)
find "$BACKUP_DIR" -name "*.tar.gz" -mtime +30 -delete

echo "Backup completed: $BACKUP_PATH.tar.gz"
```

This database schema documentation provides a comprehensive overview of the data structures and patterns used throughout the OpenUSP platform.